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
	CustomAWSEndpointURL = "http://localhost:4566"
	AWSRegion            = "ap-southeast-2"
	TestQueue            = "test-queue"
	TestMessage          = "test message"
)

// helper functions

func setup() {
	// setup environment variables to access LocalStack
	os.Setenv("AWS_REGION", AWSRegion)
	os.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("CUSTOM_AWS_ENDPOINT_URL", CustomAWSEndpointURL)

	// create queue
	if err := exec.Command(
		"aws", "sqs",
		"create-queue",
		"--queue-name", TestQueue,
		"--endpoint-url", CustomAWSEndpointURL,
		"--region", AWSRegion).Run(); err != nil {

		panic(err)
	}
}

func teardown() {
	if err := exec.Command(
		"aws", "sqs",
		"delete-queue",
		"--queue-url", awsCmdQueueURL(),
		"--region", AWSRegion,
		"--endpoint-url", CustomAWSEndpointURL).Run(); err != nil {

		panic(err)
	}
}

func awsCmdQueueURL() string {
	if out, err := exec.Command(
		"aws", "sqs",
		"get-queue-url",
		"--queue-name", TestQueue,
		"--region", AWSRegion,
		"--endpoint-url", CustomAWSEndpointURL).CombinedOutput(); err != nil {

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
		"--message-body", TestMessage,
		"--region", AWSRegion,
		"--endpoint-url", CustomAWSEndpointURL).Run(); err != nil {

		panic(err)
	}
}

func awsCmdReceiveMessage() string {
	if out, err := exec.Command(
		"aws", "sqs",
		"receive-message",
		"--queue-url", awsCmdQueueURL(),
		"--attribute-names", "body",
		"--region", AWSRegion,
		"--endpoint-url", CustomAWSEndpointURL).CombinedOutput(); err != nil {

		panic(err)
	} else {
		var payload map[string][]map[string]string
		json.Unmarshal(out, &payload)
		return payload["Messages"][0]["Body"]
	}
}

func awsCmdQueueCount() int {
	if out, err := exec.Command(
		"aws", "sqs",
		"get-queue-attributes",
		"--queue-url", awsCmdQueueURL(),
		"--attribute-name", "ApproximateNumberOfMessages",
		"--region", AWSRegion,
		"--endpoint-url", CustomAWSEndpointURL).CombinedOutput(); err != nil {

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
	assert.Equal(t, TestMessage, receivedMessage.Body)
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
	err = client.Send(awsCmdQueueURL(), TestMessage)

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, TestMessage, awsCmdReceiveMessage())
}

func TestSQSSendWithDelay(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating sqs client: %v", err))

	// ACTION
	err = client.SendWithDelay(awsCmdQueueURL(), TestMessage, 4) // delay seconds
	time.Sleep(2 * time.Second)
	messageCountBeforeDelay := awsCmdQueueCount()
	time.Sleep(3 * time.Second)
	messageCountAfterDelay := awsCmdQueueCount()

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, 0, messageCountBeforeDelay)
	assert.Equal(t, 1, messageCountAfterDelay)
}

func TestGetQueueUrl(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating sqs client: %v", err))

	// ACTION
	queueURL, err := client.GetQueueUrl(TestQueue)

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, awsCmdQueueURL(), queueURL)
}
