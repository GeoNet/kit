//go:build localstack
// +build localstack

package sqs

import (
	"context"
	"encoding/json"
	"errors"
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

func awsCmdReceiveMessages() []string {
	if out, err := exec.Command( //nolint:gosec
		"aws", "sqs",
		"receive-message",
		"--queue-url", awsCmdQueueURL(),
		"--attribute-names", "body",
		"--region", awsRegion,
		"--max-number-of-messages", "10", // AWS SQS allows up to 10 messages at a time
	).CombinedOutput(); err != nil {

		panic(err)
	} else {
		var payload map[string][]map[string]string
		_ = json.Unmarshal(out, &payload)

		var bodies []string
		for _, msg := range payload["Messages"] {
			if body, ok := msg["Body"]; ok {
				bodies = append(bodies, body)
			}
		}
		return bodies
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

func awsCmdQueueInFlightCount() int {
	out, err := exec.Command(
		"aws", "sqs",
		"get-queue-attributes",
		"--queue-url", awsCmdQueueURL(),
		"--attribute-name", "ApproximateNumberOfMessagesNotVisible",
		"--region", awsRegion).CombinedOutput()

	if err != nil {
		panic(err)
	}

	var payload map[string]map[string]string
	json.Unmarshal(out, &payload)

	rvalue, _ := strconv.Atoi(payload["Attributes"]["ApproximateNumberOfMessagesNotVisible"])
	return rvalue
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

func TestSQSReceiveBatch(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	awsCmdSendMessage()
	awsCmdSendMessage()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("error creating sqs client: %v", err))

	// ACTION
	receivedMessages, err := client.ReceiveBatch(context.TODO(), awsCmdQueueURL(), 30)

	// ASSERT
	assert.Nil(t, err)
	for _, receivedMessage := range receivedMessages {
		assert.Equal(t, testMessage, receivedMessage.Body)
	}
}

func TestSQSDelete(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	awsCmdSendMessage()
	require.Equal(t, 1, awsCmdQueueCount())

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating sqs client: %v", err))

	receivedMessage, err := client.Receive(awsCmdQueueURL(), 30)
	require.Nil(t, err, fmt.Sprintf("Error receiving test message: %v", err))
	require.Equal(t, 1, awsCmdQueueInFlightCount())

	// ACTION
	err = client.Delete(awsCmdQueueURL(), receivedMessage.ReceiptHandle)

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, 0, awsCmdQueueCount())
	assert.Equal(t, 0, awsCmdQueueInFlightCount())
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

func TestSendBatch(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("error creating sqs client: %v", err))

	// ACTION
	var maxBytes int = 1048576
	maxSizeSingleMessage := strings.Repeat("a", maxBytes)
	err = client.SendBatch(context.TODO(), awsCmdQueueURL(), []string{maxSizeSingleMessage})

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, maxSizeSingleMessage, awsCmdReceiveMessage())

	// ACTION
	tooLargeSingleMessage := maxSizeSingleMessage + "a"
	err = client.SendBatch(context.TODO(), awsCmdQueueURL(), []string{tooLargeSingleMessage})

	// ASSERT
	assert.NotNil(t, err)

	// ACTION
	var maxHalfBytes int = 524288
	maxHalfSizeMessage := strings.Repeat("a", maxHalfBytes)
	err = client.SendBatch(context.TODO(), awsCmdQueueURL(), []string{maxHalfSizeMessage, maxHalfSizeMessage})

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, maxHalfSizeMessage, awsCmdReceiveMessage())
	assert.Equal(t, maxHalfSizeMessage, awsCmdReceiveMessage())

	// ACTION
	tooLargeHalfSizeMessage := maxHalfSizeMessage + "a"
	err = client.SendBatch(context.TODO(), awsCmdQueueURL(), []string{maxHalfSizeMessage, tooLargeHalfSizeMessage})

	// ASSERT
	assert.NotNil(t, err)

	var sbe *SendBatchError
	if errors.As(err, &sbe) {
		assert.Equal(t, 2, len(sbe.Info))
		assert.Equal(t, 0, sbe.Info[0].Index)
		assert.Equal(t, 1, sbe.Info[1].Index)
	} else {
		t.Error("unexpected error type")
	}

	// ACTION
	validMessage := "test"
	invalidMessage := "\u0000"

	err = client.SendBatch(context.TODO(), awsCmdQueueURL(), []string{validMessage, invalidMessage})

	// ASSERT
	assert.NotNil(t, err)

	sbe = nil
	if errors.As(err, &sbe) {
		assert.Equal(t, 1, len(sbe.Info))
		assert.Equal(t, 1, sbe.Info[0].Index)
	} else {
		t.Error("unexpected error type")
	}

	assert.Equal(t, validMessage, awsCmdReceiveMessage())
}

