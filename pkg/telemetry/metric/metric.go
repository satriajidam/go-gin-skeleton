package metric

import (
	"context"
	"time"
)

// HTTPMetricProperty stores properties for the HTTP metrics based on incoming requests.
type HTTPMetricProperty struct {
	// Service is the service that has served the request.
	Service string
	// Endpoint is the endpoint of the request handler.
	Endpoint string
	// Method is the method of the request.
	Method string
	// Status is the response code of the request.
	Status string
}

// HTTPMetricRecorder records and measures the HTTP metrics.
// This interface has the required methods to be used with the HTTP middlewares.
type HTTPMetricRecorder interface {
	// RecordRequestDuration measures the duration of an HTTP request.
	RecordRequestDuration(ctx context.Context, prop HTTPMetricProperty, duration time.Duration)
	// RecordRequestSize measures the size of an HTTP request in bytes.
	RecordRequestSize(ctx context.Context, prop HTTPMetricProperty, sizeBytes int64)
	// RecordResponseSize measures the size of an HTTP response in bytes.
	RecordResponseSize(ctx context.Context, prop HTTPMetricProperty, sizeBytes int64)
	// AddCompletedRequests increments the number of completed requests.
	AddCompletedRequests(ctx context.Context, prop HTTPMetricProperty, quantity int)
	// AddInflightRequests increments and decrements the number of inflight requests.
	AddInflightRequests(ctx context.Context, prop HTTPMetricProperty, quantity int)
}
