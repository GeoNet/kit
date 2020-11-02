package seis

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestSLState(t *testing.T) {

	t.Run("check empty files", func(t *testing.T) {
		var state SLState
		if err := state.ReadFile(""); err != nil {
			t.Errorf("an empty read path should not return an error")
		}
		if err := state.WriteFile(""); err != nil {
			t.Errorf("an empty write path should not return an error")
		}
	})

	t.Run("check read then write", func(t *testing.T) {
		raw, err := ioutil.ReadFile("testdata/state.json")
		if err != nil {
			t.Fatal(err)
		}
		var state SLState
		if err := state.Unmarshal(raw); err != nil {
			t.Fatal(err)
		}

		data, err := state.Marshal()
		if err != nil {
			t.Fatal(err)
		}

		if a, b := strings.TrimSpace(string(raw)), strings.TrimSpace(string(data)); a != b {
			t.Errorf("marshal and unmarshal of state file do not match, expected\n%s\ngot\n%s\n", a, b)
		}
	})

	t.Run("check read then write file", func(t *testing.T) {
		tmpfile, err := ioutil.TempFile("", "test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile.Name())

		raw, err := ioutil.ReadFile("testdata/state.json")
		if err != nil {
			t.Fatal(err)
		}

		var state SLState
		if err := state.ReadFile("testdata/state.json"); err != nil {
			t.Error(err)
		}
		if err := state.WriteFile(tmpfile.Name()); err != nil {
			t.Errorf("an empty write path should not return an error")
		}

		data, err := ioutil.ReadFile(tmpfile.Name())
		if err != nil {
			t.Fatal(err)
		}

		if a, b := strings.TrimSpace(string(raw)), strings.TrimSpace(string(data)); a != b {
			t.Errorf("read and write of state file do not match, expected\n%s\ngot\n%s\n", a, b)
		}
	})

}
