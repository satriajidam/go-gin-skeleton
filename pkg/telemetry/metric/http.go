package metric

import (
	"context"
	"time"
)

// HTTPProperty stores properties for the HTTP metrics based on incoming requests.
type HTTPProperty struct {
	// Host is the host that has served the request.
	Host string
	// Endpoint is the endpoint of the request handler.
	Endpoint string
	// Method is the method of the request.
	Method string
	// Status is the response code of the request.
	Status int
}

// HTTPRecorder records and measures the HTTP metrics.
// This interface has the required methods to be used with the HTTP middlewares.
type HTTPRecorder interface {
	// RecordRequestDuration measures the duration of an HTTP request.
	RecordRequestDuration(ctx context.Context, prop HTTPProperty, duration time.Duration)
	// RecordRequestSize measures the size of an HTTP request in bytes.
	RecordRequestSize(ctx context.Context, prop HTTPProperty, sizeBytes int64)
	// RecordResponseSize measures the size of an HTTP response in bytes.
	RecordResponseSize(ctx context.Context, prop HTTPProperty, sizeBytes int64)
	// AddTotalRequests increments the total of completed requests.
	AddTotalRequests(ctx context.Context, prop HTTPProperty, quantity int64)
	// AddInflightRequests increments and decrements the number of inflight requests.
	AddInflightRequests(ctx context.Context, prop HTTPProperty, quantity int64)
}
