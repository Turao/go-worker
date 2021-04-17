package worker

import (
	"errors"
	"sync"
)

type pool struct {
	mx   *sync.RWMutex
	jobs map[string]*job
}

func NewPool(size int) *pool {
	return &pool{
		mx:   &sync.RWMutex{},
		jobs: make(map[string]*job, size),
	}
}

func (q *pool) Put(key string, value *job) error {
	q.mx.Lock()
	defer q.mx.Unlock()

	// prevent adding the same job twice
	if _, found := q.jobs[key]; found {
		return ErrIdAlreadyTaken
	}

	// could there be a job in this position already (likely not)
	q.jobs[key] = value
	return nil
}

func (q *pool) Get(id string) (*job, error) {
	q.mx.RLock()
	defer q.mx.RUnlock()
	job, found := q.jobs[id]
	if !found {
		return nil, ErrNotExists
	}
	return job, nil
}

func (q *pool) Remove(id string) error {
	q.mx.Lock()
	defer q.mx.Unlock()

	if _, found := q.jobs[id]; !found {
		return ErrNotExists
	}

	delete(q.jobs, id) // I don't like no-ops
	return nil
}

var ErrNotExists error = errors.New("job does not exist")
var ErrIdAlreadyTaken error = errors.New("id has already been taken")
