package worker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJobStartOnce(t *testing.T) {
	job := NewJob("ls", "-lah")

	err := job.Start()
	assert.Nil(t, err)
}

func TestJobStartTwice(t *testing.T) {
	job := NewJob("ls", "-lah")
	job.Start()

	err := job.Start()
	assert.NotNil(t, err)
}

func TestJobStopOnce(t *testing.T) {
	job := NewJob("sleep", "10")
	job.Start()
	time.Sleep(1 * time.Second)

	err := job.Stop()
	assert.Nil(t, err)
}

func TestJobStopTwice(t *testing.T) {
	job := NewJob("sleep", "10")
	job.Start()
	time.Sleep(2 * time.Second)
	job.Stop()
	time.Sleep(2 * time.Second)

	// os.Process should have been terminated
	// by this point
	err := job.Stop()
	assert.Equal(t, ErrAlreadyFinished, err)
}

func TestJobStopWhileStopping(t *testing.T) {
	job := NewJob("sleep", "10")
	job.Start()
	time.Sleep(2 * time.Second)
	job.Stop()

	// kill signal has been sent,
	// os.Process has not terminated yet
	err := job.Stop()
	assert.Equal(t, ErrStopping, err)
}

func TestJobStopBeforeStart(t *testing.T) {
	job := NewJob("ls", "-lah")

	err := job.Stop()
	assert.Equal(t, ErrNotStarted, err)
}
