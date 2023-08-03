package utils

import (
	"context"
)

type ValueKey string

const ContextValueKey_WaitGroup ValueKey = "WaitGroup"

type Request struct {
	RequestID string
	WorkSize  int
}

type Response struct {
	Status string
}

// RequestResponse contains a request response and the respective request that was used.
type RequestResponse struct {
	Request  Request
	Response *Response
}

// DoRequests executes a list of requests in parallel, and will fail-fast.
func DoRequests(ctx context.Context, reqs []Request, doWork func(context.Context, Request) (*Response, error)) ([]RequestResponse, error) {
	// If one of the requests fail, we want to be able to cancel the rest of them.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Initialize the Response Channel and Error Channel
	respChan, errChan := make(chan RequestResponse), make(chan error)

	for _, req := range reqs {
		// loopclosure: check references to loop variables from within nested functions
		// https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/loopclosure
		go func(req Request) {
			resp, err := doWork(ctx, req)
			if err != nil {
				errChan <- err
			} else {
				respChan <- RequestResponse{req, resp}
			}
		}(req)
	}

	// Initialize the response
	resps := make([]RequestResponse, 0, len(reqs))
	var firstErr error

	// Collect the response and fail-fast
	for range reqs {
		select {
		case resp := <-respChan:
			resps = append(resps, resp)
		case err := <-errChan:
			if firstErr == nil && err != nil {
				cancel()
				firstErr = err

				// Discard the response and return the first error
				return nil, firstErr
			}
		}
	}

	// firstErr will always be nil
	return resps, firstErr
}
