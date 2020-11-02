package seis

import (
	"bufio"
	"bytes"
	"fmt"
	//	"math"
	"strconv"
	"strings"
	"time"
)

type MSEncoding uint8

const (
	EncodingASCII      MSEncoding = 0
	EncodingInt16      MSEncoding = 1
	EncodingInt24      MSEncoding = 2
	EncodingInt32      MSEncoding = 3
	EncodingIEEEFloat  MSEncoding = 4
	EncodingIEEEDouble MSEncoding = 5

	EncodingSTEIM1            MSEncoding = 10
	EncodingSTEIM2            MSEncoding = 11
	EncodingGEOSCOPE24bit     MSEncoding = 12
	EncodingGEOSCOPE16bit3exp MSEncoding = 13
	EncodingGEOSCOPE16bit4exp MSEncoding = 14
	EncodingUSNN              MSEncoding = 15
	EncodingCDSN              MSEncoding = 16
	EncodingGraefenberg       MSEncoding = 17
	EncodingIPG               MSEncoding = 18
	EncodingSTEIM3            MSEncoding = 19

	EncodingSRO    MSEncoding = 30
	EncodingHGLP   MSEncoding = 31
	EncodingDWWSSN MSEncoding = 32
	EncodingRSTN   MSEncoding = 33
)

type WordOrder uint8

const (
	LittleEndian WordOrder = 0
	BigEndian    WordOrder = 1
)

type MSSampleType byte

const (
	UnknownType MSSampleType = 0
	ByteType    MSSampleType = 'a'
	IntegerType MSSampleType = 'i'
	FloatType   MSSampleType = 'f'
	DoubleType  MSSampleType = 'd'
)

// MSRecord represents the raw miniseed record.
type MSRecord struct {
	Header MSDataHeader
	B1000  Blockette1000 //If Present
	B1001  Blockette1001 //If Present
	Data   []byte

	stats   msRecordStats
	samples msRecordSamples
}

// msRecordStats holds decoded miniseed properties.
type msRecordStats struct {
	DataFlag   uint8
	B1001Flag  uint8
	SampleType byte
}

// msRecordSamples holds the decoded miniseed data when
// using the Unpack function with the decode data samples
// option is true.
type msRecordSamples struct {
	Int32s   []int32
	Float32s []float32
	Float64s []float64
}

// NewMSRecord decodes and unpacks the record samples from a byte slice and
// returns an MSRecord pointer, or an empty pointer and an error if it
// could not be decoded.
func NewMSRecord(buf []byte) (*MSRecord, error) {

	var ms MSRecord
	if err := ms.Unpack(buf, true); err != nil {
		return nil, err
	}

	return &ms, nil
}

// String implements the Stringer interface and provides a short summary of the
// miniseed record header.
func (m *MSRecord) String() string {
	var parts []string

	parts = append(parts, m.SrcName(false))
	parts = append(parts, fmt.Sprintf("%06d", m.SequenceNumber()))
	parts = append(parts, string(m.DataQuality()))
	parts = append(parts, fmt.Sprintf("%d", m.BlockSize()))
	parts = append(parts, fmt.Sprintf("%d samples", m.SampleCount()))
	parts = append(parts, fmt.Sprintf("%g Hz", m.SampleRate()))
	parts = append(parts, m.StartTime().Format("2006,002,15:04:05.000000"))

	return strings.Join(parts, ", ")
}

