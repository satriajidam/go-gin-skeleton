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

// HTTPRecorderConfig stores configurations for the HTTP metrics recorder.
type HTTPRecorderConfig struct {
	// DurationBuckets are the buckets used for the HTTP request duration metrics,
	// by default uses default buckets (from 5ms to 10s).
	DurationBuckets []float64
	// SizeBuckets are the buckets for the HTTP request/response size metrics,
	// by default uses a exponential buckets from 100B to 1GB.
	SizeBuckets []float64
	// HostLabel is the name that will be set to the host label, by default is `host`.
	HostLabel string
	// EndpointLabel is the name that will be set to the endpoint label, by default is `endpoint`.
	EndpointLabel string
	// MethodLabel is the name that will be set to the method label, by default is `method`.
	MethodLabel string
	// StatusLabel is the name that will be set to the response code label, by default is `status`.
	StatusLabel string
}

// Defaults sets default values for HTTP metrics recorder configurations.
func (c *HTTPRecorderConfig) Defaults() {
	if len(c.DurationBuckets) == 0 {
		c.DurationBuckets = durationBuckets
	}

	if len(c.SizeBuckets) == 0 {
		c.SizeBuckets = sizeBuckets
	}

	if c.HostLabel == "" {
		c.HostLabel = "host"
	}

	if c.EndpointLabel == "" {
		c.EndpointLabel = "endpoint"
	}

	if c.MethodLabel == "" {
		c.MethodLabel = "method"
	}

	if c.StatusLabel == "" {
		c.StatusLabel = "status"
	}
}
