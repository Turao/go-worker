package server

import (
	"context"

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
		req := request.(StopRequest)
		err := service.Stop(ctx, req.ID)
		return nil, err
	}
}

type QueryInfoRequest struct {
	ID string `json:"id"`
}
type QueryInfoResponse struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	ExitCode int    `json:"exitCode"`
	Output   string `json:"output"`
	Errors   string `json:"errors"`
}

func makeQueryInfoEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(QueryInfoRequest)
		info, err := service.QueryInfo(ctx, req.ID)
		if err != nil {
			return nil, err
		}

		return QueryInfoResponse{
			ID:       info.ID,
			Status:   info.Status,
			ExitCode: info.ExitCode,
			Output:   info.Output,
			Errors:   info.Errors,
		}, nil
	}
}
