package server

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	v1 "github.com/turao/go-worker/api/v1"
)

func makeDispatchEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(v1.DispatchRequest)
		id, err := service.Dispatch(ctx, req.Name, req.Args...)
		if err != nil {
			return nil, err
		}

		return v1.DispatchResponse{ID: id}, nil
	}
}

func makeStopEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(v1.StopRequest)
		err := service.Stop(ctx, req.ID)
		return nil, err
	}
}

func makeQueryInfoEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(v1.QueryInfoRequest)
		info, err := service.QueryInfo(ctx, req.ID)
		if err != nil {
			return nil, err
		}

		return v1.QueryInfoResponse{
			JobInfo: *info,
		}, nil
	}
}
