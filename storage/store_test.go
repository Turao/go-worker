package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type item struct {
	key   string
	value string
}

func givenItem() *item {
	return &item{key: "key", value: "value"}
}

func givenStoreWithJobAt(t *testing.T, key string) *store {
	store := NewStore(100)
	job := givenItem()

	err := store.Put(key, job)
	if err != nil {
		t.Fatal("unable to store mock data in store")
	}
	return store
}

func TestStorePutOnce(t *testing.T) {
	job := givenItem()
	store := NewStore(100)

	err := store.Put("id", job)
	assert.Nil(t, err)
}

func TestStorePutTwice(t *testing.T) {
	job := givenItem()
	store := NewStore(100)
	store.Put("id", job)

	err := store.Put("id", job)
	assert.Equal(t, ErrKeyAlreadyTaken, err)
}

func TestStoreGet(t *testing.T) {
	store := givenStoreWithJobAt(t, "id")

	job, err := store.Get("id")
	assert.Nil(t, err)
	assert.NotNil(t, job)
}

func TestStoreGetNotExists(t *testing.T) {
	store := givenStoreWithJobAt(t, "id")

	_, err := store.Get("other")
	assert.Equal(t, ErrNotExists, err)
}

func TestStoreRemove(t *testing.T) {
	store := givenStoreWithJobAt(t, "id")

	err := store.Remove("id")
	assert.Nil(t, err)
}

func TestStoreRemoveNotExists(t *testing.T) {
	store := givenStoreWithJobAt(t, "id")

	err := store.Remove("other")
	assert.Equal(t, ErrNotExists, err)
}
