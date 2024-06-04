// Package keyspaces is for working with AWS Keyspaces (for Apache Cassandra).
package keyspaces

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sigv4-auth-cassandra-gocql-driver-plugin/sigv4"
	"github.com/gocql/gocql"
)

type Keyspaces struct {
	session   *gocql.Session
	cluster   *gocql.ClusterConfig
	refreshes bool // Whether client has been told to refresh regularly.
	mutex     *sync.Mutex
}

type Logger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

// Default logger discards all logs.
type discarder struct{}

func (d discarder) Println(...interface{})        {}
func (d discarder) Printf(string, ...interface{}) {}

// Keyspaces Logger
var ksLogger Logger = discarder{}

// SetLogger allows calling code to set a desired Logger
// to receive logs from this package (eg: log.Default())
func SetLogger(l Logger) {
	ksLogger = l
}

// New returns a Keyspaces struct which wraps a Keyspaces session.
// host to connect to eg: cassandra.ap-southeast-2.amazonaws.com:9142
// certPath is the path of the required Starfield digital certificate to connect
// to Keyspaces using SSL/TLS
// (found here: https://certs.secureserver.net/repository/sf-class2-root.crt)
func New(host, certPath string) (Keyspaces, error) {

	ksClient := Keyspaces{}

	// Add the Amazon Keyspaces service endpoint
	cluster := gocql.NewCluster(host)
	ksClient.cluster = cluster

	// When host is localhost, for example in a test environment, we don't need these settings.
	if host != "localhost" {

		// Set port.
		cluster.Port = 9142

		authenticator, err := generateAuthenticator()
		if err != nil {
			return ksClient, err
		}
		cluster.Authenticator = authenticator

		// Provide the path to the certificate
		cluster.SslOpts = &gocql.SslOptions{
			CaPath: certPath,
		}
		// Override default Consistency to LocalQuorum
		cluster.Consistency = gocql.LocalQuorum

		// Enable initial host lookup.
		// see https://github.com/gocql/gocql/issues/915
		// When set to true, we were seeing this error intermittently.
		cluster.DisableInitialHostLookup = false
	}

	// Create session.
	session, err := cluster.CreateSession()
	if err != nil {
		return ksClient, err
	}
	ksClient.session = session
	ksClient.mutex = &sync.Mutex{}

	return ksClient, nil
}

// RefreshSessionEvery starts up a loop that replaces the client's
// session with a new one regularly. This can be only set once.
func (k *Keyspaces) RefreshSessionEvery(d time.Duration) {

	k.mutex.Lock()
	defer k.mutex.Unlock()

	if k.refreshes {
		return // Go routine already exists to refresh the Session.
	}
	k.refreshes = true

	// Start a loop that creates a new session regularly.
	go func(client *Keyspaces, interval time.Duration) {
		ticker := time.NewTicker(interval)
		for range ticker.C {
			if client == nil {
				return
			}
			err := client.refreshSession()
			if err != nil {
				ksLogger.Printf("error refreshing session: %v", err)
			}
		}
	}(k, d)
}

// generateAuthenticator returns a gocql Authenticator based on
// available AWS credentials.
func generateAuthenticator() (gocql.Authenticator, error) {

	cfg, err := getConfig()
	if err != nil {
		return nil, err
	}

	// Get credentials from config (required for authentication to Keyspaces)
	creds, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		return nil, err
	}
	var auth sigv4.AwsAuthenticator = sigv4.NewAwsAuthenticator()
	auth.Region = cfg.Region
	auth.AccessKeyId = creds.AccessKeyID
	auth.SecretAccessKey = creds.SecretAccessKey
	auth.SessionToken = creds.SessionToken

	ksLogger.Println("New credentials retrieved.")
	if creds.CanExpire {
		ksLogger.Printf("Credentials expire in %v. Session: ****%v\n", time.Until(creds.Expires), creds.SessionToken[len(creds.SessionToken)-6:])
	}

	return auth, nil
}

// refreshSession retrieves new credentials and creates a new Session
// for the Keyspaces client. It closes the old Session after a delay.
func (k *Keyspaces) refreshSession() error {

	if k.cluster == nil {
		return errors.New("cluster doesn't exist")
	}
	authenticator, err := generateAuthenticator()
	if err != nil {
		return err
	}
	k.cluster.Authenticator = authenticator

	// Create new Session.
	newSession, err := k.cluster.CreateSession()
	if err != nil {
		return err
	}

	// Close old session after delay, to make sure nothing is still using it.
	go func(oldSession *gocql.Session) {

		sleepFor := time.Second * 10
		ksLogger.Printf("New session created. Old Session set to close in %v\n", sleepFor)
		defer ksLogger.Println("Session closed.")

		time.Sleep(sleepFor)

		oldSession.Close()

	}(k.getSession())

	// Set the new Session
	k.setSession(newSession)

	return nil
}

// getSession returns the Keyspaces client's current gocql Session.
func (k *Keyspaces) getSession() *gocql.Session {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	return k.session
}

// setSession sets the Keyspaces client to the provided gocql Session.
func (k *Keyspaces) setSession(s *gocql.Session) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	k.session = s
}

// QueryDatabase queries the database with provided query, and returns a Scanner to iterate through the results.
func (k *Keyspaces) QueryDatabase(query string, values []interface{}) gocql.Scanner {
	return k.getSession().Query(query, values...).Iter().Scanner()
}

// ExecuteQuery executes the provided query on the database.
func (k *Keyspaces) ExecuteQuery(query string, values []interface{}) error {
	return k.getSession().Query(query, values...).Exec()
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
