package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

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

func decodeStopResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var response struct {
		Error error `json:"error,omitempty"`
	}

	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}

	return response, nil
}

func decodeQueryResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var response struct {
		ID       string `json:"id"`
		Status   string `json:"status"`
		ExitCode int    `json:"exitCode,omitempty"`
		Output   string `json:"output,omitempty"`
		Errors   string `json:"errors,omitempty"`
	}

	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
