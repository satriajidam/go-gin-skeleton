package metric

import (
	"context"
	"time"
)

// HTTPRequestProperty stores properties for the HTTP metrics of an incoming request.
type HTTPRequestProperty struct {
	Host     string
	Endpoint string
	Method   string
	Status   string
}

// HTTPInflightProperty stores properties for the HTTP metrics of an inflight request.
type HTTPInflightProperty struct {
	Host     string
	Endpoint string
	Method   string
}

// HTTPRecorder records and measures the HTTP metrics.
// This interface has the required methods to be implemented by the HTTP metrics backend.
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
	DurationBuckets []float64
	SizeBuckets     []float64
	HostLabel       string
	EndpointLabel   string
	MethodLabel     string
	StatusLabel     string
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

// HTTPRequestDuration returns HTTP request duration metric metadata.
func HTTPRequestDuration() metadata {
	return metadata{
		Name:        "http_request_duration_seconds",
		Description: "The latency of the HTTP request in seconds.",
	}
}

// HTTPRequestSize returns HTTP request size metric metadata.
func HTTPRequestSize() metadata {
	return metadata{
		Name:        "http_request_size_bytes",
		Description: "The size of the HTTP request in bytes.",
	}
}

// HTTPResponseSize returns HTTP response size metric metadata.
func HTTPResponseSize() metadata {
	return metadata{
		Name:        "http_response_size_bytes",
		Description: "The size of the HTTP response in bytes.",
	}
}

// HTTPRequestsTotal returns HTTP requests total metric metadata.
func HTTPRequestsTotal() metadata {
	return metadata{
		Name:        "http_requests_total",
		Description: "The total number of completed HTTP requests.",
	}
}

// HTTPRequestsInflight returns HTTP requests inflight metric metadata.
func HTTPRequestsInflight() metadata {
	return metadata{
		Name:        "http_requests_inflight",
		Description: "The number of inflight requests being processed at the same time.",
	}
}
