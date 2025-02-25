package cloudwatch

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type CloudWatch struct {
	client *cloudwatchlogs.Client
}

// New creates a new CloudWatch Logs client
func New() (CloudWatch, error) {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		return CloudWatch{}, err
	}
	// Create a CloudWatch Logs client
	client := cloudwatchlogs.NewFromConfig(cfg)
	return CloudWatch{client}, nil
}

// get latest log stream with specific log group name
func (cw *CloudWatch) GetLatestLogStream(logGroupName string) (types.LogStream, error) {
	var latestLogStream types.LogStream
	// Fetch log streams, sorted by LastIngestionTime
	input := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: &logGroupName,
		OrderBy:      "LastEventTime",
		Descending:   aws.Bool(true), // Get the latest first
		Limit:        aws.Int32(1),   // Only fetch the latest stream
	}
	result, err := cw.client.DescribeLogStreams(context.TODO(), input)
	if err != nil {
		return latestLogStream, err
	}
	// Check if we found a log stream
	if len(result.LogStreams) == 0 {
		return latestLogStream, fmt.Errorf("no log streams found")
	}
	// Get the latest log
	if len(result.LogStreams) > 0 {
		latestLogStream = result.LogStreams[0]
	}
	return latestLogStream, nil
}

// GetLogEvents fetches log events from a log stream
func (cw *CloudWatch) GetLogEvents(logGroupName string,
	logStreamName string, limit int32) (*cloudwatchlogs.GetLogEventsOutput, error) {
	input := &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  &logGroupName,
		LogStreamName: &logStreamName,
		Limit:         aws.Int32(limit),
	}
	return cw.client.GetLogEvents(context.TODO(), input)
}

// GetLogEventsWithFilter fetches log events from a log stream with a filter pattern
func (cw *CloudWatch) GetLogEventsWithFilter(logGroupName string,
	logStreamName string, filterPattern string) (*cloudwatchlogs.FilterLogEventsOutput, error) {
	input := &cloudwatchlogs.FilterLogEventsInput{
		LogGroupName:   &logGroupName,
		LogStreamNames: []string{logStreamName},
		FilterPattern:  aws.String(filterPattern),
	}
	return cw.client.FilterLogEvents(context.TODO(), input)
}

// HealthCheck checks health by checking cloudwatch logs
// logGroupName: log group name, e.g. tf-dev-tilde-meta-cron
// errorFilter: filter to apply on logs, e.g. "level=ERROR"
// maxAgeHour: maximum age of log in hours
func (cw *CloudWatch) HealthCheck(logGroupName, errorFilter string, maxAgeHour float64) error {
	// Get latest log stream
	latestStream, err := cw.GetLatestLogStream(logGroupName)
	if err != nil {
		return fmt.Errorf("err: %v", err)
	}
	tm := time.Unix(0, *latestStream.LastEventTimestamp*int64(time.Millisecond))
	//check log age
	logAgeHours := time.Since(tm).Hours()
	if logAgeHours > maxAgeHour {
		return fmt.Errorf("last log event age %v is older than 7 hour ", logAgeHours)
	}
	// Get filtered log events
	result, err := cw.GetLogEventsWithFilter(logGroupName, *latestStream.LogStreamName, errorFilter)
	if err != nil {
		return fmt.Errorf("err: %v", err)
	}
	if len(result.Events) > 0 {
		return fmt.Errorf("error messages found in log stream: %v", len(result.Events))
	}

	return nil
}
