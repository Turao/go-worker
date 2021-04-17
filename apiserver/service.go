package apiserver

import "github.com/turao/kami-go/worker"

type Service interface {
	Dispatch(name string, args ...string) (string, error)
	Stop(jobId string) error
	QueryInfo(jobId string) (JobInfo, error)
	QueryLogs(jobId string) (JobLogs, error)
}

type JobInfo struct {
	Id       string `json:"id"`
	Status   string `json:"status"`
	ExitCode *int   `json:"exitCode,omitempty"`
}

type JobLogs struct {
	Output string `json:"output"`
	Errors string `json:"errors"`
}

type workerservice struct {
	worker *worker.Worker
}

func NewWorkerService() *workerservice {
	return &workerservice{
		worker: worker.NewWorker(),
	}
}

func (s *workerservice) Dispatch(name string, args ...string) (string, error) {
	return s.worker.Dispatch(name, args...)
}

func (s *workerservice) Stop(jobId string) error {
	return s.worker.Stop(jobId)
}

func (s *workerservice) QueryInfo(jobId string) (*JobInfo, error) {
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

func (s *workerservice) QueryLogs(jobId string) (*JobLogs, error) {
	logs, err := s.worker.QueryLogs(jobId)
	if err != nil {
		return nil, err
	}

	return &JobLogs{
		Output: logs.Output,
		Errors: logs.Errors,
	}, nil
}
