package apiserver

import (
	"context"
	"encoding/json"
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

	r.Handle("/job", dispatchHandler).Methods("POST")
	r.Handle("/job", queryInfoHandler).Methods("GET")

	return r
}

func decodeDispatchRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body DispatchRequest // todo: go-kit proposes having a temporary struct for decoupling, but I find this overkill
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func decodeQueryInfoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body QueryInfoRequest // todo: go-kit proposes having a temporary struct for decoupling, but I find this overkill
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
