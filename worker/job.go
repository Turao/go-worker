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

	cmd                 *exec.Cmd
	onProcessStart      chan bool
	onProcessCompletion chan bool
	onProcessStop       chan bool
}

func makeJob(stdout io.Reader, stderr io.Reader, cmd *exec.Cmd) *job {
	return &job{
		id:     -1,
		stdout: stdout,
		stderr: stderr,
		state:  &state{mx: &sync.Mutex{}, status: SCHEDULED},

		cmd:                 cmd,
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

func (j *job) start() error {
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

func (j *job) stop() error {
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
		log.Println(err.Error())
		return err
	}

	go j.waitUntilCompleted()
	return nil
}

func (j *job) onProcessCompleted() error {
	log.Println("process completed")

	err := j.state.completed()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (j *job) onProcessStopped() error {
	log.Println("process stopped")

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

func (j *job) hasStarted() bool {
	j.state.mx.Lock()
	defer j.state.mx.Unlock()
	return j.state.status == RUNNING
}

func (j *job) hasFinished() bool {
	j.state.mx.Lock()
	defer j.state.mx.Unlock()
	return j.state.status == COMPLETED || j.state.status == STOPPED
}

var ErrAlreadyStarted error = errors.New("job has already been started")
var ErrAlreadyFinished error = errors.New("job has already finished (either completed or stopped)")
var ErrNotStarted error = errors.New("job has not started yet")
var ErrNotFinished error = errors.New("job has not finished yet")

// change state status
var ErrNotScheduled error = errors.New("job is not scheduled")
var ErrNotRunning error = errors.New("job is not running")

func (s *state) running() error {
	s.mx.Lock()
	defer s.mx.Unlock()
	if s.status != SCHEDULED {
		return ErrNotScheduled
	}
	s.status = RUNNING
	return nil
}

func (s *state) completed() error {
	s.mx.Lock()
	defer s.mx.Unlock()
	if s.status != RUNNING {
		return ErrNotRunning
	}
	s.status = COMPLETED
	return nil
}

func (s *state) stopped() error {
	s.mx.Lock()
	defer s.mx.Unlock()
	if s.status != RUNNING {
		return ErrNotRunning
	}
	s.status = STOPPED
	return nil
}
