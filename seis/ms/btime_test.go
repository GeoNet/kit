package ms

import "testing"

func TestBTime(t *testing.T) {
	raw := BTime{2017, 105, 8, 13, 45, 0, 250}

	t.Run("encode/decode", func(t *testing.T) {
		res := DecodeBTime(EncodeBTime(raw))

		if raw != res {
			t.Errorf("encode/decode error, epected %v but got %v", raw, res)
		}
	})

	t.Run("marshal/unmarshal", func(t *testing.T) {

		data, err := raw.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		var res BTime
		if err := res.Unmarshal(data); err != nil {
			t.Fatal(err)
		}

		if raw != res {
			t.Errorf("encode/decode error, epected %v but got %v", raw, res)
		}

	})

	t.Run("encode/unmarshal", func(t *testing.T) {

		var res BTime
		if err := res.Unmarshal(EncodeBTime(raw)); err != nil {
			t.Fatal(err)
		}

		if raw != res {
			t.Errorf("encode/unmarshal error, epected %v but got %v", raw, res)
		}
	})

	t.Run("marshal/decode", func(t *testing.T) {

		data, err := raw.Marshal()
		if err != nil {
			t.Fatal(err)
		}

		res := DecodeBTime(data)

		if raw != res {
			t.Errorf("marshal/decode error, epected %v but got %v", raw, res)
		}
	})
}
