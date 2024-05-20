/*
map180 draws SVG maps on EPSG3857.  It handles maps that cross the 180 meridian and allows
high zoom level data in a specified region.  Land and lake layer queries are cached for  speed.
*/
package map180

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/golang/groupcache"
)

const (
	width3857    = 40075016.6855785
	Width3857    = 40075016.6855785
	left3857     = width3857 / 2
	right3857    = left3857 * -1.0
	svgPolyQuery = `select array_to_string(array_agg(
		ST_AsSVG(st_transScale(ST_Intersection(ST_MakeEnvelope($1,$2,$3,$4, 3857), geom), $5, $6, $7, $8 ),0,1)
		), ' ')
	 from public.map180_layers where zoom = $9 and type = $10 and region = $11`
)

var (
	empty          *regexp.Regexp
	db             *sql.DB
	zoomRegion     bbox
	mapBounds      []bbox
	namedMapBounds map[string]bbox
	mapLayers      *groupcache.Group
	landLayers     *groupcache.Group
	lakeLayers     *groupcache.Group
	dummyCtx       = context.TODO()
)

type Map180 struct {
}

type Raw struct {
	Land, Lakes    string  // land or lake polygons as strings for use in the d property of an SVG path.
	Height, Width  int     // the height and width for the map.
	DX             float64 // for scaling.
	XShift, YShift float64 // for translation.
	LLX            float64 // for translation.
	CrossesCentral bool    // true if the map crosses the 180 meridian.
}

type map3857 struct {
	llx, lly, urx, ury float64
	dx                 float64
	xshift, yshift     float64
	width, height      int
	crossesCentral     bool
	zoom               int
	region             int
}

func init() {
	var err error
	empty, err = regexp.Compile(`^\s*$`)
	if err != nil {
		log.Fatal(err)
	}
}

/*
Init returns an initialised Map180. d must have access to the map180 tables in the
public schema.  cacheBytes is the size of the memory cache
for land and lake layers.  region must be a valid Region.

Example:

	wm, err = map180.Init(db, map180.Region(`newzealand`), 256000000)
*/
func Init(d *sql.DB, region Region, cacheBytes int64) (*Map180, error) {
	w := &Map180{}

	var err error
	if _, ok := allZoomRegions[region]; !ok {
		err = fmt.Errorf("invalid region (allZoomRegions): %s", region)
		return w, err
	}

	if _, ok := allMapBounds[region]; !ok {
		err = fmt.Errorf("invalid region (allMapBounds): %s", region)
		return w, err
	}

	if _, ok := allNamedMapBounds[region]; !ok {
		err = fmt.Errorf("invalid region (allNamedMapBounds): %s", region)
		return w, err
	}

	zoomRegion = allZoomRegions[region]
	mapBounds = allMapBounds[region]
	namedMapBounds = allNamedMapBounds[region]

	db = d

	mapLayers = groupcache.NewGroup("mapLayers", cacheBytes, groupcache.GetterFunc(layerGetter))
	landLayers = groupcache.NewGroup("landLayers", cacheBytes, groupcache.GetterFunc(landGetter))
	lakeLayers = groupcache.NewGroup("lakeLayers", cacheBytes, groupcache.GetterFunc(lakeGetter))

	return w, nil
}

