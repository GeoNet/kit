package shake

import (
	"math"
)

// algorithm to convert peak velocity in m/s into MMI intensity
type IntensityEquation interface {
	RawIntensity(vel float64) float64
}

// convert peak velocity in m/s into integer MMI intensity
func Intensity(ie IntensityEquation, vel float64) int32 {

	raw := ie.RawIntensity(vel)

	switch {
	case raw <= 1.0:
		return 1
	case raw >= 12.0:
		return 12
	default:
		// TODO - round up or not?
		return (int32)(math.Floor(raw))
	}
}

// David J. Wald, Vincent Quitoriano, Thomas H. Heaton, and Hiroo Kanamori (1999),
// "Relationships between Peak Ground Acceleration, Peak Ground Velocity, and
// Modified Mercalli Intensity in California", Earthquake Spectra, Volume 15, No. 3, August 1999.
type WaldQuitorianoHeatonKanamori1999 struct{}

func (fn WaldQuitorianoHeatonKanamori1999) RawIntensity(vel float64) float64 {
	return 2.35 + 3.47*math.Log10(100.0*math.Abs(vel)+1.0e-9)
}

// L. Faenza and A. Michelini (2010),
// "Regression analysis of MCS Intensity and ground motion parameters in Italy and its application
// in ShakeMap", Geophysical Journal International, 180: 1138–1152.
type FaenzaMichelini2010 struct{}

func (fn FaenzaMichelini2010) RawIntensity(vel float64) float64 {
	return 5.11 + 2.35*math.Log10(100.0*math.Abs(vel)+1.0e-9)
}

// "New Ground Motion to Intensity Conversion Equations (GMICEs) for New Zealand" - Jose M. Moratalla et. al.
// Seismological Research Letters 92 (2020), 448–459, doi: 10.1785/0220200156
type Moratalla2020 struct{}

func (fn Moratalla2020) RawIntensity(vel float64) float64 {
	switch v := math.Log10(100.0 * math.Abs(vel)); {
	case v < 1.0024:
		return 1.6323*v + 4.107
	default:
		return 3.837*v + 1.8970
	}
}
