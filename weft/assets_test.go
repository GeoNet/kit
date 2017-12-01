package weft

import "testing"

func TestLoadAsset(t *testing.T) {
	a, err := loadAsset("testdata/leaflet.css", "testdata")
	if err != nil {
		t.Error(err)
	}

	if a.path != "/leaflet.css" {
		t.Errorf("expected path /leaflet.css got %s", a.path)
	}

	if a.hashedPath != "/07800b98-leaflet.css" {
		t.Errorf("expected hashed path /07800b98-leaflet.css got %s", a.hashedPath)
	}

	if a.mime != "text/css" {
		t.Errorf("expected mime text/css got %s", a.mime)
	}
}
