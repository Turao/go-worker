package server

import (
	"context"

	v1 "github.com/turao/go-worker/api/v1"
	"github.com/turao/go-worker/pkg/worker"
)

type Service interface {
	Dispatch(ctx context.Context, name string, args ...string) (v1.JobID, error)
	Stop(ctx context.Context, jobId v1.JobID) error
	QueryInfo(ctx context.Context, jobId v1.JobID) (*v1.JobInfo, error)
}

type workerservice struct {
	worker *worker.Worker
}

func newWorkerService() Service {
	return workerservice{
		worker: worker.NewWorker(),
	}
}

func (s workerservice) Dispatch(ctx context.Context, name string, args ...string) (v1.JobID, error) {
	return s.worker.Dispatch(name, args...)
}

func (s workerservice) Stop(ctx context.Context, jobId v1.JobID) error {
	return s.worker.Stop(jobId)
}

func (s workerservice) QueryInfo(ctx context.Context, jobId v1.JobID) (*v1.JobInfo, error) {
	info, err := s.worker.QueryInfo(jobId)
	if err != nil {
		return nil, err
	}

	return info, nil
}
