package slogger

import (
	"testing"
	"time"
)

func TestSmartLogger(t *testing.T) {
	//define a logger with message prefix
	logger := NewSmartLogger(2*time.Second, "error connecting to")
	n := 0

	// Log some messages for testing
	if n = logger.Log("error connecting to AVLN"); n != 1 {
		t.Errorf("expected 1 but got %v", n)
	}
	if n = logger.Log("error connecting to AUCK"); n != 2 {
		t.Errorf("expected 2 but got %v", n)
	}
	if n = logger.Log("error connecting to DUND"); n != 3 {
		t.Errorf("expected 3 but got %v", n)
	}
	if n = logger.Log("error connecting to WARK"); n != 4 {
		t.Errorf("expected 4 but got %v", n)
	}
	time.Sleep(3 * time.Second)
	if n = logger.Log("Another message"); n != 0 {
		t.Errorf("expected 0 but got %v", n)
	}
	//define a logger without message prefix
	logger = NewSmartLogger(2*time.Second, "")

	if n = logger.Log("Another message"); n != 1 {
		t.Errorf("expected 1 but got %v", n)
	}
	if n = logger.Log("Another message"); n != 2 {
		t.Errorf("expected 2 but got %v", n)
	}
	time.Sleep(3 * time.Second)
	if n = logger.Log("Goodbye!"); n != 1 {
		t.Errorf("expected 1 but got %v", n)
	}
}
