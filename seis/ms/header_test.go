package ms

import (
	"testing"
	"time"
)

func TestRecord_Header(t *testing.T) {
	raw := RecordHeader{
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

	t.Run("encode/decode", func(t *testing.T) {
		res := DecodeRecordHeader(EncodeRecordHeader(raw))

		if raw != res {
			t.Errorf("encode/decode error, expected %v but got %v", raw, res)
		}
	})

	t.Run("marshal/unmarshal", func(t *testing.T) {

		data, err := raw.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		var res RecordHeader
		if err := res.Unmarshal(data); err != nil {
			t.Fatal(err)
		}

		if raw != res {
			t.Errorf("encode/decode error, expected %v but got %v", raw, res)
		}

	})

	t.Run("encode/unmarshal", func(t *testing.T) {

		var res RecordHeader
		if err := res.Unmarshal(EncodeRecordHeader(raw)); err != nil {
			t.Fatal(err)
		}

		if raw != res {
			t.Errorf("encode/unmarshal error, expected %v but got %v", raw, res)
		}
	})

	t.Run("marshal/decode", func(t *testing.T) {

		data, err := raw.Marshal()
		if err != nil {
			t.Fatal(err)
		}

		res := DecodeRecordHeader(data)

		if raw != res {
			t.Errorf("marshal/decode error, expected %v but got %v", raw, res)
		}
	})
}

func TestRecord_Get(t *testing.T) {
	raw := RecordHeader{
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

	t.Run("seq", func(t *testing.T) {
		if v, r := raw.SeqNumber(), 123456; v != r {
			t.Errorf("seq error, expected %v but got %v", r, v)
		}
	})

	t.Run("sta", func(t *testing.T) {
		if v, r := raw.Station(), "WAIM"; v != r {
			t.Errorf("sta error, expected %v but got %v", r, v)
		}
	})

	t.Run("loc", func(t *testing.T) {
		if v, r := raw.Location(), "20"; v != r {
			t.Errorf("loc error, expected %v but got %v", r, v)
		}
	})

	t.Run("cha", func(t *testing.T) {
		if v, r := raw.Channel(), "HNZ"; v != r {
			t.Errorf("cha error, expected %v but got %v", r, v)
		}
	})

	t.Run("net", func(t *testing.T) {
		if v, r := raw.Network(), "NZ"; v != r {
			t.Errorf("net error, expected %v but got %v", r, v)
		}
	})

	t.Run("src", func(t *testing.T) {
		if v, r := raw.SrcName(false), "NZ_WAIM_20_HNZ"; v != r {
			t.Errorf("src error, expected %v but got %v", r, v)
		}
	})

}

func TestRecord_Set(t *testing.T) {
	raw := RecordHeader{
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

	t.Run("seq", func(t *testing.T) {
		raw := raw

		r := 234567
		raw.SetSeqNumber(r)
		if v := raw.SeqNumber(); v != r {
			t.Errorf("seq error, expected %v but got %v", r, v)
		}
	})

	t.Run("sta", func(t *testing.T) {
		raw := raw

		r := "TEST"
		raw.SetStation(r)
		if v := raw.Station(); v != r {
			t.Errorf("sta error, expected %v but got %v", r, v)
		}
	})

	t.Run("loc", func(t *testing.T) {
		r := "XX"
		raw.SetLocation(r)
		if v := raw.Location(); v != r {
			t.Errorf("loc error, expected %v but got %v", r, v)
		}
	})

	t.Run("cha", func(t *testing.T) {
		raw := raw

		r := "YYY"
		raw.SetChannel(r)
		if v := raw.Channel(); v != r {
			t.Errorf("cha error, expected %v but got %v", r, v)
		}
	})

	t.Run("net", func(t *testing.T) {
		raw := raw

		r := "ZZ"
		raw.SetNetwork(r)
		if v := raw.Network(); v != r {
			t.Errorf("net error, expected %v but got %v", r, v)
		}
	})

	t.Run("start", func(t *testing.T) {
		raw := raw

		// no time correction has been applied so add it
		r := raw.RecordStartTime.Time().Add(raw.Correction())

		if v := raw.StartTime(); !v.Equal(r) {
			t.Errorf("start error, expected %v but got %v", r, v)
		}
	})

	t.Run("not applied", func(t *testing.T) {
		raw := raw

		// reset correction and flags
		raw.ActivityFlags = clearBit(raw.ActivityFlags, 1)
		raw.TimeCorrection = 0

		b, r := BTime{2018, 106, 9, 14, 46, 0, 251}.Time(), time.Second
		raw.SetStartTime(b)
		if v := raw.StartTime(); !v.Equal(b) {
			t.Errorf("not applied error, expected %v but got %v", b, v)
		}

		raw.SetCorrection(r, false)
		if isBitSet(raw.ActivityFlags, 1) {
			t.Errorf("not applied error, activity flag should not be set")
		}
		if v := raw.Correction(); v != r {
			t.Errorf("not applied error, expected %v but got %v", r, v)
		}

		if v := raw.StartTime().Sub(b); v != r {
			t.Errorf("not applied error, expected %v but got %v", r, v)
		}

	})

	t.Run("applied", func(t *testing.T) {
		raw := raw

		b, r := BTime{2018, 106, 9, 14, 46, 0, 251}.Time(), time.Second

		raw.SetStartTime(b)
		raw.SetCorrection(r, true)

		if !isBitSet(raw.ActivityFlags, 1) {
			t.Errorf("applied error, activity flag should be set")
		}

		if v := raw.StartTime(); !v.Equal(b) {
			t.Errorf("applied error, expected no changed but got %v", v)
		}

		if v := raw.Correction(); v != r {
			t.Errorf("applied error, expected %v but got %v", r, v)
		}
	})

}

func TestBits(t *testing.T) {

	for i := 0; i < 8; i++ {
		r := byte(0x01 << i)

		//nolint:gosec
		if v := setBit(0, uint8(i)); v != r {
			t.Errorf("invalid bit set: expected %v but got %v for %d", r, v, i)
		}
	}

	for i := 0; i < 8; i++ {
		//nolint:gosec
		v := setBit(0, uint8(i))

		//nolint:gosec
		if !isBitSet(v, uint8(i)) {
			t.Errorf("invalid bit test: expected bit to be set %d", i)
		}

		//nolint:gosec
		v = clearBit(v, uint8(i))

		//nolint:gosec
		if isBitSet(v, uint8(i)) {
			t.Errorf("invalid bit test: expected bit to not be set %d", i)
		}
	}

	for i := 0; i < 8; i++ {
		//nolint:gosec
		v := setBit(0xff, uint8(i))
		
		//nolint:gosec
		if !isBitSet(v, uint8(i)) {
			t.Errorf("invalid bit test: expected bit to be set %d", i)
		}

		//nolint:gosec
		v = clearBit(v, uint8(i))

		for j := 0; j < i; j++ {
			//nolint:gosec
			if !isBitSet(v, uint8(j)) {
				t.Errorf("invalid bit test: expected bit to be set %d/%d", i, j)
			}
		}

		for j := i + 1; j < 8; j++ {
			//nolint:gosec
			if !isBitSet(v, uint8(j)) {
				t.Errorf("invalid bit test: expected bit to be set %d/%d", i, j)
			}
		}
	}
}
