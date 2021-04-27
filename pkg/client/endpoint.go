package client

import (
	"context"
	"fmt"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	httpkit "github.com/go-kit/kit/transport/http"

	v1 "github.com/turao/go-worker/api/v1"
)

func makeDispatchEndpoint(url *url.URL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(v1.DispatchRequest)
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

func makeStopEndpoint(url *url.URL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(v1.StopRequest)
		url.Path = fmt.Sprintf("/job/%s/stop", req.ID)

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

func makeQueryEndpoint(url *url.URL) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(v1.QueryInfoRequest)
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
