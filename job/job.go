package job

import (
	"bytes"
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

	cmd                 *exec.Cmd
	onProcessStart      signalOnce
	onProcessCompletion signalOnce
	onProcessStop       signalOnce
}

// signalOnce provides a way to cleanly close channels while having concurrent event handler calls (e.g. stop)
type signalOnce struct {
	once sync.Once
	ch   chan bool
}

func NewJob(name string, args ...string) *job {
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
		onProcessStart:      signalOnce{once: sync.Once{}, ch: make(chan bool, 1)},
		onProcessCompletion: signalOnce{once: sync.Once{}, ch: make(chan bool, 1)},
		onProcessStop:       signalOnce{once: sync.Once{}, ch: make(chan bool, 1)},
	}
}

func (j *job) ID() string {
	return j.id
}

func (j *job) Status() string {
	return string(j.state.Status())
}

func (j *job) ExitCode() int {
	return j.state.ExitCode()
}

func (j *job) Logs() (string, string) {
	return j.logs.Output(), j.logs.Errors()
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

func (j *job) waitUntilCompleted() error {
	if !j.hasStarted() {
		return ErrNotStarted
	}

	if j.hasFinished() {
		return ErrAlreadyFinished
	}

	err := ErrAlreadyWaiting
	j.onProcessCompletion.once.Do(func() {
		log.Println("waiting for job process to finish")
		j.cmd.Process.Wait()
		log.Println("process completed, signaling app")

		j.onProcessCompletion.ch <- true
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
	case <-j.onProcessCompletion.ch:
		return j.onProcessCompleted()

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

func (j *job) onProcessCompleted() error {
	log.Println("process completed")

	exitCode := j.cmd.ProcessState.ExitCode()
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
