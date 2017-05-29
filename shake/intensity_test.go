package shake

import (
	"testing"
)

func TestIntensity(t *testing.T) {

	// taken from standard shakemap plots .... (converted to m/s)
	var intensities = []struct {
		v float64 // Velocity m/s
		i int32   // ShakeMap
		w int32   // WaldQuitorianoHeatonKanamori1999
		f int32   // FaenzaMichelini2010
	}{
		// original tests
		{0.01 / 100.0, 1, 1, 1},
		{0.09 / 100.0, 2, 1, 2},
		{1.9 / 100.0, 5, 3, 5},
		{5.8 / 100.0, 6, 4, 6},
		{11.0 / 100.0, 7, 5, 7},
		{22.0 / 100.0, 8, 7, 8},
		{43.0 / 100.0, 8, 8, 8},
		{83.0 / 100.0, 9, 9, 9},
		{161.0 / 100.0, 10, 10, 10},
		// bounds tests
		{0.0, 1, 1, 1},
		{-0.09 / 100.0, 2, 1, 2},
		{-100.0 / 100.0, 1, 9, 9},
		{1.0e+10, 12, 12, 12},
	}

	for _, intensity := range intensities {
		if i := Intensity(WaldQuitorianoHeatonKanamori1999{}, intensity.v); i != intensity.w {
			t.Errorf("invalid WaldQuitorianoHeatonKanamori1999 rawintensity [%g cm/s]: %d (calculated) != %d (expected)", 100.0*intensity.v, i, intensity.w)
		}
		if i := Intensity(FaenzaMichelini2010{}, intensity.v); i != intensity.f {
			t.Errorf("invalid FaenzaMichelini2010 rawintensity [%g cm/s]: %d (calculated) != %d (expected)", 100.0*intensity.v, i, intensity.f)
		}
	}
}
