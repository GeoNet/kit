package wgs84

import (
	"math"
	"sort"
	"testing"
)

func TestClosestNZ(t *testing.T) {
	in := []struct {
		id                  string
		longitude, latitude float64
		name                string
		bearing, distance   float64
	}{
		{id: loc(), longitude: 176.0479775, latitude: -39.62710223, name: "Taihape", distance: 22.077, bearing: 74.651},
		{id: loc(), longitude: 174.8382288, latitude: -40.90474677, name: "Paraparaumu", distance: 13.733522, bearing: 277.031977},
	}

	for _, v := range in {
		l, err := ClosestNZ(v.latitude, v.longitude)
		if err != nil {
			t.Error(err)
		}

		if v.name != l.Name {
			t.Errorf("%s got name %s expected %s", v.id, l.Name, v.name)
		}

		if math.Abs(v.bearing-l.Bearing) > 0.001 {
			t.Errorf("%s got bearing %f expected %f", v.id, l.Bearing, v.bearing)
		}

		if math.Abs(v.distance-l.Distance) > 0.001 {
			t.Errorf("%s got distance %f expected %f", v.id, l.Distance, v.distance)
		}
	}
}

func TestLocalitiesNZ(t *testing.T) {
	in := []struct {
		id                  string
		longitude, latitude float64
		name                string
		bearing, distance   float64
	}{
		{id: loc(), longitude: 176.0479775, latitude: -39.62710223, name: "Taihape", distance: 22.077, bearing: 74.651},
		{id: loc(), longitude: 174.8382288, latitude: -40.90474677, name: "Paraparaumu", distance: 13.733522, bearing: 277.031977},
	}

	for _, v := range in {
		l, err := LocalitiesNZ(v.latitude, v.longitude)
		if err != nil {
			t.Error(err)
		}

		sort.Sort(ByDistance(l))

		if v.name != l[0].Name {
			t.Errorf("%s got name %s expected %s", v.id, l[0].Name, v.name)
		}

		if math.Abs(v.bearing-l[0].Bearing) > 0.001 {
			t.Errorf("%s got bearing %f expected %f", v.id, l[0].Bearing, v.bearing)
		}

		if math.Abs(v.distance-l[0].Distance) > 0.001 {
			t.Errorf("%s got distance %f expected %f", v.id, l[0].Distance, v.distance)
		}
	}
}

func TestCompass(t *testing.T) {
	in := []struct {
		id string
		b  float64
		c  string
	}{
		{id: loc(), b: 1.0, c: "north"},
		{id: loc(), b: 45.0, c: "north-east"},
		{id: loc(), b: 95.0, c: "east"},
		{id: loc(), b: 125.0, c: "south-east"},
		{id: loc(), b: 160.0, c: "south"},
		{id: loc(), b: 220.0, c: "south-west"},
		{id: loc(), b: 270.0, c: "west"},
		{id: loc(), b: 295.0, c: "north-west"},
		{id: loc(), b: 340.0, c: "north"},
		{id: loc(), b: 9999.0, c: "north"},
	}

	for _, v := range in {
		if Compass(v.b) != v.c {
			t.Errorf("%s for %f got compass %s expected %s", v.id, v.b, Compass(v.b), v.c)
		}
	}
}

func TestCompassShort(t *testing.T) {
	in := []struct {
		id string
		b  float64
		c  string
	}{
		{id: loc(), b: 1.0, c: "N"},
		{id: loc(), b: 45.0, c: "NE"},
		{id: loc(), b: 95.0, c: "E"},
		{id: loc(), b: 125.0, c: "SE"},
		{id: loc(), b: 160.0, c: "S"},
		{id: loc(), b: 220.0, c: "SW"},
		{id: loc(), b: 270.0, c: "W"},
		{id: loc(), b: 295.0, c: "NW"},
		{id: loc(), b: 340.0, c: "N"},
		{id: loc(), b: 9999.0, c: "N"},
	}

	for _, v := range in {
		if CompassShort(v.b) != v.c {
			t.Errorf("%s for %f got compass %s expected %s", v.id, v.b, CompassShort(v.b), v.c)
		}
	}
}

func TestLocality_Description(t *testing.T) {
	in := []struct {
		id string
		l  Locality
		e  string
	}{
		{id: loc(), l: Locality{Name: "test", Bearing: 1.0, Distance: 4.5}, e: "Within 5 km of test"},
		{id: loc(), l: Locality{Name: "test", Bearing: 1.0, Distance: 10.5}, e: "10 km north of test"},
		{id: loc(), l: Locality{Name: "test", Bearing: 1.0, Distance: 14.5}, e: "10 km north of test"},
		{id: loc(), l: Locality{Name: "test", Bearing: 1.0, Distance: 15.5}, e: "15 km north of test"},
	}

	for _, v := range in {
		if v.l.Description() != v.e {
			t.Errorf("%s expected description '%s' got '%s'", v.id, v.e, v.l.Description())
		}
	}
}

func TestLocality_DescriptionShort(t *testing.T) {
	in := []struct {
		id string
		l  Locality
		e  string
	}{
		{id: loc(), l: Locality{Name: "test", Bearing: 1.0, Distance: 4.5}, e: "Within 5 km of test"},
		{id: loc(), l: Locality{Name: "test", Bearing: 1.0, Distance: 10.5}, e: "10 km N of test"},
		{id: loc(), l: Locality{Name: "test", Bearing: 1.0, Distance: 14.5}, e: "10 km N of test"},
		{id: loc(), l: Locality{Name: "test", Bearing: 1.0, Distance: 15.5}, e: "15 km N of test"},
	}

	for _, v := range in {
		if v.l.DescriptionShort() != v.e {
			t.Errorf("%s expected description '%s' got '%s'", v.id, v.e, v.l.DescriptionShort())
		}
	}
}
