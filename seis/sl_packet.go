package seis

import (
	"fmt"
)

const (
	SLPacketSize = 8 + 512
)

type SLPacket struct {
	SL   [2]byte   // ASCII String == "SL"
	Seq  [6]byte   // ASCII sequence number
	Data [512]byte // Fixed size payload
}

func NewSLPacket(data []byte) (*SLPacket, error) {
	if l := len(data); l < SLPacketSize {
		return nil, fmt.Errorf("invalid packet data length: %d", l)
	}
	if data[0] != 'S' || data[1] != 'L' {
		return nil, fmt.Errorf("invalid packet header tag: %v", string(data[0:2]))
	}

	var pkt SLPacket

	copy(pkt.SL[:], data[0:2])
	copy(pkt.Seq[:], data[2:8])
	copy(pkt.Data[:], data[8:])

	return &pkt, nil
}
