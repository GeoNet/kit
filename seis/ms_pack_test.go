package seis

import (
	"io/ioutil"
	"testing"
)

func TestMSUnpack(t *testing.T) {
	raw, err := ioutil.ReadFile("testdata/NZ.AUCT.40.BTT.mseed")
	if err != nil {
		t.Fatal(err)
	}
	msr, err := NewMSRecord(raw)
	if err != nil {
		t.Fatal(err)
	}

	hdr, data := NewMSRecordStream(msr), msr.ToInt32s()

	if err := hdr.PackInt32(msr.StartTime(), 100, false, data, func(m *MSRecord, d []byte, l bool) error {
		t.Log(m.String(), len(d), m.EndTime(), l)
		return nil
	}); err != nil {
		t.Error(err)
	}

}