func TestSendNBatch(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("error creating sqs client: %v", err))

	// ACTION
	var maxBytes int = 1048576
	maxSizeSingleMessage := strings.Repeat("a", maxBytes)
	batchesSent, err := client.SendNBatch(context.TODO(), awsCmdQueueURL(), []string{maxSizeSingleMessage, maxSizeSingleMessage})

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, 2, batchesSent)
	assert.Equal(t, maxSizeSingleMessage, awsCmdReceiveMessage())
	assert.Equal(t, maxSizeSingleMessage, awsCmdReceiveMessage())

	assert.Equal(t, 0, awsCmdQueueCount())

	// ACTION
	tooLargeSingleMessage := maxSizeSingleMessage + "a"
	batchesSent, err = client.SendNBatch(context.TODO(), awsCmdQueueURL(), []string{maxSizeSingleMessage, tooLargeSingleMessage})

	// ASSERT
	assert.NotNil(t, err)
	assert.Equal(t, 1, batchesSent)

	var sbe *SendNBatchError
	if errors.As(err, &sbe) {
		assert.Equal(t, 1, len(sbe.Info))
		assert.True(t, indexIsPresent(sbe.Info, 1))
	} else {
		t.Error("unexpected error type")
	}

	assert.Equal(t, maxSizeSingleMessage, awsCmdReceiveMessage())
	assert.Equal(t, 0, awsCmdQueueCount())

	// ACTION
	smallMessageText := "small"
	smallMessageCount := 21
	smallMessages := make([]string, smallMessageCount)
	for i := range smallMessageCount {
		smallMessages[i] = smallMessageText
	}
	batchesSent, err = client.SendNBatch(context.TODO(), awsCmdQueueURL(), smallMessages)

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, 3, batchesSent)

	receiveCount := 0
	for range batchesSent {
		messages := awsCmdReceiveMessages()
		for _, m := range messages {
			if m == smallMessageText {
				receiveCount++
			}
		}
	}
	assert.Equal(t, smallMessageCount, receiveCount)

	// ACTION
	invalidMessage := "\u0000"

	batchesSent, err = client.SendNBatch(context.TODO(), awsCmdQueueURL(), []string{
		tooLargeSingleMessage,
		maxSizeSingleMessage,
		smallMessageText,
		maxSizeSingleMessage,
		invalidMessage,
		smallMessageText,
		smallMessageText,
		invalidMessage,
		tooLargeSingleMessage,
	})

	// ASSERT
	assert.NotNil(t, err)
	assert.Equal(t, 4, batchesSent)

	sbe = nil
	if errors.As(err, &sbe) {
		assert.Equal(t, 4, len(sbe.Info))
		assert.True(t, indexIsPresent(sbe.Info, 0))
		assert.True(t, indexIsPresent(sbe.Info, 4))
		assert.True(t, indexIsPresent(sbe.Info, 7))
		assert.True(t, indexIsPresent(sbe.Info, 8))
	} else {
		t.Error("unexpected error type")
	}

	assert.Equal(t, 5, len(awsCmdReceiveMessages()))
}

func indexIsPresent(info []SendBatchErrorEntry, index int) bool {
	for _, entry := range info {
		if entry.Index == index {
			return true
		}
	}
	return false
}

func TestDeleteBatch(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("error creating sqs client: %v", err))

	// Send and receive messages to get receipt handles
	messages := []string{"message1", "message2", "message3"}
	err = client.SendBatch(context.TODO(), awsCmdQueueURL(), messages)
	require.Nil(t, err)
	require.Equal(t, 3, awsCmdQueueCount())
	require.Equal(t, 0, awsCmdQueueInFlightCount())

	receivedMessages, err := client.ReceiveBatch(context.TODO(), awsCmdQueueURL(), 30)
	require.Nil(t, err)
	require.Equal(t, 3, len(receivedMessages))
	require.Equal(t, 3, awsCmdQueueInFlightCount())
	receiptHandles := make([]string, 0)
	for _, rm := range receivedMessages {
		receiptHandles = append(receiptHandles, rm.ReceiptHandle)
	}

	// ACTION
	err = client.DeleteBatch(context.TODO(), awsCmdQueueURL(), receiptHandles)

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, 0, awsCmdQueueCount())
	assert.Equal(t, 0, awsCmdQueueInFlightCount())

	// ACTION
	invalidReceiptHandle := "invalid-receipt-handle"
	err = client.DeleteBatch(context.TODO(), awsCmdQueueURL(), []string{invalidReceiptHandle})

	// ASSERT
	assert.NotNil(t, err)

	var dbe *DeleteBatchError
	if errors.As(err, &dbe) {
		assert.Equal(t, 1, len(dbe.Info))
		assert.Equal(t, 0, dbe.Info[0].Index)
	} else {
		t.Error("unexpected error type")
	}
}

