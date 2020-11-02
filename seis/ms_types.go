package seis

import (
	"encoding/binary"
	"time"
)

const (
	MSDataHeaderSize    = 48
	BTimeSize           = 10
	BlocketteHeaderSize = 4
	Blockette1000Size   = 4
	Blockette1001Size   = 4
)

type BTime struct { //SEED Representation of Time
	Year   uint16
	Doy    uint16
	Hour   uint8
	Minute uint8
	Second uint8
	Unused byte   //Required for "alignment"
	S0001  uint16 //.0001 of a second 0-9999
}

func (b BTime) Time() time.Time {
	return time.Date(
		int(b.Year),
		1,
		1,
		int(b.Hour),
		int(b.Minute),
		int(b.Second),
		int(b.S0001)*100000,
		time.UTC,
	).AddDate(0, 0, int(b.Doy-1))
}

func NewBTime(t time.Time) BTime {
	return BTime{
		Year:   uint16(t.Year()),
		Doy:    uint16(t.YearDay()),
		Hour:   uint8(t.Hour()),
		Minute: uint8(t.Minute()),
		Second: uint8(t.Second()),
		S0001:  uint16(t.Nanosecond() / 100000),
	}
}

func UnmarshalBTime(b [BTimeSize]byte) (r BTime) {
	r.Year = binary.BigEndian.Uint16(b[0:2])
	r.Doy = binary.BigEndian.Uint16(b[2:4])
	r.Hour = b[4]
	r.Minute = b[5]
	r.Second = b[6]
	r.Unused = b[7]
	r.S0001 = binary.BigEndian.Uint16(b[8:10])

	return
}

func MarshalBTime(r BTime) (b [BTimeSize]byte) {
	binary.BigEndian.PutUint16(b[0:2], r.Year)
	binary.BigEndian.PutUint16(b[2:4], r.Doy)
	b[4] = r.Hour
	b[5] = r.Minute
	b[6] = r.Second
	b[7] = r.Unused
	binary.BigEndian.PutUint16(b[8:10], r.S0001)

	return
}

type MSDataHeader struct {
	SeqNumber   [6]byte //ASCII String representing a 7 digit number
	DataQuality byte    //ASCI: D, R, Q or M
	Reserved    byte

	//These are ascii strings
	Station  [5]byte
	Location [2]byte
	Channel  [3]byte
	Network  [2]byte

	StartTime            BTime
	SampleCount          uint16 // Number of Samples in the data block which may or may not be unpacked.
	SampleRateFactor     int16  // >0: Samples/Second <0: Second/Samples 0: Seconds/Sample, ASCII/OPAQUE DATA records
	SampleRateMultiplier int16  // >0: Multiplication Factor <0: Division Factor

	//flags are bit flags
	ActivityFlags    byte
	IOClockFlags     byte
	DataQualityFlags byte

	BlockettesToFollow uint8
	TimeCorrection     int32 // 0.0001 second units
	BeginningOfData    uint16
	FirstBlockette     uint16
}

func UnmarshalMSDataHeader(b [MSDataHeaderSize]byte) (r MSDataHeader) {
	copy(r.SeqNumber[:], b[0:6])
	r.DataQuality = b[6]
	r.Reserved = b[7]

	copy(r.Station[:], b[8:13])
	copy(r.Location[:], b[13:15])
	copy(r.Channel[:], b[15:18])
	copy(r.Network[:], b[18:20])

	var btime [10]byte
	copy(btime[:], b[20:30])
	r.StartTime = UnmarshalBTime(btime)
	r.SampleCount = binary.BigEndian.Uint16(b[30:32])
	r.SampleRateFactor = int16(binary.BigEndian.Uint16(b[32:34]))
	r.SampleRateMultiplier = int16(binary.BigEndian.Uint16(b[34:36]))

	r.ActivityFlags = b[36]
	r.IOClockFlags = b[37]
	r.DataQualityFlags = b[38]

	r.BlockettesToFollow = b[39]
	r.TimeCorrection = int32(binary.BigEndian.Uint32(b[40:44]))
	r.BeginningOfData = binary.BigEndian.Uint16(b[44:46])
	r.FirstBlockette = binary.BigEndian.Uint16(b[46:48])

	return
}

func MarshalMSDataHeader(r MSDataHeader) (b [MSDataHeaderSize]byte) {
	copy(b[0:6], r.SeqNumber[:])
	b[6] = r.DataQuality
	b[7] = r.Reserved

	copy(b[8:13], r.Station[:])
	copy(b[13:15], r.Location[:])
	copy(b[15:18], r.Channel[:])
	copy(b[18:20], r.Network[:])

	btime := MarshalBTime(r.StartTime)
	copy(b[20:30], btime[:])
	binary.BigEndian.PutUint16(b[30:32], r.SampleCount)
	binary.BigEndian.PutUint16(b[32:34], uint16(r.SampleRateFactor))
	binary.BigEndian.PutUint16(b[34:36], uint16(r.SampleRateMultiplier))

	b[36] = r.ActivityFlags
	b[37] = r.IOClockFlags
	b[38] = r.DataQualityFlags

	b[39] = r.BlockettesToFollow
	binary.BigEndian.PutUint32(b[40:44], uint32(r.TimeCorrection))
	binary.BigEndian.PutUint16(b[44:46], r.BeginningOfData)
	binary.BigEndian.PutUint16(b[46:48], r.FirstBlockette)

	return
}

type BlocketteHeader struct {
	BlocketteType uint16
	NextBlockette uint16 //Byte of next blockette, 0 if last blockette
}

func UnmarshalBlocketteHeader(b [BlocketteHeaderSize]byte) (r BlocketteHeader) {
	r.BlocketteType = binary.BigEndian.Uint16(b[0:2])
	r.NextBlockette = binary.BigEndian.Uint16(b[2:4])

	return
}

func MarshalBlocketteHeader(r BlocketteHeader) (b [BlocketteHeaderSize]byte) {
	binary.BigEndian.PutUint16(b[0:2], r.BlocketteType)
	binary.BigEndian.PutUint16(b[2:4], r.NextBlockette)

	return
}

type Blockette1000 struct { //"Data Only Seed Blockette" (excluding header)
	Encoding     uint8
	WordOrder    uint8
	RecordLength uint8
	Reserved     uint8
}

func UnmarshalBlockette1000(b [Blockette1000Size]byte) (r Blockette1000) {
	r.Encoding = b[0]
	r.WordOrder = b[1]
	r.RecordLength = b[2]
	r.Reserved = b[3]

	return
}

func MarshalBlockette1000(r Blockette1000) (b [Blockette1000Size]byte) {
	b[0] = uint8(r.Encoding)
	b[1] = uint8(r.WordOrder)
	b[2] = uint8(r.RecordLength)
	b[3] = uint8(r.Reserved)

	return
}

type Blockette1001 struct { //"Data Extension Blockette" (excluding header)
	TimingQuality uint8
	MicroSec      int8 //Increased accuracy for starttime
	Reserved      uint8
	FrameCount    uint8
}

func UnmarshalBlockette1001(b [Blockette1001Size]byte) (r Blockette1001) {
	r.TimingQuality = b[0]
	r.MicroSec = int8(b[1])
	r.Reserved = b[2]
	r.FrameCount = b[3]

	return
}

func MarshalBlockette1001(r Blockette1001) (b [Blockette1001Size]byte) {
	b[0] = r.TimingQuality
	b[1] = uint8(r.MicroSec)
	b[2] = r.Reserved
	b[3] = r.FrameCount

	return
}

//TODO: Implement Other Blockette Types?
