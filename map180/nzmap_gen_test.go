//go:build generate
// +build generate

package map180

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/GeoNet/kit/cfg"
	_ "github.com/lib/pq"
)

/*
Use these tests to generate files in nzmap. Firstly, start up a Postgres DB (recommended to use the one
defined in the haz repo: https://github.com/GeoNet/haz/tree/main/etc). e.g.:

	docker build -t haz-db .
	docker run -p 5432:5432 -e POSTGRES_PASSWORD=test haz-db:latest

Start psql:

	psql -h localhost -U postgres -d hazard

Then copy and paste the following commands:

    drop table public.map180_layers;
    drop table public.map180_labels;

    CREATE TABLE public.map180_layers (
            mapPK SERIAL PRIMARY KEY,
            region INT NOT NULL,
            zoom INT NOT NULL,
            type INT NOT NULL
    );

    SELECT addgeometrycolumn('public', 'map180_layers', 'geom', 3857, 'MULTIPOLYGON', 2);

    CREATE INDEX ON public.map180_layers (zoom);
    CREATE INDEX ON public.map180_layers (region);
    CREATE INDEX ON public.map180_layers (type);
    CREATE INDEX ON public.map180_layers USING gist (geom);

    GRANT SELECT ON public.map180_layers TO PUBLIC;

    CREATE TABLE public.map180_labels (
            labelPK SERIAL PRIMARY KEY,
            zoom INT NOT NULL,
            type INT NOT NULL,
            name text
    );

    SELECT addgeometrycolumn('public', 'map180_labels', 'geom', 3857, 'POINT', 2);

    CREATE INDEX ON public.map180_labels (zoom);
    CREATE INDEX ON public.map180_labels USING gist (geom);

    GRANT SELECT ON public.map180_labels TO PUBLIC;


Then add the coast line data:

	mkdir work
    cp ./data/new_zealand_map_layers.ddl.gz ./work/new_zealand_map_layers.ddl.gz
    gunzip ./work/new_zealand_map_layers.ddl.gz
    psql -h localhost -U postgres -d hazard -f ./work/new_zealand_map_layers.ddl
    rm -rf ./work


Export the environment variables in env.list, then generate and format the nzmap files:

    go test -tags generate
    goimports -w nzmap/
*/

const (
	iconWidth   = 100
	mediumWidth = 500
	// used for drawing landmasses.
	landPath = "<path class=\"nzmap-land\" fill=\"#f0f5f5\" stroke=\"#778899\" stroke-width=\"0.75\" stroke-linejoin=\"round\" d=\"%s\"/>"
	// used for drawing lakes ontop of landmasses.
	lakePath = "<path class=\"nzmap-lake\" fill=\"#ffffff\" stroke=\"#778899\" stroke-width=\"0.75\" stroke-linejoin=\"round\" d=\"%s\"/>"
	// fixed size SVG images (icons)
	fixed = "<?xml version=\"1.0\"?><svg height=\"%d\" width=\"%d\" xmlns=\"http://www.w3.org/2000/svg\"><title>Map of New Zealand.</title>"
	// responsive
	responsive = "<?xml version=\"1.0\"?><svg viewBox=\"0 0 %d %d\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\"><title>Map of New Zealand.</title>"
)

var wm *Map180

func setup(t *testing.T) {
	p, err := cfg.PostgresEnv()

	if err != nil {
		log.Fatalf("error reading DB config from the environment vars: %s", err)
	}

	db, err = sql.Open("postgres", p.Connection())
	if err != nil {
		log.Fatalf("ERROR: problem with DB config: %s", err)
	}

	db.SetMaxIdleConns(p.MaxIdle)
	db.SetMaxOpenConns(p.MaxOpen)

	err = db.Ping()
	if err != nil {
		log.Fatal("can't ping DB")
	}

	if wm == nil {
		wm, err = Init(db, Region(`newzealand`), 256000000)
		if err != nil {
			log.Fatalf("ERROR: problem with map180 init: %s", err)
		}
	}

}

func teardown() {
	db.Close()
}

