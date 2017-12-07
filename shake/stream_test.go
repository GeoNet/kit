package shake

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"
)

const GRAVITY = 9.80665

func TestStreams(t *testing.T) {

	// Check real data for a couple of sites, one, KIKS, has only strong motion
	// data, whereas the other two have both weak and strong motion continuous data.
	// The estimates based on weak motion data use scwfparam on an independent sensor.
	// Arbitarily uses 0.15 as the pga/pgv relative error thresholds.

	// For reference, sensor FBA-ES-T (2g) has gain 1.01971 Volts per M/S**S
	// and the Q330 datalogger has gain 419430.4 counts per Volt.
	// The Basalt (at KIKS) has the same overall gain as the Q330 with FBA-ES-T.

	// in the tests pga and pgv returned from Peaks is converted from m/s/s and m/s to
	// %g and cm/s so that it can be compared to values calculated externally (from scwfparam).

	var tests = []struct {
		path   string  // path to raw samples
		length int     // expected length
		sps    float64 // sampling rate
		q      float64 // filter factor
		pga    string  // formatted pga to check calculations
		pgv    string  // formatted pgv to check calculations
		extPga float64 // external pga to check algorithm
		extPgv float64 // external pgv to check algorithm
	}{
		// based on v1a analysis
		// 2433.7 mm/s/s and 353.33 mm/s
		{
			"./testdata/NZKIKS_HNZ20.txt",
			120000,
			200.0,
			0.98829,
			"25.8213",
			"30.1618",
			100.0 * 2.4337 / GRAVITY,
			35.333,
		},
		// scwfparam 2016p969664 TCW L4C-3D (short-period)
		//<acc value="0.2050960949" flag="0"/> units of %g [uses 9.81 as g]
		//<vel value="0.0544451529" flag="0"/> units of cm/s
		{
			"./testdata/NZTCW__HNZ20.txt",
			10000,
			200.0,
			0.98829,
			"0.2087",
			"0.0602",
			0.2050960949,
			0.0544451529,
		},
		// scwfparam 2016p969664 DUWZ LE-3DliteMkII (short-period)
		// <acc value="0.0481134810" flag="0"/> units of %g [uses 9.81 as g]
		// <vel value="0.0304631505" flag="0"/> units of cm/s
		{
			"./testdata/NZDUWZ_HNZ20.txt",
			10000,
			200.0,
			0.98829,
			"0.0481",
			"0.0271",
			0.0481134810,
			0.0304631505,
		},
		// scwfparam 2016p981371 TSZ CMG-3ESP (broadband)
		// <acc value="0.0229067493" flag="0"/> units of %g [uses 9.81 as g]
		// <vel value="0.0185129043" flag="0"/> units of cm/s
		{
			"./testdata/NZTSZ__BNZ20.txt",
			30000,
			50.0,
			0.95395,
			"0.0254",
			"0.0163",
			0.0229067493,
			0.0185129043,
		},
	}

	for _, x := range tests {
		s := Stream{
			HighPass:   NewHighPass(427697.373184, x.q),
			Integrator: NewIntegrator(1.0, 1.0/x.sps, x.q),
		}
		d, err := func(path string) ([]int32, error) {
			file, err := os.Open(path)
			if err != nil {
				return nil, err
			}
			defer file.Close()

			var samples []int32
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				parts := strings.Fields(scanner.Text())
				if len(parts) > 1 {
					s, err := strconv.ParseInt(parts[1], 10, 32)
					if err != nil {
						return nil, err
					}
					samples = append(samples, int32(s))
				}
			}

			if err := scanner.Err(); err != nil {
				return nil, err
			}

			return samples, nil
		}(x.path)
		if err != nil {
			t.Fatalf("unable to load test file: %v", err)
		}
		if len(d) != x.length {
			t.Fatalf("unable to decode test data file, length mismatch %s: found %d, expected %d", x.path, len(d), x.length)
		}
		s.Reset()
		s.Condition(d)

		pga, pgv := s.Peaks(d)

		if v := strconv.FormatFloat(100*pga/GRAVITY, 'f', 4, 64); v != x.pga {
			t.Errorf("invalid pga %s: found %s, expected %s", x.path, v, x.pga)
		}
		if v := strconv.FormatFloat(pgv*100, 'f', 4, 64); v != x.pgv {
			t.Errorf("invalid pgv %s: found %s, expected %s", x.path, v, x.pgv)
		}

		if rpga := math.Abs(((100 * pga / GRAVITY) - x.extPga) / x.extPga); rpga > 0.15 {
			t.Errorf("large pga error %s:, expected: %g got %g, relative error: %g", x.path, x.extPga, pga, rpga)
		}
		if rpgv := math.Abs(((100 * pgv) - x.extPgv) / x.extPgv); rpgv > 0.15 {
			t.Errorf("large pgv error %s:, expected: %g got %g, relative error: %g", x.path, x.extPgv, pgv, rpgv)
		}
	}
}
