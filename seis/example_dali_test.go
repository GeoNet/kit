// +build dlink

package seis

import (
	"strings"
	"testing"
)

func TestDLConnect(t *testing.T) {

	dlink := NewDLink("localhost:16000")

	conn, err := dlink.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	if s := conn.Id(); !strings.HasPrefix(s, "seis:seis") {
		t.Errorf("expected id \"seis:seis\" but got %s", s)
	}
	if !conn.Writable() {
		t.Errorf("expected to be writeable")
	}
	if d := conn.Size(); d != 512 {
		t.Errorf("expected to have 512 size but got %d", d)
	}
}
