package seis

import "testing"

func TestMarshallUnmarshalBtime(t *testing.T) {
	raw := BTime{2017, 105, 8, 13, 45, 0, 250}
	res := UnmarshalBTime(MarshalBTime(raw))

	if raw != res {
		t.Errorf("unmarshall error, epected %v but got %v", raw, res)
	}
}

func TestMarshallUnmarshalMSDataHeader(t *testing.T) {
	raw := MSDataHeader{
		[6]byte{'1', '2', '3', '4', '5', '6'},
		'D',
		uint8(0),
		[5]byte{' ', 'W', 'A', 'I', 'M'},
		[2]byte{'2', '0'},
		[3]byte{'H', 'N', 'Z'},
		[2]byte{'N', 'Z'},
		BTime{2017, 105, 8, 13, 45, 0, 250},
		565,
		10,
		-1,
		uint8(1), uint8(2), uint8(3),
		2,
		25,
		64,
		42,
	}

	res := UnmarshalMSDataHeader(MarshalMSDataHeader(raw))

	if raw != res {
		t.Errorf("data header unmarshall error, epected %v but got %v", raw, res)
	}
}

func TestMarshallUnmarshalBlocketteHeader(t *testing.T) {
	raw := BlocketteHeader{1000, 95}

	res := UnmarshalBlocketteHeader(MarshalBlocketteHeader(raw))

	if raw != res {
		t.Errorf("blockette header unmarshall error, epected %v but got %v", raw, res)
	}
}

func TestMarshallUnmarshalBlockette1000(t *testing.T) {
	raw := Blockette1000{11, 1, 67, 0}

	res := UnmarshalBlockette1000(MarshalBlockette1000(raw))

	if raw != res {
		t.Errorf("blockette 1000 unmarshall error, epected %v but got %v", raw, res)
	}
}

func TestMarshallUnmarshalBlockette1001(t *testing.T) {
	raw := Blockette1001{100, 10, 0, 14}

	res := UnmarshalBlockette1001(MarshalBlockette1001(raw))

	if res != raw {
		t.Errorf("blockette 1001 unmarshall error, epected %v but got %v", raw, res)
	}
}
