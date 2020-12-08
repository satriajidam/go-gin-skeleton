package opentelemetry

import (
	"context"
	"time"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/satriajidam/go-gin-skeleton/pkg/telemetry/metric"
)

// Config has the dependencies and values of the recorder.
type HTTPRecorderConfig struct {
	// Prefix is the prefix that will be set on the metrics, by default it will be empty.
	Prefix string
	// DurationBuckets are the buckets used for the HTTP request duration metrics,
	// by default uses default buckets (from 5ms to 10s).
	DurationBuckets []float64
	// SizeBuckets are the buckets for the HTTP request/response size metrics,
	// by default uses a exponential buckets from 100B to 1GB.
	SizeBuckets []float64
	// ServiceLabel is the name that will be set to the service label, by default is `service`.
	ServiceLabel string
	// EndpointLabel is the name that will be set to the endpoint label, by default is `endpoint`.
	EndpointLabel string
	// MethodLabel is the name that will be set to the method label, by default is `method`.
	MethodLabel string
	// StatusLabel is the name that will be set to the response code label, by default is `status`.
	StatusLabel string
}

func (c *HTTPRecorderConfig) defaults() {
	if len(c.DurationBuckets) == 0 {
		c.DurationBuckets = prom.DefBuckets
	}

	if len(c.SizeBuckets) == 0 {
		c.SizeBuckets = prom.ExponentialBuckets(100, 10, 8)
	}

	if c.EndpointLabel == "" {
		c.EndpointLabel = "endpoint"
	}

	if c.StatusLabel == "" {
		c.StatusLabel = "status"
	}

	if c.MethodLabel == "" {
		c.MethodLabel = "method"
	}

	if c.ServiceLabel == "" {
		c.ServiceLabel = "service"
	}
}

type httpRecorder struct {
}

// NewHTTPRecorder returns a new Recorder that uses OpenTelemetry as the backend.
func NewHTTPRecorder(cfg HTTPRecorderConfig) (metric.HTTPMetricRecorder, error) {
	cfg.defaults()

	r := &httpRecorder{}

	return r, nil
}

func (r *httpRecorder) RecordRequestDuration(
	ctx context.Context, prop metric.HTTPMetricProperty, duration time.Duration,
) {
}

func (r *httpRecorder) RecordRequestSize(
	ctx context.Context, prop metric.HTTPMetricProperty, sizeBytes int64,
) {
}

func (r *httpRecorder) RecordResponseSize(
	ctx context.Context, prop metric.HTTPMetricProperty, sizeBytes int64,
) {
}

func (r *httpRecorder) AddCompletedRequests(
	ctx context.Context, prop metric.HTTPMetricProperty, quantity int,
) {
}

func (r *httpRecorder) AddInflightRequests(
	ctx context.Context, prop metric.HTTPMetricProperty, quantity int,
) {
}
