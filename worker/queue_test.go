package worker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func givenJob() *job {
	return NewJob("ls", "-lah")
}

func givenQueueWithJobAt(t *testing.T, key string) *queue {
	queue := NewQueue(100)
	job := givenJob()

	err := queue.Put(key, job)
	if err != nil {
		t.Fatal("unable to store mock data in queue")
	}
	return queue
}

func TestQueuePutOnce(t *testing.T) {
	job := givenJob()
	queue := NewQueue(100)

	err := queue.Put("id", job)
	assert.Nil(t, err)
}

func TestQueuePutTwice(t *testing.T) {
	job := givenJob()
	queue := NewQueue(100)
	queue.Put("id", job)

	err := queue.Put("id", job)
	assert.Equal(t, ErrIdAlreadyTaken, err)
}

func TestQueueGet(t *testing.T) {
	queue := givenQueueWithJobAt(t, "id")

	job, err := queue.Get("id")
	assert.Nil(t, err)
	assert.NotNil(t, job)
}

func TestQueueGetNotExists(t *testing.T) {
	queue := givenQueueWithJobAt(t, "id")

	_, err := queue.Get("other")
	assert.Equal(t, ErrNotExists, err)
}

func TestQueueRemove(t *testing.T) {
	queue := givenQueueWithJobAt(t, "id")

	err := queue.Remove("id")
	assert.Nil(t, err)
}

func TestQueueRemoveNotExists(t *testing.T) {
	queue := givenQueueWithJobAt(t, "id")

	err := queue.Remove("other")
	assert.Equal(t, ErrNotExists, err)
}
