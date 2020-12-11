// IMPORTANT:
// This package is still a work in progress. It isn't recommended for production usage.
// Currently, the OpenTelemetry dev team is still in the middle of migrating opencensus-go
// APIs to opentelemetry-go library. So for production usage please use the OpenCensus
// backend instead.
package opentelemetry

import (
	"context"
	"time"

	"github.com/satriajidam/go-gin-skeleton/pkg/telemetry/metric"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/unit"
)

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
func NewHTTPRecorder(cfg metric.HTTPRecorderConfig) metric.HTTPRecorder {
	cfg.Defaults()

	r := &httpRecorder{}

	r.initLabelKeys(cfg)
	r.initMeasurements(cfg)

	// TODO: Use OpenCensus View APIs (WIP) to configure metric outputs.
	// Please refer to the following article to learn more about this topic:
	// https://github.com/open-telemetry/opentelemetry-go/issues/689

	return r
}

func (r *httpRecorder) initLabelKeys(cfg metric.HTTPRecorderConfig) {
	r.hostKey = label.Key(cfg.HostLabel)
	r.endpointKey = label.Key(cfg.EndpointLabel)
	r.methodKey = label.Key(cfg.MethodLabel)
	r.statusKey = label.Key(cfg.StatusLabel)
}

func (r *httpRecorder) initMeasurements(cfg metric.HTTPRecorderConfig) {
	meter := otel.Meter("http")

	requestDuration := otelmetric.Must(meter).NewFloat64ValueRecorder(
		metric.HTTPRequestDuration().Name,
		otelmetric.WithDescription(metric.HTTPRequestDuration().Description),
		otelmetric.WithUnit(unit.Milliseconds),
	)
	requestSize := otelmetric.Must(meter).NewFloat64ValueRecorder(
		metric.HTTPRequestSize().Name,
		otelmetric.WithDescription(metric.HTTPRequestSize().Description),
		otelmetric.WithUnit(unit.Bytes),
	)
	responseSize := otelmetric.Must(meter).NewFloat64ValueRecorder(
		metric.HTTPResponseSize().Name,
		otelmetric.WithDescription(metric.HTTPResponseSize().Description),
		otelmetric.WithUnit(unit.Bytes),
	)
	requestsTotal := otelmetric.Must(meter).NewInt64Counter(
		metric.HTTPRequestsTotal().Name,
		otelmetric.WithDescription(metric.HTTPRequestsTotal().Description),
		otelmetric.WithUnit(unit.Dimensionless),
	)
	requestsInflight := otelmetric.Must(meter).NewInt64UpDownCounter(
		metric.HTTPRequestsInflight().Name,
		otelmetric.WithDescription(metric.HTTPRequestsInflight().Description),
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
