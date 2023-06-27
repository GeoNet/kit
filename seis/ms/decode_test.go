package ms

import (
	"os"
	"testing"
)

func TestRecord_Bytes(t *testing.T) {

	files := map[string]struct {
		first byte
		last  byte
	}{
		"geonet-seedlink-info-ascii.mseed": {'<', '>'},
	}

	for k, v := range files {
		t.Run("decode data: "+k, func(t *testing.T) {
			raw, err := os.ReadFile("testdata/" + k)
			if err != nil {
				t.Fatal(err)
			}
			var record Record
			if err := record.Unpack(raw); err != nil {
				t.Error(err)
			}
			data, err := record.Bytes()
			if err != nil {
				t.Fatal(err)
			}
			if n := record.SampleCount(); n != len(data) {
				t.Fatalf("invalid number of samples, expected %d, got %d", n, len(data))

			}
			if d := data[0]; d != v.first {
				t.Errorf("invalid first samples, expected %v, got %v", v.first, d)
			}
			if d := data[len(data)-1]; d != v.last {
				t.Errorf("invalid first samples, expected %v, got %v", v.last, d)
			}
		})
	}
}

func TestRecord_Int32s(t *testing.T) {

	files := map[string]struct {
		first int32
		last  int32
	}{
		"basic.mseed":          {-3022, -3027},
		"empty_location.mseed": {4440, 5717},
		"NZ.AUCT.40.BTT.mseed": {33901, 33836},
		"NZ.CHIT.40.BTT.mseed": {557, 1608},
		"steim1.mseed":         {-16268, -15555},
		"wel2000.mseed":        {1, -1},
	}

	for k, v := range files {
		t.Run("decode data: "+k, func(t *testing.T) {
			raw, err := os.ReadFile("testdata/" + k)
			if err != nil {
				t.Fatal(err)
			}
			var record Record
			if err := record.Unpack(raw); err != nil {
				t.Error(err)
			}
			data, err := record.Int32s()
			if err != nil {
				t.Fatal(err)
			}
			if n := record.SampleCount(); n != len(data) {
				t.Fatalf("invalid number of samples, expected %d, got %d", n, len(data))

			}
			if d := data[0]; d != v.first {
				t.Errorf("invalid first samples, expected %v, got %v", v.first, d)
			}
			if d := data[len(data)-1]; d != v.last {
				t.Errorf("invalid first samples, expected %v, got %v", v.last, d)
			}
		})
	}
}

func TestRecord_Float32s(t *testing.T) {

	files := map[string]struct {
		first float32
		last  float32
	}{
		"4096_float.mseed": {-0.00068550772, -0.00080318749},
	}

	for k, v := range files {
		t.Run("decode data: "+k, func(t *testing.T) {
			raw, err := os.ReadFile("testdata/" + k)
			if err != nil {
				t.Fatal(err)
			}
			var record Record
			if err := record.Unpack(raw); err != nil {
				t.Error(err)
			}
			data, err := record.Float64s()
			if err != nil {
				t.Fatal(err)
			}
			if n := record.SampleCount(); n != len(data) {
				t.Fatalf("invalid number of samples, expected %d, got %d", n, len(data))

			}
			if d := float32(data[0]); d != v.first {
				t.Errorf("invalid first samples, expected %v, got %v", v.first, d)
			}
			if d := float32(data[len(data)-1]); d != v.last {
				t.Errorf("invalid first samples, expected %v, got %v", v.last, d)
			}
		})
	}
}
