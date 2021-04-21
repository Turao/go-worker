package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
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

	queryInfoHandler := kithttp.NewServer(
		makeQueryInfoEndpoint(svc),
		decodeQueryInfoRequest,
		encodeResponse,
		opts...,
	)

	queryLogsHandler := kithttp.NewServer(
		makeQueryLogsEndpoint(svc),
		decodeQueryLogsRequest,
		encodeResponse,
		opts...,
	)

	r.Handle("/job", dispatchHandler).Methods("POST")

	r.Handle("/job/{id}/info", queryInfoHandler).Methods("GET")
	r.Handle("/job/{id}/logs", queryLogsHandler).Methods("GET")

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

	return DispatchRequest{
		Name: body.Name,
		Args: body.Args,
	}, nil
}

func decodeQueryInfoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	jobID, found := params["id"]
	if !found {
		return nil, errors.New("unable to find id in URL params")
	}

	return QueryInfoRequest{
		ID: jobID,
	}, nil
}

func decodeQueryLogsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	jobID, found := params["id"]
	if !found {
		return nil, errors.New("unable to find id in URL params")
	}

	return QueryLogsRequest{
		ID: jobID,
	}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
