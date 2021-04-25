package client

import (
	"net/url"

	"github.com/go-kit/kit/endpoint"
	httpkit "github.com/go-kit/kit/transport/http"
)

func makeDispatchEndpoint(url *url.URL) endpoint.Endpoint {
	var opts []httpkit.ClientOption

	return httpkit.NewClient(
		"POST",
		url,
		httpkit.EncodeJSONRequest,
		decodeDispatchResponse,
		opts...,
	).Endpoint()
}

func makeStopEndpoint(url *url.URL) endpoint.Endpoint {
	var opts []httpkit.ClientOption

	return httpkit.NewClient(
		"POST",
		url,
		httpkit.EncodeJSONRequest,
		decodeStopResponse,
		opts...,
	).Endpoint()
}

func makeQueryEndpoint(url *url.URL) endpoint.Endpoint {
	var opts []httpkit.ClientOption

	return httpkit.NewClient(
		"GET",
		url,
		httpkit.EncodeJSONRequest,
		decodeQueryResponse,
		opts...,
	).Endpoint()
}
