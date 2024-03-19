package s3

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestS3GetAllConcurrently(t *testing.T) {

	// ARRANGE
	setup()
	defer teardown()

	// ASSERT parameter errors.
	_, err := NewConcurrent(0, 100, 1000)
	assert.NotNil(t, err)
	_, err = NewConcurrent(100, 0, 1000)
	assert.NotNil(t, err)
	_, err = NewConcurrent(100, 100, 0)
	assert.NotNil(t, err)
	_, err = NewConcurrent(100, 10, 99)
	assert.NotNil(t, err)
	_, err = NewConcurrent(100, 101, 1000)
	assert.NotNil(t, err)

	client, err := NewConcurrent(100, 10, 1000)
	require.Nil(t, err, fmt.Sprintf("error creating s3 client concurrency manager: %v", err))

	// ASSERT computed fields.
	assert.Equal(t, 100, len(client.manager.workerPool.channel))
	assert.Equal(t, 100, len(client.manager.memoryPool.channel))
	assert.Equal(t, int64(10), client.manager.memoryChunkSize)
	assert.Equal(t, int64(10*100), client.manager.memoryTotalSize)
	assert.Equal(t, 10, client.manager.maxWorkersPerRequest)

	// ASSERT memory chunk size is correct in memory pool.
	chunk := <-client.manager.memoryPool.channel
	assert.Equal(t, int64(10), chunk)
	client.manager.memoryPool.channel <- chunk

	// ASSERT worker get/release methods work expectedly.
	w := client.manager.getWorkers(1)
	assert.Equal(t, 99, len(client.manager.workerPool.channel))
	client.manager.returnWorker(w[0])
	assert.Equal(t, 100, len(client.manager.workerPool.channel))

	// ASSERT memory get/release methods work expectedly.
	elevenByteFile := types.Object{Size: aws.Int64(11)} // requires 2 memory chunks.
	client.manager.secureMemory([]types.Object{elevenByteFile})
	assert.Equal(t, 98, len(client.manager.memoryPool.channel))
	client.manager.releaseMemory(20)
	assert.Equal(t, 100, len(client.manager.memoryPool.channel))

	// ARRANGE bucket with test objects.
	total := 20
	keys := make([]string, total)
	for i := 0; i < total; i++ {
		keys[i] = fmt.Sprintf("%s-%v", testObjectKey, i)
	}
	awsCmdPutKeys(keys)

	// ACTION
	objects, _ := client.ListAllObjects(testBucket, "")
	tooManyBytes := make([]types.Object, 10*len(objects))
	for _, o := range objects {
		for i := 0; i < 10; i++ {
			tooManyBytes = append(tooManyBytes, o)
		}
	}
	output := client.GetAllConcurrently(testBucket, "", tooManyBytes)

	// ASSERT error returned
	for hf := range output {
		assert.NotNil(t, hf.Error)
	}

	// ACTION
	objects, _ = client.ListAllObjects(testBucket, "")
	output = client.GetAllConcurrently(testBucket, "", objects)
	outputKeys := make([]string, 0)
	for hf := range output {
		outputKeys = append(outputKeys, hf.Key)
	}

	// ASSERT input and output order is the same.
	require.Equal(t, len(outputKeys), total)
	for i := 0; i < total; i++ {
		assert.Equal(t, aws.ToString(objects[i].Key), outputKeys[i])
	}

	// ASSERT all workers and memory returned to pools.
	time.Sleep(2 * time.Second)
	assert.Equal(t, 100, len(client.manager.workerPool.channel))
	assert.Equal(t, 100, len(client.manager.memoryPool.channel))

	// ASSERT that process blocked when all memory secured.
	tenByteFile := types.Object{Size: aws.Int64(10)}
	oneThousandBytesOfFiles := make([]types.Object, 100)
	for i := 0; i < 100; i++ {
		oneThousandBytesOfFiles[i] = tenByteFile
	}
	client.manager.secureMemory(oneThousandBytesOfFiles)
	ch := make(chan chan HydratedFile)
	go func() {
		ch <- client.GetAllConcurrently(testBucket, "", objects)
	}()

	for {
		select {
		case <-ch:
			t.Error("process was not blocked")
		case <-time.After(time.Second):
			// Timed out as expected
			return
		}
	}
}
