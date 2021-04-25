package client

import (
	"context"
	"encoding/json"
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
		ID string `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
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
