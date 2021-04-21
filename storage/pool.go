package storage

import (
	"errors"
	"sync"
)

type pool struct {
	mx    *sync.RWMutex
	items map[string]interface{}
}

func NewPool(size int) *pool {
	return &pool{
		mx:    &sync.RWMutex{},
		items: make(map[string]interface{}, size),
	}
}

func (q *pool) Put(key string, value interface{}) error {
	q.mx.Lock()
	defer q.mx.Unlock()

	// prevent adding the same job twice
	if _, found := q.items[key]; found {
		return ErrIdAlreadyTaken
	}

	// could there be a job in this position already (likely not)
	q.items[key] = value
	return nil
}

func (q *pool) Get(id string) (interface{}, error) {
	q.mx.RLock()
	defer q.mx.RUnlock()
	job, found := q.items[id]
	if !found {
		return nil, ErrNotExists
	}
	return job, nil
}

func (q *pool) Remove(id string) error {
	q.mx.Lock()
	defer q.mx.Unlock()

	if _, found := q.items[id]; !found {
		return ErrNotExists
	}

	delete(q.items, id) // I don't like no-ops
	return nil
}

var ErrNotExists error = errors.New("job does not exist")
var ErrIdAlreadyTaken error = errors.New("id has already been taken")