// New Zealand icon map - lon lat grid at 0.1 degrees
func TestIconNewZealand(t *testing.T) {
	setup(t)
	defer teardown()

	// New Zealand icon map - lon lat grid at 0.1 degrees
	b, err := newBbox("165,-48,180,-34")
	if err != nil {
		t.Fatal("Getting bbox " + err.Error())
	}

	m, err := b.newMap3857(iconWidth)
	if err != nil {
		t.Fatal("Getting map " + err.Error())
	}

	land, err := m.nePolySVG(m.zoom, 0)
	if err != nil {
		t.Fatal("Getting land " + err.Error())
	}

	lakes, err := m.nePolySVG(m.zoom, 1)
	if err != nil {
		t.Fatal("Getting coast")
	}

	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf(fixed, m.height, m.width))
	buf.WriteString(fmt.Sprintf(landPath, land))
	buf.WriteString(fmt.Sprintf(lakePath, lakes))

	var out bytes.Buffer

	out.WriteString("package nzmap\n")
	out.WriteString("var nzIcon = `" + buf.String() + "`\n\n")
	out.WriteString("func init() {\n")

	for x := 165.0; x <= 180.0; x = x + 0.1 {
		for y := -48.0; y <= -34.0; y = y + 0.1 {
			p := NewMarker(x, y, "", "", "")
			m.marker3857(&p)
			out.WriteString(fmt.Sprintf("nzIconPts[%.f][%.f] = pt{x:%d, y:%d}\n", x*10-1650, y*10+480, int(p.x), int(p.y)))
		}
	}
	out.WriteString("}\n")

	err = os.WriteFile("nzmap/nzicon.go", out.Bytes(), 0644)
	if err != nil {
		t.Fatal(err)
	}
}

// New Zealand, Raoul, Chathams icon map.  1 degree grid.
// the bbox is slightly larger than the grid to make the height the same as other
// icon maps.
func TestIconNewZealandRaoulChathams(t *testing.T) {
	setup(t)
	defer teardown()

	b, err := newBbox("165,-48,-167.4,-20")
	if err != nil {
		t.Fatal("Getting bbox " + err.Error())
	}

	m, err := b.newMap3857(iconWidth)
	if err != nil {
		t.Fatal("Getting map " + err.Error())
	}

	land, err := m.nePolySVG(m.zoom, 0)
	if err != nil {
		t.Fatal("Getting land " + err.Error())
	}

	lakes, err := m.nePolySVG(m.zoom, 1)
	if err != nil {
		t.Fatal("Getting coast " + err.Error())
	}

	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf(fixed, m.height, m.width))
	buf.WriteString(fmt.Sprintf(landPath, land))
	buf.WriteString(fmt.Sprintf(lakePath, lakes))

	var out bytes.Buffer

	out.WriteString("package nzmap\n")
	out.WriteString("var nzrcIcon = `" + buf.String() + "`\n\n")
	out.WriteString("func init() {\n")

	var xs float64
	for x := 165.0; x <= 192; x++ {
		xs = x
		if x > 180 {
			xs = xs - 360
		}
		for y := -48.0; y <= -20.0; y++ {
			p := NewMarker(xs, y, "", "", "")
			m.marker3857(&p)
			out.WriteString(fmt.Sprintf("nzrcIconPts[%.f][%.f] = pt{x:%d, y:%d}\n", x-165, y+48, int(p.x), int(p.y)))
		}
	}
	out.WriteString("}\n")

	err = os.WriteFile("nzmap/nzrcicon.go", out.Bytes(), 0644)
	if err != nil {
		t.Fatal(err)
	}
}