/** TODO: investigate best option for packing if required -- dependent on block size etc.
// Pack the record into a byte slice using the given encoding format.
func (m *MSRecord) Pack(encoding int) ([]byte, error) {
	//var err error
	var output []byte

	//encode the data so we can put its details in the header
	if m.Stats.DataFlag != 0 {
		if encoding >= 0 { //If encoding is set use the specified encoding
			m.B1000.Encoding = uint8(encoding)
		}

		switch m.B1000.Encoding {
		case EncodingASCII:
			if m.Stats.SampleType != 'a' {
				return output, fmt.Errorf("pack: EncodingASCII is only available for text data")
			}
		case EncodingSTEIM1:
			if m.Stats.SampleType == 'i' {
				steimOutput, err := encodeSteim(1, m.Data)
				if err != nil {
					return output, fmt.Errorf("pack: encoding failed: %v", err)
				}
				m.B1001.FrameCount = steimOutput.FrameCount //TODO: Should we set the B1001Flag here? Will value in a 0 clock quality value if that hasn't been set already
				m.B1000.WordOrder = steimOutput.WordOrder
				m.RawData = steimOutput.EncodedData
			} else {
				return output, fmt.Errorf("pack: EncodingSTEIM1 is only available for integer data")
			}
		case EncodingSTEIM2:
			if m.Stats.SampleType == 'i' {
				steimOutput, err := encodeSteim(2, m.Data)
				if err != nil {
					return output, fmt.Errorf("pack: encoding failed: %v", err)
				}
				m.B1001.FrameCount = steimOutput.FrameCount //TODO: Should we set the B1001Flag here? Will value in a 0 clock quality value if that hasn't been set already
				m.B1000.WordOrder = steimOutput.WordOrder
				m.RawData = steimOutput.EncodedData
			} else {
				return output, fmt.Errorf("pack: EncodingSTEIM2 is only available for integer data")
			}
		default:
			return output, fmt.Errorf("pack: data encoding %v is not supported", m.B1000.Encoding)
		}
	}

	m.Header.BlockettesToFollow = 1 //Blockette1000 is required

	//calculate the size of our record so we appropriately size output
	outputSize := MSDataHeaderSize + BlocketteHeaderSize + Blockette1000Size
	if m.Stats.B1001Flag != 0 {
		m.Header.BlockettesToFollow++
		outputSize += BlocketteHeaderSize + Blockette1001Size
	}

	m.Header.BeginningOfData = uint16(outputSize) //The byte location of the actual data
	outputSize += len(m.RawData)                  //The final size of the packed record

	if outputSize > MSeedSize {
		return output, fmt.Errorf("packing this record would result in a %v byte record", outputSize)
	}

	output = make([]byte, MSeedSize)

	//write the header
	pointer := uint16(MSDataHeaderSize)
	bmsh := MarshalMSDataHeader(m.Header)
	copy(output[:pointer], bmsh[:])

	//set the location of the first blockette (directly after the header)
	m.Header.FirstBlockette = pointer

	blockettesToWrite := m.Header.BlockettesToFollow - 1 //Keep track so we know which blockette is the last
	var macroMore = func(i uint8) uint16 {
		if i > 0 {
			return 1
		}
		return 0
	} //A macro to determine if we need to set the NextBlockette value

	//write the 1000 blockette
	bh := BlocketteHeader{
		BlocketteType: 1000,
		NextBlockette: (pointer + BlocketteHeaderSize + Blockette1000Size) * macroMore(blockettesToWrite),
	}
	bbh := MarshalBlocketteHeader(bh)
	copy(output[pointer:pointer+BlocketteHeaderSize], bbh[:])
	pointer += BlocketteHeaderSize

	bb1000 := MarshalBlockette1000(m.B1000)
	copy(output[pointer:pointer+Blockette1000Size], bb1000[:])
	pointer += Blockette1000Size
	blockettesToWrite--

	//write 'optional' blockettes
	//TODO: The wel2000.mseed test file does not have a blockette 1001 but leaves blank bytes where it would be, is it ok to not include it at all?
	if m.Stats.B1001Flag != 0 {
		bh = BlocketteHeader{
			BlocketteType: 1001,
			NextBlockette: (pointer + BlocketteHeaderSize + Blockette1000Size) * macroMore(blockettesToWrite),
		}
		bbh = MarshalBlocketteHeader(bh)
		copy(output[pointer:pointer+BlocketteHeaderSize], bbh[:])
		pointer += BlocketteHeaderSize

		bb1001 := MarshalBlockette1001(m.B1001)
		copy(output[pointer:pointer+Blockette1001Size], bb1001[:])
		pointer += Blockette1001Size
		blockettesToWrite--
	}

	//write the encoded data
	copy(output[pointer:pointer+uint16(len(m.RawData))], m.RawData)
	pointer += uint16(len(m.RawData))

	return output, nil
}
**/

