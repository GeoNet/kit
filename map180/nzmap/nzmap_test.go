package nzmap

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

// these tests also output SVG to svg_test/ for visual inspection.
func init() {
	if err := os.Mkdir("svg_test", 0755); err != nil {
		fmt.Println("Error creating svg_test directory", err)
	}
}

func TestIconWellington(t *testing.T) {
	var b bytes.Buffer
	p := Point{Longitude: 174.7,
		Latitude: -41.2,
		Value:    -1,
	}
	p.Icon(&b)

	if !p.Visible() {
		t.Error("point should be on the icon map")
	}

	b.WriteString(fmt.Sprintf("<circle cx=\"%d\" cy=\"%d\" r=\"5\" fill=\"black\"/></svg>", p.X(), p.Y()))

	if err := ioutil.WriteFile("svg_test/nzicon-wellington.svg", b.Bytes(), 0644); err != nil { // nolint: gosec
		t.Fatal(err)
	}
}

func TestIconRaoul(t *testing.T) {
	var b bytes.Buffer
	p := Point{Longitude: -177.9286,
		Latitude: -29.2684,
		Value:    -1,
	}
	p.Icon(&b)

	if !p.Visible() {
		t.Error("point should be on the icon map")
	}

	b.WriteString(fmt.Sprintf("<circle cx=\"%d\" cy=\"%d\" r=\"5\" fill=\"black\"/></svg>", p.X(), p.Y()))

	if err := ioutil.WriteFile("svg_test/nzicon-raoul.svg", b.Bytes(), 0644); err != nil { // nolint: gosec
		t.Fatal(err)
	}
}

func TestIconAucklandIsland(t *testing.T) {
	var b bytes.Buffer
	p := Point{Longitude: 166.102,
		Latitude: -50.72,
		Value:    -1,
	}
	p.Icon(&b)

	if !p.Visible() {
		t.Error("point should be on the icon map")
	}

	b.WriteString(fmt.Sprintf("<circle cx=\"%d\" cy=\"%d\" r=\"5\" fill=\"black\"/></svg>", p.X(), p.Y()))

	if err := ioutil.WriteFile("svg_test/nzicon-aucklandisland.svg", b.Bytes(), 0644); err != nil { // nolint: gosec
		t.Fatal(err)
	}
}

func TestIconCanberra(t *testing.T) {
	var b bytes.Buffer
	p := Point{Longitude: 149.1300,
		Latitude: -35.2809,
		Value:    -1,
	}
	p.Icon(&b)

	if p.Visible() {
		t.Error("point should not be on the icon map")
	}

	b.WriteString(fmt.Sprintf("<circle cx=\"%d\" cy=\"%d\" r=\"5\" fill=\"black\"/></svg>", p.X(), p.Y()))

	if err := ioutil.WriteFile("svg_test/nzicon-canberra.svg", b.Bytes(), 0644); err != nil { // nolint: gosec
		t.Fatal(err)
	}
}

func TestMediumWellington(t *testing.T) {
	var b bytes.Buffer

	var pt Points

	pt = append(pt, Point{Longitude: 174.7,
		Latitude: -41.2,
		Value:    -1,
	})

	pt.Medium(&b)

	for _, p := range pt {
		if !p.Visible() {
			t.Error("point should be on the map")
		}

		b.WriteString(fmt.Sprintf("<circle cx=\"%d\" cy=\"%d\" r=\"5\" fill=\"black\"/>", p.X(), p.Y()))
	}

	b.WriteString("</svg>")

	if err := ioutil.WriteFile("svg_test/nzmedium-wellington.svg", b.Bytes(), 0644); err != nil { // nolint: gosec
		t.Fatal(err)
	}
}

func TestMediumIob(t *testing.T) {
	var b bytes.Buffer

	var pt Points

	// should be on the map but not plotable (interpolation). Use value to id.  should not iob
	pt = append(pt, Point{Longitude: 180.0,
		Latitude: -34.0,
		Value:    -1,
	})

	pt.Medium(&b)

	for _, p := range pt {
		if p.Visible() {
			t.Error("point should not be on the map")
		}

		b.WriteString(fmt.Sprintf("<circle cx=\"%d\" cy=\"%d\" r=\"5\" fill=\"black\"/>", p.X(), p.Y()))
	}

	b.WriteString("</svg>")

	if err := ioutil.WriteFile("svg_test/nzmedium-iob.svg", b.Bytes(), 0644); err != nil { // nolint: gosec
		t.Fatal(err)
	}
}

func TestMediumRegionWellington(t *testing.T) {
	var b bytes.Buffer

	var pt Points = Points{
		// ~ Wellington
		Point{Longitude: 174.7,
			Latitude: -41.2,
			Value:    -1,
		},
		// ~ Raoul
		Point{Longitude: -177.9286,
			Latitude: -29.2684,
			Value:    -10,
		},

		// Chathams (north east point) - should be off the map. Use value to id
		Point{Longitude: -176.254,
			Latitude: -43.751,
			Value:    -10,
		},
		// Canberra - should be off the map. Use value to id
		Point{Longitude: 149.1300,
			Latitude: -35.2809,
			Value:    -10,
		},
		// Top right - should be on the map but not plottable (interpolation). Use value to id.  should not iob
		Point{Longitude: 180.0,
			Latitude: -34.0,
			Value:    -10,
		},
	}

	pt.Medium(&b)

	for _, p := range pt {
		if p.Value == -1 {
			if !p.Visible() {
				t.Error("point should be on the map")
			}
		} else {
			if p.Visible() {
				t.Error("point should not be on the map")
			}
		}

		b.WriteString(fmt.Sprintf("<circle cx=\"%d\" cy=\"%d\" r=\"2\" fill=\"black\"/>", p.X(), p.Y()))
	}

	b.WriteString("</svg>")

	if err := ioutil.WriteFile("svg_test/nzmediumregion-wellington.svg", b.Bytes(), 0644); err != nil { // nolint: gosec
		t.Fatal(err)
	}
}

func TestMediumRegionRaoul(t *testing.T) {
	var b bytes.Buffer

	var pt Points = Points{
		// ~ Raoul
		Point{Longitude: -177.9286,
			Latitude: -29.2684,
			Value:    -1,
		},
		// ~ Wellington
		Point{Longitude: 174.7,
			Latitude: -41.2,
			Value:    -1,
		},
		// Chathams (north east point)
		Point{Longitude: -176.254,
			Latitude: -43.751,
			Value:    -1,
		},
		// Top right - should be on the map but not plottable (interpolation). Use value to id.  should not iob
		Point{Longitude: -167.4,
			Latitude: -20.0,
			Value:    -10,
		},
		// Canberra - should be off the map. Use value to id
		Point{Longitude: 149.1300,
			Latitude: -35.2809,
			Value:    -10,
		},
	}

	pt.Medium(&b)

	for _, p := range pt {
		if p.Value == -1 {
			if !p.Visible() {
				t.Error("point should be on the map")
			}
		} else {
			if p.Visible() {
				t.Error("point should not be on the map")
			}
		}

		b.WriteString(fmt.Sprintf("<circle cx=\"%d\" cy=\"%d\" r=\"2\" fill=\"black\"/>", p.X(), p.Y()))
	}

	b.WriteString("</svg>")

	if err := ioutil.WriteFile("svg_test/nzmediumregion-raoul.svg", b.Bytes(), 0644); err != nil { // nolint: gosec
		t.Fatal(err)
	}
}