// New Zealand, Southern Ocean icon map.  1 degree grid.
func TestIconNewZealandSouth(t *testing.T) {
	setup(t)
	defer teardown()

	b, err := newBbox("156,-55,180,-34")
	if err != nil {
		t.Fatal("Getting bbox " + err.Error())
	}

	m, err := b.newMap3857(iconWidth)
	if err != nil {
		t.Fatal("Getting map " + err.Error())
	}

	land, err := m.nePolySVG(m.zoom, 0)
	if err != nil {
		t.Fatal("Getting land " + err.Error())
	}

	lakes, err := m.nePolySVG(m.zoom, 1)
	if err != nil {
		t.Fatal("Getting coast " + err.Error())
	}

	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf(fixed, m.height, m.width))
	buf.WriteString(fmt.Sprintf(landPath, land))
	buf.WriteString(fmt.Sprintf(lakePath, lakes))

	var out bytes.Buffer

	out.WriteString("package nzmap\n")
	out.WriteString("var nzsIcon = `" + buf.String() + "`\n\n")
	out.WriteString("func init() {\n")

	for x := 156.0; x <= 180.0; x++ {
		for y := -55.0; y <= -34.0; y++ {
			p := NewMarker(x, y, "", "", "")
			m.marker3857(&p)
			out.WriteString(fmt.Sprintf("nzsIconPts[%.f][%.f] = pt{x:%d, y:%d}\n", x-156, y+55, int(p.x), int(p.y)))
		}
	}
	out.WriteString("}\n")

	err = os.WriteFile("nzmap/nzsicon.go", out.Bytes(), 0644)
	if err != nil {
		t.Fatal(err)
	}
}

// New Zealand map - lon lat grid at 0.1 degrees
func TestNewZealand(t *testing.T) {
	setup(t)
	defer teardown()

	b, err := newBbox("165,-48,180,-34")
	if err != nil {
		t.Fatal("Getting bbox " + err.Error())
	}

	m, err := b.newMap3857(mediumWidth)
	if err != nil {
		t.Fatal("Getting map " + err.Error())
	}

	land, err := m.nePolySVG(m.zoom, 0)
	if err != nil {
		t.Fatal("Getting land " + err.Error())
	}

	lakes, err := m.nePolySVG(m.zoom, 1)
	if err != nil {
		t.Fatal("Getting coast")
	}

	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf(responsive, m.width, m.height))
	buf.WriteString(fmt.Sprintf(landPath, land))
	buf.WriteString(fmt.Sprintf(lakePath, lakes))

	var out bytes.Buffer

	out.WriteString("package nzmap\n")
	out.WriteString("var nzMedium = `" + buf.String() + "`\n\n")
	out.WriteString("func init() {\n")

	for x := 165.0; x <= 180.0; x = x + 0.1 {
		for y := -48.0; y <= -34.0; y = y + 0.1 {
			p := NewMarker(x, y, "", "", "")
			m.marker3857(&p)
			out.WriteString(fmt.Sprintf("nzMediumPts[%.f][%.f] = pt{x:%d, y:%d}\n", x*10-1650, y*10+480, int(p.x), int(p.y)))
		}
	}
	out.WriteString("}\n")

	err = os.WriteFile("nzmap/nzmedium.go", out.Bytes(), 0644)
	if err != nil {
		t.Fatal(err)
	}
}

// New Zealand, Raoul, Chathams map.
// the bbox is slightly larger than the grid to make the height the same as other
// icon maps.
func TestNewZealandRaoulChathams(t *testing.T) {
	setup(t)
	defer teardown()

	b, err := newBbox("165,-48,-167.4,-20")
	if err != nil {
		t.Fatal("Getting bbox " + err.Error())
	}

	m, err := b.newMap3857(mediumWidth)
	if err != nil {
		t.Fatal("Getting map " + err.Error())
	}

	land, err := m.nePolySVG(m.zoom, 0)
	if err != nil {
		t.Fatal("Getting land " + err.Error())
	}

	lakes, err := m.nePolySVG(m.zoom, 1)
	if err != nil {
		t.Fatal("Getting coast")
	}

	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf(responsive, m.width, m.height))
	buf.WriteString(fmt.Sprintf(landPath, land))
	// TODO no lakes are turning up?
	buf.WriteString(fmt.Sprintf(lakePath, lakes))

	var out bytes.Buffer

	out.WriteString("package nzmap\n")
	out.WriteString("var nzrcMedium = `" + buf.String() + "`\n\n")
	out.WriteString("func init() {\n")

	var xs float64
	for x := 165.0; x <= 192; x++ {
		xs = x
		if x > 180 {
			xs = xs - 360
		}
		for y := -48.0; y <= -20.0; y++ {
			p := NewMarker(xs, y, "", "", "")
			m.marker3857(&p)
			out.WriteString(fmt.Sprintf("nzrcMediumPts[%.f][%.f] = pt{x:%d, y:%d}\n", x-165, y+48, int(p.x), int(p.y)))
		}
	}
	out.WriteString("}\n")

	err = os.WriteFile("nzmap/nzrcmedium.go", out.Bytes(), 0644)
	if err != nil {
		t.Fatal(err)
	}
}

