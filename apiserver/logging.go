package apiserver

import (
	"context"
	"log"
)

type loggingMiddleware struct {
	Service
}

func (l *loggingMiddleware) Dispatch(ctx context.Context, name string, args ...string) (string, error) {
	log.Println("[Dispatch]", "called")
	defer log.Println("[Dispatch]", "finished")
	return l.Service.Dispatch(ctx, name, args...)
}

func (l *loggingMiddleware) Stop(ctx context.Context, jobId string) error {
	log.Println("[Stop]", "called")
	defer log.Println("[Stop]", "finished")
	return l.Service.Stop(ctx, jobId)
}

func (l *loggingMiddleware) QueryInfo(ctx context.Context, jobId string) (JobInfo, error) {
	log.Println("[QueryInfo]", "called")
	defer log.Println("[QueryInfo]", "finished")
	return l.Service.QueryInfo(ctx, jobId)
}

func (l *loggingMiddleware) QueryLogs(ctx context.Context, jobId string) (JobLogs, error) {
	log.Println("[QueryLogs]", "called")
	defer log.Println("[QueryLogs]", "finished")
	return l.Service.QueryLogs(ctx, jobId)
}
