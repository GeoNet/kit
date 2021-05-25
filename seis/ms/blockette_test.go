package ms

import (
	"testing"
)

func TestBlockette_Header(t *testing.T) {
	raw := BlocketteHeader{1000, 95}

	t.Run("encode/decode", func(t *testing.T) {
		res := DecodeBlocketteHeader(EncodeBlocketteHeader(raw))

		if raw != res {
			t.Errorf("encode/decode error, epected %v but got %v", raw, res)
		}
	})

	t.Run("marshal/unmarshal", func(t *testing.T) {

		data, err := raw.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		var res BlocketteHeader
		if err := res.Unmarshal(data); err != nil {
			t.Fatal(err)
		}

		if raw != res {
			t.Errorf("encode/decode error, epected %v but got %v", raw, res)
		}

	})

	t.Run("encode/unmarshal", func(t *testing.T) {

		var res BlocketteHeader
		if err := res.Unmarshal(EncodeBlocketteHeader(raw)); err != nil {
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

		res := DecodeBlocketteHeader(data)

		if raw != res {
			t.Errorf("marshal/decode error, epected %v but got %v", raw, res)
		}
	})
}

func TestBlockette_1000(t *testing.T) {
	raw := Blockette1000{11, 1, 67, 0}

	t.Run("encode/decode", func(t *testing.T) {
		res := DecodeBlockette1000(EncodeBlockette1000(raw))

		if raw != res {
			t.Errorf("encode/decode error, epected %v but got %v", raw, res)
		}
	})

	t.Run("marshal/unmarshal", func(t *testing.T) {

		data, err := raw.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		var res Blockette1000
		if err := res.Unmarshal(data); err != nil {
			t.Fatal(err)
		}

		if raw != res {
			t.Errorf("encode/decode error, epected %v but got %v", raw, res)
		}

	})

	t.Run("encode/unmarshal", func(t *testing.T) {

		var res Blockette1000
		if err := res.Unmarshal(EncodeBlockette1000(raw)); err != nil {
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

		res := DecodeBlockette1000(data)

		if raw != res {
			t.Errorf("marshal/decode error, epected %v but got %v", raw, res)
		}
	})
}

func TestBlockette_1001(t *testing.T) {
	raw := Blockette1001{100, 10, 0, 14}

	t.Run("encode/decode", func(t *testing.T) {
		res := DecodeBlockette1001(EncodeBlockette1001(raw))

		if raw != res {
			t.Errorf("encode/decode error, epected %v but got %v", raw, res)
		}
	})

	t.Run("marshal/unmarshal", func(t *testing.T) {

		data, err := raw.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		var res Blockette1001
		if err := res.Unmarshal(data); err != nil {
			t.Fatal(err)
		}

		if raw != res {
			t.Errorf("encode/decode error, epected %v but got %v", raw, res)
		}

	})

	t.Run("encode/unmarshal", func(t *testing.T) {

		var res Blockette1001
		if err := res.Unmarshal(EncodeBlockette1001(raw)); err != nil {
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

		res := DecodeBlockette1001(data)

		if raw != res {
			t.Errorf("marshal/decode error, epected %v but got %v", raw, res)
		}
	})
}
