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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	customAWSEndpoint = "http://localhost:4566"
	testRegion        = "ap-southeast-2"
	testBucket        = "test-bucket"

	testObjectKey  = "test-key"
	testObjectData = "some data"

	testMetaKey   = "test-meta-key"
	testMetaValue = "test-meta-value"

	testPrefix          = "test-prefix"
	testPrefixDelimiter = "_"

	testNewKey = "test-new-key"
)

// helper functions

func setAwsEnv() {
	os.Setenv("AWS_REGION", testRegion)
	os.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_ENDPOINT_URL", customAWSEndpoint)
}

func setup() {
	// setup environment variable to run AWS CLI/SDK
	setAwsEnv()

	// create bucket
	if err := exec.Command( //nolint:gosec
		"aws", "s3api",
		"create-bucket",
		"--bucket", testBucket,
		"--create-bucket-configuration", fmt.Sprintf(
			"{\"LocationConstraint\": \"%v\"}", testRegion),
	).Run(); err != nil {
		panic(err)
	}
}

func teardown() {
	setAwsEnv()

	if err := exec.Command( //nolint:gosec
		"aws", "s3",
		"rb", fmt.Sprintf("s3://%v", testBucket),
		"--force",
	).Run(); err != nil {

		panic(err)
	}
}

func awsCmdPopulateBucket() {
	// create test data
	tmpDir, _ := os.MkdirTemp("", "")
	defer os.RemoveAll(tmpDir)

	testDataFilepath := filepath.Join(tmpDir, "data.txt")
	testFile, _ := os.Create(testDataFilepath)
	_, _ = testFile.WriteString(testObjectData)
	testFile.Close()

	// populate bucket
	if err := exec.Command(
		"aws", "s3api",
		"put-object",
		"--bucket", testBucket,
		"--key", testObjectKey,
		"--body", testDataFilepath,
		"--metadata", fmt.Sprintf("%v=%v", testMetaKey, testMetaValue),
	).Run(); err != nil {

		panic(err)
	}
}

func awsCmdBucketExists(bucket string) bool {
	if err := exec.Command(
		"aws", "s3api",
		"head-bucket",
		"--bucket", bucket,
	).Run(); err != nil {
		return false
	}
	return true
}

func awsCmdExists(key string) bool {
	if err := exec.Command(
		"aws", "s3api",
		"head-object",
		"--bucket", testBucket,
		"--key", key,
	).Run(); err != nil {

		return false
	}
	return true
}

func awsCmdPutKey(key string) {
	if err := exec.Command(
		"aws", "s3api",
		"put-object",
		"--bucket", testBucket,
		"--key", key,
	).Run(); err != nil {

		panic(err)
	}
}

func awsCmdPutKeys(keys []string) {
	// create test data
	tmpDir, _ := os.MkdirTemp("", "")
	defer os.RemoveAll(tmpDir)

	for _, k := range keys {
		testDataFilepath := filepath.Join(tmpDir, k)
		testFile, _ := os.Create(testDataFilepath)
		_, _ = testFile.WriteString(testObjectData)
		testFile.Close()
	}
	// sync to bucket
	if err := exec.Command(
		"aws", "s3",
		"sync", tmpDir, fmt.Sprintf("s3://%v", testBucket),
	).Run(); err != nil {

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
		"--bucket", testBucket,
		"--key", testObjectKey)

	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	var metaData map[string]interface{}
	_ = json.Unmarshal(out, &metaData)
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
		"--bucket", testBucket,
		"--key", testObjectKey,
		testDataFilepath).Run(); err != nil {
		panic(err)
	}

	testFileContents, _ := os.ReadFile(testDataFilepath)
	return string(testFileContents)
}

// THE TESTS

func TestCheckBucket(t *testing.T) {
	setup()
	defer teardown()

	// test
	client, err := New()
	assert.Nil(t, err)
	//test existing bucket
	err = client.CheckBucket(testBucket)
	assert.Nil(t, err)

	//test none existing bucket
	testBucket1 := "test1"
	err = client.CheckBucket(testBucket1)
	assert.NotNil(t, err)

}

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

