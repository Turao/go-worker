package storage

import (
	"errors"
	"sync"
)

type store struct {
	mx    *sync.RWMutex
	items map[string]interface{}
}

var ErrNotExists error = errors.New("job does not exist")
var ErrKeyAlreadyTaken error = errors.New("key has already been taken")

func New() *store {
	return &store{
		mx:    &sync.RWMutex{},
		items: make(map[string]interface{}),
	}
}

func (q *store) Put(key string, value interface{}) error {
	q.mx.Lock()
	defer q.mx.Unlock()

	// prevent adding the same job twice
	if _, found := q.items[key]; found {
		return ErrKeyAlreadyTaken
	}

	// could there be a job in this position already (likely not)
	q.items[key] = value
	return nil
}

func (q *store) Get(key string) (interface{}, error) {
	q.mx.RLock()
	defer q.mx.RUnlock()
	job, found := q.items[key]
	if !found {
		return nil, ErrNotExists
	}
	return job, nil
}

func (q *store) Remove(key string) error {
	q.mx.Lock()
	defer q.mx.Unlock()

	if _, found := q.items[key]; !found {
		return ErrNotExists
	}

	delete(q.items, key) // I don't like no-ops
	return nil
}
