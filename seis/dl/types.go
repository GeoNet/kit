package dl

import (
	"fmt"
	"time"
)

const (
	PreheaderSize = 3
)

type Packet struct {
	Preheader Preheader
	Header    []byte
	Body      []byte
}

func (d Packet) header() string {
	return string(d.Header)
}

func (d Packet) body() string {
	return string(d.Body)
}

type Preheader struct {
	DL           [2]byte //ASCII String == "DL"
	HeaderLength uint8   //1 byte describing the length of rest of the header
}

func UnmarshalPreheader(b [PreheaderSize]byte) (r Preheader) {
	copy(r.DL[:], b[0:2])
	r.HeaderLength = b[2]
	return
}

func MarshalPreheader(r Preheader) (b [PreheaderSize]byte) {
	copy(b[0:2], r.DL[:])
	b[2] = r.HeaderLength
	return
}

func packetToBytes(dlp Packet) ([]byte, error) {

	if len(dlp.Header) > 255 {
		return []byte{}, fmt.Errorf("cannot send a header larger than 225 bytes (uint8)")
	}

	dlp.Preheader = Preheader{
		DL:           [2]byte{'D', 'L'},
		HeaderLength: uint8(len(dlp.Header)),
	}

	out := make([]byte, 0, PreheaderSize+len(dlp.Header)+len(dlp.Body))

	mdlp := MarshalPreheader(dlp.Preheader)

	out = append(out, mdlp[:]...)
	out = append(out, dlp.Header...)
	out = append(out, dlp.Body...)

	return out, nil
}

func packetFromBytes(in []byte) (Packet, error) {
	var dlp Packet
	var pointer int

	var phIn [PreheaderSize]byte
	copy(phIn[:], in[:PreheaderSize])
	dlp.Preheader = UnmarshalPreheader(phIn)
	pointer += PreheaderSize

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