func TestCreateS3ClientWithOptions(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	awsCmdPopulateBucket()

	// NOTE: This test case is to make sure setting endpoint to nil can cause an error.
	//       However, BaseEndpoint is deprecated in V2 and new way is quite complicated.
	//       Not worth the efforts in testing this.
	// // ACTION
	// s3Client, err := NewWithOptions(func(options *s3.Options) {
	// 	options.BaseEndpoint = nil
	// })

	// // ASSERT
	// assert.Nil(t, err)

	// // ACTION
	// _, err := s3Client.ListAll(testBucket, "")

	// // ASSERT
	// assert.NotNil(t, err)

	// ACTION
	s3Client, err := NewWithOptions(func(options *s3.Options) {
		options.Region = testRegion
	})

	// ASSERT
	assert.Nil(t, err)

	// ACTION
	_, err = s3Client.ListAll(testBucket, "")

	// ASSERT
	assert.Nil(t, err)
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
	err = client.Get(testBucket, testObjectKey, "", &dataObject)

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, testObjectData, dataObject.String())
}

func TestS3GetByteRange(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	awsCmdPopulateBucket()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating s3 client: %v", err))

	// ACTION
	testByteRange := "bytes=0-2"
	dataObject := bytes.Buffer{}
	err = client.GetByteRange(testBucket, testObjectKey, "", testByteRange, &dataObject)

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, testObjectData[:3], dataObject.String())
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
		testBucket, testObjectKey, "", &dataObject)

	// ASSERT
	assert.Nil(t, err)

	// object is what we expect
	assert.Equal(t, testObjectData, dataObject.String())

	// modified time is what we expect
	assert.Equal(t, meta.lastModified, lastModified)

	// bonus test, LastModified function
	lastModified, err = client.LastModified(testBucket, testObjectKey, "")
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
	meta, err := client.GetMeta(testBucket, testObjectKey, "")

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, testMetaValue, meta[testMetaKey])

	// for non existing object, we want an empty result instead of not found error
	meta, err = client.GetMeta(testBucket, fmt.Sprintf("%s.%d", testObjectKey, time.Now().Unix()), "")

	// ASSERT
	assert.Nil(t, err)
	assert.Empty(t, meta)
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
	contentLength, lastModified, err := client.GetContentSizeTime(testBucket, testObjectKey)

	// ASSERT
	assert.Nil(t, err)
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
	err = client.Put(testBucket, testObjectKey, []byte(testObjectData))

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, testObjectData, awsCmdGetTestObject())
}

func TestS3PutWithMetadata(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating s3 client: %v", err))

	// ACTION
	err = client.PutWithMetadata(
		testBucket,
		testObjectKey,
		[]byte(testObjectData),
		map[string]string{testMetaKey: testMetaValue})

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, testObjectData, awsCmdGetTestObject())

	// test meta data
	metaData := awsCmdMeta()
	assert.Equal(t, testMetaValue, metaData.meta[testMetaKey])
}

func TestS3Exists(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	awsCmdPopulateBucket()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating s3 client: %v", err))

	// ACTION
	exists, err := client.Exists(testBucket, testObjectKey)

	// ASSERT
	assert.Nil(t, err)
	assert.True(t, exists)

	// ACTION
	exists, err = client.Exists(testBucket, "thisdoesntexists")

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
		listing, err := client.List(testBucket, tt.prefix, 1000)

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
		listAll, err := client.ListAll(testBucket, tt.prefix)
		assert.Nil(t, err)
		assert.Equal(t, listing, listAll)
	}
}

func TestS3ListAll(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	// populate bucket with over 1000 objects (the limit at
	// which ListAll's continuation token functionality is required).
	numObjects := 1005
	keys := make([]string, 0)
	keyGroups := [11]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K"}
	for i := 0; i < numObjects; i++ {
		keyGroup := keyGroups[i/100]
		keys = append(keys, fmt.Sprintf("%s%s%s%06d", testPrefix, keyGroup, testPrefixDelimiter, i))
	}
	awsCmdPutKeys(keys)

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating s3 client: %v", err))

	// ACTION
	listing, err := client.ListAll(testBucket, testPrefix)

	// ASSERT

	// got listing ok
	assert.Nil(t, err)

	// expected objects in list in the correct order.
	assert.Equal(t, listing, keys)

	// ACTION
	testPrefixes := make([]string, 0)
	for _, keyGroup := range keyGroups {
		testPrefixes = append(testPrefixes, testPrefix+keyGroup)
	}
	listingConcurrent, err := client.ListAllObjectsConcurrently(testBucket, testPrefixes)
	var keysConcurrent []string
	for _, object := range listingConcurrent {
		keysConcurrent = append(keysConcurrent, aws.ToString(object.Key))
	}

	// ASSERT

	// got listing ok
	assert.Nil(t, err)

	// expected objects in list in the correct order.
	assert.Equal(t, keysConcurrent, keys)
}

