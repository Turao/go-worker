package job

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJobStartOnce(t *testing.T) {
	job := New("ls")

	err := job.Start()
	assert.Nil(t, err)
}

func TestJobStartWhileStarting(t *testing.T) {
	job := New("ls")
	job.Start()

	err := job.Start()
	assert.Equal(t, ErrStarting, err)
}

func TestJobStartAlreadyStarted(t *testing.T) {
	job := New("sleep", "3")
	job.Start()
	time.Sleep(1 * time.Second)

	err := job.Start()
	assert.Equal(t, ErrAlreadyStarted, err)
}

func TestJobStopOnce(t *testing.T) {
	job := New("sleep", "5")
	job.Start()
	time.Sleep(1 * time.Second)

	err := job.Stop()
	assert.Nil(t, err)
}

func TestJobStopWhileStopping(t *testing.T) {
	job := New("sleep", "5")
	job.Start()
	time.Sleep(1 * time.Second)
	job.Stop()

	// kill signal has been sent,
	// os.Process has not terminated yet
	err := job.Stop()
	assert.Equal(t, ErrStopping, err)
}

func TestJobStopAlreadyFinished(t *testing.T) {
	job := New("ls")
	job.Start()
	time.Sleep(1 * time.Second)

	// os.Process should have been terminated
	// by this point
	err := job.Stop()
	assert.Equal(t, ErrAlreadyFinished, err)
}

func TestJobStopOnJobNotStarted(t *testing.T) {
	job := New("ls")

	err := job.Stop()
	assert.Equal(t, ErrNotStarted, err)
}

func TestJobWaitOnJobNotStarted(t *testing.T) {
	job := New("sleep", "3")

	err := job.waitUntilCompleted()
	assert.Equal(t, ErrNotStarted, err)
}

func TestJobWaitOnJobAlreadyFinished(t *testing.T) {
	job := New("ls")
	job.Start()
	time.Sleep(1 * time.Second)
	// os.Process should have been completed
	// by this point
	err := job.waitUntilCompleted()
	assert.Equal(t, ErrAlreadyFinished, err)
}

func TestJobWaitWhileAlreadyWaiting(t *testing.T) {
	job := New("sleep", "3")
	job.Start()
	time.Sleep(1 * time.Second)
	// os.Process should have been started
	// by this point
	go job.waitUntilCompleted()

	err := job.waitUntilCompleted()
	assert.Equal(t, ErrAlreadyWaiting, err)
}
