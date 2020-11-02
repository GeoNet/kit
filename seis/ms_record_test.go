package seis

import (
	"io/ioutil"
	"testing"
)

func TestMSRecord_Unpack(t *testing.T) {

	files := []string{
		"basic.mseed",
		"empty_location.mseed",
		"geonet-seedlink-info-ascii.mseed",
		"NZ.AUCT.40.BTT.mseed",
		"NZ.CHIT.40.BTT.mseed",
		"steim1.mseed",
		"wel2000.mseed",
		"4096_float.mseed",
	}

	for _, k := range files {
		t.Run("unpack header: "+k, func(t *testing.T) {
			raw, err := ioutil.ReadFile("testdata/" + k)
			if err != nil {
				t.Fatal(err)
			}
			var header MSRecord
			if err := header.Unpack(raw, false); err != nil {
				t.Error(err)
			}
			var data MSRecord
			if err := data.Unpack(raw, true); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestMSRecord_File(t *testing.T) {

	files := map[string]string{
		"4096_float.mseed":                 "NZ_WAIS_20_HNZ, 000001, D, 4096, 1008 samples, 200 Hz, 2020,298,09:36:24.733165",
		"basic.mseed":                      "NZ_TDHS_20_BN1, 175910, D, 512, 646 samples, 50 Hz, 2016,245,16:36:51.051896",
		"empty_location.mseed":             "AU_MOO__BHE, 000001, D, 512, 370 samples, 40 Hz, 2019,001,00:00:08.319536",
		"geonet-seedlink-info-ascii.mseed": "SL_INFO__INF, 000001, D, 512, 216 samples, 0 Hz, 2019,074,03:45:47.508800",
		"NZ.AUCT.40.BTT.mseed":             "NZ_AUCT_40_BTT, 000001, D, 512, 358 samples, 10 Hz, 2019,099,01:52:28.069500",
		"NZ.CHIT.40.BTT.mseed":             "NZ_CHIT_40_BTT, 000001, D, 512, 256 samples, 10 Hz, 2019,099,01:52:55.069500",
		"steim1.mseed":                     "AU_MILA_00_BHZ, 136425, D, 512, 412 samples, 40 Hz, 2019,148,00:00:09.025000",
		"wel2000.mseed":                    "NZ_WEL_20_BNE, 000001, D, 512, 720 samples, 50 Hz, 2000,021,13:43:00.180000",
	}

	for k, v := range files {
		t.Run("unpack header: "+k, func(t *testing.T) {
			raw, err := ioutil.ReadFile("testdata/" + k)
			if err != nil {
				t.Fatal(err)
			}
			ms, err := NewMSRecord(raw)
			if err != nil {
				t.Fatal(err)
			}
			if s := ms.String(); s != v {
				t.Errorf("invalid block, expected \"%s\", got \"%s\"", v, s)
			}
		})
	}
}
