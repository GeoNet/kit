// Package mmi is for Modificed Mercalli Intensity calculations in New Zealand.
package mmi

import "math"

// MMI calculates the maximum Modificed Mercalli Intensity for the quake.
// depth is in km.
func MMI(depth, magnitude float64) float64 {
	var w, m float64
	d := math.Abs(depth)
	rupture := d

	if d < 100 {
		w = math.Min(0.5*math.Pow(10, magnitude-5.39), 30.0)
		rupture = math.Max(d-0.5*w*0.85, 0.0)
	}

	if d < 70.0 {
		m = 4.40 + 1.26*magnitude - 3.67*math.Log10(rupture*rupture*rupture+1634.691752)/3.0 + 0.012*d + 0.409
	} else {
		m = 3.76 + 1.48*magnitude - 3.50*math.Log10(rupture*rupture*rupture)/3.0 + 0.0031*d
	}

	if m < 3.0 {
		m = -1.0
	}

	return m
}

// MMIDistance calculates the MMI at distance for New Zealand.  Distance and depth are in km.
func MMIDistance(depth, magnitude, distance float64) float64 {
	// Minimum depth of 5 for numerical instability.
	d := math.Max(math.Abs(depth), 5.0)
	s := math.Hypot(d, distance)

	return math.Max(MMI(depth, magnitude)-1.18*math.Log(s/d)-0.0044*(s-d), -1.0)
}

// MMIIntensity returns the string describing mmi.
func MMIIntensity(mmi float64) string {
	switch {
	case mmi >= 7:
		return "severe"
	case mmi >= 6:
		return "strong"
	case mmi >= 5:
		return "moderate"
	case mmi >= 4:
		return "light"
	case mmi >= 3:
		return "weak"
	default:
		return "unnoticeable"
	}
}

// IntensityMMI returns the minimum MMI for the intensity.
func IntensityMMI(Intensity string) float64 {
	switch Intensity {
	case "severe":
		return 7
	case "strong":
		return 6
	case "moderate":
		return 5
	case "light":
		return 4
	case "weak":
		return 3
	default:
		return -9
	}
}

// Severity returns the CAP severity for mmi.
func Severity(mmi float64) string {
	switch {
	case mmi >= 8:
		return `Extreme`
	case mmi >= 7:
		return `Severe`
	case mmi >= 6:
		return `Moderate`
	default:
		return "Minor"
	}
}
