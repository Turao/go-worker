package server

import (
	"context"

	"github.com/turao/go-worker/pkg/worker"
)

type Service interface {
	Dispatch(ctx context.Context, name string, args ...string) (string, error)
	Stop(ctx context.Context, jobId string) error
	QueryInfo(ctx context.Context, jobId string) (*JobInfo, error)
}

type JobInfo struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	ExitCode int    `json:"exitCode"`
	Output   string `json:"output"`
	Errors   string `json:"errors"`
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

func (s workerservice) QueryInfo(ctx context.Context, jobId string) (*JobInfo, error) {
	info, err := s.worker.QueryInfo(jobId)
	if err != nil {
		return nil, err
	}

	return &JobInfo{
		ID:       info.ID,
		Status:   info.Status,
		ExitCode: info.ExitCode,
		Output:   info.Output,
		Errors:   info.Errors,
	}, nil
}
