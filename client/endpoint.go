package client

import (
	"context"
	"fmt"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	httpkit "github.com/go-kit/kit/transport/http"
)

type DispatchRequest struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}

func makeDispatchEndpoint(url *url.URL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DispatchRequest)

		url.Path = "/job"
		var opts []httpkit.ClientOption

		return httpkit.NewClient(
			"POST",
			url,
			httpkit.EncodeJSONRequest,
			decodeDispatchResponse,
			opts...,
		).Endpoint()(ctx, req)
	}
}

type StopRequest struct {
	ID string `json:"id"`
}

func makeStopEndpoint(url *url.URL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(StopRequest)
		url.Path = fmt.Sprintf("/job/%s", req.ID)

		var opts []httpkit.ClientOption

		return httpkit.NewClient(
			"POST",
			url,
			httpkit.EncodeJSONRequest,
			decodeStopResponse,
			opts...,
		).Endpoint()(ctx, req)
	}

}

type QueryRequest struct {
	ID string `json:"id"`
}

func makeQueryEndpoint(url *url.URL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(QueryRequest)
		url.Path = fmt.Sprintf("/job/%s/info", req.ID)

		var opts []httpkit.ClientOption

		return httpkit.NewClient(
			"GET",
			url,
			httpkit.EncodeJSONRequest,
			decodeQueryResponse,
			opts...,
		).Endpoint()(ctx, req)
	}
}
