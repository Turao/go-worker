package client

import (
	"context"
	"net/url"
	"time"

	"github.com/go-kit/kit/endpoint"
	v1 "github.com/turao/go-worker/api/v1"
)

// client wraps an http client and add a bunch of stuff to it
type client struct {
	// dependencies
	// auth
	// server
	// logger (?)

	dispatch endpoint.Endpoint
	stop     endpoint.Endpoint
	query    endpoint.Endpoint
}

func New(url *url.URL) *client {
	return &client{
		dispatch: makeDispatchEndpoint(url),
		stop:     makeStopEndpoint(url),
		query:    makeQueryEndpoint(url),
	}
}

func (c *client) Dispatch(name string, args ...string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	request := v1.DispatchRequest{
		Name: name,
		Args: args,
	}

	res, err := c.dispatch(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *client) Stop(jobID string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	request := v1.StopRequest{
		ID: jobID,
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

	request := v1.QueryInfoRequest{
		ID: jobID,
	}

	res, err := c.query(ctx, request)
	if err != nil {
		return nil, err
	}

	return res, nil
}
