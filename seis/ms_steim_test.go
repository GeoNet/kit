package seis

import (
	"bytes"
	"testing"
)

func TestSteim_Nibble(t *testing.T) {

	var NibbleTests = []struct {
		w0    []byte
		index int
		value uint8
	}{
		{
			[]byte{0x0, 0x0, 0x30, 0x0}, 9, 3,
		},
	}

	for _, w := range NibbleTests {
		t.Run("get nibble", func(t *testing.T) {
			r := getNibble(w.w0, w.index)
			if r != w.value {
				t.Errorf("Expected value %v, got %v", w.value, r)
			}
		})
		t.Run("write nibble", func(t *testing.T) {
			b := make([]byte, 4)
			writeNibble(b, w.index, w.value)
			if !bytes.Equal(b, w.w0) {
				t.Errorf("Expected value %032b, got %032b", w.w0, b)
			}
		})
	}
}

func TestSteim_Int(t *testing.T) {

	var varIntTest = []struct {
		v    uint32
		i    int32
		bits uint8
	}{
		{1, 1, 2},
		{3, -1, 2},
		{25, 25, 6},
		{39, -25, 6},
		{5906, 5906, 14},
		{10478, -5906, 14},
		{25603942, 25603942, 26},
		{41504922, -25603942, 26},
		{292392304, 292392304, 32},
		{4002574992, -292392304, 32},

		{1, 1, 30},
		{1073741823, -1, 30},

		{1, 1, 32},
		{4294967295, -1, 32},
	}

	for _, x := range varIntTest {
		t.Run("uint to int32", func(t *testing.T) {

			i := uintVarToInt32(x.v, x.bits)
			if i != x.i {
				t.Errorf("input %v %vbits: expected output %v got %v", x.v, x.bits, x.i, i)
			}
		})
		t.Run("int32 to uint", func(t *testing.T) {
			v := int32ToUintVar(x.i, x.bits)
			if v != x.v {
				t.Errorf("input %v %vbits: expected output %v got %v", x.i, x.bits, x.v, v)
			}
		})
	}
}
