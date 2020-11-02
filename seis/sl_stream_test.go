// +build link

package seis

import (
	"context"
	"io/ioutil"
	"net"
	"os"
	"testing"
	"time"
)

func TestSLink_Refresh(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	t.Run("first", func(t *testing.T) {
		sl := NewSLink("link.geonet.org.nz:18000",
			SetStreams("NZ_WEL:HHZ HNZ,NZ_CAW:EHZ"),
			SetStateFile(tmpfile.Name()),
			SetRefresh(time.Second),
		)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := sl.CollectWithContext(ctx, func(seq string, data []byte) (bool, error) {
			t.Log(seq)
			return false, nil
		}); err != nil {
			switch e, ok := err.(net.Error); {
			case !ok || !e.Timeout():
				t.Error(err)
			}
		}
	})

	first, err := ioutil.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	t.Log("->", string(first))

	t.Run("second", func(t *testing.T) {
		sl := NewSLink("link.geonet.org.nz:18000",
			SetStreams("NZ_WEL:HHZ HNZ,NZ_CAW:EHZ"),
			SetStateFile(tmpfile.Name()),
			SetRefresh(time.Second),
		)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := sl.CollectWithContext(ctx, func(seq string, data []byte) (bool, error) {
			t.Log(seq)
			return false, nil
		}); err != nil {
			switch e, ok := err.(net.Error); {
			case !ok || !e.Timeout():
				t.Error(err)
			}
		}
	})

	second, err := ioutil.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	t.Log("->", string(second))
}
