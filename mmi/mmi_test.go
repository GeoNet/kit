package mmi

import (
	"math"
	"runtime"
	"strconv"
	"testing"
)

func TestMMI(t *testing.T) {
	in := []struct {
		id               string
		depth, magnitude float64
		mmi              float64
	}{
		{id: loc(), depth: 5.0, magnitude: 6.3, mmi: 8.86},
		{id: loc(), depth: 40.0, magnitude: 6.8, mmi: 8.19},
		{id: loc(), depth: 11.0, magnitude: 7.1, mmi: 9.96},
		{id: loc(), depth: 150.0, magnitude: 1.5, mmi: -1.0},
		{id: loc(), depth: 150.0, magnitude: 6.5, mmi: 6.23},
		{id: loc(), depth: 7.0, magnitude: 4.4, mmi: 6.41},
	}

	for _, v := range in {
		if math.Abs(v.mmi-(MMI(v.depth, v.magnitude))) > 0.05 {
			t.Errorf("%s incorrect MMI expected %f got %f", v.id, v.mmi, MMI(v.depth, v.magnitude))
		}
	}
}

func TestMMIDistance(t *testing.T) {
	in := []struct {
		id                         string
		depth, magnitude, distance float64
		mmid                       float64
	}{
		{id: loc(), depth: 27.4, magnitude: 3.9, distance: 110.0, mmid: 2.65},
		{id: loc(), depth: 22.2, magnitude: 4.2, distance: 5.0, mmid: 5.27},
		{id: loc(), depth: 22.2, magnitude: 4.2, distance: 0.0, mmid: 5.27},
	}

	for _, v := range in {
		if math.Abs(v.mmid-(MMIDistance(v.depth, v.magnitude, v.distance))) > 0.1 {
			t.Errorf("%s expected MMI ditance %f got %f", v.id, v.mmid, MMIDistance(v.depth, v.magnitude, v.distance))
		}
	}
}

func loc() string {
	_, _, l, _ := runtime.Caller(1)
	return "L" + strconv.Itoa(l)
}
