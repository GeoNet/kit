package map180

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type bbox struct {
	llx, lly, urx, ury float64 // EPSG:4327
	region             int     // If not known set to -1
	crosses180         bool
	title              string
}

/*
ValidBbox returns a nil error for a valid boundingBox.  Valid options are one of:
  - an empty string: ""
  - a string definition of a bounding box using ',' separated
    longitude latitude (float) on EPSG4327.  This should be lower
    left and upper right corners.  It may cross 180.  Longitude can be -180 to 180
    or 0 to 360.  Latitude must be <= 85.0 && >= -85.0  Examples:
    "165,-48,179,-34"  // New Zealand
    "165,-48,-175,-34" // New Zealand and Chathams
    "165,-48,185,-34" // New Zealand and Chathams
*/
func ValidBbox(boundingBox string) error {
	if boundingBox == "" {
		return nil
	}

	_, err := newBbox(boundingBox)
	return err
}

func BboxToWKTPolygon(boundingBox string) (p string, err error) {
	if boundingBox == "" {
		err = fmt.Errorf("valid but empty boundingBox")
		return
	}

	b, err := newBbox(boundingBox)
	if err != nil {
		return
	}

	return fmt.Sprintf("POLYGON((%f %f,%f %f,%f %f,%f %f,%f %f))",
		b.llx, b.lly,
		b.llx, b.ury,
		b.urx, b.ury,
		b.urx, b.lly,
		b.llx, b.lly), nil
}

// parses boundingBox and returns a bbox
func newBbox(boundingBox string) (b bbox, err error) {
	// If it's a named bbox return that
	b, ok := namedMapBounds[boundingBox]
	if ok {
		return
	}

	s := strings.Split(boundingBox, ",")

	b.llx, err = strconv.ParseFloat(s[0], 64)
	if err != nil {
		err = fmt.Errorf("Invalid boundingBox: %s", boundingBox)
		return
	}

	b.lly, err = strconv.ParseFloat(s[1], 64)
	if err != nil {
		err = fmt.Errorf("Invalid boundingBox: %s", boundingBox)
		return
	}

	b.urx, err = strconv.ParseFloat(s[2], 64)
	if err != nil {
		err = fmt.Errorf("Invalid boundingBox: %s", boundingBox)
		return
	}

	b.ury, err = strconv.ParseFloat(s[3], 64)
	if err != nil {
		err = fmt.Errorf("Invalid boundingBox: %s", boundingBox)
		return
	}

	if b.lly < -85.0 {
		err = fmt.Errorf("bbox out of bounds %s", boundingBox)
		return
	}

	if b.ury > 85.0 {
		err = fmt.Errorf("bbox out of bounds %s", boundingBox)
		return
	}

	b.setCrosses180()

	err = b.setRegion()

	return
}

// finds a bbox for the markers.  Uses the default mapBounds for the set Region.
func newBboxFromMarkers(m []Marker) (bbox, error) {
	var geom string
	switch len(m) {
	case 0:
		return bbox{}, fmt.Errorf("zero length markers, can't determine map bounds")
	case 1:
		geom = fmt.Sprintf("POINT(%f %f)", m[0].longitude, m[0].latitude)
	default:
		geom = "LINESTRING("
		for _, mr := range m {
			geom = geom + fmt.Sprintf("%f %f,", mr.longitude, mr.latitude)
		}
		geom = strings.TrimSuffix(geom, ",")
		geom = geom + ")"
	}

	for _, b := range mapBounds {
		var in bool
		err := db.QueryRow(`select ST_Within(ST_ShiftLongitude(st_setsrid(ST_GeomFromText($1), 4326)),
		ST_ShiftLongitude(ST_MakeEnvelope($2,$3,$4,$5, 4326)))`, geom, b.llx, b.lly, b.urx, b.ury).Scan(&in)
		if err != nil {
			return bbox{}, err
		}

		if in {
			return b, nil
		}
	}

	return world, nil
}

func (b *bbox) setCrosses180() {
	if (b.llx >= 0.0 && b.llx < 180.0 && b.urx > -180.0 && b.urx <= 0.0) ||
		(b.llx >= 0.0 && b.llx < 180.0 && b.urx > 180.0 && b.urx <= 360.0) {
		b.crosses180 = true
	}
}

func (b *bbox) setRegion() error {
	// allow for 0-360 bbox queries
	if b.urx > 180 {
		b.urx = -360.0 + b.urx
	}
	if b.llx > 180 {
		b.llx = -360.0 + b.llx
	}

	var in bool
	err := db.QueryRow(`select ST_Within(ST_ShiftLongitude(ST_MakeEnvelope($1,$2,$3,$4, 4326)),
		ST_ShiftLongitude(ST_MakeEnvelope($5,$6,$7,$8, 4326)))`, b.llx, b.lly, b.urx, b.ury,
		zoomRegion.llx, zoomRegion.lly, zoomRegion.urx, zoomRegion.ury).Scan(&in)
	if err != nil {
		return err
	}

	if in {
		b.region = zoomRegion.region
	} else {
		b.region = 0 // the world region
	}

	return nil
}

// map3857 contains information for drawing map on EPSG3857
func (b *bbox) newMap3857(width int) (m map3857, err error) {
	// bbox on 3857
	// tried using st_MakeEnvelope so that only needed to hit DB once
	// but it does not do what I for crossing 180
	err = db.QueryRow(`with p as (
		select st_transform(st_setsrid(st_makepoint($1, $2), 4326), 3857) as pt 
		)
	select ST_X(pt), ST_Y(pt) from p;`, b.llx, b.lly).Scan(&m.llx, &m.lly)
	if err != nil {
		return
	}

	err = db.QueryRow(`with p as (
		select st_transform(st_setsrid(st_makepoint($1, $2), 4326), 3857) as pt 
		)
	select ST_X(pt), ST_Y(pt) from p;`, b.urx, b.ury).Scan(&m.urx, &m.ury)
	if err != nil {
		return
	}

	// Minimum across the top of the bbox on 3857
	// Allows for crossing 180
	var x float64
	err = db.QueryRow(`select ST_Distance(
	st_transform(st_setsrid(st_makepoint($1, $2), 4326), 3857),
	st_transform(st_setsrid(st_makepoint($3, $4), 4326), 3857)
	)`, b.llx, b.ury, b.urx, b.ury).Scan(&x)
	if err != nil {
		return
	}

	if b.crosses180 {
		m.crossesCentral = true
		x = width3857 - x
	}

	m.width = width
	m.dx = 1 / (x / float64(m.width))
	m.height = int(math.Abs(m.ury-m.lly)*m.dx) + 1

	// shift the map area into SVG view based on upper left corner.
	m.xshift = m.llx * -1
	m.yshift = m.ury * -1

	m.width = width

	m.region = b.region

	mPix := x / float64(m.width) // m per pixel

	// region 1 has higher zoom data.
	switch m.region {
	case 1:
		switch {
		case mPix < 50:
			m.zoom = 3
		case mPix < 500:
			m.zoom = 2
		case mPix < 5000.0:
			m.zoom = 1
		default:
			m.zoom = 0
		}

	case 0:
		m.zoom = 0
	}

	return
}
