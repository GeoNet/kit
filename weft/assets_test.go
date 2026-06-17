package weft

import (
	"io/fs"
	"os"
	"path/filepath"
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
			"testdata/test.js",
			&asset{
				path:       "/test.js",
				hashedPath: "/7ae97332-test.js",
				mime:       "text/javascript",
				fileType:   "js",
				sri:        "sha384-QpJfAj2w6B/9M/RPFCW5SdxSs8wf4DRuer8K06bMu8cqj0Cu91WZYh4spHDPmKO/",
			},
			56,
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

func TestUpdateAsset(t *testing.T) {

	testData := []struct {
		testName             string
		filename             string
		append               string
		expectedInitial      *asset
		expectedInitialBytes int
		expectedResult       *asset
	}{
		{
			"Update CSS file",
			"testdata/leaflet.css",
			"abc",
			&asset{
				path:       "/leaflet.css",
				hashedPath: "/07800b98-leaflet.css",
				mime:       "text/css",
				fileType:   "css",
				sri:        "sha384-9oKBsxAYdVVBJcv3hwG8RjuoJhw9GwYLqXdQRDxi2q0t1AImNHOap8y6Qt7REVd4",
			},
			13429,
			&asset{
				path:       "/leaflet.css",
				hashedPath: "/35aea7ae-leaflet.css",
				mime:       "text/css",
				fileType:   "css",
				sri:        "sha384-pQdxLofki9LA7dW8kunwJTtCD/uhhLglB46EU576cEgXCtj7bJqASfVDb7IVDxnC",
			},
		},
	}

	for _, d := range testData {

		t.Run(d.testName, func(t *testing.T) {

			// Make a copy of test data into temp directory, and count number of files.
			tmpDir := t.TempDir()

			count := 0
			err := filepath.WalkDir("testdata", func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if !d.IsDir() {
					destPath := filepath.Join(tmpDir, d.Name())

					input, err := os.ReadFile(path) //nolint:gosec
					if err != nil {
						t.Fatalf("failed to read source file: %v", err)
					}
					if err := os.WriteFile(destPath, input, 0600); err != nil { //nolint: gosec
						t.Fatalf("failed to copy file to temp dir: %v", err)
					}
					count++
				}
				return nil
			})
			if err != nil {
				t.Error(err)
			}
			if count < 1 {
				t.Fatal("should be at least one test file in testdata")
			}

			err = InitAssets(tmpDir, tmpDir)
			if err != nil {
				t.Error(err)
			}

			assetsLength := len(assetStore.assets)
			hashesLength := len(assetStore.hashes)

			if assetsLength != count*2 {
				t.Errorf("expected %v files in asset store, found %v", count*2, assetsLength)
			}
			if hashesLength != count {
				t.Errorf("expected %v files in asset store hashes, found %v", count, hashesLength)
			}

			// Append to end of file to make a change
			destPath := filepath.Join(tmpDir, d.expectedInitial.path)
			f, err := os.OpenFile(destPath, os.O_APPEND|os.O_WRONLY, 0600) //nolint:gosec
			if err != nil {
				t.Fatal(err)
			}
			_, err = f.WriteString(d.append)
			if err != nil {
				t.Fatal(err)
			}
			err = f.Close()
			if err != nil {
				t.Fatal(err)
			}

			// Action
			if err := UpdateAsset(destPath); err != nil {
				t.Errorf("failed to update asset: %v", err)
			}

			// Assert
			got, ok := assetStore.assets[d.expectedResult.path]
			if !ok {
				t.Fatalf("path %s not found in store", d.expectedResult.path)
			}
			gotHashed, ok := assetStore.assets[d.expectedResult.hashedPath]
			if !ok {
				t.Fatalf("hashed path %s not found in store", d.expectedResult.path)
			}
			if got.hashedPath != gotHashed.hashedPath || got.sri != gotHashed.sri {
				t.Fatalf("expected asset for path and hashedPath to be the same")
			}

			expectedLength := d.expectedInitialBytes + len(d.append)
			if len(got.b) != expectedLength {
				t.Errorf("expected %d bytes, got %d", expectedLength, len(got.b))
			}
			if got.hashedPath != d.expectedResult.hashedPath {
				t.Errorf("expected hashed path %s instead got %s", d.expectedResult.hashedPath, got.hashedPath)
			}
			if got.sri != d.expectedResult.sri {
				t.Errorf("expected sri hash %s instead got %s", d.expectedResult.sri, got.sri)
			}
			newAssetsLength := len(assetStore.assets)
			newHashesLength := len(assetStore.hashes)

			if newAssetsLength != assetsLength {
				t.Errorf("asset store assets unexpected length (expected no change). Expected: %v Found: %v", assetsLength, newAssetsLength)
			}
			if newHashesLength != hashesLength {
				t.Errorf("asset store hashes unexpected length (expected no change). Expected: %v Found: %v", hashesLength, newHashesLength)
			}
		})
	}
}

func TestCreateSubResourceTag(t *testing.T) {
	err := InitAssets("testdata", "testdata")
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
			"testdata/test.js",
			`<script src="/7ae97332-test.js" type="text/javascript" integrity="sha384-QpJfAj2w6B/9M/RPFCW5SdxSs8wf4DRuer8K06bMu8cqj0Cu91WZYh4spHDPmKO/" nonce="abcdefgh" async></script>`,
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
	err := InitAssets("testdata", "testdata")
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

func TestCreateImportTag(t *testing.T) {
	err := InitAssets("testdata", "testdata")
	if err != nil {
		t.Error(err)
	}

	work := []struct {
		testName      string
		nonce         string
		importMapping map[string]string
		expected      string
	}{
		{
			"No nonce, one module file",
			"",
			map[string]string{
				"test.mjs": "/assets/js/hashprefix-test.mjs",
			},
			`<script type="importmap">
{
	"imports":{
		"test.mjs":"/assets/js/hashprefix-test.mjs"
	}
}
</script>`,
		},
		{
			"Nonce present, two module files",
			"abcdefg",
			map[string]string{
				"test1.mjs": "/assets/js/hashprefix-test1.mjs",
				"test2.mjs": "/assets/js/hashprefix-test2.mjs",
			},
			`<script type="importmap" nonce="abcdefg">
{
	"imports":{
		"test1.mjs":"/assets/js/hashprefix-test1.mjs",
		"test2.mjs":"/assets/js/hashprefix-test2.mjs"
	}
}
</script>`,
		},
	}

	for _, w := range work {
		t.Run(w.testName, func(t *testing.T) {
			tag := createImportMapTag(w.importMapping, w.nonce)
			if tag != w.expected {
				t.Errorf("import map tag\n '%v' did not equal expected\n '%v'", tag, w.expected)
			}
		})
	}
}
