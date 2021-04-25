package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	httpkit "github.com/go-kit/kit/transport/http"
)

func makeStartEndpoint(url *url.URL) endpoint.Endpoint {
	var opts []httpkit.ClientOption

	return httpkit.NewClient(
		"POST",
		url,
		encodeRequest,
		decodeResponse,
		opts...,
	).Endpoint()
}

func makeStopEndpoint(url *url.URL) endpoint.Endpoint {
	var opts []httpkit.ClientOption

	return httpkit.NewClient(
		"POST",
		url,
		encodeRequest,
		decodeResponse,
		opts...,
	).Endpoint()
}

func makeQueryEndpoint(url *url.URL) endpoint.Endpoint {
	var opts []httpkit.ClientOption

	return httpkit.NewClient(
		"GET",
		url,
		encodeRequest,
		decodeResponse,
		opts...,
	).Endpoint()
}

func encodeRequest(ctx context.Context, r *http.Request, body interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return err
	}
	r.Body = io.NopCloser(&buf)
	return nil
}

func decodeResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	return r, nil
}
