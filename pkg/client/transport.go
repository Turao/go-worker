package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	v1 "github.com/turao/go-worker/api/v1"
)

func decodeDispatchResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var response v1.DispatchResponse

	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func decodeStopResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var response v1.StopResponse

	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}

	return response, nil
}

func decodeQueryResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var response v1.QueryInfoResponse

	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
