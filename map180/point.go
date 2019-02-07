package map180

type Point struct {
	Longitude, Latitude   float64
	x, y                  int
	Value                 float64 // Points can be sorted by Value.
	Stroke, Fill, Opacity string  // Optional for any later drawing funcs
	Size                  int     // Optional for any later drawing funcs
}

type Points []Point

func (p Points) Len() int           { return len(p) }
func (p Points) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Points) Less(i, j int) bool { return p[i].Value < p[j].Value }

/*
Returns the SVG x coord for p.  Map() must be called to set this value.
*/
func (p *Point) X() int {
	return p.x
}

/*
Returns the SVG y coord for p.  Map() must be called to set this value.
*/
func (p *Point) Y() int {
	return p.y
}

// sets x,y on 3857 from long lat on 4326.  Allows for crossing 180.
func (m map3857) point3857(p *Point) (err error) {
	// map crosses 180 and the point is to the right
	switch m.crossesCentral && p.Longitude > -180.0 && p.Longitude < 0.0 {
	case true:
		err = db.QueryRow(svgPointQuery, p.Longitude, p.Latitude, width3857-m.llx, m.yshift, m.dx, m.dx).Scan(&p.x, &p.y)
	case false:
		err = db.QueryRow(svgPointQuery, p.Longitude, p.Latitude, m.xshift, m.yshift, m.dx, m.dx).Scan(&p.x, &p.y)
	}

	return
}
