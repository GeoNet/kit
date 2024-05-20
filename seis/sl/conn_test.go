//go:build link
// +build link

package sl

import (
	"testing"
	"time"
)

// use "go test -tags link" to test
func TestConn(t *testing.T) {
	conn, err := NewConn("link.geonet.org.nz", 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	_, err = conn.GetInfo("streams")
	if err != nil {
		t.Error(err)
	}
}
