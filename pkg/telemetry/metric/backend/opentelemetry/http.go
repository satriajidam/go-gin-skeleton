package opentelemetry

import (
	"context"
	"time"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/satriajidam/go-gin-skeleton/pkg/telemetry/metric"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/unit"
)

// HTTPRecorderConfig stores configurations for the HTTP metrics recorder.
type HTTPRecorderConfig struct {
	// InstrumentationName must be the name of the library providing instrumentation.
	// This name may be the same as the instrumented code only if that code provides
	// built-in instrumentation. If the instrumentationName is empty, it will be set
	// to `http`.
	InstrumentationName string
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

func (c *HTTPRecorderConfig) defaults() {
	if c.InstrumentationName == "" {
		c.InstrumentationName = "http"
	}

	if len(c.DurationBuckets) == 0 {
		c.DurationBuckets = prom.DefBuckets
	}

	if len(c.SizeBuckets) == 0 {
		c.SizeBuckets = prom.ExponentialBuckets(100, 10, 8)
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

type httpRecorder struct {
	// Label keys.
	hostKey     label.Key
	endpointKey label.Key
	methodKey   label.Key
	statusKey   label.Key

	// Measurements.
	requestDuration  *otelmetric.Float64ValueRecorder
	requestSize      *otelmetric.Float64ValueRecorder
	responseSize     *otelmetric.Float64ValueRecorder
	requestsTotal    *otelmetric.Int64Counter
	requestsInflight *otelmetric.Int64UpDownCounter
}

// NewHTTPRecorder returns a new Recorder that uses OpenTelemetry as the backend.
func NewHTTPRecorder(cfg HTTPRecorderConfig) metric.HTTPRecorder {
	cfg.defaults()

	r := &httpRecorder{}

	r.initLabelKeys(cfg)
	r.initMeasurements(cfg)

	return r
}

func (r *httpRecorder) initLabelKeys(cfg HTTPRecorderConfig) {
	r.hostKey = label.Key(cfg.HostLabel)
	r.endpointKey = label.Key(cfg.EndpointLabel)
	r.methodKey = label.Key(cfg.MethodLabel)
	r.statusKey = label.Key(cfg.StatusLabel)
}

func (r *httpRecorder) initMeasurements(cfg HTTPRecorderConfig) {
	meter := otel.Meter(cfg.InstrumentationName)

	requestDuration := otelmetric.Must(meter).NewFloat64ValueRecorder(
		"http_request_duration_milliseconds",
		otelmetric.WithDescription("The latency of the HTTP request."),
		otelmetric.WithUnit(unit.Milliseconds),
	)
	requestSize := otelmetric.Must(meter).NewFloat64ValueRecorder(
		"http_request_size_bytes",
		otelmetric.WithDescription("The size of the HTTP request."),
		otelmetric.WithUnit(unit.Bytes),
	)
	responseSize := otelmetric.Must(meter).NewFloat64ValueRecorder(
		"http_response_size_bytes",
		otelmetric.WithDescription("The size of the HTTP response."),
		otelmetric.WithUnit(unit.Bytes),
	)
	requestsTotal := otelmetric.Must(meter).NewInt64Counter(
		"http_requests_total",
		otelmetric.WithDescription("The total number of completed HTTP requests."),
		otelmetric.WithUnit(unit.Dimensionless),
	)
	requestsInflight := otelmetric.Must(meter).NewInt64UpDownCounter(
		"http_requests_inflight",
		otelmetric.WithDescription("The number of inflight requests being processed at the same time."),
		otelmetric.WithUnit(unit.Dimensionless),
	)

	r.requestDuration = &requestDuration
	r.requestSize = &requestSize
	r.responseSize = &responseSize
	r.requestsTotal = &requestsTotal
	r.requestsInflight = &requestsInflight
}

func (r *httpRecorder) reqPropToLabelPairs(prop metric.HTTPRequestProperty) []label.KeyValue {
	return []label.KeyValue{
		r.hostKey.String(prop.Host),
		r.endpointKey.String(prop.Endpoint),
		r.methodKey.String(prop.Method),
		r.statusKey.String(prop.Status),
	}
}

func (r *httpRecorder) infPropToLabelPairs(prop metric.HTTPInflightProperty) []label.KeyValue {
	return []label.KeyValue{
		r.hostKey.String(prop.Host),
		r.endpointKey.String(prop.Endpoint),
	}
}

func (r *httpRecorder) RecordRequestDuration(
	ctx context.Context, prop metric.HTTPRequestProperty, duration time.Duration,
) {
	r.requestDuration.Record(ctx, float64(duration.Milliseconds()), r.reqPropToLabelPairs(prop)...)
}

func (r *httpRecorder) RecordRequestSize(
	ctx context.Context, prop metric.HTTPRequestProperty, sizeBytes int64,
) {
	r.requestSize.Record(ctx, float64(sizeBytes), r.reqPropToLabelPairs(prop)...)
}

func (r *httpRecorder) RecordResponseSize(
	ctx context.Context, prop metric.HTTPRequestProperty, sizeBytes int64,
) {
	r.responseSize.Record(ctx, float64(sizeBytes), r.reqPropToLabelPairs(prop)...)
}

func (r *httpRecorder) AddTotalRequests(
	ctx context.Context, prop metric.HTTPRequestProperty, quantity int64,
) {
	r.requestsTotal.Add(ctx, quantity, r.reqPropToLabelPairs(prop)...)
}

func (r *httpRecorder) AddInflightRequests(
	ctx context.Context, prop metric.HTTPInflightProperty, quantity int64,
) {
	r.requestsInflight.Add(ctx, quantity, r.infPropToLabelPairs(prop)...)
}
