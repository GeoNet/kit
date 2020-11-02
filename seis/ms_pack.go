package seis

import (
	//	"encoding/binary"
	"fmt"
	//	"strconv"
	"time"
)

// MSRecordStream is used as a template to encode miniseed blocks.
type MSRecordStream struct {
	// stream details
	Station     string
	Location    string
	Channel     string
	Network     string
	DataQuality byte

	// sampling rate
	Factor     int
	Multiplier int

	// block details
	WordOrder    WordOrder
	RecordLength int

	// running details
	SeqNumber int
}

// samplePeriod returns the expected sample interval for the given stream.
func (m MSRecordStream) samplePeriod() time.Duration {
	return samplePeriod(int(m.Factor), int(m.Multiplier))
}

// blockLength returns the number of bytes expected for the full record.
func (m MSRecordStream) blockLength() int {
	return 1 << m.RecordLength
}

// beginningOfData returns the offset to the data for the given block.
func (m *MSRecordStream) beginningOfData() int {
	return MSDataHeaderSize + 2*BlocketteHeaderSize + Blockette1000Size + Blockette1001Size
}

// dataLength returns the available space for data in the block.
func (m *MSRecordStream) dataLength() int {
	return m.blockLength() - m.beginningOfData()
}

// NewMSRecordStream builds a MSRecordStream pointer from an MSRecord pointer.
func NewMSRecordStream(msr *MSRecord) *MSRecordStream {
	return &MSRecordStream{
		Station:     string(msr.Header.Station[:]),
		Location:    string(msr.Header.Location[:]),
		Channel:     string(msr.Header.Channel[:]),
		Network:     string(msr.Header.Network[:]),
		DataQuality: msr.Header.DataQuality,

		Factor:       int(msr.Header.SampleRateFactor),
		Multiplier:   int(msr.Header.SampleRateMultiplier),
		RecordLength: int(msr.B1000.RecordLength),
		WordOrder:    WordOrder(msr.B1000.WordOrder),
	}
}

// MSRecordFunc is used to process packed miniseed blocks.
type MSRecordFunc func(*MSRecord, []byte, bool) error

// MSRecordTrace holds the running encoded data needed together with a MSRecordStream to
// build a miniseed block.
type msRecordTrace struct {
	// sequence details
	SeqNumber int

	// timing information
	StartTime     time.Time
	TimingQuality int
	ClockUnlocked bool

	// sample details
	SampleCount   int
	ResidualCount int
	SampleData    []byte

	// encoding details
	Encoding MSEncoding

	// record details
	FrameCount int
}

// Last returns whether the current MSRecordTrace has no residual counts.
func (r msRecordTrace) Last() bool {
	return !(r.ResidualCount > 0)
}

// msRecordFunc is used to convert raw data into a slice of MSRecordTrace.
type msRecordFunc func(*MSRecordStream) ([]msRecordTrace, error)

func (m *MSRecordStream) PackInt32(start time.Time, quality int, locked bool, data []int32, fn MSRecordFunc) error {
	return m.pack(encodeStream(start, quality, locked, encodeInt32(data, m.dataLength(), m.WordOrder)), fn)
}

