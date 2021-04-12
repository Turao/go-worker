package worker

import (
	"bytes"
	"errors"
	"log"
	"os/exec"
	"sync"

	"github.com/google/uuid"
)

type job struct {
	id string

	state *state
	logs  *logs

	cmd                 *exec.Cmd
	onProcessStart      chan bool
	onProcessCompletion chan bool
	onProcessStop       chan bool
}

func NewJob(name string, args ...string) *job {
	command := exec.Command(name, args...)

	logs := logs{
		stdout: &threadSafeBuffer{mx: sync.RWMutex{}, buf: bytes.Buffer{}},
		stderr: &threadSafeBuffer{mx: sync.RWMutex{}, buf: bytes.Buffer{}},
	}
	command.Stdout = logs.stdout
	command.Stderr = logs.stderr

	return &job{
		id:    uuid.New().String(),
		state: &state{mx: &sync.RWMutex{}, status: SCHEDULED},
		logs:  &logs,

		cmd:                 command,
		onProcessStart:      make(chan bool, 1),
		onProcessCompletion: make(chan bool, 1),
		onProcessStop:       make(chan bool, 1),
	}
}

// watch reacts on process signals
func (j *job) watch() error {
	go func() {
		<-j.onProcessStart
		j.onProcessStarted()
	}()

	select {
	case <-j.onProcessCompletion:
		return j.onProcessCompleted()

	case <-j.onProcessStop:
		return j.onProcessStopped()
	}
}

func (j *job) Start() error {
	if j.hasStarted() {
		return ErrAlreadyStarted
	}

	if j.hasFinished() {
		return ErrAlreadyFinished
	}

	log.Println("starting job...")
	err := j.cmd.Start()
	if err != nil {
		return err
	}
	go j.watch()

	j.onProcessStart <- true
	close(j.onProcessStart)

	return nil
}

func (j *job) Stop() error {
	if !j.hasStarted() {
		return ErrNotStarted
	}

	if j.hasFinished() {
		return ErrAlreadyFinished
	}

	log.Println("stopping job...")
	err := j.cmd.Process.Kill()
	if err != nil {
		return err
	}

	j.onProcessStop <- true
	close(j.onProcessStop)

	return nil
}

func (j *job) waitUntilCompleted() {
	log.Println("waiting for job process to finish")
	j.cmd.Process.Wait()
	log.Println("process completed, signaling app")

	j.onProcessCompletion <- true
	close(j.onProcessCompletion)
}

func (j *job) onProcessStarted() error {
	log.Println("process started")

	err := j.state.running()
	if err != nil {
		return err
	}

	go j.waitUntilCompleted()
	return nil
}

func (j *job) onProcessCompleted() error {
	log.Println("process completed")

	exitCode := j.cmd.ProcessState.ExitCode()
	err := j.state.completed(&exitCode)
	if err != nil {
		return err
	}
	return nil
}

func (j *job) onProcessStopped() error {
	log.Println("process stopped")

	exitCode := j.cmd.ProcessState.ExitCode()
	err := j.state.stopped(&exitCode)
	if err != nil {
		return err
	}
	return nil
}

// state provides thread-safe operations for reading/updating the job current state
type state struct {
	mx       *sync.RWMutex
	status   status
	exitCode *int
}

func (s *state) Status() status {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.status
}

func (s *state) ExitCode() *int {
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

var ErrAlreadyStarted error = errors.New("job has already been started")
var ErrAlreadyFinished error = errors.New("job has already finished (either completed or stopped)")
var ErrNotStarted error = errors.New("job has not started yet")
var ErrNotFinished error = errors.New("job has not finished yet")

var ErrNotScheduled error = errors.New("job is not scheduled")
var ErrNotRunning error = errors.New("job is not running")
var ErrNoExitCode error = errors.New("job terminated without exit code")

func (s *state) running() error {
	s.mx.Lock()
	defer s.mx.Unlock()
	if s.status != SCHEDULED {
		return ErrNotScheduled
	}

	s.status = RUNNING
	return nil
}

func (s *state) completed(exitCode *int) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	if s.status != RUNNING {
		return ErrNotRunning
	}

	if exitCode == nil {
		return ErrNoExitCode
	}

	s.status = COMPLETED
	s.exitCode = exitCode
	return nil
}

func (s *state) stopped(exitCode *int) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	if s.status != RUNNING {
		return ErrNotRunning
	}

	if exitCode == nil {
		return ErrNoExitCode
	}

	s.status = STOPPED
	s.exitCode = exitCode
	return nil
}

// logs provides thread-safe operations for reading job logs
type logs struct {
	stdout *threadSafeBuffer
	stderr *threadSafeBuffer
}

func (l *logs) Output() string {
	return l.stdout.String()
}

func (l *logs) Errors() string {
	return l.stderr.String()
}
