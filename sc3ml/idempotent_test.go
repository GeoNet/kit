package sc3ml

import (
	"testing"
	"time"
)

func TestIdpQuake(t *testing.T) {
	idpq := IdpQuake{}

	q := Quake{
		Time:     time.Now().UTC().Add(time.Duration(-15 * time.Minute)),
		PublicID: "1234",
	}

	if idpq.Seen(q) != false {
		t.Error("should not have seen quake 1234")
	}

	idpq.Add(q)

	if idpq.Seen(q) != true {
		t.Error("should have seen quake 1234")
	}

	// an old quake
	o := Quake{
		Time:     time.Now().UTC().Add(time.Duration(-100 * time.Minute)),
		PublicID: "1234567",
	}

	idpq.Add(o)

	if idpq.Seen(o) != false {
		t.Error("should not have seen quake 1234567 - it's old and should have been removed.")
	}

	if idpq.Seen(q) != true {
		t.Error("should have seen quake 1234")
	}

}
