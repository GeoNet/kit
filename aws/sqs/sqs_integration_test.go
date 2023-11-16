//go:build localstack
// +build localstack

package sqs

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	cutomAWSEndpointURL = "http://localhost:4566"
	awsRegion           = "ap-southeast-2"
	testQueue           = "test-queue"
	testMessage         = "test message"

	testMessageAttributeKey   = "test-key"
	testMessageAttributeValue = "test-value"
)

// helper functions

func setup() {
	// setup environment variables to access LocalStack
	os.Setenv("AWS_REGION", awsRegion)
	os.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("CUSTOM_AWS_ENDPOINT_URL", cutomAWSEndpointURL)

	// create queue
	if err := exec.Command(
		"aws", "sqs",
		"create-queue",
		"--queue-name", testQueue,
		"--endpoint-url", cutomAWSEndpointURL,
		"--region", awsRegion).Run(); err != nil {

		panic(err)
	}
}

func teardown() {
	if err := exec.Command(
		"aws", "sqs",
		"delete-queue",
		"--queue-url", awsCmdQueueURL(),
		"--region", awsRegion,
		"--endpoint-url", cutomAWSEndpointURL).Run(); err != nil {

		panic(err)
	}
}

func awsCmdQueueURL() string {
	if out, err := exec.Command(
		"aws", "sqs",
		"get-queue-url",
		"--queue-name", testQueue,
		"--region", awsRegion,
		"--endpoint-url", cutomAWSEndpointURL).CombinedOutput(); err != nil {

		panic(err)
	} else {
		var payload map[string]string
		json.Unmarshal(out, &payload)
		return payload["QueueUrl"]
	}
}

func awsCmdSendMessage() {
	if err := exec.Command(
		"aws", "sqs",
		"send-message",
		"--queue-url", awsCmdQueueURL(),
		"--message-body", testMessage,
		"--message-attributes", fmt.Sprintf(
			"%s={DataType=String, StringValue=\"%s\"}",
			testMessageAttributeKey, testMessageAttributeValue),
		"--region", awsRegion,
		"--endpoint-url", cutomAWSEndpointURL).Run(); err != nil {

		panic(err)
	}
}

func awsCmdReceiveMessage() string {
	if out, err := exec.Command(
		"aws", "sqs",
		"receive-message",
		"--queue-url", awsCmdQueueURL(),
		"--attribute-names", "body",
		"--region", awsRegion,
		"--endpoint-url", cutomAWSEndpointURL).CombinedOutput(); err != nil {

		panic(err)
	} else {
		var payload map[string][]map[string]string
		json.Unmarshal(out, &payload)
		return payload["Messages"][0]["Body"]
	}
}

func awsCmdReceiveMessageWithAttributes() (string, map[string]string) {
	if out, err := exec.Command(
		"aws", "sqs",
		"receive-message",
		"--queue-url", awsCmdQueueURL(),
		"--message-attribute-names", "All",
		"--region", awsRegion,
		"--endpoint-url", cutomAWSEndpointURL).CombinedOutput(); err != nil {

		panic(err)
	} else {
		var payload map[string][]map[string]interface{}
		json.Unmarshal(out, &payload)

		msgAttributes := map[string]string{}
		for k, v := range payload["Messages"][0]["MessageAttributes"].(map[string]interface{}) {
			msgAttributes[k] = v.(map[string]interface{})["StringValue"].(string)
		}

		return payload["Messages"][0]["Body"].(string), msgAttributes
	}
}

func awsCmdQueueCount() int {
	if out, err := exec.Command(
		"aws", "sqs",
		"get-queue-attributes",
		"--queue-url", awsCmdQueueURL(),
		"--attribute-name", "ApproximateNumberOfMessages",
		"--region", awsRegion,
		"--endpoint-url", cutomAWSEndpointURL).CombinedOutput(); err != nil {

		panic(err)
	} else {
		var payload map[string]map[string]string
		json.Unmarshal(out, &payload)
		rvalue, _ := strconv.Atoi(payload["Attributes"]["ApproximateNumberOfMessages"])
		return rvalue
	}
}

func TestSQSNewAndReady(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	// two tests, good and bad

	// GOOD

	// ACTION
	client, err := New()
	ready := client.Ready()

	// ASSERT
	assert.Nil(t, err)
	assert.True(t, ready)

	// BAD

	os.Unsetenv("AWS_REGION")

	// ACTION
	client, err = New()
	ready = client.Ready()

	// ASSERT
	assert.NotNil(t, err)
	assert.False(t, ready)
}

func TestSQSNewWithRetries(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	// ACTION
	_, err := NewWithMaxRetries(2)

	// ASSERT
	assert.Nil(t, err)

	// ARRAGE
	os.Unsetenv("AWS_REGION")

	// ACTION
	_, err = NewWithMaxRetries(2)

	// ASSERT
	assert.NotNil(t, err)
}

func TestSQSReceive(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	awsCmdSendMessage()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating sqs client: %v", err))

	// ACTION
	receivedMessage, err := client.Receive(awsCmdQueueURL(), 1)

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, testMessage, receivedMessage.Body)
	assert.Equal(t, testMessageAttributeValue, receivedMessage.MessageAttributes[testMessageAttributeKey])
}

func TestSQSReceiveWithAttributes(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	awsCmdSendMessage()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating sqs client: %v", err))

	// ACTION
	receivedMessage, err := client.ReceiveWithAttributes(
		awsCmdQueueURL(), 1, []types.QueueAttributeName{"All"})

	// ASSERT
	assert.Nil(t, err)
	assert.True(t, len(receivedMessage.Attributes) > 0)
}

func TestSQSDelete(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	awsCmdSendMessage()
	require.Equal(t, 1, awsCmdQueueCount())

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating sqs client: %v", err))

	receivedMessage, err := client.Receive(awsCmdQueueURL(), 1)
	require.Nil(t, err, fmt.Sprintf("Error receiving test message: %v", err))

	// ACTION
	err = client.Delete(awsCmdQueueURL(), receivedMessage.ReceiptHandle)

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, 0, awsCmdQueueCount())
}

func TestSQSSend(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating sqs client: %v", err))

	// ACTION
	err = client.Send(awsCmdQueueURL(), testMessage)

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, testMessage, awsCmdReceiveMessage())
}

func TestSQSSendMessageWithAttributes(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating sqs client: %v", err))

	// ACTION
	err = client.SendWithAttributes(
		awsCmdQueueURL(),
		testMessage,
		map[string]string{testMessageAttributeKey: testMessageAttributeValue})

	// ASSERT
	assert.Nil(t, err)

	msg, attributes := awsCmdReceiveMessageWithAttributes()
	assert.Equal(t, testMessage, msg)
	assert.Equal(t, testMessageAttributeValue, attributes[testMessageAttributeKey])
}

func TestSQSSendWithDelay(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating sqs client: %v", err))

	// ACTION
	err = client.SendWithDelay(awsCmdQueueURL(), testMessage, 5) // delay seconds
	assert.Nil(t, err)                                           // fail fast here to avoid loop queue check

	start := time.Now()
	timeout := 10 * time.Second
	for awsCmdQueueCount() < 1 {
		time.Sleep(500 * time.Millisecond)
		if time.Since(start) > timeout {
			break
		}
	}
	timeElapsed := time.Since(start)

	// ASSERT
	assert.True(t, timeElapsed > 5*time.Second)
	assert.True(t, timeElapsed < timeout)
}

func TestGetQueueUrl(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating sqs client: %v", err))

	// ACTION
	queueURL, err := client.GetQueueUrl(testQueue)

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, awsCmdQueueURL(), queueURL)
}
