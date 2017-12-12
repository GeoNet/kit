package mseed_test

import (
	"github.com/GeoNet/kit/mseed"
	"io"
	"os"
	"testing"
	"time"
)

func TestMSR(t *testing.T) {
	msr := mseed.NewMSRecord()
	defer mseed.FreeMSRecord(msr)

	record := make([]byte, 512)

	r, err := os.Open("etc/NZ.ABAZ.10.EHE.D.2016.079")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	_, err = io.ReadFull(r, record)
	if err != nil {
		t.Fatal(err)
	}

	err = msr.Unpack(record, 512, 1, 0)
	if err != nil {
		t.Fatal(err)
	}

	if msr.Network() != "NZ" {
		t.Errorf("expected network code NZ got %s", msr.Network())
	}

	// Station has null termination in the test file.
	if msr.Station() != "ABAZ" {
		t.Errorf("expected station code ABAZ got %s", msr.Station())
	}

	if msr.Location() != "10" {
		t.Errorf("expected location code 10 got %s", msr.Location())
	}

	if msr.Channel() != "EHE" {
		t.Errorf("expected channel code EHE got %s", msr.Channel())
	}

	if msr.Numsamples() != 397 {
		t.Errorf("expected 397 samples got %d", msr.Numsamples())
	}

	s, err := time.Parse(time.RFC3339Nano, "2016-03-19T00:00:01.968393Z")
	if err != nil {
		t.Error(err)
	}

	if !s.Equal(msr.Starttime()) {
		t.Errorf("start time does not match expected 2016-03-19T00:00:01.968393Z got %s",
			msr.Starttime().Format(time.RFC3339Nano))
	}

	e, err := time.Parse(time.RFC3339Nano, "2016-03-19T00:00:05.928393Z")
	if err != nil {
		t.Error(err)
	}

	if !e.Equal(msr.Endtime()) {
		t.Errorf("end time does not match expected 2016-03-19T00:00:05.928393Z got %s",
			msr.Endtime().Format(time.RFC3339Nano))
	}

	d, err := msr.DataSamples()
	if err != nil {
		t.Error(err)
	}

	if len(d) != 397 {
		t.Errorf("expected 397 data samples got %d", len(d))
	}

	if d[0] != 227 {
		t.Errorf("expected first data value 227 got %d", d[0])
	}
}
