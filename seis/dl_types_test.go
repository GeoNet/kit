package seis

import "testing"

func TestMarshallUnmarshalDLPreheader(t *testing.T) {
	o := DLPreheader{
		[2]byte{'S', 'L'},
		255,
	}

	b := MarshalDLPreheader(o)

	o2 := UnmarshalDLPreheader(b)

	if o2 != o {
		t.Errorf("unmarshalled does not match unmarshalled\nORIG:\n%v\nUNMARSHALLED:\n%v", o, o2)
	}
}
