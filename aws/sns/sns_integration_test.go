//go:build localstack
// +build localstack

package sns

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	customAWSEndpointURL = "http://localhost:4566"
	awsRegion            = "ap-southeast-2"
	testTopic            = "test-topic"
	testQueue            = "test-queue"
	testMessage          = "test message"
)

type SNSMessage struct {
	Message  string `json:"Message"`
	TopicArn string `json:"TopicArn"`
}

func setup(createQueue, createTopic bool) (string, string) {
	os.Setenv("AWS_REGION", awsRegion)
	os.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_ENDPOINT_URL", customAWSEndpointURL)

	var queueArn, topicArn string
	if createQueue {
		queueArn = awsCmdCreateQueue(testQueue)
	}
	if createTopic {
		topicArn = awsCmdCreateTopic(testTopic)
	}
	return queueArn, topicArn
}

func teardown(queueArn, topicArn string) {
	if queueArn != "" {
		awsCmdDeleteQueue(strings.Split(queueArn, ":")[5])
	}
	if topicArn != "" {
		awsCmdDeleteTopic(topicArn)
	}
	os.Unsetenv("AWS_REGION")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_ENDPOINT_URL")
}

func TestCheckTopic(t *testing.T) {
	_, topicArn := setup(false, true)
	defer teardown("", topicArn)

	client, err := New()
	ready := client.Ready()
	// ASSERT
	require.Nil(t, err)
	assert.True(t, ready)

	//check existing topic
	err = client.CheckTopic(topicArn)
	assert.Nil(t, err)

	//check none existing topic
	err = client.CheckTopic(topicArn + "_1")
	assert.NotNil(t, err)

}

func TestSNSNewAndReady(t *testing.T) {
	// ARRANGE
	setup(false, false)
	defer teardown("", "")

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

func TestSNSNewWithRetries(t *testing.T) {
	// ARRANGE
	setup(false, false)
	defer teardown("", "")

	// ACTION
	_, err := NewWithMaxRetries(2)

	// ASSERT
	assert.Nil(t, err)

	// ARRANGE
	os.Unsetenv("AWS_REGION")

	// ACTION
	_, err = NewWithMaxRetries(2)

	// ASSERT
	assert.NotNil(t, err)
}

func TestSNSClient(t *testing.T) {
	// ARRANGE
	setup(false, false)
	defer teardown("", "")

	client, err := New()
	ready := client.Ready()

	require.Nil(t, err)
	require.True(t, ready)

	// ACTION
	underlyingClient := client.Client()

	// ASSERT
	assert.NotNil(t, underlyingClient)
}

func TestSNSPublish(t *testing.T) {
	// ARRANGE
	queueArn, topicArn := setup(true, true)
	defer teardown(queueArn, topicArn)
	awsCmdSubscribeQueueToTopic(topicArn, queueArn)

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("error creating sns client: %v", err))

	// ACTION
	err = client.Publish(topicArn, []byte(testMessage))

	// ASSERT
	assert.Nil(t, err)
	var m SNSMessage
	message := awsCmdReceiveMessage(testQueue)
	err = json.Unmarshal([]byte(message), &m)
	assert.Nil(t, err)
	assert.Equal(t, testMessage, m.Message)
}

func TestSNSCreateTopic(t *testing.T) {
	// ARRANGE
	setup(false, false)

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("error creating sns client: %v", err))

	// ACTION
	topicArn, err := client.CreateTopic(testTopic, nil)
	defer teardown("", topicArn)

	// ASSERT
	assert.Nil(t, err)
	assert.True(t, awsCmdCheckTopicExists(topicArn))
}

func TestSNSCreateTopicWithAttributes(t *testing.T) {
	// ARRANGE
	setup(false, false)

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("error creating sns client: %v", err))

	// ACTION
	testAttribute := "test-display-name"
	topicArn, err := client.CreateTopic(testTopic, map[string]string{
		"DisplayName": testAttribute,
	})
	defer teardown("", topicArn)

	// ASSERT
	assert.Nil(t, err)
	assert.True(t, awsCmdCheckTopicExists(topicArn))
	assert.True(t, awsCmdCheckTopicAttribute(topicArn, "DisplayName", testAttribute))
}

func TestSNSDeleteTopic(t *testing.T) {
	// ARRANGE
	_, topicArn := setup(false, true)
	defer teardown("", topicArn)

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("error creating sns client: %v", err))

	// ACTION
	err = client.DeleteTopic(topicArn)

	// ASSERT
	assert.Nil(t, err)
	assert.False(t, awsCmdCheckTopicExists(topicArn))
}