/*
SVG draws an SVG image showing a map of markers.  The returned map uses EPSG3857.
Width is the SVG image width in pixels (height is calculated).
If boundingBox is the empty string then the map bounds are calculated from the markers.
See ValidBbox for boundingBox options.
*/
func (w *Map180) SVG(boundingBox string, width int, markers []Marker, insetBbox string) (buf bytes.Buffer, err error) {
	// If the bbox is zero type then figure it out from the markers.
	var b bbox
	if boundingBox == "" {
		b, err = newBboxFromMarkers(markers)
		if err != nil {
			return
		}
	} else {
		b, err = newBbox(boundingBox)
		if err != nil {
			return
		}
	}

	m, err := b.newMap3857(width)
	if err != nil {
		return
	}

	buf.WriteString(`<?xml version="1.0"?>`)
	buf.WriteString(fmt.Sprintf("<svg height=\"%d\" width=\"%d\" xmlns=\"http://www.w3.org/2000/svg\">",
		m.height, m.width))
	if b.title != "" {
		buf.WriteString(`<title>Map of ` + b.title + `.</title>`)
	} else {
		buf.WriteString(`<title>Map of ` + boundingBox + `.</title>`)
	}

	// Get the land and lakes layers from the cache.  This creates them
	// if they haven't been cached already.
	var landLakes string

	err = mapLayers.Get(dummyCtx, m.toKey(), groupcache.StringSink(&landLakes))
	if err != nil {
		return
	}

	buf.WriteString(landLakes)

	if insetBbox != "" {
		var inset bbox
		inset, err = newBbox(insetBbox)
		if err != nil {
			return
		}
		var in map3857
		in, err = inset.newMap3857(80)
		if err != nil {
			return
		}

		var insetMap string
		err = mapLayers.Get(dummyCtx, in.toKey(), groupcache.StringSink(&insetMap))
		if err != nil {
			return
		}

		// use 2 markers to put a the main map bbox as a rect
		ibboxul := NewMarker(b.llx, b.ury, ``, ``, ``)
		err = in.marker3857(&ibboxul)
		if err != nil {
			return
		}

		ibboxlr := NewMarker(b.urx, b.lly, ``, ``, ``)
		err = in.marker3857(&ibboxlr)
		if err != nil {
			return
		}

		// if bbox rect is tiny make it bigger and shift it a little.
		iw := int(ibboxlr.x - ibboxul.x)
		if iw < 5 {
			iw = 5
			ibboxul.x = ibboxul.x - 2
		}

		ih := int(ibboxlr.y - ibboxul.y)
		if ih < 5 {
			ih = 5
			ibboxul.y = ibboxul.y - 2
		}

		buf.WriteString(fmt.Sprintf("<g transform=\"translate(10,10)\"><rect x=\"-3\" y=\"-3\" width=\"%d\" height=\"%d\" rx=\"10\" ry=\"10\" fill=\"white\"/>",
			in.width+6, in.height+6))

		buf.WriteString(insetMap)

		buf.WriteString(fmt.Sprintf("<rect x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" fill=\"red\" opacity=\"0.5\"/>",
			int(ibboxul.x), int(ibboxul.y), iw, ih) + `</g>`)

	} // end of inset

	err = m.drawMarkers(markers, &buf)
	if err != nil {
		return
	}

	var short bool
	if m.width < 250 {
		short = true
	}
	labelMarkers(markers, m.width-5, m.height-5, `end`, 12, short, &buf)

	buf.WriteString("</svg>")
	return
}

/*
Map draws an SVG image to buf for the bbox regions.  The returned map uses EPSG3857.
Width is the SVG image width in pixels (height is calculated). pts haz X Y and values
initialised for later drawing.  The SVG in buf is not closed.  See ValidBbox for boundingBox options.
*/
func (w *Map180) Map(boundingBox string, width int, pts Points, insetBbox string, buf *bytes.Buffer) (err error) {
	// If the bbox is zero type then figure it out from the markers.
	var b bbox
	b, err = newBbox(boundingBox)
	if err != nil {
		return
	}

	m, err := b.newMap3857(width)
	if err != nil {
		return
	}

	buf.WriteString(`<?xml version="1.0"?>`)
	buf.WriteString(fmt.Sprintf("<svg  viewBox=\"0 0 %d %d\"  xmlns=\"http://www.w3.org/2000/svg\">",
		m.width, m.height))
	if b.title != "" {
		buf.WriteString(`<title>Map of ` + b.title + `.</title>`)
	} else {
		buf.WriteString(`<title>Map of ` + boundingBox + `.</title>`)
	}

	// Get the land and lakes layers from the cache.  This creates them
	// if they haven't been cached already.
	var landLakes string

	err = mapLayers.Get(dummyCtx, m.toKey(), groupcache.StringSink(&landLakes))
	if err != nil {
		return
	}

	buf.WriteString(landLakes)

	if insetBbox != "" {
		var inset bbox
		inset, err = newBbox(insetBbox)
		if err != nil {
			return
		}
		var in map3857
		in, err = inset.newMap3857(80)
		if err != nil {
			return
		}

		var insetMap string
		err = mapLayers.Get(dummyCtx, in.toKey(), groupcache.StringSink(&insetMap))
		if err != nil {
			return
		}

		// use 2 markers to put a the main map bbox as a rect
		ibboxul := NewMarker(b.llx, b.ury, ``, ``, ``)
		err = in.marker3857(&ibboxul)
		if err != nil {
			return
		}

		ibboxlr := NewMarker(b.urx, b.lly, ``, ``, ``)
		err = in.marker3857(&ibboxlr)
		if err != nil {
			return
		}

		// if bbox rect is tiny make it bigger and shift it a little.
		iw := int(ibboxlr.x - ibboxul.x)
		if iw < 5 {
			iw = 5
			ibboxul.x = ibboxul.x - 2
		}

		ih := int(ibboxlr.y - ibboxul.y)
		if ih < 5 {
			ih = 5
			ibboxul.y = ibboxul.y - 2
		}

		buf.WriteString(fmt.Sprintf("<g transform=\"translate(10,10)\"><rect x=\"-3\" y=\"-3\" width=\"%d\" height=\"%d\" rx=\"10\" ry=\"10\" fill=\"white\"/>",
			in.width+6, in.height+6))

		buf.WriteString(insetMap)

		buf.WriteString(fmt.Sprintf("<rect x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" fill=\"red\" opacity=\"0.5\"/>",
			int(ibboxul.x), int(ibboxul.y), iw, ih) + `</g>`)

	} // end of inset

	for i := range pts {
		if pts[i].Latitude <= 85.0 && pts[i].Latitude >= -85.0 {
			if err = m.point3857(&pts[i]); err != nil {
				return
			}
		}
	}

	return
}

