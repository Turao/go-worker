package worker

import (
	"errors"
	"sync"
)

type queue struct {
	mx   *sync.Mutex
	jobs map[int]*job
}

func makeQueue(size int) *queue {
	return &queue{
		mx:   &sync.Mutex{},
		jobs: make(map[int]*job, size),
	}
}

func (q *queue) addJob(job *job) (int, error) {
	q.mx.Lock()
	defer q.mx.Unlock()

	// prevent adding the same job twice
	if job.id != -1 {
		return -1, errors.New("job already has an id assigned, thus was already added")
	}

	id := len(q.jobs)
	job.id = id

	q.jobs[id] = job // could there be a job in this position already (likely not)
	return id, nil
}

func (q *queue) getJob(id int) (*job, error) {
	q.mx.Lock()
	defer q.mx.Unlock()
	job, found := q.jobs[id]
	if !found {
		return nil, errors.New("job does not exist")
	}
	return job, nil
}

func (q *queue) removeJob(id int) error {
	q.mx.Lock()
	defer q.mx.Unlock()
	_, found := q.jobs[id]
	if !found {
		return errors.New("job does not exist")
	}

	delete(q.jobs, id)
	return nil
}
