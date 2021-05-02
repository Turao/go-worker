package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	v1 "github.com/turao/go-worker/api/v1"
)

func makeHandler(svc Service) http.Handler {
	r := mux.NewRouter()

	var opts []kithttp.ServerOption
	dispatchHandler := kithttp.NewServer(
		makeDispatchEndpoint(svc),
		decodeDispatchRequest,
		encodeResponse,
		opts...,
	)

	stopHandler := kithttp.NewServer(
		makeStopEndpoint(svc),
		decodeStopRequest,
		encodeResponse,
		opts...,
	)

	queryInfoHandler := kithttp.NewServer(
		makeQueryInfoEndpoint(svc),
		decodeQueryInfoRequest,
		encodeResponse,
		opts...,
	)

	r.Handle("/job", dispatchHandler).Methods("POST")
	r.Handle("/job/{id}/stop", stopHandler).Methods("POST")

	r.Handle("/job/{id}/info", queryInfoHandler).Methods("GET")

	return r
}

func decodeDispatchRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Name string   `json:"name"`
		Args []string `json:"args"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	return v1.DispatchRequest{
		Name: body.Name,
		Args: body.Args,
	}, nil
}

func decodeStopRequest(_ context.Context, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	jobID, found := params["id"]
	if !found {
		return nil, errors.New("unable to find id in URL params")
	}

	return v1.StopRequest{
		ID: v1.JobID(jobID),
	}, nil
}

func decodeQueryInfoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	jobID, found := params["id"]
	if !found {
		return nil, errors.New("unable to find id in URL params")
	}

	return v1.QueryInfoRequest{
		ID: v1.JobID(jobID),
	}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
