//go:build localstack
// +build localstack

package s3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	CustomAWSEndpointURL = "http://localhost:4566"
	AWSRegion            = "ap-southeast-2"
	TestBucket           = "test-bucket"

	TestObjectKey  = "test-key"
	TestObjectData = "some data"

	TestMetaKey   = "test-meta-key"
	TestMetaValue = "test-meta-value"

	TestPrefix          = "test-prefix"
	TestPrefixDelimiter = "_"

	TestNewKey = "test-new-key"
)

// helper functions

func setup() {
	// setup environment variables to access LocalStack
	os.Setenv("AWS_REGION", AWSRegion)
	os.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("CUSTOM_AWS_ENDPOINT_URL", CustomAWSEndpointURL)

	// create bucket
	if err := exec.Command(
		"aws", "s3api",
		"create-bucket",
		"--bucket", TestBucket,
		"--create-bucket-configuration", fmt.Sprintf(
			"{\"LocationConstraint\": \"%v\"}", AWSRegion),
		"--endpoint-url", CustomAWSEndpointURL).Run(); err != nil {

		panic(err)
	}
}

func teardown() {
	if err := exec.Command(
		"aws", "s3",
		"rb", fmt.Sprintf("s3://%v", TestBucket),
		"--force",
		"--endpoint-url", CustomAWSEndpointURL).Run(); err != nil {

		panic(err)
	}
}

func awsCmdPopulateBucket() {
	// create test data
	tmpDir, _ := os.MkdirTemp("", "")
	defer os.RemoveAll(tmpDir)

	testDataFilepath := filepath.Join(tmpDir, "data.txt")
	testFile, _ := os.Create(testDataFilepath)
	testFile.WriteString(TestObjectData)
	testFile.Close()

	// populate bucket
	if err := exec.Command(
		"aws", "s3api",
		"put-object",
		"--bucket", TestBucket,
		"--key", TestObjectKey,
		"--body", testDataFilepath,
		"--metadata", fmt.Sprintf("%v=%v", TestMetaKey, TestMetaValue),
		"--endpoint-url", CustomAWSEndpointURL).Run(); err != nil {

		panic(err)
	}
}

func awsCmdExists(key string) bool {
	if err := exec.Command(
		"aws", "s3api",
		"head-object",
		"--bucket", TestBucket,
		"--key", key,
		"--endpoint-url", CustomAWSEndpointURL).Run(); err != nil {

		return false
	}
	return true
}

func awsCmdPutKey(key string) {
	if err := exec.Command(
		"aws", "s3api",
		"put-object",
		"--bucket", TestBucket,
		"--key", key,
		"--endpoint-url", CustomAWSEndpointURL).Run(); err != nil {

		panic(err)
	}
}

type awsMeta struct {
	lastModified  time.Time
	contentLength int64
	meta          map[string]string
}

func awsCmdMeta() awsMeta {
	var rvalue = awsMeta{
		meta: make(map[string]string),
	}

	cmd := exec.Command(
		"aws", "s3api",
		"head-object",
		"--bucket", TestBucket,
		"--key", TestObjectKey,
		"--endpoint-url", CustomAWSEndpointURL)

	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	var metaData map[string]interface{}
	json.Unmarshal(out, &metaData)
	testLastModified, err := time.Parse(
		"2006-01-02T15:04:05", metaData["LastModified"].(string)[0:19])
	if err != nil {
		panic(err)
	}

	rvalue.lastModified = testLastModified
	rvalue.contentLength = int64(metaData["ContentLength"].(float64))
	// note: assumes string=string type aws meta data
	for k, v := range metaData["Metadata"].(map[string]interface{}) {
		rvalue.meta[k] = v.(string)
	}

	return rvalue
}

func awsCmdGetTestObject() string {
	tmpDir, _ := os.MkdirTemp("", "")
	defer os.RemoveAll(tmpDir)

	testDataFilepath := filepath.Join(tmpDir, "data.txt")

	if err := exec.Command(
		"aws", "s3api",
		"get-object",
		"--bucket", TestBucket,
		"--key", TestObjectKey,
		"--endpoint-url", CustomAWSEndpointURL,
		testDataFilepath).Run(); err != nil {
		panic(err)
	}

	testFileContents, _ := os.ReadFile(testDataFilepath)
	return string(testFileContents)
}

// THE TESTS

func TestCreateS3ClientAndReady(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	// test good case

	// ACTION
	client, err := New()

	// ASSERT
	assert.Nil(t, err)
	assert.True(t, client.Ready())

	// test bad case

	// ARRANGE
	os.Unsetenv("AWS_REGION")

	// ACTION
	client, err = New()

	// ASSERT
	assert.NotNil(t, err)
	assert.False(t, client.Ready())
}

func TestCreateS3ClientWithMaxRetries(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	// ACTION
	_, err := NewWithMaxRetries(2)

	// ASSERT
	assert.Nil(t, err)

	// test bad case

	// ARRANGE
	os.Unsetenv("AWS_REGION")

	// ACTION
	_, err = NewWithMaxRetries(2)

	// ASSERT
	assert.NotNil(t, err)
}
func TestS3Get(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	awsCmdPopulateBucket()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating s3 client: %v", err))

	// ACTION
	dataObject := bytes.Buffer{}
	err = client.Get(TestBucket, TestObjectKey, "", &dataObject)

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, TestObjectData, dataObject.String())
}

