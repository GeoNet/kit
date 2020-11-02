// +build link

package seis

import (
	"testing"
	"time"
)

// use "go test -tags link" to test
func TestSLConn(t *testing.T) {
	conn, err := NewSLConn("link.geonet.org.nz", 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	_, err = conn.GetSLInfo("streams")
	if err != nil {
		t.Error(err)
	}
}
