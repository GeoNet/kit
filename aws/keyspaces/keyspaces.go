// Package keyspaces is for working with AWS Keyspaces (for Apache Cassandra).
package keyspaces

import (
	"errors"
	"fmt"
	"os"

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

	var region, service_username, service_password string
	if region = os.Getenv("AWS_REGION"); region == "" {
		return ksClient, errors.New("AWS_REGION is not set")
	}
	if service_username = os.Getenv("KEYSPACE_USER"); service_username == "" {
		return ksClient, errors.New("KEYSPACE_USER is not set")
	}
	if service_password = os.Getenv("KEYSPACE_PW"); service_password == "" {
		return ksClient, errors.New("KEYSPACE_PW is not set")
	}

	// Add the Amazon Keyspaces service endpoint
	cluster := gocql.NewCluster(fmt.Sprintf("cassandra.%s.amazonaws.com:9142", region))

	// Add your service specific credentials
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: service_username,
		Password: service_password,
	}
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
