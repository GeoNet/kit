package seis

import (
	"encoding/binary"
)

type encodedData struct {
	PackedData    []byte
	SampleCount   int
	ResidualCount int
	FrameCount    int
}

func (e encodedData) Last() bool {
	return !(e.ResidualCount > 0)
}

type encodingFunc func() ([]encodedData, error)

func encodeInt32(data []int32, maxData int, wordOrder WordOrder) encodingFunc {
	return func() ([]encodedData, error) {
		var res []encodedData

		count, remaining := 0, len(data)

		for remaining > 0 {
			buf := make([]byte, maxData)

			var samples int
			for p := 0; remaining > 0 && (p+4) < maxData; p += 4 {
				switch wordOrder {
				case 0:
					binary.LittleEndian.PutUint32(buf[p:p+4], uint32(data[count]))
				default:
					binary.BigEndian.PutUint32(buf[p:p+4], uint32(data[count]))
				}

				count, remaining = count+1, remaining-1

				samples++
			}

			res = append(res, encodedData{
				PackedData:    buf,
				SampleCount:   samples,
				ResidualCount: remaining,
			})
		}

		return res, nil
	}
}
