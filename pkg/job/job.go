package job

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/google/uuid"
)

type job struct {
	id string

	state *state
	logs  *logs

	cmd *exec.Cmd

	onProcessStart      signalOnce
	onProcessCompletion signalOnce
	onProcessStop       signalOnce
}

// signalOnce provides a way to cleanly send messages/close channels while having concurrent event handler calls (e.g. stop)
type signalOnce struct {
	once sync.Once
	ch   chan interface{}
}

type JobInfo struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	ExitCode int    `json:"exitCode"`
	Output   string `json:"output"`
	Errors   string `json:"errors"`
}

var ErrNotStarted error = errors.New("job has not started yet")
var ErrStarting error = errors.New("job is starting")
var ErrAlreadyStarted error = errors.New("job has already been started")
var ErrAlreadyWaiting error = errors.New("job is already waiting for process to complete")
var ErrStopping error = errors.New("job is stopping")
var ErrAlreadyFinished error = errors.New("job has already finished (either completed or stopped)")

func New(name string, args ...string) *job {
	command := exec.Command(name, args...)

	logs := logs{
		stdout: &threadSafeBuffer{mx: sync.RWMutex{}, buf: bytes.Buffer{}},
		stderr: &threadSafeBuffer{mx: sync.RWMutex{}, buf: bytes.Buffer{}},
	}
	command.Stdout = io.MultiWriter(os.Stdout, logs.stdout)
	command.Stderr = io.MultiWriter(os.Stderr, logs.stderr)

	return &job{
		id:    uuid.New().String(),
		state: &state{mx: &sync.RWMutex{}, status: SCHEDULED, exitCode: UnknownOrTerminated},
		logs:  &logs,

		cmd:                 command,
		onProcessStart:      signalOnce{once: sync.Once{}, ch: make(chan interface{}, 1)},
		onProcessCompletion: signalOnce{once: sync.Once{}, ch: make(chan interface{}, 1)},
		onProcessStop:       signalOnce{once: sync.Once{}, ch: make(chan interface{}, 1)},
	}
}

func (j *job) ID() string {
	return j.id
}

func (j *job) Start() error {
	if j.hasStarted() {
		return ErrAlreadyStarted
	}

	if j.hasFinished() {
		return ErrAlreadyFinished
	}

	err := ErrStarting
	j.onProcessStart.once.Do(func() {
		log.Println("starting job...")
		err = j.cmd.Start()
		if err != nil {
			return
		}
		go j.watch()

		j.onProcessStart.ch <- true
		close(j.onProcessStart.ch)
		err = nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (j *job) Stop() error {
	if !j.hasStarted() {
		return ErrNotStarted
	}

	if j.hasFinished() {
		return ErrAlreadyFinished
	}

	err := ErrStopping
	j.onProcessStop.once.Do(func() {
		log.Println("stopping job...")
		err = j.cmd.Process.Kill()
		if err != nil {
			return
		}

		j.onProcessStop.ch <- true
		close(j.onProcessStop.ch)
		err = nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (j *job) Info() *JobInfo {
	return &JobInfo{
		ID:       j.id,
		Status:   string(j.state.Status()),
		ExitCode: j.state.ExitCode(),
		Output:   j.logs.Output(),
		Errors:   j.logs.Errors(),
	}
}

func (j *job) waitUntilCompleted() error {
	if !j.hasStarted() {
		return ErrNotStarted
	}

	if j.hasFinished() {
		return ErrAlreadyFinished
	}

	err := ErrAlreadyWaiting
	j.onProcessCompletion.once.Do(func() {
		log.Println("waiting for job process to complete")
		// wait releases all process resources (and creates a new os.ProcessState on return)
		// we need to propagate the process state value somehow...
		ps, err := j.cmd.Process.Wait()
		if err != nil {
			return
		}
		log.Println("process completed, signaling job")

		j.onProcessCompletion.ch <- ps
		close(j.onProcessCompletion.ch)
		err = nil
	})

	if err != nil {
		return err
	}

	return nil
}

// watch reacts on process signals
func (j *job) watch() error {
	go func() {
		<-j.onProcessStart.ch
		j.onProcessStarted()
	}()

	select {
	case ps := <-j.onProcessCompletion.ch:
		return j.onProcessCompleted(ps.(*os.ProcessState))

	case <-j.onProcessStop.ch:
		return j.onProcessStopped()
	}
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

func (j *job) onProcessCompleted(ps *os.ProcessState) error {
	log.Println("process completed with status", ps)

	exitCode := ps.ExitCode()
	err := j.state.completed(exitCode)
	if err != nil {
		return err
	}
	return nil
}

func (j *job) onProcessStopped() error {
	log.Println("process stopped")

	exitCode := j.cmd.ProcessState.ExitCode()
	err := j.state.stopped(exitCode)
	if err != nil {
		return err
	}
	return nil
}

func (j *job) hasStarted() bool {
	return j.state.Status() != SCHEDULED
}

func (j *job) hasFinished() bool {
	return j.state.Status() == COMPLETED || j.state.Status() == STOPPED
}
