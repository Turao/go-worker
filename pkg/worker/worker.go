package worker

import (
	"log"

	v1 "github.com/turao/go-worker/api/v1"
	"github.com/turao/go-worker/pkg/job"
	"github.com/turao/go-worker/pkg/storage"
)

type Storage interface {
	Put(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Remove(key string) error
}

type Job interface {
	ID() v1.JobID
	Start() error
	Stop() error
	Info() *v1.JobInfo
}

type Worker struct {
	store Storage
}

func NewWorker() *Worker {
	return &Worker{store: storage.New()}
}

func (w *Worker) Dispatch(name string, args ...string) (v1.JobID, error) {
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

func (w *Worker) Stop(jobID v1.JobID) error {
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

func (w *Worker) QueryInfo(jobID v1.JobID) (*v1.JobInfo, error) {
	item, err := w.store.Get(string(jobID))
	if err != nil {
		return nil, err
	}

	job := item.(Job) // need casting as we don't have generics yet...

	return job.Info(), nil
}
