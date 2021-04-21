package apiserver

import (
	"context"

	"github.com/turao/kami-go/worker"
)

type Service interface {
	Dispatch(ctx context.Context, name string, args ...string) (string, error)
	Stop(ctx context.Context, jobId string) error
	QueryInfo(ctx context.Context, jobId string) (*JobInfo, error)
	QueryLogs(ctx context.Context, jobId string) (*JobLogs, error)
}

type JobInfo struct {
	Id       string `json:"id"`
	Status   string `json:"status"`
	ExitCode int    `json:"exitCode"`
}

type JobLogs struct {
	Output string `json:"output"`
	Errors string `json:"errors"`
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
		Id:       info.Id,
		Status:   info.Status,
		ExitCode: info.ExitCode,
	}, nil
}

func (s workerservice) QueryLogs(ctx context.Context, jobId string) (*JobLogs, error) {
	logs, err := s.worker.QueryLogs(jobId)
	if err != nil {
		return nil, err
	}

	return &JobLogs{
		Output: logs.Output,
		Errors: logs.Errors,
	}, nil
}
