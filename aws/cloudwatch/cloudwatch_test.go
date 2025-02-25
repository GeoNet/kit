//go:build devtest
// +build devtest

package cloudwatch

import (
	"log"
	"os"
	"testing"
	"time"
)

// go test -tags devtest -v -run TestExtractLatestLogs
func TestExtractLatestLogs(t *testing.T) {
	os.Setenv("AWS_REGION", "ap-southeast-2")
	logGroupName := "tf-dev-tilde-meta-cron"
	client, err := New()
	if err != nil {
		t.Error("Test failed, ", err)
	}
	latestStream, err := client.GetLatestLogStream(logGroupName)
	if err != nil {
		t.Error("Test failed, ", err)
	}
	tm := time.Unix(0, *latestStream.LastEventTimestamp*int64(time.Millisecond))
	logAgeHours := time.Since(tm).Hours()
	t.Log("logAgeHours: ", logAgeHours)

	// Get the latest log events
	result, err := client.GetLogEvents(logGroupName, *latestStream.LogStreamName, 10)
	if err != nil {
		t.Error("Test failed, ", err)
	}

	for _, event := range result.Events {
		tm := time.Unix(0, *event.Timestamp*int64(time.Millisecond))
		log.Printf("log messages[%s] %s\n", tm, *event.Message)
	}

	// Get the latest error log events
	result1, err := client.GetLogEventsWithFilter(logGroupName, *latestStream.LogStreamName, "ERROR")
	if err != nil {
		t.Error("Test failed, ", err)
	}

	for _, event := range result1.Events {
		tm := time.Unix(0, *event.Timestamp*int64(time.Millisecond))
		log.Printf("error messages[%s] %s\n", tm, *event.Message)
	}
}