func TestSNSSubscribeQueue(t *testing.T) {
	// ARRANGE
	queueArn, topicArn := setup(true, true)
	defer teardown(queueArn, topicArn)

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("error creating sns client: %v", err))

	// ACTION
	subscriptionArn, err := client.SubscribeQueue(topicArn, queueArn)

	// ASSERT
	assert.Nil(t, err)
	assert.True(t, awsCmdCheckQueueSubscribedToTopic(topicArn, queueArn, subscriptionArn))
}

// helper functions

func awsCmdCheckQueueSubscribedToTopic(topicArn, queueArn, subscriptionArn string) bool {

	if out, err := exec.Command(
		"aws", "sns",
		"list-subscriptions-by-topic",
		"--topic-arn", topicArn,
		"--region", awsRegion,
		"--output", "json").CombinedOutput(); err != nil {

		panic(err)
	} else {
		var payload map[string][]map[string]string
		json.Unmarshal(out, &payload)
		subscriptions := payload["Subscriptions"]
		for _, subscription := range subscriptions {
			if endpoint, exists := subscription["Endpoint"]; exists && endpoint == queueArn {
				if subArn, exists := subscription["SubscriptionArn"]; exists && subArn == subscriptionArn {
					return true
				}
			}
		}
	}
	return false
}

func awsCmdQueueURL(name string) string {
	if out, err := exec.Command(
		"aws", "sqs",
		"get-queue-url",
		"--queue-name", name,
		"--region", awsRegion,
	).CombinedOutput(); err != nil {

		panic(err)
	} else {
		var payload map[string]string
		json.Unmarshal(out, &payload)
		return payload["QueueUrl"]
	}
}

func awsCmdReceiveMessage(name string) string {
	if out, err := exec.Command(
		"aws", "sqs",
		"receive-message",
		"--queue-url", awsCmdQueueURL(name),
		"--attribute-names", "body",
		"--region", awsRegion,
	).CombinedOutput(); err != nil {

		panic(err)
	} else {
		var payload map[string][]map[string]string
		json.Unmarshal(out, &payload)
		if _, ok := payload["Messages"]; !ok {
			return ""
		}
		return payload["Messages"][0]["Body"]
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
	return awsCmdGetQueueArn(strings.TrimSpace(string(url)))
}

func awsCmdCreateTopic(name string) string {
	arn, err := exec.Command(
		"aws", "sns",
		"create-topic",
		"--name", name,
		"--region", awsRegion,
		"--query", "TopicArn",
		"--output", "text").CombinedOutput()

	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(arn))
}

func awsCmdDeleteQueue(name string) {
	if err := exec.Command(
		"aws", "sqs",
		"delete-queue",
		"--queue-url", awsCmdQueueURL(name),
		"--region", awsRegion).Run(); err != nil {

		panic(err)
	}
}

func awsCmdDeleteTopic(arn string) {
	if err := exec.Command(
		"aws", "sns",
		"delete-topic",
		"--topic-arn", arn,
		"--region", awsRegion).Run(); err != nil {

		panic(err)
	}
}

func awsCmdSubscribeQueueToTopic(topicArn, queueArn string) {
	if err := exec.Command(
		"aws", "sns",
		"subscribe",
		"--topic-arn", topicArn,
		"--protocol", "sqs",
		"--notification-endpoint", queueArn,
		"--region", awsRegion).Run(); err != nil {

		panic(err)
	}
}

func awsCmdCheckTopicAttribute(arn, attribute, expectedValue string) bool {

	if out, err := exec.Command(
		"aws", "sns",
		"get-topic-attributes",
		"--topic-arn", arn,
		"--region", awsRegion,
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

func awsCmdCheckTopicExists(arn string) bool {
	if out, err := exec.Command(
		"aws", "sns",
		"list-topics",
		"--region", awsRegion,
		"--output", "json").CombinedOutput(); err != nil {

		panic(err)
	} else {
		var payload map[string][]map[string]string
		json.Unmarshal(out, &payload)
		topicsMap := payload["Topics"]
		if len(topicsMap) == 0 {
			return false
		}
		for range topicsMap {
			a := topicsMap[0]["TopicArn"]
			if a == arn {
				return true
			}
		}
	}
	return false
}
