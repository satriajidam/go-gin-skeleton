package metric

import (
	"context"
	"time"
)

// HTTPRequestProperty stores properties for the HTTP metrics of an incoming request.
type HTTPRequestProperty struct {
	// Host is the host that has served the request.
	Host string
	// Endpoint is the endpoint of the request handler.
	Endpoint string
	// Method is the method of the request.
	Method string
	// Status is the response code of the request.
	Status string
}

// HTTPInflightProperty stores properties for the HTTP metrics of an inflight request.
type HTTPInflightProperty struct {
	// Host is the host that has served the request.
	Host string
	// Endpoint is the endpoint of the request handler.
	Endpoint string
}

// HTTPRecorder records and measures the HTTP metrics.
// This interface has the required methods to be implemented by the HTTP metrics backend
// and used by the middleware.
type HTTPRecorder interface {
	// RecordRequestDuration measures the duration of an HTTP request.
	RecordRequestDuration(ctx context.Context, prop HTTPRequestProperty, duration time.Duration)
	// RecordRequestSize measures the size of an HTTP request in bytes.
	RecordRequestSize(ctx context.Context, prop HTTPRequestProperty, sizeBytes int64)
	// RecordResponseSize measures the size of an HTTP response in bytes.
	RecordResponseSize(ctx context.Context, prop HTTPRequestProperty, sizeBytes int64)
	// AddTotalRequests increments the total of completed requests.
	AddTotalRequests(ctx context.Context, prop HTTPRequestProperty, quantity int64)
	// AddInflightRequests increments and decrements the number of inflight requests.
	AddInflightRequests(ctx context.Context, prop HTTPInflightProperty, quantity int64)
}
