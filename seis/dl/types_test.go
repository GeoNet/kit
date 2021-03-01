package dl

import "testing"

func TestMarshallUnmarshalPreheader(t *testing.T) {
	o := Preheader{
		[2]byte{'S', 'L'},
		255,
	}

	b := MarshalPreheader(o)

	o2 := UnmarshalPreheader(b)

	if o2 != o {
		t.Errorf("unmarshalled does not match unmarshalled\nORIG:\n%v\nUNMARSHALLED:\n%v", o, o2)
	}
}
