package server

import (
	"context"

	"github.com/turao/go-worker/pkg/job"
	"github.com/turao/go-worker/pkg/worker"
)

type Service interface {
	Dispatch(ctx context.Context, name string, args ...string) (string, error)
	Stop(ctx context.Context, jobId string) error
	QueryInfo(ctx context.Context, jobId string) (*job.JobInfo, error)
}

type workerservice struct {
	worker *worker.Worker
}

func newWorkerService() Service {
	return workerservice{
		worker: worker.NewWorker(),
	}
}

func (s workerservice) Dispatch(ctx context.Context, name string, args ...string) (string, error) {
	return s.worker.Dispatch(name, args...)
}

func (s workerservice) Stop(ctx context.Context, jobId string) error {
	return s.worker.Stop(jobId)
}

func (s workerservice) QueryInfo(ctx context.Context, jobId string) (*job.JobInfo, error) {
	info, err := s.worker.QueryInfo(jobId)
	if err != nil {
		return nil, err
	}

	return info, nil
}
