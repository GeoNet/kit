package weft

import (
	"testing"
)

func TestLoadAssets(t *testing.T) {

	testData := []struct {
		testName       string
		filename       string
		prefix         string
		expectedResult *asset
	}{
		{
			"Load CSS file",
			"testdata/leaflet.css",
			"testdata",
			&asset{
				path:       "/leaflet.css",
				hashedPath: "/07800b98-leaflet.css",
				mime:       "text/css",
				sri:        "sha384-9oKBsxAYdVVBJcv3hwG8RjuoJhw9GwYLqXdQRDxi2q0t1AImNHOap8y6Qt7REVd4",
			},
		},
		{
			"Load JS file",
			"testdata/gnss-map-plot.js",
			"testdata",
			&asset{
				path:       "/gnss-map-plot.js",
				hashedPath: "/e83a0b0f-gnss-map-plot.js",
				mime:       "text/javascript",
				sri:        "sha384-haxRijtRHhpn6nbt+JNpioqOj0AwB+THIaVdUZ34R9sQrQL2vmf/pn6GPnQq+AI1",
			},
		},
	}
	// SRI hash calculated with `openssl dgst -sha384 -binary leaflet.css | openssl base64 -A`
	// from https://www.srihash.org/

	for _, d := range testData {

		t.Run(d.testName, func(t *testing.T) {

			a, err := loadAsset(d.filename, d.prefix)
			if err != nil {
				t.Error(err)
			}

			if a.path != d.expectedResult.path {
				t.Errorf("expected path %s instead got %s", d.expectedResult.path, a.path)
			}

			if a.hashedPath != d.expectedResult.hashedPath {
				t.Errorf("expected hashed path %s instead got %s", d.expectedResult.hashedPath, a.hashedPath)
			}

			if a.mime != d.expectedResult.mime {
				t.Errorf("expected mime %s instead got %s", d.expectedResult.mime, a.mime)
			}

			if a.sri != d.expectedResult.sri {
				t.Errorf("expected sri hash %s instead got %s", d.expectedResult.sri, a.sri)
			}
		})
	}
}

func TestCreateSubResourceTag(t *testing.T) {
	err := initAssets("testdata", "testdata")
	if err != nil {
		t.Error(err)
	}

	work := []struct {
		nonce    string
		path     string
		expected string
	}{
		{"",
			"testdata/leaflet.css",
			`<link rel="stylesheet" href="/07800b98-leaflet.css" integrity="sha384-9oKBsxAYdVVBJcv3hwG8RjuoJhw9GwYLqXdQRDxi2q0t1AImNHOap8y6Qt7REVd4">`,
		},
		{"abcdefgh",
			"testdata/gnss-map-plot.js",
			`<script src="/e83a0b0f-gnss-map-plot.js" type="text/javascript" integrity="sha384-haxRijtRHhpn6nbt+JNpioqOj0AwB+THIaVdUZ34R9sQrQL2vmf/pn6GPnQq+AI1" nonce="abcdefgh"></script>`,
		},
	}

	for _, w := range work {
		t.Run(w.path, func(t *testing.T) {

			a, err := loadAsset(w.path, "testdata")
			if err != nil {
				t.Error(err)
			}

			tag, err := createSubResourceTag(a, w.nonce)
			if err != nil {
				t.Fatalf("Couldn't create embedded resource tag: %v", err)
			}

			if tag != w.expected {
				t.Fatalf("output tag '%v' did not equal epected '%v'", tag, w.expected)
			}
		})
	}
}
