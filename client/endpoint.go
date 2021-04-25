package client

import (
	"context"
	"encoding/json"
	"net/http"
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

func decodeDispatchResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var response struct {
		ID string `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
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

func decodeStopResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var response struct {
		ID string `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
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

func decodeQueryResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var response struct {
		ID       string `json:"id"`
		Status   string `json:"status"`
		ExitCode int    `json:"exitCode"`
		Output   string `json:"output"`
		Errors   string `json:"errors"`
	}

	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