// Unpack the record form a byte slice, the dataflag can be used to suppress decoding the waveform data
// for efficency if only the header information is required.
func (m *MSRecord) Unpack(buf []byte, dataflag bool) error {
	var err error

	//Unpack The Fixed Header
	if len(buf) < MSDataHeaderSize {
		return fmt.Errorf("unpack: given %v bytes; not enough to parse header", len(buf))
	}
	var h [MSDataHeaderSize]byte
	copy(h[:], buf[:MSDataHeaderSize])
	m.Header = UnmarshalMSDataHeader(h)
	if !isValidMSHeader(m.Header) {
		return fmt.Errorf("unpack: input is not a valid MSEED record: incorrect header")
	}

	//Unpack Blockettes
	pointer := m.Header.FirstBlockette //TODO: This could be replaced with bytes.Reader()
	for i := 0; i < int(m.Header.BlockettesToFollow); i++ {
		if pointer == 0 {
			return fmt.Errorf(
				"unpack: next blockette pointer == 0 after %v blockettes but BlockettesToFollow = %v",
				i-1, m.Header.BlockettesToFollow)
		}

		//Get the blockette header
		if len(buf) < int(pointer+BlocketteHeaderSize) {
			return fmt.Errorf("unpack: given %v bytes; not enough to parse blockette header at %v", len(buf), pointer)
		}
		var h [BlocketteHeaderSize]byte
		copy(h[:], buf[pointer:pointer+BlocketteHeaderSize])
		bhead := UnmarshalBlocketteHeader(h)
		bpointer := pointer + BlocketteHeaderSize //start of blockette content

		switch bhead.BlocketteType {
		case 1000:
			if len(buf) < int(bpointer+Blockette1000Size) {
				return fmt.Errorf("unpack: given %v bytes; not enough to parse blockette 1000 at %v", len(buf), bpointer)
			}
			var h [Blockette1000Size]byte
			copy(h[:], buf[bpointer:bpointer+Blockette1000Size])
			m.B1000 = UnmarshalBlockette1000(h)
		case 1001:
			if len(buf) < int(bpointer+Blockette1001Size) {
				return fmt.Errorf("unpack: given %v bytes; not enough to parse blockette 1001 at %v", len(buf), bpointer)
			}
			var h [Blockette1001Size]byte
			copy(h[:], buf[bpointer:bpointer+Blockette1001Size])
			m.B1001 = UnmarshalBlockette1001(h)
			m.stats.B1001Flag = 1
		default:
			return fmt.Errorf("unpack: unsupported blockette type: %v", err)
		}

		pointer = bhead.NextBlockette
	}

	pointer = m.Header.BeginningOfData
	m.Data = buf[pointer:]
	if dataflag {
		switch MSEncoding(m.B1000.Encoding) {
		case EncodingASCII:
			m.stats.SampleType = 'a'
			m.stats.DataFlag = 1
		case EncodingInt32:
			m.stats.SampleType = 'i'
			m.stats.DataFlag = 1
			m.samples.Int32s, err = decodeInt32(m.Data, m.B1000.WordOrder, m.Header.SampleCount)
			if err != nil {
				return fmt.Errorf("unpack: %v", err)
			}
		case EncodingIEEEFloat:
			m.stats.SampleType = 'f'
			m.stats.DataFlag = 1
			m.samples.Float32s, err = decodeFloat32(m.Data, m.B1000.WordOrder, m.Header.SampleCount)
			if err != nil {
				return fmt.Errorf("unpack: %v", err)
			}
		case EncodingIEEEDouble:
			m.stats.SampleType = 'd'
			m.stats.DataFlag = 1
			m.samples.Float64s, err = decodeFloat64(m.Data, m.B1000.WordOrder, m.Header.SampleCount)
			if err != nil {
				return fmt.Errorf("unpack: %v", err)
			}
		case EncodingSTEIM1:
			framecount := uint8((len(buf) - int(pointer)) / 64)
			if m.B1001.FrameCount != 0 {
				framecount = m.B1001.FrameCount
			}
			if int(framecount)*64 > (len(buf) - int(pointer)) { //make sure the decoding doesn't overrun the buffer
				return fmt.Errorf("unpack: header reported more bytes then are present in data packet: %v > %v", framecount*64, len(buf)-int(pointer))
			}
			m.stats.SampleType = 'i'
			m.samples.Int32s, err = decodeSteim(1, m.Data, m.B1000.WordOrder, framecount, m.Header.SampleCount)
			if err != nil {
				return fmt.Errorf("unpack: %v", err)
			}

			m.stats.DataFlag = 1
			if len(m.samples.Int32s) != int(m.Header.SampleCount) {
				return fmt.Errorf("unpack: expected %v samples, decoding returned %v", m.Header.SampleCount, len(m.samples.Int32s))
			}
		case EncodingSTEIM2: //STEIM2
			framecount := uint8((len(buf) - int(pointer)) / 64)
			if m.B1001.FrameCount != 0 {
				framecount = m.B1001.FrameCount
			}
			if int(framecount)*64 > (len(buf) - int(pointer)) { //make sure the decoding doesn't overrun the buffer
				return fmt.Errorf("unpack: header reported more bytes then are present in data packet: %v > %v", framecount*64, len(buf)-int(pointer))
			}
			m.stats.SampleType = 'i'
			m.samples.Int32s, err = decodeSteim(2, m.Data, m.B1000.WordOrder, framecount, m.Header.SampleCount)
			if err != nil {
				return fmt.Errorf("unpack: %v", err)
			}

			m.stats.DataFlag = 1
			if len(m.samples.Int32s) != int(m.Header.SampleCount) {
				return fmt.Errorf("unpack: expected %v samples, decoding returned %v", m.Header.SampleCount, len(m.samples.Int32s))
			}
		default:
			return fmt.Errorf("unpack: data encoding %v is not supported", m.B1000.Encoding)
		}
	}

	return nil
}

