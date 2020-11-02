// +build link

package seis

import (
	"context"
	"testing"
	"time"
)

func TestMSRecord_Example(t *testing.T) {

	slink := NewSLink("link.geonet.org.nz:18000")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := slink.CollectWithContext(ctx, func(seq string, data []byte) (bool, error) {
		if ms, err := NewMSRecord(data); err == nil {
			t.Log(ms.SrcName(false), time.Since(ms.EndTime()))
		}
		return false, nil
	}); err != nil {
		t.Fatal(err)
	}
}