// Pack uses an encoding function and an output function to build miniseed data blocks.
func (m *MSRecordStream) pack(recFn msRecordFunc, packFn MSRecordFunc) error {
	recs, err := recFn(m)
	if err != nil {
		return err
	}

	for _, rec := range recs {
		maxLength := 1 << m.RecordLength

		header := MSDataHeader{
			DataQuality:          m.DataQuality,
			StartTime:            NewBTime(rec.StartTime),
			SampleCount:          uint16(rec.SampleCount), // Number of Samples in the data block which may or may not be unpacked.
			SampleRateFactor:     int16(m.Factor),         //int16  // >0: Samples/Second <0: Second/Samples 0: Seconds/Sample, ASCII/OPAQUE DATA records
			SampleRateMultiplier: int16(m.Multiplier),     // int16  // >0: Multiplication Factor <0: Division Factor

			//TODO: handle finer details - e.g. locked
			//flags are bit flags
			//ActivityFlags    byte
			//IOClockFlags     byte
			//DataQualityFlags byte
			//TimeCorrection     int32 // 0.0001 second units //

			BlockettesToFollow: 2,
			BeginningOfData:    uint16(MSDataHeaderSize + 2*BlocketteHeaderSize + Blockette1000Size + Blockette1001Size), FirstBlockette: uint16(MSDataHeaderSize),
		}

		// update the header arrays
		copy(header.SeqNumber[:], []byte(fmt.Sprintf("%06d", rec.SeqNumber)))
		copy(header.Station[:], []byte(m.Station))
		copy(header.Location[:], []byte(m.Location))
		copy(header.Channel[:], []byte(m.Channel))
		copy(header.Network[:], []byte(m.Network))

		b1000 := Blockette1000{
			Encoding:     uint8(rec.Encoding),
			WordOrder:    uint8(m.WordOrder),
			RecordLength: uint8(m.RecordLength),
		}
		b1001 := Blockette1001{
			TimingQuality: uint8(rec.TimingQuality),
			//TODO: update the timing code.
			//MicroSec  :     int8 //Increased accuracy for starttime
			FrameCount: uint8(rec.FrameCount),
		}

		// pack samples into a byte slice using the encoded header and blockettes
		data, err := packRecord(header, b1000, b1001, rec.SampleData, maxLength)
		if err != nil {
			return err
		}

		// no output function?
		if packFn == nil {
			continue
		}

		// unpack as a check ...
		msr, err := NewMSRecord(data)
		if err != nil {
			return err
		}

		// as well as being passed into the packing function.
		if err := packFn(msr, data, rec.Last()); err != nil {
			return err
		}
	}

	return nil
}

func encodeStream(start time.Time, quality int, locked bool, fn encodingFunc) msRecordFunc {
	return func(m *MSRecordStream) ([]msRecordTrace, error) {
		var res []msRecordTrace

		var count int

		seqno := m.SeqNumber
		delta := m.samplePeriod()

		blks, err := fn()
		if err != nil {
			return nil, err
		}

		for _, blk := range blks {

			if seqno < 1 || seqno > 999999 {
				seqno = 1
			}

			res = append(res, msRecordTrace{
				SeqNumber:     seqno,
				StartTime:     start.Add(delta * time.Duration(count)),
				TimingQuality: quality,
				ClockUnlocked: !locked,
				SampleData:    blk.PackedData,
				SampleCount:   blk.SampleCount,
				ResidualCount: blk.ResidualCount,
				Encoding:      EncodingInt32,
			})

			count += blk.SampleCount

			seqno++
		}

		m.SeqNumber = seqno

		return res, nil
	}
}

func packRecord(hdr MSDataHeader, b1000 Blockette1000, b1001 Blockette1001, data []byte, maxlen int) ([]byte, error) {

	if p := MSDataHeaderSize + 2*BlocketteHeaderSize + Blockette1000Size + Blockette1001Size + len(data); p > maxlen {
		return nil, fmt.Errorf("record too large, found %d but needs to be at most %d", p, maxlen)
	}

	output := make([]byte, maxlen)

	// write the header
	pointer := uint16(MSDataHeaderSize)
	bmsh := MarshalMSDataHeader(hdr)
	copy(output[:pointer], bmsh[:])

	// write the 1000 blockette
	bh := BlocketteHeader{
		BlocketteType: 1000,
		NextBlockette: (pointer + BlocketteHeaderSize + Blockette1000Size),
	}
	bbh := MarshalBlocketteHeader(bh)
	copy(output[pointer:pointer+BlocketteHeaderSize], bbh[:])
	pointer += BlocketteHeaderSize

	bb1000 := MarshalBlockette1000(b1000)
	copy(output[pointer:pointer+Blockette1000Size], bb1000[:])
	pointer += Blockette1000Size

	// write the 1001 blockette
	bh = BlocketteHeader{
		BlocketteType: 1001,
	}
	bbh = MarshalBlocketteHeader(bh)
	copy(output[pointer:pointer+BlocketteHeaderSize], bbh[:])
	pointer += BlocketteHeaderSize

	bb1001 := MarshalBlockette1001(b1001)
	copy(output[pointer:pointer+Blockette1001Size], bb1001[:])
	pointer += Blockette1001Size

	// write the encoded data
	copy(output[pointer:pointer+uint16(len(data))], data)

	return output, nil
}