func TestS3PrefixExists(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	awsCmdPutKey(testPrefix + "/" + testObjectKey)

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating s3 client: %v", err))

	// ACTION
	exists, err := client.PrefixExists(testBucket, testPrefix)

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, true, exists)

	// ACTION
	exists, err = client.PrefixExists(testBucket, "thisdoesnotexists")

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
		testPrefix + testPrefixDelimiter + "file1",
		testPrefix + testPrefixDelimiter + "file2",
		testPrefix + testPrefixDelimiter + "file3",
		testPrefix + testPrefixDelimiter + "file4",
	}

	//populate bucket with several objects, some with common prefix
	for _, key := range testKeys {
		awsCmdPutKey(key)
	}

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating s3 client: %v", err))

	// ACTION
	listing, err := client.ListCommonPrefixes(testBucket, testPrefix, "_")

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, []string{testPrefix + testPrefixDelimiter}, listing)

	// test ListObjects

	// ACTION
	objects, err := client.ListObjects(testBucket, testPrefix, 1000)
	var keys []string
	for _, object := range objects {
		keys = append(keys, *object.Key)
	}

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, testKeys, keys)

	// test ListObjects with limit

	// ACTION
	objectsLimited, err := client.ListObjects(testBucket, testPrefix, 2)
	var keysLimited []string
	for _, object := range objectsLimited {
		keysLimited = append(keysLimited, *object.Key)
	}

	// ASSERT
	assert.Nil(t, err)
	assert.Equal(t, testKeys[:2], keysLimited)
}

func TestS3Delete(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	awsCmdPopulateBucket()
	require.True(t, awsCmdExists(testObjectKey))

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating s3 client: %v", err))

	// ACTION
	err = client.Delete(testBucket, testObjectKey)

	// ASSERT
	assert.Nil(t, err)

	// should no longer exists
	assert.False(t, awsCmdExists(testObjectKey))
}

func TestS3Copy(t *testing.T) {
	// ARRANGE
	setup()
	defer teardown()

	awsCmdPopulateBucket()

	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating s3 client: %v", err))

	// ACTION
	err = client.Copy(testBucket, testNewKey, testBucket+"/"+testObjectKey)

	// ASSERT
	assert.Nil(t, err)

	// new object exists
	assert.True(t, awsCmdExists(testNewKey))
}

func TestS3CreateBucket(t *testing.T) {
	// ARRANGE
	setAwsEnv()
	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating s3 client: %v", err))

	bucket := "new-bucket"
	require.False(t, awsCmdBucketExists(bucket), "error arranging test, bucket already exists")

	// ACTION
	err = client.CreateBucket(bucket)
	t.Cleanup(func() {
		if err := client.DeleteBucket(bucket); err != nil {
			t.Fatalf("Failed to delete bucket during cleanup: %v", err)
		}
	})

	// ASSERT
	assert.Nil(t, err)
	assert.True(t, awsCmdBucketExists(bucket))
}

func TestS3DeleteBucket(t *testing.T) {
	// ARRANGE
	setAwsEnv()
	client, err := New()
	require.Nil(t, err, fmt.Sprintf("Error creating s3 client: %v", err))

	bucket := "bucket-to-delete"
	require.Nil(t, client.CreateBucket(bucket), "error arranging test, couldn't create bucket to delete")
	require.True(t, awsCmdBucketExists(bucket), "error arranging test, doesn't exist")

	// ACTION
	err = client.DeleteBucket(bucket)

	// ASSERT
	assert.Nil(t, err)
	assert.False(t, awsCmdBucketExists(bucket))
}
