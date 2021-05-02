package worker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/turao/go-worker/pkg/job"
	"github.com/turao/go-worker/pkg/storage"
)

func givenDispatchedJob(t *testing.T, worker *Worker) string {
	jobID, err := worker.Dispatch("sleep", "10")
	if err != nil {
		t.Fatal("unable to dispatch mock job")
	}
	time.Sleep(2 * time.Second)
	return jobID
}

func TestDispatch(t *testing.T) {
	worker := NewWorker()

	jobID, err := worker.Dispatch("ls", "-lah")
	assert.Nil(t, err)
	assert.NotNil(t, jobID)
}

func TestStopOnce(t *testing.T) {
	worker := NewWorker()
	jobID := givenDispatchedJob(t, worker)

	err := worker.Stop(jobID)
	assert.Nil(t, err)
}

func TestStopAlreadyStoppedJob(t *testing.T) {
	worker := NewWorker()
	jobID := givenDispatchedJob(t, worker)
	worker.Stop(jobID)
	time.Sleep(1 * time.Second)

	// os.Process should have been terminated
	// by this point
	err := worker.Stop(jobID)
	assert.Equal(t, job.ErrAlreadyFinished, err)
}

func TestStopWhileJobIsStopping(t *testing.T) {
	worker := NewWorker()
	jobID := givenDispatchedJob(t, worker)
	worker.Stop(jobID)

	// kill signal has been sent,
	// os.Process has not terminated yet
	err := worker.Stop(jobID)
	assert.Equal(t, job.ErrStopping, err)
}

func TestStopNotExists(t *testing.T) {
	worker := NewWorker()
	givenDispatchedJob(t, worker)

	err := worker.Stop("other")
	assert.Equal(t, storage.ErrNotExists, err)
}
