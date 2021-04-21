package apiserver

import (
	"context"
	"log"

	"github.com/go-kit/kit/endpoint"
)

type DispatchRequest struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}

type DispatchResponse struct {
	ID string `json:"id"`
}

func makeDispatchEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		log.Println("[dispatch endpoint]", "called")
		req := request.(DispatchRequest)
		id, err := service.Dispatch(ctx, req.Name, req.Args...)
		if err != nil {
			return nil, err
		}

		return DispatchResponse{ID: id}, nil
	}
}

type StopRequest struct {
	ID string `json:"id"`
}

type StopResponse struct{}

func makeStopEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		log.Println("[stop endpoint]", "called")
		req := request.(StopRequest)
		err := service.Stop(ctx, req.ID)
		return nil, err
	}
}

type QueryInfoRequest struct {
	ID string
}
type QueryInfoResponse struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	ExitCode int    `json:"exitCode"`
}

func makeQueryInfoEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		log.Println("[query info endpoint]", "called")
		req := request.(QueryInfoRequest)
		info, err := service.QueryInfo(ctx, req.ID)
		if err != nil {
			return nil, err
		}

		return QueryInfoResponse{
			ID:       info.Id,
			Status:   info.Status,
			ExitCode: info.ExitCode,
		}, nil
	}
}

type QueryLogsRequest struct {
	ID string `json:"id"`
}

type QueryLogsResponse struct {
	Output string `json:"output,omitempty"`
	Errors string `json:"errors,omitempty"`
}

func makeQueryLogsEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		log.Println("[query logs endpoint]", "called")
		req := request.(QueryLogsRequest)
		logs, err := service.QueryLogs(ctx, req.ID)
		if err != nil {
			return nil, err
		}

		return QueryLogsResponse{
			Output: logs.Output,
			Errors: logs.Errors,
		}, nil
	}
}
