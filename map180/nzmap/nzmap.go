/*
Package nzmap is for drawing large scale context maps in SVG format of the New Zealand region.
There are no external dependencies.

Calling code can draw markers on the maps as required.  The closing </svg> tag must also be added.

  var b bytes.Buffer
  var pts nzmap.Points

  pts.Medium(b)

	for _, p := range pts {
		if p.Visible() {
			b.WriteString(fmt.Sprintf("<path d=\"M%d %d l5 0 l-5 -8 l-5 8 Z\" stroke-width=\"0\" fill=\"blue\" opacity=\"0.7\"></path>",
			p.X(), p.Y()))
		}
	}

 b.WriteString("</svg>")

*/
package nzmap

import (
	"bytes"
	"math"
)

var nzIconPts [151][141]pt
var nzrcIconPts [29][29]pt
var nzsIconPts [25][22]pt
var nzMediumPts [151][140]pt
var nzrcMediumPts [28][29]pt
var nzrMediumPts [22][22]pt
var nzsMediumPts [25][22]pt

type pt struct {
	x, y int
}

type Point struct {
	Longitude, Latitude   float64
	x, y                  int
	visible               bool
	Value                 float64 // Points can be sorted by Value.
	Stroke, Fill, Opacity string  // Optional for any later drawing funcs
	Size                  int     // Optional for any later drawing funcs
}

type Points []Point

func (p Points) Len() int {
	return len(p)
}
func (p Points) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p Points) Less(i, j int) bool {
	return p[i].Value < p[j].Value
}

/*
Returns the SVG x coord for p.  Icon() or Map() must be called to set this value.
*/
func (p *Point) X() int {
	return p.x
}

/*
Returns the SVG y coord for p.  Icon() or Map() must be called to set this value.
*/
func (p *Point) Y() int {
	return p.y
}

/*
Visible returns true if p is visible on the SVG canvas.  Icon() or Map() must be called to set this value.
*/
func (p *Point) Visible() bool {
	return p.visible
}

/*
Icon selects a small low resolution coastline map of New Zealand, and possibly surrounding regions, that is
suitable for icon size maps.
*/
func (p *Point) Icon(b *bytes.Buffer) {
	longitude := p.Longitude
	latitude := p.Latitude

	if longitude < 0 {
		longitude = longitude + 360.0
	}

	switch {
	// New Zealand.
	case longitude >= 165.0 && longitude <= 180.0 && latitude >= -48.0 && latitude <= -34.0:
		x := int(math.Floor((longitude+0.05)*10)) - 1650
		y := int(math.Floor((latitude+0.05)*10)) + 480
		if x >= 0 && x <= 151 && y >= 0 && y <= 141 {
			b.WriteString(nzIcon)
			p.x = nzIconPts[x][y].x
			p.y = nzIconPts[x][y].y
			p.visible = true
		}
	// New Zealand, Raoul, Chathams.
	case longitude >= 165.0 && longitude <= 193.0 && latitude >= -48.0 && latitude <= -20.0:
		x := int(math.Floor(longitude+0.5)) - 165
		y := int(math.Floor(latitude+0.5)) + 48
		if x >= 0 && x <= 29 && y >= 0 && y <= 29 {
			b.WriteString(nzrcIcon)
			p.x = nzrcIconPts[x][y].x
			p.y = nzrcIconPts[x][y].y
			p.visible = true
		}
	// New Zealand, South.
	case longitude >= 156.0 && longitude <= 180.0 && latitude >= -55.0 && latitude <= -34.0:
		x := int(math.Floor(longitude+0.5)) - 156
		y := int(math.Floor(latitude+0.5)) + 55
		if x >= 0 && x <= 21 && y >= 0 && y <= 22 {
			b.WriteString(nzsIcon)
			p.x = nzsIconPts[x][y].x
			p.y = nzsIconPts[x][y].y
			p.visible = true
		}
	default:
		b.WriteString(nzIcon)
		p.x = -1000
		p.y = -1000
	}
}

