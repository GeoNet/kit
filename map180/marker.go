package map180

import (
	"bytes"
	"fmt"
)

const (
	svgPointQuery = `with p as (
		select st_transScale(st_transform(st_setsrid(st_makepoint($1, $2), 4326), 3857), $3, $4, $5, $6) as pt
		)
	select round(ST_X(pt)), round(ST_Y(pt)*-1) from p`
)

type Marker struct {
	longitude, latitude   float64
	drawSVG               SVGMarker
	x, y                  float64
	svgColour             string
	size                  int
	id, label, shortLabel string
}

// NewMarker returns a Marker to be drawn at longitude, latitude (EPSG:4327).
// Latitude must be in the range -85 to 85 otherwise the Marker will not be drawn.
// Defaults to a red triangle size 10.
func NewMarker(longitude, latitude float64, id, label, shortLabel string) Marker {
	m := Marker{
		longitude:  longitude,
		latitude:   latitude,
		drawSVG:    SVGTriangle,
		svgColour:  "red",
		size:       10,
		id:         id,
		label:      label,
		shortLabel: shortLabel,
	}

	return m
}

func (m *Marker) SetSvgColour(colour string) {
	m.svgColour = colour
}

func (m *Marker) SetSize(size int) {
	m.size = size
}

// SetSVGMarker sets the SVGMarker func that will be used
// to draw the marker.
func (m *Marker) SetSVGMarker(f SVGMarker) {
	m.drawSVG = f
}

type SVGMarker func(Marker, *bytes.Buffer)

// SVGTriangle is an SVGMarker func.
func SVGTriangle(m Marker, b *bytes.Buffer) {
	w := int(m.size / 2)
	h := w * 2
	c := int(float64(h) * 0.33)

	b.WriteString(fmt.Sprintf("<g id=\"%s\"><path d=\"M%d %d l%d 0 l-%d -%d l-%d %d Z\" fill=\"%s\" opacity=\"0.5\">",
		m.id, int(m.x), int(m.y)+c, w, w, h, w, h, m.svgColour))
	b.WriteString(`<desc>` + m.label + `.</desc>`)
	b.WriteString(fmt.Sprint(`<set attributeName="opacity" from="0.5" to="1" begin="mouseover" end="mouseout"  dur="2s"/></path>`))
	b.WriteString(fmt.Sprintf("<path d=\"M%d %d l%d 0 l-%d -%d l-%d %d Z\" stroke=\"%s\" stroke-width=\"1\" fill=\"none\" opacity=\"1\" /></g>",
		int(m.x), int(m.y)+c, w, w, h, w, h, m.svgColour))
}

// puts the label or short label on the SVG image all at the same place.
// labels are made visible using mouseover on the marker with the same id.
func labelMarkers(m []Marker, x, y int, anchor string, fontSize int, short bool, b *bytes.Buffer) {
	b.WriteString(`<g id="marker_labels">`)
	for _, mr := range m {
		b.WriteString(fmt.Sprintf("<text x=\"%d\" y=\"%d\" font-size=\"%d\" visibility=\"hidden\" text-anchor=\"%s\">", x, y, fontSize, anchor))
		if short {
			b.WriteString(mr.shortLabel)
		} else {
			b.WriteString(mr.label)
		}
		b.WriteString(fmt.Sprintf("<set attributeName=\"visibility\" from=\"hidden\" to=\"visible\" begin=\"%s.mouseover\" end=\"%s.mouseout\" dur=\"2s\"/>",
			mr.id, mr.id))
		b.WriteString(`</text>`)
	}
	b.WriteString(`</g>`)
}

func (m map3857) drawMarkers(markers []Marker, b *bytes.Buffer) (err error) {
	for _, mr := range markers {
		if mr.latitude <= 85.0 && mr.latitude >= -85.0 {
			err = m.marker3857(&mr)
			if err != nil {
				return
			}
			mr.drawSVG(mr, b)
		}
	}
	return
}

// sets x,y on 3857 from long lat on 4326.  Allows for crossing 180.
func (m map3857) marker3857(marker *Marker) (err error) {
	// map crosses 180 and the point is to the right
	switch m.crossesCentral && marker.longitude > -180.0 && marker.longitude < 0.0 {
	case true:
		err = db.QueryRow(svgPointQuery, marker.longitude, marker.latitude, width3857-m.llx, m.yshift, m.dx, m.dx).Scan(&marker.x, &marker.y)
	case false:
		err = db.QueryRow(svgPointQuery, marker.longitude, marker.latitude, m.xshift, m.yshift, m.dx, m.dx).Scan(&marker.x, &marker.y)
	}

	return
}
