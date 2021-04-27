package sl

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

func mustParse(str string) time.Time {
	t, err := time.Parse("2006,002,15:04:05.999999", str)
	if err != nil {
		return time.Time{}
	}
	return t
}

func TestStation(t *testing.T) {

	headers := []Station{
		{
			Network:   "NZ",
			Station:   "TDHS",
			Sequence:  175910,
			Timestamp: mustParse("2016,245,16:36:51.0518"),
		},
		{
			Network:   "AU",
			Station:   "MOO",
			Sequence:  1,
			Timestamp: mustParse("2019,001,00:00:08.3195"),
		},
		{
			Network:   "SL",
			Station:   "INFO",
			Sequence:  1,
			Timestamp: mustParse("2019,074,03:45:47.508800"),
		},
		{
			Network:   "NZ",
			Station:   "AUCT",
			Sequence:  1,
			Timestamp: mustParse("2019,099,01:52:28.069500"),
		},
		{
			Network:   "NZ",
			Station:   "CHIT",
			Sequence:  1,
			Timestamp: mustParse("2019,099,01:52:55.069500"),
		},
		{
			Network:   "AU",
			Station:   "MILA",
			Sequence:  136425,
			Timestamp: mustParse("2019,148,00:00:09.025000"),
		},
		{
			Network:   "NZ",
			Station:   "WEL",
			Sequence:  1,
			Timestamp: mustParse("2000,021,13:43:00.180000"),
		},
	}

	raw, err := ioutil.ReadFile("testdata/test.ms")
	if err != nil {
		t.Fatal(err)
	}

	blocks := len(raw) / 512
	if n := len(headers); n != blocks {
		t.Fatalf("invalid number of blocks, expected %d, but got %d", n, blocks)
	}

	for n, h := range headers {
		s := UnpackStation(fmt.Sprintf("%06X", h.Sequence), raw[n*512:(n+1)*512])
		if s.Network != h.Network {
			t.Errorf("%d invalid network, expected %s, but got %s", n, h.Network, s.Network)
		}
		if s.Station != h.Station {
			t.Errorf("%d invalid station, expected %s, but got %s", n, h.Station, s.Station)
		}
		if s.Sequence != h.Sequence {
			t.Errorf("%d invalid sequence, expected %d, but got %d", n, h.Sequence, s.Sequence)
		}
		if !s.Timestamp.Equal(h.Timestamp) {
			t.Errorf("%d invalid timestamp, expected %v, but got %v", n, h.Timestamp, s.Timestamp)
		}
	}
}
