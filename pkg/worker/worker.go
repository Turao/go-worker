package worker

import (
	"log"

	"github.com/turao/go-worker/pkg/job"
	"github.com/turao/go-worker/pkg/storage"
)

type Storage interface {
	Put(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Remove(key string) error
}

type Job interface {
	ID() string
	Start() error
	Stop() error
	Info() *job.JobInfo
}

type worker struct {
	store Storage
}

func NewWorker() *worker {
	return &worker{store: storage.New()}
}

func (w *worker) Dispatch(name string, args ...string) (string, error) {
	log.Println("dispatching new job for command:", name, args)

	job := job.New(name, args...)
	err := w.store.Put(string(job.ID()), job)
	if err != nil {
		log.Println("unable to store command", err.Error())
		return "", err
	}

	err = job.Start()
	if err != nil {
		log.Println("unable to dispatch command", err.Error())
		return "", err
	}

	return job.ID(), nil
}

func (w *worker) Stop(jobID string) error {
	item, err := w.store.Get(string(jobID))
	if err != nil {
		log.Println("unable to retrieve job", jobID, err.Error())
		return err
	}

	job := item.(Job) // need casting as we don't have generics yet...
	err = job.Stop()
	if err != nil {
		log.Println("unable to stop job", jobID, err.Error())
		return err
	}

	return nil
}

func (w *worker) QueryInfo(jobID string) (*job.JobInfo, error) {
	item, err := w.store.Get(string(jobID))
	if err != nil {
		return nil, err
	}

	job := item.(Job) // need casting as we don't have generics yet...

	return job.Info(), nil
}
