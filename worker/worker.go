package worker

import (
	"fmt"
	"log"

	"github.com/turao/kami-go/job"
	"github.com/turao/kami-go/storage"
)

type Storage interface {
	Put(key string, value interface{}) error
	Get(id string) (interface{}, error)
	Remove(id string) error
}

type Job interface {
	ID() string
	Start() error
	Stop() error
	Status() string
	ExitCode() int
	Logs() (string, string)
}

type Worker struct {
	store Storage
}

func NewWorker() *Worker {
	defaultPoolSize := 100
	return &Worker{store: storage.NewPool(defaultPoolSize)}
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
	Id       string `json:"id"`
	Status   string `json:"status"`
	ExitCode int    `json:"exitCode"`
}

func (w *Worker) QueryInfo(jobId string) (*JobInfo, error) {
	item, err := w.store.Get(jobId)
	if err != nil {
		return nil, err
	}

	job := item.(Job) // need casting as we don't have generics yet...
	return &JobInfo{
		Id:       fmt.Sprint(job.ID()),
		Status:   string(job.Status()),
		ExitCode: job.ExitCode(),
	}, nil
}

type JobLogs struct {
	Output string
	Errors string
}

func (w *Worker) QueryLogs(jobId string) (*JobLogs, error) {
	item, err := w.store.Get(jobId)
	if err != nil {
		return nil, err
	}

	job := item.(Job) // need casting as we don't have generics yet...
	stdout, stderr := job.Logs()

	return &JobLogs{
		Output: stdout,
		Errors: stderr,
	}, nil
}
