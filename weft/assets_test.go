package weft

import (
	"testing"
)

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

	// Comparison calculated with `openssl dgst -sha384 -binary leaflet.css | openssl base64 -A` from https://www.srihash.org/
	if a.sri != "sha384-9oKBsxAYdVVBJcv3hwG8RjuoJhw9GwYLqXdQRDxi2q0t1AImNHOap8y6Qt7REVd4" {
		t.Errorf("got sri hash '%v' expected 'sha384-9oKBsxAYdVVBJcv3hwG8RjuoJhw9GwYLqXdQRDxi2q0t1AImNHOap8y6Qt7REVd4'", a.sri)
	}
}

func TestCreateSubResourceTag(t *testing.T) {
	err := initAssets("testdata", "testdata")
	if err != nil {
		t.Error(err)
	}

	work := []struct {
		path     string
		expected string
	}{
		{
			"testdata/leaflet.css",
			`<link rel="stylesheet" href="/07800b98-leaflet.css" integrity="sha384-9oKBsxAYdVVBJcv3hwG8RjuoJhw9GwYLqXdQRDxi2q0t1AImNHOap8y6Qt7REVd4">`,
		},
		{
			"testdata/gnss-map-plot.js",
			`<script src="/e83a0b0f-gnss-map-plot.js" type="text/javascript" integrity="sha384-haxRijtRHhpn6nbt+JNpioqOj0AwB+THIaVdUZ34R9sQrQL2vmf/pn6GPnQq+AI1"></script>`,
		},
	}

	for _, w := range work {
		t.Run(w.path, func(t *testing.T) {

			a, err := loadAsset(w.path, "testdata")
			if err != nil {
				t.Error(err)
			}

			tag, err := createSubResourceTag(a)
			if err != nil {
				t.Fatalf("Couldn't create embedded resource tag: %v", err)
			}

			if tag != w.expected {
				t.Fatalf("output tag '%v' did not equal epected '%v'", tag, w.expected)
			}
		})
	}
}