func TestS3GetWithLastModified(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	awsCmdPopulateBucket()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating s3 client: %v", err))

	meta := awsCmdMeta()

	// ACTION
	dataObject := bytes.Buffer{}
	lastModified, err := client.GetWithLastModified(
		TestBucket, TestObjectKey, "", &dataObject)

	// ASSERT
	assert.Nil(t, err)

	// object is what we expect
	assert.Equal(t, TestObjectData, dataObject.String())

	// modified time is what we expect
	assert.Equal(t, meta.lastModified, lastModified)

	// bonus test, LastModified function
	lastModified, err = client.LastModified(TestBucket, TestObjectKey, "")
	assert.Nil(t, err)
	assert.Equal(t, meta.lastModified, lastModified)
}

func TestS3GetMeta(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	awsCmdPopulateBucket()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating s3 client: %v", err))

	// ACTION
	meta, err := client.GetMeta(TestBucket, TestObjectKey, "")

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, TestMetaValue, meta[TestMetaKey])
}

func TestS3GetContentSizeTime(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	awsCmdPopulateBucket()

	meta := awsCmdMeta()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating s3 client: %v", err))

	// ACTION
	contentLength, lastModified, err := client.GetContentSizeTime(TestBucket, TestObjectKey)

	// ASSERT
	assert.Equal(t, meta.contentLength, contentLength)
	assert.Equal(t, meta.lastModified, lastModified)
}

func TestS3Put(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating s3 client: %v", err))

	// ACTION
	err = client.Put(TestBucket, TestObjectKey, []byte(TestObjectData))

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, TestObjectData, awsCmdGetTestObject())
}

func TestS3PutWithMetadata(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating s3 client: %v", err))

	// ACTION
	err = client.PutWithMetadata(
		TestBucket,
		TestObjectKey,
		[]byte(TestObjectData),
		map[string]string{TestMetaKey: TestMetaValue})

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, TestObjectData, awsCmdGetTestObject())

	// test meta data
	metaData := awsCmdMeta()
	assert.Equal(t, TestMetaValue, metaData.meta[TestMetaKey])
}

func TestS3Exists(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	awsCmdPopulateBucket()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating s3 client: %v", err))

	// ACTION
	exists, err := client.Exists(TestBucket, TestObjectKey)

	// ASSERT
	assert.Nil(t, err)
	assert.True(t, exists)

	// ACTION
	exists, err = client.Exists(TestBucket, "thisdoesntexists")

	// ASSERT
	assert.Nil(t, err)
	assert.False(t, exists)
}

func TestS3List(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	// create several test inputs
	var tests = []struct {
		prefix     string
		numObjects int
	}{
		{"prefix1", 2},
		{"prefix2", 3},
	}

	//populate bucket with several objects, some with common prefix
	for _, tt := range tests {
		for i := 0; i < tt.numObjects; i++ {
			awsCmdPutKey(fmt.Sprintf("%v-%v", tt.prefix, i))
		}
	}

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating s3 client: %v", err))

	for _, tt := range tests {
		// ACTION
		listing, err := client.List(TestBucket, tt.prefix, 1000)

		// ASSERT

		// got listing ok
		assert.Nil(t, err)

		// expected number of objects in list
		assert.Equal(t, tt.numObjects, len(listing))

		// expected prefix
		for _, item := range listing {
			assert.True(t, strings.HasPrefix(item, tt.prefix))
		}

		// sneak in a ListAll test as well
		listAll, err := client.ListAll(TestBucket, tt.prefix)
		assert.Nil(t, err)
		assert.Equal(t, listing, listAll)
	}
}

func TestS3PrefixExists(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	awsCmdPutKey(TestPrefix + "/" + TestObjectKey)

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating s3 client: %v", err))

	// ACTION
	exists, err := client.PrefixExists(TestBucket, TestPrefix)

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, true, exists)

	// ACTION
	exists, err = client.PrefixExists(TestBucket, "thisdoesnotexists")

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, false, exists)
}

func TestS3ListCommonPrefixes(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	// create several test inputs
	var testKeys = []string{
		// prefix_1 etc.
		TestPrefix + TestPrefixDelimiter + "file1",
		TestPrefix + TestPrefixDelimiter + "file2",
		TestPrefix + TestPrefixDelimiter + "file3",
		TestPrefix + TestPrefixDelimiter + "file4",
	}

	//populate bucket with several objects, some with common prefix
	for _, key := range testKeys {
		awsCmdPutKey(key)
	}

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating s3 client: %v", err))

	// ACTION
	listing, err := client.ListCommonPrefixes(TestBucket, TestPrefix, "_")

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, []string{TestPrefix + TestPrefixDelimiter}, listing)

	// test ListObjects

	// ACTION
	objects, err := client.ListObjects(TestBucket, TestPrefix)
	var keys []string
	for _, object := range objects {
		keys = append(keys, *object.Key)
	}

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, testKeys, keys)
}

func TestS3Delete(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	awsCmdPopulateBucket()
	require.True(t, awsCmdExists(TestObjectKey))

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating s3 client: %v", err))

	// ACTION
	err = client.Delete(TestBucket, TestObjectKey)

	// ASSERT
	assert.Nil(t, err)

	// should no longer exists
	assert.False(t, awsCmdExists(TestObjectKey))
}

func TestS3Copy(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	awsCmdPopulateBucket()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating s3 client: %v", err))

	// ACTION
	err = client.Copy(TestBucket, TestNewKey, TestBucket+"/"+TestObjectKey)

	// ASSERT
	assert.Nil(t, err)

	// new object exists
	assert.True(t, awsCmdExists(TestNewKey))
}