// SequenceNumber returns the record sequence number.
func (m *MSRecord) SequenceNumber() int {
	s := string(m.Header.SeqNumber[:]) //Convert the byte array to a slice so we can convert to a string
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0 //TODO: Maybe better error handling, we're swallowing an error here
	}
	return int(i)
}

// Network returns the cleaned miniseed header network string.
func (m *MSRecord) Network() string {
	return strings.TrimSpace(string(m.Header.Network[:]))
}

// Station returns the cleaned miniseed header station string.
func (m *MSRecord) Station() string {
	return strings.TrimSpace(string(m.Header.Station[:]))
}

// Location returns the cleaned miniseed header location string.
func (m *MSRecord) Location() string {
	return strings.TrimSpace(string(m.Header.Location[:]))
}

// Channel returns the cleaned miniseed header channel string.
func (m *MSRecord) Channel() string {
	return strings.TrimSpace(string(m.Header.Channel[:]))
}

// DataQuality returns the miniseed header data quality flag.
func (m *MSRecord) DataQuality() byte {
	return m.Header.DataQuality
}

// StartTime returns the corrected header start time.
func (m *MSRecord) StartTime() time.Time {
	st := time.Date(int(m.Header.StartTime.Year),
		1,
		int(m.Header.StartTime.Doy),
		int(m.Header.StartTime.Hour),
		int(m.Header.StartTime.Minute),
		int(m.Header.StartTime.Second),
		int(m.Header.StartTime.S0001)*100000,
		time.UTC)

	//Check Time Correction
	if !flagCheck(m.Header.ActivityFlags, 1) { //bit 1 indicates whether time correction has been applied
		st = st.Add(time.Microsecond * 100 * time.Duration(m.Header.TimeCorrection))
	}

	//Get Enhanced Timing From Blockette 1001
	if m.B1001.MicroSec != 0 {
		st = st.Add(time.Microsecond * time.Duration(m.B1001.MicroSec))
	}

	return st
}

