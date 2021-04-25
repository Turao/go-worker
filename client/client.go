package client

import (
	"context"
	"net/url"
	"time"

	"github.com/go-kit/kit/endpoint"
)

// client wraps an http client and add a bunch of stuff to it
type client struct {
	// dependencies
	// auth
	// server
	// logger (?)

	start endpoint.Endpoint
	stop  endpoint.Endpoint
	query endpoint.Endpoint
}

func New(url *url.URL) *client {
	return &client{
		start: makeDispatchEndpoint(url),
		stop:  makeStopEndpoint(url),
		query: makeQueryEndpoint(url),
	}
}

func (c *client) Start(name string, args ...string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	request := struct {
		Name string   `json:"name"`
		Args []string `json:"args"`
	}{
		Name: name,
		Args: args,
	}

	res, err := c.start(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *client) Stop(jobID string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	request := struct {
		JobID string `json:"id"`
	}{
		JobID: jobID,
	}

	res, err := c.stop(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *client) Query(jobID string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	request := struct {
		JobID string `json:"id"`
	}{
		JobID: jobID,
	}

	res, err := c.query(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}
