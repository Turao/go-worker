package job

import (
	"errors"
	"sync"
)

// state provides thread-safe operations for reading/updating the job current state
type state struct {
	mx       *sync.RWMutex
	status   status
	exitCode int
}

const UnknownOrTerminated = -1

func (s *state) Status() status {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.status
}

func (s *state) ExitCode() int {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.exitCode
}

type status string

const (
	SCHEDULED = "scheduled"
	RUNNING   = "running"
	COMPLETED = "completed"
	STOPPED   = "stopped"
)

func (j *job) hasStarted() bool {
	j.state.mx.RLock()
	defer j.state.mx.RUnlock()
	return j.state.status != SCHEDULED
}

func (j *job) hasFinished() bool {
	j.state.mx.RLock()
	defer j.state.mx.RUnlock()
	return j.state.status == COMPLETED || j.state.status == STOPPED
}

func (s *state) running() error {
	s.mx.Lock()
	defer s.mx.Unlock()
	if s.status != SCHEDULED {
		return ErrNotScheduled
	}

	s.status = RUNNING
	return nil
}

func (s *state) completed(exitCode int) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	if s.status != RUNNING {
		return ErrNotRunning
	}

	s.status = COMPLETED
	s.exitCode = exitCode
	return nil
}

func (s *state) stopped(exitCode int) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	if s.status != RUNNING {
		return ErrNotRunning
	}

	s.status = STOPPED
	s.exitCode = exitCode
	return nil
}

var ErrNotStarted error = errors.New("job has not started yet")
var ErrStarting error = errors.New("job is starting")
var ErrAlreadyStarted error = errors.New("job has already been started")
var ErrAlreadyWaiting error = errors.New("job is already waiting for process to complete")
var ErrStopping error = errors.New("job is stopping")
var ErrAlreadyFinished error = errors.New("job has already finished (either completed or stopped)")

var ErrNotScheduled error = errors.New("job is not scheduled")
var ErrNotRunning error = errors.New("job is not running")
