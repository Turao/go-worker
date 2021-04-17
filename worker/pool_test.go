package worker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func givenJob() *job {
	return NewJob("ls", "-lah")
}

func givenPoolWithJobAt(t *testing.T, key string) *pool {
	pool := NewPool(100)
	job := givenJob()

	err := pool.Put(key, job)
	if err != nil {
		t.Fatal("unable to store mock data in pool")
	}
	return pool
}

func TestPoolPutOnce(t *testing.T) {
	job := givenJob()
	pool := NewPool(100)

	err := pool.Put("id", job)
	assert.Nil(t, err)
}

func TestPoolPutTwice(t *testing.T) {
	job := givenJob()
	pool := NewPool(100)
	pool.Put("id", job)

	err := pool.Put("id", job)
	assert.Equal(t, ErrIdAlreadyTaken, err)
}

func TestPoolGet(t *testing.T) {
	pool := givenPoolWithJobAt(t, "id")

	job, err := pool.Get("id")
	assert.Nil(t, err)
	assert.NotNil(t, job)
}

func TestPoolGetNotExists(t *testing.T) {
	pool := givenPoolWithJobAt(t, "id")

	_, err := pool.Get("other")
	assert.Equal(t, ErrNotExists, err)
}

func TestPoolRemove(t *testing.T) {
	pool := givenPoolWithJobAt(t, "id")

	err := pool.Remove("id")
	assert.Nil(t, err)
}

func TestPoolRemoveNotExists(t *testing.T) {
	pool := givenPoolWithJobAt(t, "id")

	err := pool.Remove("other")
	assert.Equal(t, ErrNotExists, err)
}
