// Package sns is for messaging with AWS SNS.
package sns

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SNS struct {
	client *sns.Client
}

// New returns an SNS struct which wraps an SNS client using the default AWS credentials chain.
// This consults (in order) environment vars, config files, ec2 and ecs roles.
// Requests with recoverable errors will be retried with the default retrier.
// It is an error if the AWS_REGION environment variable is not set.
func New() (SNS, error) {
	cfg, err := getConfig()
	if err != nil {
		return SNS{}, err
	}
	return SNS{client: sns.NewFromConfig(cfg)}, nil
}

// NewWithMaxRetries returns the same as New(), but with the
// back off set to up to maxRetries times.
func NewWithMaxRetries(maxRetries int) (SNS, error) {
	cfg, err := getConfig()
	if err != nil {
		return SNS{}, err
	}
	client := sns.NewFromConfig(cfg, func(options *sns.Options) {
		options.Retryer = retry.AddWithMaxAttempts(options.Retryer, maxRetries)
	})
	return SNS{client: client}, nil
}

// getConfig returns the default AWS Config struct.
func getConfig() (aws.Config, error) {
	if os.Getenv("AWS_REGION") == "" {
		return aws.Config{}, errors.New("AWS_REGION is not set")
	}

	var cfg aws.Config
	var err error

	if awsEndpoint := os.Getenv("AWS_ENDPOINT_URL"); awsEndpoint != "" {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:   "aws",
				SigningRegion: region,
				URL:           awsEndpoint,
			}, nil
		})

		cfg, err = config.LoadDefaultConfig(
			context.TODO(),
			config.WithEndpointResolverWithOptions(customResolver))
	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO())
	}

	if err != nil {
		return aws.Config{}, err
	}
	return cfg, nil
}

// Ready returns whether the SNS client has been initialised.
func (s *SNS) Ready() bool {
	return s.client != nil
}

// Publish publishes message to topicArn.
func (s *SNS) Publish(topicArn string, message []byte) error {
	input := sns.PublishInput{
		Message:  aws.String(string(message)),
		TopicArn: aws.String(topicArn),
	}

	result, err := s.client.Publish(context.TODO(), &input)
	if err != nil {
		return err
	}

	if aws.ToString(result.MessageId) == "" {
		return errors.New("empty message ID in SNS publish result")
	}

	return nil
}

// CreateTopic creates an SNS topic with the specified name. You can provide a
// map of topic attributes if needed (can be set to nil if not).
// Returns the ARN of the created topic.
func (s *SNS) CreateTopic(topicName string, topicAttributes map[string]string) (string, error) {
	topic, err := s.client.CreateTopic(context.TODO(), &sns.CreateTopicInput{
		Name:       aws.String(topicName),
		Attributes: topicAttributes,
	})
	if err != nil {
		return "", err
	}
	return aws.ToString(topic.TopicArn), nil
}

// DeleteTopic delete an SNS topic.
func (s *SNS) DeleteTopic(topicArn string) error {
	_, err := s.client.DeleteTopic(context.TODO(), &sns.DeleteTopicInput{TopicArn: aws.String(topicArn)})
	return err
}

// CheckTopic checks if an SNS topic exists and is accessible given its ARN.
func (s *SNS) CheckTopic(topicArn string) error {
	_, err := s.client.GetTopicAttributes(context.TODO(), &sns.GetTopicAttributesInput{
		TopicArn: aws.String(topicArn),
	})
	return err
}

// SubscribeQueue subscribes an SQS queue to an SNS topic.
func (s *SNS) SubscribeQueue(topicArn string, queueArn string) (string, error) {
	output, err := s.client.Subscribe(context.TODO(), &sns.SubscribeInput{
		Protocol:              aws.String("sqs"),
		TopicArn:              aws.String(topicArn),
		Endpoint:              aws.String(queueArn),
		ReturnSubscriptionArn: true,
	})
	if err != nil {
		return "", err
	}
	return aws.ToString(output.SubscriptionArn), err
}
