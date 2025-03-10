//go:build localstack
// +build localstack

package sqs

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	customAWSEndpointURL = "http://localhost:4566"
	awsRegion            = "ap-southeast-2"
	testQueue            = "test-queue"
	testMessage          = "test message"
)

// helper functions

func setup() {
	// setup environment variables to access LocalStack
	os.Setenv("AWS_REGION", awsRegion)
	os.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_ENDPOINT_URL", customAWSEndpointURL)

	// create queue
	awsCmdCreateQueue(testQueue)
}

func teardown() {
	if err := exec.Command(
		"aws", "sqs",
		"delete-queue",
		"--queue-url", awsCmdQueueURL(),
		"--region", awsRegion,
	).Run(); err != nil {

		panic(err)
	}
}

func awsCmdCreateQueue(name string) string {
	url, err := exec.Command(
		"aws", "sqs",
		"create-queue",
		"--queue-name", name,
		"--region", awsRegion,
		"--query", "QueueUrl",
		"--output", "text").CombinedOutput()

	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(url))
}

func awsCmdQueueURL() string {
	if out, err := exec.Command(
		"aws", "sqs",
		"get-queue-url",
		"--queue-name", testQueue,
		"--region", awsRegion,
	).CombinedOutput(); err != nil {

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
		"--region", awsRegion,
	).Run(); err != nil {

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
	).CombinedOutput(); err != nil {

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
		"--region", awsRegion).CombinedOutput(); err != nil {

		panic(err)
	} else {
		var payload map[string]map[string]string
		json.Unmarshal(out, &payload)
		rvalue, _ := strconv.Atoi(payload["Attributes"]["ApproximateNumberOfMessages"])
		return rvalue
	}
}

func awsCmdGetQueueArn(url string) string {
	arn, err := exec.Command(
		"aws", "sqs",
		"get-queue-attributes",
		"--queue-url", url,
		"--attribute-names", "QueueArn",
		"--region", awsRegion,
		"--query", "Attributes.QueueArn",
		"--output", "text").CombinedOutput()

	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(arn))
}

func awsCmdDeleteQueue(url string) {
	if err := exec.Command(
		"aws", "sqs",
		"delete-queue",
		"--queue-url", url,
		"--region", awsRegion).Run(); err != nil {

		panic(err)
	}
}

func awsCmdCheckSQSAttribute(url, attribute, expectedValue string) bool {

	if out, err := exec.Command(
		"aws", "sqs",
		"get-queue-attributes",
		"--queue-url", url,
		"--region", awsRegion,
		"--attribute-names", attribute,
		"--output", "json").CombinedOutput(); err != nil {
		panic(err)
	} else {
		var payload map[string]map[string]string
		err = json.Unmarshal(out, &payload)
		if err != nil {
			panic(err)
		}
		attributes, ok := payload["Attributes"]
		if !ok {
			panic("attributes not found in payload")
		}
		if value, exists := attributes[attribute]; exists && value == expectedValue {
			return true
		}
	}
	return false
}

func TestCheckQueue(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	client, err := New()

	// ASSERT
	require.Nil(t, err)

	//test existing queue
	queue, err := client.GetQueueUrl(testQueue)

	assert.Nil(t, err)

	err = client.CheckQueue(queue)
	assert.Nil(t, err)

	//test none existing queue
	err = client.CheckQueue(queue + "_1")
	assert.NotNil(t, err)

}

func TestSQSClient(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	client, err := New()
	ready := client.Ready()

	require.Nil(t, err)
	require.True(t, ready)

	// ACTION
	underlyingClient := client.Client()

	// ASSERT
	assert.NotNil(t, underlyingClient)
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

func TestGetQueueARN(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating sqs client: %v", err))

	// ACTION
	queueARN, err := client.GetQueueARN(awsCmdQueueURL())

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, awsCmdGetQueueArn(awsCmdQueueURL()), queueARN)
}

func TestCreateQueue(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating sqs client: %v", err))

	// ACTION
	newQueue := "test-queue-2"
	url, err := client.CreateQueue(newQueue, false)
	defer awsCmdDeleteQueue(url)

	// ASSERT
	assert.Nil(t, err)
	assert.NotPanics(t, func() { awsCmdGetQueueArn(url) })
}

func TestCreateFifoQueue(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating sqs client: %v", err))

	// ACTION
	newQueue := "test-queue-2.fifo"
	url, err := client.CreateQueue(newQueue, true)
	defer awsCmdDeleteQueue(url)

	// ASSERT
	assert.Nil(t, err)
	assert.NotPanics(t, func() { awsCmdGetQueueArn(url) })
	assert.True(t, awsCmdCheckSQSAttribute(url, "FifoQueue", "true"))
}

func TestDeleteQueue(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating sqs client: %v", err))
	newQueue := "test-queue-2"
	newQueueUrl := awsCmdCreateQueue(newQueue)

	// ACTION
	err = client.DeleteQueue(newQueueUrl)

	// ASSERT
	assert.Nil(t, err)
	assert.Panics(t, func() { awsCmdDeleteQueue(newQueueUrl) })
}
