package worker

import (
	"errors"
	"io"
	"log"
	"os/exec"
	"sync"
)

type job struct {
	id     int
	stdout io.Reader
	stderr io.Reader
	state  *state

	cmd exec.Cmd
}

func makeJob(stdout io.Reader, stderr io.Reader, cmd exec.Cmd) *job {
	return &job{
		id:     -1,
		stdout: stdout,
		stderr: stderr,
		state:  &state{mx: &sync.Mutex{}, status: SCHEDULED},
		cmd:    cmd,
	}
}

func (j *job) start() error {
	log.Println("starting job...")
	err := j.state.running()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (j *job) complete() error {
	log.Println("completing job...")
	err := j.state.completed()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (j *job) stop() error {
	log.Println("stopping job...")
	err := j.state.stopped()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// job statuses
type state struct {
	mx     *sync.Mutex
	status status
}

type status string

const (
	SCHEDULED = "scheduled"
	RUNNING   = "running"
	COMPLETED = "completed"
	STOPPED   = "stopped"
)

func (s *state) scheduled() error {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.status = SCHEDULED
	return nil
}

func (s *state) running() error {
	s.mx.Lock()
	defer s.mx.Unlock()
	if s.status != SCHEDULED {
		return errors.New("cannot change process to running state")
	}
	s.status = RUNNING
	return nil
}

func (s *state) completed() error {
	s.mx.Lock()
	defer s.mx.Unlock()
	if s.status != RUNNING {
		return errors.New("process is not running thus cannot be completed")
	}
	s.status = COMPLETED
	return nil
}

func (s *state) stopped() error {
	s.mx.Lock()
	defer s.mx.Unlock()
	if s.status != RUNNING {
		return errors.New("process is not running thus cannot be stopped")
	}
	s.status = STOPPED
	return nil
}
