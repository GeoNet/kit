package slink

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestNewSLCD(t *testing.T) {
	slconn := NewSLCD()
	defer FreeSLCD(slconn)
	if slconn == nil {
		t.Error("unable to create new SLCD")
	}
}

func TestSetNetDly(t *testing.T) {
	slconn := NewSLCD()
	defer FreeSLCD(slconn)
	if slconn == nil {
		t.Error("unable to create new SLCD")
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
		t.Error("unable to create new SLCD")
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
		t.Error("unable to create new SLCD")
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
		t.Error("unable to create new SLCD")
	}
	tf, err := ioutil.TempFile("", "slconn_test")
	if err != nil {
		t.Error("unable to open temporary file")
	}
	defer os.Remove(tf.Name())
	ioutil.WriteFile(tf.Name(), ([]byte)("# A comment\nGE ISP  BH?.D\nNL HGN\nMN AGU BH? HH?\n"), 0644)
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

func TestParseStreamList(t *testing.T) {
	slconn := NewSLCD()
	defer FreeSLCD(slconn)
	if slconn == nil {
		t.Error("unable to create new SLCD")
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
	count, err = slconn.ParseStreamList("IU__KONO:BHE BHN,GE_WLF,MN_AQU:HH?.D", "")
	if err == nil {
		t.Error("shouldn't be able to parse stream list, invalid string")
	}
}