/*
MapRaw returns a Raw struct which can be used for drawing SVG maps e.g.,

	raw, err := wm.MapRaw(bbox, width)
	b := bytes.Buffer
	b.WriteString(`<?xml version="1.0"?>`)
	b.WriteString(fmt.Sprintf("<svg  viewBox=\"0 0 %d %d\"  xmlns=\"http://www.w3.org/2000/svg\">",
	raw.Width, raw.Height))
	b.WriteString(fmt.Sprintf("<rect x=\"0\" y=\"0\" width=\"%d\" height=\"%d\" style=\"fill: azure\"/>", raw.Width, raw.Height))
	b.WriteString(fmt.Sprintf("<path style=\"fill: wheat; stroke-width: 1; stroke-linejoin: round; stroke: lightslategrey\" d=\"%s\"/>", raw.Land))
	b.WriteString(fmt.Sprintf("<path style=\"fill: azure; stroke-width: 1; stroke-linejoin: round; stroke: lightslategrey\" d=\"%s\"/>", raw.Lakes))
	b.WriteString("</svg>")

The other properties can be used to scale and translate for drawing on the map e.g.,

	    type point struct {
		latitude, longitude float64
		x, y                float64
	    }

	    // create pts []point with x,y, set to EPSG3857 and latitude longitude EPSG4326
	    // range of pts ...

	    switch raw.CrossesCentral && p.longitude > -180.0 && p.longitude < 0.0 {
		case true:
			p.x = (p.x + map180.Width3857 - raw.LLX) * raw.DX
			p.y = (p.y - math.Abs(raw.YShift)) * raw.DX
		case false:
			p.x = (p.x - math.Abs(raw.XShift)) * raw.DX
			p.y = (p.y - math.Abs(raw.YShift)) * raw.DX
	    }

	    // draw p on SVG.
*/
func (w *Map180) MapRaw(boundingBox string, width int) (mr Raw, err error) {
	var b bbox
	if b, err = newBbox(boundingBox); err != nil {
		return
	}

	var m map3857
	if m, err = b.newMap3857(width); err != nil {
		return
	}

	mr.Height = m.height
	mr.Width = m.width
	mr.DX = m.dx
	mr.CrossesCentral = m.crossesCentral
	mr.LLX = m.llx
	mr.YShift = m.yshift
	mr.XShift = m.xshift

	// Get the land and lakes layers from the cache.  This creates them
	// if they haven't been cached already.
	if err = landLayers.Get(dummyCtx, m.toKey(), groupcache.StringSink(&mr.Land)); err != nil {
		return
	}

	if err = lakeLayers.Get(dummyCtx, m.toKey(), groupcache.StringSink(&mr.Lakes)); err != nil {
		return
	}

	return
}

