// Package keyspaces is for working with AWS Keyspaces (for Apache Cassandra).
package keyspaces

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sigv4-auth-cassandra-gocql-driver-plugin/sigv4"
	"github.com/gocql/gocql"
)

type Keyspaces struct {
	Session *gocql.Session
}

// New returns a Keyspaces struct which wraps a Keyspaces session. certPath
// is the path of the required Starfield digital certificate to connect
// to Keyspaces using SSL/TLS
// (found here: https://certs.secureserver.net/repository/sf-class2-root.crt)
func New(certPath string) (Keyspaces, error) {

	ksClient := Keyspaces{}

	cfg, err := getConfig()
	if err != nil {
		return ksClient, err
	}

	region := cfg.Region

	// Add the Amazon Keyspaces service endpoint
	cluster := gocql.NewCluster(fmt.Sprintf("cassandra.%s.amazonaws.com:9142", region))

	// Get credentials from config (required for authentication to Keyspaces)
	creds, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		return ksClient, err
	}
	var auth sigv4.AwsAuthenticator = sigv4.NewAwsAuthenticator()
	auth.Region = region
	auth.AccessKeyId = creds.AccessKeyID
	auth.SecretAccessKey = creds.SecretAccessKey
	auth.SessionToken = creds.SessionToken
	cluster.Authenticator = auth

	// Provide the path to the certificate
	cluster.SslOpts = &gocql.SslOptions{
		CaPath: certPath,
	}
	// Override default Consistency to LocalQuorum
	cluster.Consistency = gocql.LocalQuorum
	// Disable initial host lookup
	cluster.DisableInitialHostLookup = true

	// Create and return session.
	session, err := cluster.CreateSession()
	if err != nil {
		return ksClient, err
	}
	ksClient.Session = session

	return ksClient, nil
}

// QueryDatabase queries the database with provided query, and returns a Scanner to iterate through the results.
func (k *Keyspaces) QueryDatabase(query string, values []interface{}) gocql.Scanner {
	return k.Session.Query(query, values...).Iter().Scanner()
}

// getConfig returns the default AWS Config struct.
func getConfig() (aws.Config, error) {
	if os.Getenv("AWS_REGION") == "" {
		return aws.Config{}, errors.New("AWS_REGION is not set")
	}

	var cfg aws.Config
	var err error

	if awsEndpoint := os.Getenv("CUSTOM_AWS_ENDPOINT_URL"); awsEndpoint != "" {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:       "aws",
				URL:               awsEndpoint,
				HostnameImmutable: true,
			}, nil
		})

		cfg, err = config.LoadDefaultConfig(
			context.TODO(),
			config.WithEndpointResolverWithOptions(customResolver),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO())
	}

	if err != nil {
		return aws.Config{}, err
	}
	return cfg, nil
}
