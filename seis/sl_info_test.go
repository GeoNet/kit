package seis

import (
	"io/ioutil"
	"testing"
)

func TestStationInfo(t *testing.T) {

	var checks = map[string]string{
		"capabilites": "capabilites.xml",
		"id":          "id.xml",
		"stations":    "stations.xml",
		"streams":     "streams.xml",
	}

	for k, v := range checks {
		t.Run(k, func(t *testing.T) {
			raw, err := ioutil.ReadFile("testdata/" + v)
			if err != nil {
				t.Fatal(err)
			}
			var info SLInfo
			if err := info.Unmarshal(raw); err != nil {
				t.Error(err)
			}
		})
	}
}