// extact the sampling rate from the two seed factors.
func sampleRate(factor, multiplier int) float64 {
	var samprate float64

	switch f := float64(factor); {
	case factor > 0:
		samprate = f
	case factor < 0:
		samprate = -1.0 / f
	}

	switch m := float64(multiplier); {
	case multiplier > 0:
		samprate = samprate * m
	case multiplier < 0:
		samprate = -1.0 * (samprate / m)
	}

	return samprate
}

// flip the factors to get the sample interval
func samplePeriod(factor, multiplier int) time.Duration {
	if sps := sampleRate(factor, multiplier); sps > 0.0 {
		return time.Duration(float64(time.Second) / sps)
	}
	return 0
}

// only valid for divisible sample rates in Hz or intervals in seconds, return 0 otherwise.
/**
func sampleFactor(sps float64) int {
	switch {
	case (sps - math.Floor(sps)) < 0.000001:
		return int(math.Floor(sps))
	case (1.0/sps - math.Floor(1.0/sps)) < 0.000001:
		return int(math.Floor(-1.0 / sps))
	default:
		return 0
	}
}
**/

// SampleRate returns the decoded header sampling rate in samples per second.
func (m *MSRecord) SampleRate() float64 {
	return sampleRate(int(m.Header.SampleRateFactor), int(m.Header.SampleRateMultiplier))
}

// SampleCount returns the number of samples in the record, independent of whether they are decoded or not.
func (m *MSRecord) SampleCount() int {
	return int(m.Header.SampleCount)
}

// Encoding returns the miniseed data format encoding.
func (m *MSRecord) Encoding() MSEncoding {
	return MSEncoding(m.B1000.Encoding)
}

// ByteOrder returns the miniseed data byte order.
func (m *MSRecord) ByteOrder() WordOrder {
	switch m.B1000.WordOrder {
	case 0:
		return LittleEndian
	default:
		return BigEndian
	}
}

func trimRight(data []byte) []byte {
	return bytes.TrimRight(data, "\x00")
}

// NumSamples returns the number of samples decoded from the miniseed record.
func (m *MSRecord) NumSamples() int {
	if m.stats.DataFlag != 0 {
		switch m.stats.SampleType {
		case 'a':
			return len(trimRight(m.Data))
		case 'i':
			return len(m.samples.Int32s)
		case 'f':
			return len(m.samples.Float32s)
		case 'd':
			return len(m.samples.Float64s)
		}
	}
	return 0
}

// SampleType returns the type of samples decoded, or UnknownType if no data has been decoded.
func (m *MSRecord) SampleType() MSSampleType {
	if m.stats.DataFlag != 0 {
		return MSSampleType(m.stats.SampleType)
	}
	return UnknownType
}

// EndTime returns the calculated time of the last sample.
func (m *MSRecord) EndTime() time.Time { //TODO: Handle leap seconds?
	var d time.Duration
	sc := m.Header.SampleCount
	sr := m.SampleRate()

	if sc > 0 && sr > 0 {
		//Endtime = The number of samples (-1) * the duration of each sample (sampRate)
		d = time.Duration(m.Header.SampleCount-1) * time.Duration(float64(time.Second)/m.SampleRate()+0.5)
		//The +0.5 causes the conversion to be a "rounding" conversion rather than truncation
	}

	return m.StartTime().Add(d)
}

// SamplePeriod converts the sample rate into a time interval.
// For invalid sampling rates a zero duration is returned.
func (m *MSRecord) SamplePeriod() time.Duration {
	return samplePeriod(int(m.Header.SampleRateFactor), int(m.Header.SampleRateMultiplier))
}

// PacketSize returns the length of the packet
func (m *MSRecord) BlockSize() int {
	if n := int(m.B1000.RecordLength); n > 0 {
		return 1 << n
	}
	return 0
}