/*
Medium returns a map of New Zealand or the region around New Zealand at medium resolution.
Linear interpolation is used between grid points to estimate the location of each Point on the map.

pts[0] is used to decide which region to return.
*/
func (pts Points) Medium(b *bytes.Buffer) {
	if pts == nil || len(pts) == 0 {
		b.WriteString(nzMedium)
		return
	}

	longitude := pts[0].Longitude
	latitude := pts[0].Latitude

	if longitude < 0 {
		longitude = longitude + 360.0
	}

	switch {
	// New Zealand.
	case longitude >= 165.0 && longitude <= 180.0 && latitude >= -48.0 && latitude <= -34.0:
		// the long/lat grid is accurate to 0.1 degree below that use a linear approximation between
		// the grid values.  This removes liniations in the plot
		var p, pp pt
		for i, v := range pts {
			if v.Longitude < 0 {
				v.Longitude = v.Longitude + 360.0
			}
			xi, xf := math.Modf(v.Longitude*10 - 1650.0)
			x := int(xi)
			yi, yf := math.Modf(v.Latitude*10 + 480.0)
			y := int(yi)

			if x >= 0 && x < 150 && y >= 0 && y < 140 {
				p = nzMediumPts[int(x)][y]
				pp = nzMediumPts[x+1][y+1]
				pts[i].x = p.x + int(float64(pp.x-p.x)*xf)
				pts[i].y = p.y + int(float64(pp.y-p.y)*yf)
				pts[i].visible = true
			} else {
				pts[i].x = -1000
				pts[i].y = -1000
			}
		}
		b.WriteString(nzMedium)
		//b, err := newBbox("165,-48,-174,-27")
	// New Zealand, Raoul
	case longitude >= 165.0 && longitude <= 186.0 && latitude >= -48.0 && latitude <= -27.0:
		var p, pp pt
		for i, v := range pts {
			if v.Longitude < 0 {
				v.Longitude = v.Longitude + 360.0
			}
			xi, xf := math.Modf(v.Longitude - 165.0)
			x := int(xi)
			yi, yf := math.Modf(v.Latitude + 48.0)
			y := int(yi)

			if x >= 0 && x < 21 && y >= 0 && y < 21 {
				p = nzrMediumPts[int(x)][y]
				pp = nzrMediumPts[x+1][y+1]
				pts[i].x = p.x + int(float64(pp.x-p.x)*xf)
				pts[i].y = p.y + int(float64(pp.y-p.y)*yf)
				pts[i].visible = true
			} else {
				pts[i].x = -1000
				pts[i].y = -1000
			}
		}
		b.WriteString(nzrMedium)
	// New Zealand, Raoul, Chathams.
	case longitude >= 165.0 && longitude <= 193.0 && latitude >= -48.0 && latitude <= -20.0:
		var p, pp pt
		for i, v := range pts {
			if v.Longitude < 0 {
				v.Longitude = v.Longitude + 360.0
			}
			xi, xf := math.Modf(v.Longitude - 165.0)
			x := int(xi)
			yi, yf := math.Modf(v.Latitude + 48.0)
			y := int(yi)

			if x >= 0 && x < 27 && y >= 0 && y < 28 {
				p = nzrcMediumPts[int(x)][y]
				pp = nzrcMediumPts[x+1][y+1]
				pts[i].x = p.x + int(float64(pp.x-p.x)*xf)
				pts[i].y = p.y + int(float64(pp.y-p.y)*yf)
				pts[i].visible = true
			} else {
				pts[i].x = -1000
				pts[i].y = -1000
			}
		}
		b.WriteString(nzrcMedium)
	// New Zealand, South.
	case longitude >= 156.0 && longitude <= 180.0 && latitude >= -55.0 && latitude <= -34.0:
		var p, pp pt
		for i, v := range pts {
			if v.Longitude < 0 {
				v.Longitude = v.Longitude + 360.0
			}
			xi, xf := math.Modf(v.Longitude - 156.0)
			x := int(xi)
			yi, yf := math.Modf(v.Latitude + 55.0)
			y := int(yi)

			if x >= 0 && x < 24 && y >= 0 && y < 21 {
				p = nzsMediumPts[int(x)][y]
				pp = nzsMediumPts[x+1][y+1]
				pts[i].x = p.x + int(float64(pp.x-p.x)*xf)
				pts[i].y = p.y + int(float64(pp.y-p.y)*yf)
				pts[i].visible = true
			} else {
				pts[i].x = -1000
				pts[i].y = -1000
			}
		}
		b.WriteString(nzsMedium)
	default:
		// the first point is not a region we have so send NZ.
		// some of the later points may still be on the map.
		var p, pp pt
		for i, v := range pts {
			if v.Longitude < 0 {
				v.Longitude = v.Longitude + 360.0
			}
			xi, xf := math.Modf(v.Longitude*10 - 1650.0)
			x := int(xi)
			yi, yf := math.Modf(v.Latitude*10 + 480.0)
			y := int(yi)

			if x >= 0 && x < 150 && y >= 0 && y < 140 {
				p = nzMediumPts[int(x)][y]
				pp = nzMediumPts[x+1][y+1]
				pts[i].x = p.x + int(float64(pp.x-p.x)*xf)
				pts[i].y = p.y + int(float64(pp.y-p.y)*yf)
				pts[i].visible = true
			} else {
				pts[i].x = -1000
				pts[i].y = -1000
			}
		}
		b.WriteString(nzMedium)
	}
}
