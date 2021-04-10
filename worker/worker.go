package worker

import (
	"fmt"
	"io"
	"log"
	"os/exec"
)

type worker struct {
	queue *queue
}

func MakeWorker() *worker {
	defaultQueueSize := 100
	return &worker{queue: makeQueue(defaultQueueSize)}
}

func (w *worker) Dispatch(cmd exec.Cmd) (int, error) {
	var stdout io.Reader
	var stderr io.Reader
	job := makeJob(stdout, stderr)
	_, err := w.queue.addJob(job)
	if err != nil {
		log.Println("unable to dispatch command", err.Error())
		return -1, err
	}

	// mock job start (this should be done by a separate goroutine)
	err = job.start()
	if err != nil {
		log.Println("job started!")
	}

	return job.id, nil
}

func (w *worker) Stop(jobId int) error {
	job, err := w.queue.getJob(jobId)
	if err != nil {
		// this could be sensitive, maybe log, maybe don't ...
		log.Println("job does not exist", jobId, err.Error())
		return err
	}

	err = job.stop()
	if err != nil {
		log.Println("unable to stop job", jobId, err.Error())
		return err
	}

	err = w.queue.removeJob(jobId)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

type jobInfo struct {
	id     string
	status string
}

func (w *worker) QueryInfo(jobId int) (*jobInfo, error) {
	job, err := w.queue.getJob(jobId)
	if err != nil {
		return nil, err
	}

	return &jobInfo{
		id:     fmt.Sprint(job.id),
		status: string(job.state.status),
	}, nil
}

type jobLogs struct {
	stdout io.Reader
	stderr io.Reader
}

func (w *worker) QueryLogs(jobId int) (*jobLogs, error) {
	job, err := w.queue.getJob(jobId)
	if err != nil {
		return nil, err
	}
	return &jobLogs{
		stdout: job.stdout,
		stderr: job.stderr,
	}, nil
}