// ToStrings returns an ASCII encodied record as a slice of line strings.
func (m *MSRecord) ToStrings() []string {
	switch t, e := m.SampleType(), m.Encoding(); {
	case t == ByteType && e == EncodingASCII:
		var lines []string
		buf := bytes.NewBuffer(trimRight(m.Data[:m.Header.SampleCount]))
		scanner := bufio.NewScanner(buf)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return nil
		}
		return lines
	default:
		return nil
	}
}

// ToBytes returns a cleaned byte slice of the expected ASCII record.
func (m *MSRecord) ToBytes() []byte {
	switch t, e := m.SampleType(), m.Encoding(); {
	case t == ByteType && e == EncodingASCII:
		return trimRight(m.Data[:m.SampleCount()])
	default:
		return nil
	}
}

// ToInts returns the record numerical data converted to an int slice.
func (m *MSRecord) ToInts() []int {
	var ints []int
	switch t, e := m.SampleType(), m.Encoding(); {
	case t == IntegerType && e == EncodingSTEIM1:
		for _, i := range m.samples.Int32s {
			ints = append(ints, int(i))
		}
		return ints
	case t == IntegerType && e == EncodingSTEIM2:
		for _, i := range m.samples.Int32s {
			ints = append(ints, int(i))
		}
		return ints
	case t == FloatType && e == EncodingIEEEFloat:
		for _, i := range m.samples.Float32s {
			ints = append(ints, int(i))
		}
		return ints
	case t == DoubleType && e == EncodingIEEEDouble:
		for _, i := range m.samples.Float64s {
			ints = append(ints, int(i))
		}
	}
	return ints
}

func (m *MSRecord) ToInt32s() []int32 {
	var ints []int32
	for _, i := range m.ToInts() {
		ints = append(ints, int32(i))
	}
	return ints
}

// ToFloat64s returns the record numerical data converted to an float64 slice.
func (m *MSRecord) ToFloat64s() []float64 {
	var floats []float64
	switch t, e := m.SampleType(), m.Encoding(); {
	case t == IntegerType && e == EncodingSTEIM1:
		for _, i := range m.samples.Int32s {
			floats = append(floats, float64(i))
		}
	case t == IntegerType && e == EncodingSTEIM2:
		for _, i := range m.samples.Int32s {
			floats = append(floats, float64(i))
		}
	case t == FloatType && e == EncodingIEEEFloat:
		for _, i := range m.samples.Float32s {
			floats = append(floats, float64(i))
		}
	case t == DoubleType && e == EncodingIEEEDouble:
		return append([]float64{}, m.samples.Float64s...)
	}
	return floats
}

// SrcName returns the standard stream representation of the packet header.
func (m *MSRecord) SrcName(quality bool) string {
	if quality {
		return strings.Join([]string{m.Network(), m.Station(), m.Location(), m.Channel(), string(m.DataQuality())}, "_")
	}
	return strings.Join([]string{m.Network(), m.Station(), m.Location(), m.Channel()}, "_")
}

func isValidMSHeader(mh MSDataHeader) bool { //from libmseed.h MS_ISVALIDHEADER
	for _, b := range mh.SeqNumber {
		if !((b >= '0' && b <= '9') || (b == ' ') || (b == 0)) {
			return false
		}
	}
	if !(mh.DataQuality == 'D' || mh.DataQuality == 'R' || mh.DataQuality == 'M' || mh.DataQuality == 'Q') {
		return false
	}
	if !(mh.Reserved == ' ' || mh.Reserved == 0) {
		return false
	}
	return mh.StartTime.Hour <= 23 && mh.StartTime.Minute <= 59 && mh.StartTime.Second <= 60 //TODO: Why is this 60? Leap Second?
}

func flagCheck(b byte, index uint8) bool {
	if index > 8 {
		return false
	}
	if b&(0x1<<index) == 0 {
		return false
	}
	return true
}
