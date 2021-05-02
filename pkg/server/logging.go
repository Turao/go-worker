package server

import (
	"context"
	"log"

	v1 "github.com/turao/go-worker/api/v1"
)

type loggingMiddleware struct {
	next Service
}

func (l loggingMiddleware) Dispatch(ctx context.Context, name string, args ...string) (v1.JobID, error) {
	log.Println("[Dispatch]", "called")
	defer log.Println("[Dispatch]", "finished")
	return l.next.Dispatch(ctx, name, args...)
}

func (l loggingMiddleware) Stop(ctx context.Context, jobId v1.JobID) error {
	log.Println("[Stop]", "called")
	defer log.Println("[Stop]", "finished")
	return l.next.Stop(ctx, jobId)
}

func (l loggingMiddleware) QueryInfo(ctx context.Context, jobId v1.JobID) (*v1.JobInfo, error) {
	log.Println("[QueryInfo]", "called")
	defer log.Println("[QueryInfo]", "finished")
	return l.next.QueryInfo(ctx, jobId)
}
