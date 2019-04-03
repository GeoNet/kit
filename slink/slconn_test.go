//nolint //cgo generates code that doesn't pass linting
package slink

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestNewSLCD(t *testing.T) {
	slconn := NewSLCD()
	defer FreeSLCD(slconn)
	if slconn == nil {
		t.Fatal("unable to create new SLCD")
	}
}

func TestSetNetDly(t *testing.T) {
	slconn := NewSLCD()
	defer FreeSLCD(slconn)
	if slconn == nil {
		t.Fatal("unable to create new SLCD")
	}
	if slconn.NetDly() != 30 {
		t.Error("NetDly unexpected default value")
	}
	slconn.SetNetDly(10)
	if slconn.NetDly() != 10 {
		t.Error("NetDly not set")
	}
}

func TestSetNetTo(t *testing.T) {
	slconn := NewSLCD()
	defer FreeSLCD(slconn)
	if slconn == nil {
		t.Fatal("unable to create new SLCD")
	}
	if slconn.NetTo() != 600 {
		t.Error("NetTo unexpected default value")
	}
	slconn.SetNetTo(10)
	if slconn.NetTo() != 10 {
		t.Error("NetTo not set")
	}
}

func TestSetNetKeepAlive(t *testing.T) {
	slconn := NewSLCD()
	defer FreeSLCD(slconn)
	if slconn == nil {
		t.Fatal("unable to create new SLCD")
	}
	if slconn.KeepAlive() != 0 {
		t.Error("KeepAlive unexpected default value")
	}
	slconn.SetKeepAlive(10)
	if slconn.KeepAlive() != 10 {
		t.Error("KeepAlive not set")
	}
}

func TestReadStreamList(t *testing.T) {
	slconn := NewSLCD()
	defer FreeSLCD(slconn)
	if slconn == nil {
		t.Fatal("unable to create new SLCD")
	}
	tf, err := ioutil.TempFile("", "slconn_test")
	if err != nil {
		t.Fatal("unable to open temporary file")
	}
	defer os.Remove(tf.Name())
	if err = ioutil.WriteFile(tf.Name(), ([]byte)("# A comment\nGE ISP  BH?.D\nNL HGN\nMN AGU BH? HH?\n"), 0644); err != nil {
		t.Errorf("unable to write file: %v", err)
	}
	count, err := slconn.ReadStreamList(tf.Name(), "")
	if err != nil {
		t.Error("unable to read stream list, invalid format")
	}
	if count != 4 {
		t.Error("unable to read stream list, incorrect count")
	}
	count, err = slconn.ReadStreamList(tf.Name(), "H??")
	if err != nil {
		t.Error("unable to read stream list, invalid format")
	}
	if count != 4 {
		t.Error("unable to read stream list, incorrect count")
	}
}

func TestSetBeginTime(t *testing.T) {
	slconn := NewSLCD()
	defer FreeSLCD(slconn)
	if slconn == nil {
		t.Fatal("unable to create new SLCD")
	}
	if slconn.BeginTime() != "" {
		t.Error("BeginTime unexpected default value")
	}

	s, err := time.Parse("2006-01-02,15:04:05", "2017-10-02,10:20:30")
	if err != nil {
		t.Fatal(err)
	}

	slconn.SetBeginTime(s.Format(TimeFormat))
	if ans, res := "2017,10,02,10,20,30", slconn.BeginTime(); ans != res {
		t.Errorf("BeginTime mismatch: expected \"%s\", got \"%s\"", ans, res)
	}
}

func TestSetEndTime(t *testing.T) {
	slconn := NewSLCD()
	defer FreeSLCD(slconn)
	if slconn == nil {
		t.Fatal("unable to create new SLCD")
	}
	if slconn.EndTime() != "" {
		t.Error("EndTime unexpected default value")
	}

	s, err := time.Parse("2006-01-02,15:04:05", "2017-10-01,00:10:20")
	if err != nil {
		t.Fatal(err)
	}

	slconn.SetEndTime(s.Format(TimeFormat))
	if ans, res := "2017,10,01,00,10,20", slconn.EndTime(); ans != res {
		t.Errorf("EndTime mismatch: expected \"%s\", got \"%s\"", ans, res)
	}
}

func TestParseStreamList(t *testing.T) {

	LogInit(0, nil, nil)

	slconn := NewSLCD()
	defer FreeSLCD(slconn)
	if slconn == nil {
		t.Fatal("unable to create new SLCD")
	}
	count, err := slconn.ParseStreamList("IU_KONO:BHE BHN,GE_WLF,MN_AQU:HH?.D", "")
	if err != nil {
		t.Error("unable to parse stream list, invalid string")
	}
	if count != 3 {
		t.Error("unable to parse stream list, wrong count")
	}
	count, err = slconn.ParseStreamList("IU_KONO:BHE BHN,GE_WLF,MN_AQU:HH?.D", "H??")
	if err != nil {
		t.Error("unable to parse stream list, invalid string")
	}
	if count != 3 {
		t.Error("unable to parse stream list, wrong count")
	}
	_, err = slconn.ParseStreamList("IU__KONO:BHE BHN,GE_WLF,MN_AQU:HH?.D", "")
	if err == nil {
		t.Error("shouldn't be able to parse stream list, invalid string")
	}
}

func TestLogInitErrorMessage(t *testing.T) {

	var msg string
	LogInit(100, nil, func() func(string) {
		return func(s string) {
			msg = s
		}
	}())

	slconn := NewSLCD()
	defer FreeSLCD(slconn)
	if slconn == nil {
		t.Fatal("unable to create new SLCD")
	}

	if _, err := slconn.ParseStreamList("IU__KONO:BHE BHN,GE_WLF,MN_AQU:HH?.D", ""); err == nil {
		t.Error("shouldn't be able to parse stream list, invalid string")
	}

	if ans := "not in NET_STA format: IU__KONO"; msg != ans {
		t.Errorf("error message: expected \"%s\", received \"%s\"", ans, msg)
	}
}
