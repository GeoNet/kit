package weft

import (
	"testing"
)

func TestLoadAssets(t *testing.T) {

	testData := []struct {
		testName       string
		filename       string
		expectedResult *asset
		expectedBytes  int
	}{
		{
			"Load CSS file",
			"testdata/leaflet.css",
			&asset{
				path:       "/leaflet.css",
				hashedPath: "/07800b98-leaflet.css",
				mime:       "text/css",
				fileType:   "css",
				sri:        "sha384-9oKBsxAYdVVBJcv3hwG8RjuoJhw9GwYLqXdQRDxi2q0t1AImNHOap8y6Qt7REVd4",
			},
			13429,
		},
		{
			"Load JS file",
			"testdata/gnss-map-plot.js",
			&asset{
				path:       "/gnss-map-plot.js",
				hashedPath: "/e83a0b0f-gnss-map-plot.js",
				mime:       "text/javascript",
				fileType:   "js",
				sri:        "sha384-haxRijtRHhpn6nbt+JNpioqOj0AwB+THIaVdUZ34R9sQrQL2vmf/pn6GPnQq+AI1",
			},
			8513,
		},
		{
			"Load MJS file",
			"testdata/test.mjs",
			&asset{
				path:       "/test.mjs",
				hashedPath: "/3616e4a4-test.mjs",
				mime:       "text/javascript",
				fileType:   "mjs",
				sri:        "sha384-yL9nK0JVp9FW9oAfkQ2kQC/9CcuAMK4vmyb8q+TY2SokmBFflIxJpZJ6Nk8Xqw5r",
			},
			64,
		},
	}
	// SRI hash calculated with `openssl dgst -sha384 -binary leaflet.css | openssl base64 -A`
	// from https://www.srihash.org/

	for _, d := range testData {

		t.Run(d.testName, func(t *testing.T) {

			a, err := loadAsset(d.filename, "testdata")
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

			if a.fileType != d.expectedResult.fileType {
				t.Errorf("expected file type %s instead got %s", d.expectedResult.fileType, a.fileType)
			}

			if a.sri != d.expectedResult.sri {
				t.Errorf("expected sri hash %s instead got %s", d.expectedResult.sri, a.sri)
			}

			if len(a.b) != d.expectedBytes {
				t.Errorf("expected %v bytes instead got %v", d.expectedBytes, len(a.b))
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
		attr     string
		path     string
		expected string
	}{
		{"",
			"",
			"testdata/leaflet.css",
			`<link rel="stylesheet" href="/07800b98-leaflet.css" integrity="sha384-9oKBsxAYdVVBJcv3hwG8RjuoJhw9GwYLqXdQRDxi2q0t1AImNHOap8y6Qt7REVd4" >`,
		},
		{"abcdefgh",
			"async",
			"testdata/gnss-map-plot.js",
			`<script src="/e83a0b0f-gnss-map-plot.js" type="text/javascript" integrity="sha384-haxRijtRHhpn6nbt+JNpioqOj0AwB+THIaVdUZ34R9sQrQL2vmf/pn6GPnQq+AI1" nonce="abcdefgh" async></script>`,
		},
		{"ijklmnop",
			"defer",
			"testdata/test.mjs",
			`<script src="/3616e4a4-test.mjs" type="module" integrity="sha384-yL9nK0JVp9FW9oAfkQ2kQC/9CcuAMK4vmyb8q+TY2SokmBFflIxJpZJ6Nk8Xqw5r" nonce="ijklmnop" defer></script>`,
		},
	}

	for _, w := range work {
		t.Run(w.path, func(t *testing.T) {

			a, err := loadAsset(w.path, "testdata")
			if err != nil {
				t.Error(err)
			}

			tag, err := createSubResourceTag(a, w.nonce, w.attr)
			if err != nil {
				t.Fatalf("Couldn't create embedded resource tag: %v", err)
			}

			if tag != w.expected {
				t.Fatalf("output tag '%v' did not equal epected '%v'", tag, w.expected)
			}
		})
	}
}

func TestCreateSubResourcePreloadTag(t *testing.T) {
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
			"testdata/test.mjs",
			`<link rel="modulepreload" href="/3616e4a4-test.mjs" integrity="sha384-yL9nK0JVp9FW9oAfkQ2kQC/9CcuAMK4vmyb8q+TY2SokmBFflIxJpZJ6Nk8Xqw5r"/>`,
		},
		{"abcdefg",
			"testdata/test.mjs",
			`<link rel="modulepreload" href="/3616e4a4-test.mjs" integrity="sha384-yL9nK0JVp9FW9oAfkQ2kQC/9CcuAMK4vmyb8q+TY2SokmBFflIxJpZJ6Nk8Xqw5r" nonce="abcdefg"/>`,
		},
	}

	for _, w := range work {
		t.Run(w.path, func(t *testing.T) {

			a, err := loadAsset(w.path, "testdata")
			if err != nil {
				t.Fatal(err)
			}
			tag, err := createSubResourcePreloadTag(a, w.nonce)
			if err != nil {
				t.Errorf("Couldn't create embedded resource preload tag: %v", err)
			}
			if tag != w.expected {
				t.Errorf("output tag '%v' did not equal epected '%v'", tag, w.expected)
			}
		})
	}
}