func TestDeleteNBatch(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("error creating sqs client: %v", err))

	// Send and receive messages to get receipt handles
	messages := []string{"msg1", "msg2", "msg3", "msg4", "msg5", "msg6", "msg7", "msg8", "msg9", "msg10", "msg11"}
	batchesSent, err := client.SendNBatch(context.TODO(), awsCmdQueueURL(), messages)

	require.Nil(t, err)
	require.Equal(t, 2, batchesSent)

	receivedMessages1, err := client.ReceiveBatch(context.TODO(), awsCmdQueueURL(), 30)
	require.Nil(t, err)
	require.Equal(t, 10, len(receivedMessages1))
	receivedMessages2, err := client.ReceiveBatch(context.TODO(), awsCmdQueueURL(), 30)
	require.Nil(t, err)
	require.Equal(t, 1, len(receivedMessages2))

	require.Equal(t, 11, awsCmdQueueInFlightCount())

	receiptHandles := make([]string, 0)
	for _, rm := range receivedMessages1 {
		receiptHandles = append(receiptHandles, rm.ReceiptHandle)
	}
	for _, rm := range receivedMessages2 {
		receiptHandles = append(receiptHandles, rm.ReceiptHandle)
	}

	// ACTION
	batchesDeleted, err := client.DeleteNBatch(context.TODO(), awsCmdQueueURL(), receiptHandles)

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, 2, batchesDeleted)
	assert.Equal(t, 0, awsCmdQueueCount())
	assert.Equal(t, 0, awsCmdQueueInFlightCount())

	// ARRANGE

	// Send and receive messages to get receipt handles
	messages = []string{"msg1", "msg2", "msg3", "msg4", "msg5", "msg6", "msg7", "msg8", "msg9", "msg10", "msg11", "msg12"}
	batchesSent, err = client.SendNBatch(context.TODO(), awsCmdQueueURL(), messages)

	require.Nil(t, err)
	require.Equal(t, 2, batchesSent)

	receivedMessages1, err = client.ReceiveBatch(context.TODO(), awsCmdQueueURL(), 30)
	require.Nil(t, err)
	require.Equal(t, 10, len(receivedMessages1))
	receivedMessages2, err = client.ReceiveBatch(context.TODO(), awsCmdQueueURL(), 30)
	require.Nil(t, err)
	require.Equal(t, 2, len(receivedMessages2))
	require.Equal(t, 12, awsCmdQueueInFlightCount())

	receiptHandles = make([]string, 0)
	for _, rm := range receivedMessages1 {
		receiptHandles = append(receiptHandles, rm.ReceiptHandle)
	}
	for _, rm := range receivedMessages2 {
		receiptHandles = append(receiptHandles, rm.ReceiptHandle)
	}
	invalidReceiptHandle := "invalid-receipt-handle"
	receiptHandles[0] = invalidReceiptHandle                      // Replace a valid receipt handle with an invalid one.
	receiptHandles = append(receiptHandles, invalidReceiptHandle) // Append an invalid receipt handle (index 12)

	// ACTION
	batchesDeleted, err = client.DeleteNBatch(context.TODO(), awsCmdQueueURL(), receiptHandles)
	assert.NotNil(t, err)

	var dbe *DeleteNBatchError
	if errors.As(err, &dbe) {
		assert.Equal(t, 2, len(dbe.Info))
		assert.Equal(t, 0, dbe.Info[0].Index)
		assert.Equal(t, 12, dbe.Info[1].Index)
	} else {
		t.Error("unexpected error type")
	}
	assert.Equal(t, 2, batchesDeleted)
	assert.Equal(t, 0, awsCmdQueueCount())
	assert.Equal(t, 1, awsCmdQueueInFlightCount())
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

func TestMessageVisibility(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	client, err := New()
	require.Nil(t, err, "error with test setup")

	require.Nil(t, client.Send(awsCmdQueueURL(), testMessage), "error with test setup")

	// receive, with visibility timeout of 5min
	receivedMessage, err := client.Receive(awsCmdQueueURL(), 5*60)

	require.Nil(t, err, "error with test setup")
	require.Equal(t, awsCmdQueueCount(), 0)

	// ACTION
	// set visibility timeout to 10 s
	client.SetMessageVisibility(awsCmdQueueURL(), receivedMessage.ReceiptHandle, 10)

	// ASSERT
	assert.Equal(t, awsCmdQueueCount(), 0) // not visible yet
	time.Sleep(11 * time.Second)
	assert.Equal(t, awsCmdQueueCount(), 1) // message visible now
}
