package worker

import (
	"fmt"
	"log"
)

type worker struct {
	queue *queue
}

func MakeWorker() *worker {
	defaultQueueSize := 100
	return &worker{queue: makeQueue(defaultQueueSize)}
}

func (w *worker) Dispatch(name string, args ...string) (string, error) {
	log.Println("dispatching new job for", name)

	job := makeJob(name, args...)
	err := w.queue.put(job.id, job)
	if err != nil {
		log.Println("unable to store command", err.Error())
		return "", err
	}

	// mock job start (this should be done by a separate goroutine)
	err = job.start()
	if err != nil {
		log.Println("unable to dispatch command", err.Error())
		return "", err
	}

	return job.id, nil
}

func (w *worker) Stop(jobId string) error {
	job, err := w.queue.get(jobId)
	if err != nil {
		// this could be sensitive, maybe log, maybe don't ...
		log.Println("unable to retrieve job", jobId, err.Error())
		return err
	}

	err = job.stop()
	if err != nil {
		log.Println("unable to stop job", jobId, err.Error())
		return err
	}

	return nil
}

type JobInfo struct {
	Id     string `json:id`
	Status string `json:status`
}

func (w *worker) QueryInfo(jobId string) (*JobInfo, error) {
	job, err := w.queue.get(jobId)
	if err != nil {
		return nil, err
	}

	return &JobInfo{
		Id:     fmt.Sprint(job.id),
		Status: string(job.state.getStatus()),
	}, nil
}

type JobLogs struct {
	Output string
	Errors string
}

func (w *worker) QueryLogs(jobId string) (*JobLogs, error) {
	job, err := w.queue.get(jobId)
	if err != nil {
		return nil, err
	}

	// pull out the logs from the job, but do not expose their pointers
	logs := &JobLogs{
		Output: job.logs.getOutput(),
		Errors: job.logs.getErrors(),
	}

	return logs, nil
}