// New Zealand, Southern Ocean map.  1 degree grid.
func TestNewZealandSouth(t *testing.T) {
	setup(t)
	defer teardown()

	b, err := newBbox("156,-55,180,-34")
	if err != nil {
		t.Fatal("Getting bbox " + err.Error())
	}

	m, err := b.newMap3857(mediumWidth)
	if err != nil {
		t.Fatal("Getting map " + err.Error())
	}

	land, err := m.nePolySVG(m.zoom, 0)
	if err != nil {
		t.Fatal("Getting land " + err.Error())
	}

	lakes, err := m.nePolySVG(m.zoom, 1)
	if err != nil {
		t.Fatal("Getting coast " + err.Error())
	}

	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf(responsive, m.width, m.height))
	buf.WriteString(fmt.Sprintf(landPath, land))
	buf.WriteString(fmt.Sprintf(lakePath, lakes))

	var out bytes.Buffer

	out.WriteString("package nzmap\n")
	out.WriteString("var nzsMedium = `" + buf.String() + "`\n\n")
	out.WriteString("func init() {\n")

	for x := 156.0; x <= 180.0; x++ {
		for y := -55.0; y <= -34.0; y++ {
			p := NewMarker(x, y, "", "", "")
			m.marker3857(&p)
			out.WriteString(fmt.Sprintf("nzsMediumPts[%.f][%.f] = pt{x:%d, y:%d}\n", x-156, y+55, int(p.x), int(p.y)))
		}
	}
	out.WriteString("}\n")

	err = os.WriteFile("nzmap/nzsmedium.go", out.Bytes(), 0644)
	if err != nil {
		t.Fatal(err)
	}
}

// New Zealand, Raoul
func TestNewZealandRaoul(t *testing.T) {
	setup(t)
	defer teardown()

	b, err := newBbox("165,-48,-174,-27")
	if err != nil {
		t.Fatal("Getting bbox " + err.Error())
	}

	m, err := b.newMap3857(mediumWidth)
	if err != nil {
		t.Fatal("Getting map " + err.Error())
	}

	land, err := m.nePolySVG(m.zoom, 0)
	if err != nil {
		t.Fatal("Getting land " + err.Error())
	}

	lakes, err := m.nePolySVG(m.zoom, 1)
	if err != nil {
		t.Fatal("Getting coast")
	}

	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf(responsive, m.width, m.height))
	buf.WriteString(fmt.Sprintf(landPath, land))
	// TODO no lakes are turning up?
	buf.WriteString(fmt.Sprintf(lakePath, lakes))

	var out bytes.Buffer

	out.WriteString("package nzmap\n")
	out.WriteString("var nzrMedium = `" + buf.String() + "`\n\n")
	out.WriteString("func init() {\n")

	var xs float64
	for x := 165.0; x <= 186; x++ {
		xs = x
		if x > 180 {
			xs = xs - 360
		}
		for y := -48.0; y <= -27.0; y++ {
			p := NewMarker(xs, y, "", "", "")
			m.marker3857(&p)
			out.WriteString(fmt.Sprintf("nzrMediumPts[%.f][%.f] = pt{x:%d, y:%d}\n", x-165, y+48, int(p.x), int(p.y)))
		}
	}
	out.WriteString("}\n")

	err = os.WriteFile("nzmap/nzrmedium.go", out.Bytes(), 0644)
	if err != nil {
		t.Fatal(err)
	}
}
