package worker

import (
	"log"

	"github.com/turao/go-worker/job"
	"github.com/turao/go-worker/storage"
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
	Status() string
	ExitCode() int
	Output() string
	Errors() string
}

type Worker struct {
	store Storage
}

func NewWorker() *Worker {
	return &Worker{store: storage.NewStore()}
}

func (w *Worker) Dispatch(name string, args ...string) (string, error) {
	log.Println("dispatching new job for command:", name, args)

	job := job.NewJob(name, args...)
	err := w.store.Put(job.ID(), job)
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

func (w *Worker) Stop(jobId string) error {
	item, err := w.store.Get(jobId)
	if err != nil {
		// this could be sensitive, maybe log, maybe don't ...
		log.Println("unable to retrieve job", jobId, err.Error())
		return err
	}

	job := item.(Job) // need casting as we don't have generics yet...
	err = job.Stop()
	if err != nil {
		log.Println("unable to stop job", jobId, err.Error())
		return err
	}

	return nil
}

type JobInfo struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	ExitCode int    `json:"exitCode"`
	Output   string `json:"output"`
	Errors   string `json:"errors"`
}

func (w *Worker) QueryInfo(jobId string) (*JobInfo, error) {
	item, err := w.store.Get(jobId)
	if err != nil {
		return nil, err
	}

	job := item.(Job) // need casting as we don't have generics yet...

	return &JobInfo{
		ID:       job.ID(),
		Status:   job.Status(),
		ExitCode: job.ExitCode(),
		Output:   job.Output(),
		Errors:   job.Errors(),
	}, nil
}
