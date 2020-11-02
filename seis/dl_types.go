package seis

import (
	"fmt"
	"time"
)

const (
	DLPreheaderSize = 3
)

type DLPacket struct {
	Preheader DLPreheader
	Header    []byte
	Body      []byte
}

func (d DLPacket) header() string {
	return string(d.Header)
}

func (d DLPacket) body() string {
	return string(d.Body)
}

type DLPreheader struct {
	DL           [2]byte //ASCII String == "DL"
	HeaderLength uint8   //1 byte describing the length of rest of the header
}

func UnmarshalDLPreheader(b [DLPreheaderSize]byte) (r DLPreheader) {
	copy(r.DL[:], b[0:2])
	r.HeaderLength = b[2]
	return
}

func MarshalDLPreheader(r DLPreheader) (b [DLPreheaderSize]byte) {
	copy(b[0:2], r.DL[:])
	b[2] = r.HeaderLength
	return
}

func dlPacketToBytes(dlp DLPacket) ([]byte, error) {

	if len(dlp.Header) > 255 {
		return []byte{}, fmt.Errorf("cannot send a header larger than 225 bytes (uint8)")
	}

	dlp.Preheader = DLPreheader{
		DL:           [2]byte{'D', 'L'},
		HeaderLength: uint8(len(dlp.Header)),
	}

	out := make([]byte, 0, DLPreheaderSize+len(dlp.Header)+len(dlp.Body))

	mdlp := MarshalDLPreheader(dlp.Preheader)

	out = append(out, mdlp[:]...)
	out = append(out, dlp.Header...)
	out = append(out, dlp.Body...)

	return out, nil
}

func dlPacketFromBytes(in []byte) (DLPacket, error) {
	var dlp DLPacket
	var pointer int

	var phIn [DLPreheaderSize]byte
	copy(phIn[:], in[:DLPreheaderSize])
	dlp.Preheader = UnmarshalDLPreheader(phIn)
	pointer += DLPreheaderSize

	hLength := int(dlp.Preheader.HeaderLength)
	if pointer+hLength > len(in) {
		return dlp, fmt.Errorf("header length of %v overflows byte slice of len %v", hLength, len(in))
	}
	dlp.Header = in[pointer : pointer+hLength]
	pointer += hLength

	dlp.Body = in[pointer:]

	return dlp, nil
}

// A time as microseconds since the Unix epoch
func hpTime(t time.Time) int64 {
	return t.UnixNano() / 1e3 //TODO: OK to truncate?
}