func (m *map3857) nePolySVG(zoom int, layer int) (string, error) {
	// db errors are ignored. It is not an error for there to be no data in the bbox.
	// should be possible to check for an empty row error but the pg driver
	// seems to be confused by the gist index on map_layers when there is no data at the requested
	// zoom / layer.
	//
	// keep err in the signature in case there is a fix later

	var l string
	var r string

	switch m.crossesCentral {
	case true:
		//  things to the left of 180.
		_ = db.QueryRow(svgPolyQuery,
			m.llx, m.lly, left3857, m.ury, m.xshift, m.yshift, m.dx, m.dx, zoom, layer, m.region).Scan(&l)
		//  things to the right of 180 and shift them over.
		_ = db.QueryRow(svgPolyQuery,
			right3857, m.lly, m.urx, m.ury, width3857-m.llx, m.yshift, m.dx, m.dx, zoom, layer, m.region).Scan(&r)
	case false:
		_ = db.QueryRow(svgPolyQuery,
			m.llx, m.lly, m.urx, m.ury, m.xshift, m.yshift, m.dx, m.dx, zoom, layer, m.region).Scan(&l)
	}

	p := l + r

	// recurse if there was no data.  This allows for uneven data loading at a region + zoom.
	if zoom > 0 && empty.MatchString(p) {
		p, _ = m.nePolySVG(zoom-1, layer)
	}

	return p, nil
}

// Functions for map layers with groupcache.

func layerGetter(ctx context.Context, key string, dest groupcache.Sink) error {
	m, err := fromKey(key)
	if err != nil {
		return err
	}

	land, err := m.nePolySVG(m.zoom, 0)
	if err != nil {
		return err
	}

	lakes, err := m.nePolySVG(m.zoom, 1)
	if err != nil {
		return err
	}

	l, err := m.labels()
	if err != nil {
		return err
	}

	var b bytes.Buffer

	b.WriteString(fmt.Sprintf("<path fill=\"wheat\" stroke-width=\"1\"  stroke-linejoin=\"round\" stroke=\"lightslategrey\" d=\"%s\"/>", land))
	b.WriteString(fmt.Sprintf("<path fill=\"azure\" stroke-width=\"1\"  stroke=\"lightslategrey\" d=\"%s\"/>", lakes))
	b.WriteString(labelsToSVG(l))

	return dest.SetString(b.String())
}

func landGetter(ctx context.Context, key string, dest groupcache.Sink) error {
	m, err := fromKey(key)
	if err != nil {
		return err
	}

	land, err := m.nePolySVG(m.zoom, 0)
	if err != nil {
		return err
	}

	return dest.SetString(land)
}

func lakeGetter(ctx context.Context, key string, dest groupcache.Sink) error {
	m, err := fromKey(key)
	if err != nil {
		return err
	}

	lakes, err := m.nePolySVG(m.zoom, 1)
	if err != nil {
		return err
	}

	return dest.SetString(lakes)
}

func (m *map3857) toKey() string {
	return fmt.Sprintf("%f,%f,%f,%f,%f,%f,%f,%d,%d,%d,%d,%t",
		m.llx, m.lly, m.urx, m.ury, m.dx, m.xshift, m.yshift, m.width, m.height, m.zoom, m.region, m.crossesCentral)
}

func fromKey(key string) (m map3857, err error) {
	k := strings.Split(key, ",")
	if len(k) != 12 {
		err = fmt.Errorf("Wrong length for key exptected 12 got %d", len(k))
		return
	}

	m.llx, err = strconv.ParseFloat(k[0], 64)
	if err != nil {
		return
	}

	m.lly, err = strconv.ParseFloat(k[1], 64)
	if err != nil {
		return
	}

	m.urx, err = strconv.ParseFloat(k[2], 64)
	if err != nil {
		return
	}

	m.ury, err = strconv.ParseFloat(k[3], 64)
	if err != nil {
		return
	}

	m.dx, err = strconv.ParseFloat(k[4], 64)
	if err != nil {
		return
	}

	m.xshift, err = strconv.ParseFloat(k[5], 64)
	if err != nil {
		return
	}

	m.yshift, err = strconv.ParseFloat(k[6], 64)
	if err != nil {
		return
	}

	m.width, err = strconv.Atoi(k[7])
	if err != nil {
		return
	}

	m.height, err = strconv.Atoi(k[8])
	if err != nil {
		return
	}

	m.zoom, err = strconv.Atoi(k[9])
	if err != nil {
		return
	}

	m.region, err = strconv.Atoi(k[10])
	if err != nil {
		return
	}

	m.crossesCentral, err = strconv.ParseBool(k[11])
	if err != nil {
		return
	}

	return
}
