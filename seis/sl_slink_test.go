// +build link

package seis

import (
	"context"
	"net"
	"testing"
	"time"
)

// use "go test -tags link" to test
func TestSLink_Connect(t *testing.T) {

	work := []struct {
		streamlist, deflist string
		toGrab              int
		start, end          bool
	}{
		{
			"NZ_???T:4??TT", "", 5, false, false,
		},
		{
			"*_*:41???", "", 5, false, false,
		},
		{
			"*_*", "", 5, false, false,
		},
		{ //TODO: valid?
			"*", "", 10, false, false,
		},
		{
			"GIST", "4????", 5, true, false,
		},
		{
			"NZ_AUCT,NZ_CHIT,NZ_PUYT,KAIT,RFRT", "40BTT 41BTT 40VTT 41VTT 40LTT 41LTT", 25, true, true,
		},
	}

	for _, w := range work {
		t.Run(w.streamlist, func(t *testing.T) {
			sl := NewSLink("link.geonet.org.nz:18000",
				SetTimeout(5*time.Second),
				SetStreams(w.streamlist),
				SetSelectors(w.deflist),
				SetStartTime(func() time.Time {
					switch {
					case w.start && w.end:
						return time.Now().Add(-(time.Minute * 30))
					case w.start:
						return time.Now().Add(-time.Minute * 10)
					default:
						return time.Time{}
					}
				}()),
				SetEndTime(func() time.Time {
					switch {
					case w.start && w.end:
						return time.Now().Add(-(time.Minute * 18))
					default:
						return time.Time{}
					}
				}()),
			)

			var count int
			if err := sl.Collect(func(seq string, pkt []byte) (bool, error) {
				t.Log(seq)
				if count > w.toGrab {
					return true, nil
				}
				count++
				return false, nil
			}); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestSLink_KeepAlive(t *testing.T) {

	sl := NewSLink("link.geonet.org.nz:18000",
		SetStreams("XX_ZZZZ"),
		SetKeepAlive(200*time.Millisecond),
		SetNetTo(500*time.Millisecond),
		SetTimeout(100*time.Millisecond),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := sl.CollectWithContext(ctx, func(seq string, data []byte) (bool, error) {
		// should not reach here
		return false, nil
	}); err != nil {
		t.Error(err)
	}

}

func TestSLink_NetTo(t *testing.T) {

	sl := NewSLink("link.geonet.org.nz:18000",
		SetStreams("XX_ZZZZ"),
		SetKeepAlive(0),
		SetNetTo(500*time.Millisecond),
		SetTimeout(100*time.Millisecond),
	)

	// this should return after netto seconds ...
	err := sl.Collect(func(seq string, data []byte) (bool, error) {
		// should not reach here
		return false, nil
	})

	switch e, ok := err.(net.Error); {
	case !ok || !e.Timeout():
		t.Error(err)

	}
}
