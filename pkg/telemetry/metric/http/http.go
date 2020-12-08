package http

import (
	"context"
	"time"
)

// ReqProperties are properties for the metrics based on incoming HTTP requests.
type ReqProperties struct {
	// Service is the service that has served the request.
	Service string
	// Endpoint is the endpoint of the request handler.
	Endpoint string
	// Method is the method of the request.
	Method string
	// Status is the response code of the request.
	Status string
}

// Recorder records and measures the metrics.
// This interface has the required methods to be used with the HTTP middlewares.
type Recorder interface {
	// RecordHTTPRequestDuration measures the duration of an HTTP request.
	RecordHTTPRequestDuration(ctx context.Context, props ReqProperties, duration time.Duration)
	// RecordHTTPRequestSize measures the size of an HTTP request in bytes.
	RecordHTTPRequestSize(ctx context.Context, props ReqProperties, sizeBytes int64)
	// RecordHTTPResponseSize measures the size of an HTTP response in bytes.
	RecordHTTPResponseSize(ctx context.Context, props ReqProperties, sizeBytes int64)
	// AddCompletedRequests increments the number of completed requests.
	AddCompletedRequests(ctx context.Context, props ReqProperties, quantity int)
	// AddInflightRequests increments and decrements the number of inflight requests.
	AddInflightRequests(ctx context.Context, props ReqProperties, quantity int)
}

type dummy int

// Dummy is a dummy recorder.
const Dummy = dummy(0)

func (dummy) RecordHTTPRequestDuration(_ context.Context, _ ReqProperties, _ time.Duration) {}
func (dummy) RecordHTTPRequestSize(_ context.Context, _ ReqProperties, _ int64)             {}
func (dummy) RecordHTTPResponseSize(_ context.Context, _ ReqProperties, _ int64)            {}
func (dummy) AddCompletedRequests(_ context.Context, _ ReqProperties, _ int)                {}
func (dummy) AddInflightRequests(_ context.Context, _ ReqProperties, _ int)                 {}

var _ Recorder = Dummy
