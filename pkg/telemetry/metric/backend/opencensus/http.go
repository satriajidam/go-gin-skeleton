package opencensus

import (
	"context"
	"fmt"
	"time"

	"github.com/satriajidam/go-gin-skeleton/pkg/telemetry/metric"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

type httpRecorder struct {
	// Tag keys.
	hostKey     tag.Key
	endpointKey tag.Key
	methodKey   tag.Key
	statusKey   tag.Key

	// Measurements.
	requestDuration  *stats.Float64Measure
	requestSize      *stats.Int64Measure
	responseSize     *stats.Int64Measure
	requestsTotal    *stats.Int64Measure
	requestsInflight *stats.Int64Measure
}

// NewHTTPRecorder returns a new HTTP Recorder with OpenCensus backend.
func NewHTTPRecorder(cfg metric.HTTPRecorderConfig) metric.HTTPRecorder {
	cfg.Defaults()

	r := &httpRecorder{}

	if err := r.initTagKeys(cfg); err != nil {
		panic(fmt.Errorf("failed initializing opencensus http recorder tag keys: %v", err))
	}

	r.initMeasurements()

	if err := r.registerViews(cfg); err != nil {
		panic(fmt.Errorf("failed registering opencensus http recorder views: %v", err))
	}

	return r
}

func (r *httpRecorder) initTagKeys(cfg metric.HTTPRecorderConfig) error {
	hostKey, err := tag.NewKey(cfg.HostLabel)
	if err != nil {
		return err
	}
	r.hostKey = hostKey

	endpointKey, err := tag.NewKey(cfg.EndpointLabel)
	if err != nil {
		return err
	}
	r.endpointKey = endpointKey

	methodKey, err := tag.NewKey(cfg.MethodLabel)
	if err != nil {
		return err
	}
	r.methodKey = methodKey

	statusKey, err := tag.NewKey(cfg.StatusLabel)
	if err != nil {
		return err
	}
	r.statusKey = statusKey

	return nil
}

func (r *httpRecorder) initMeasurements() {
	r.requestDuration = stats.Float64(
		metric.HTTPRequestDuration().Name,
		metric.HTTPRequestDuration().Description,
		stats.UnitSeconds,
	)

	r.requestSize = stats.Int64(
		metric.HTTPRequestSize().Name,
		metric.HTTPRequestSize().Description,
		stats.UnitBytes,
	)

	r.responseSize = stats.Int64(
		metric.HTTPResponseSize().Name,
		metric.HTTPResponseSize().Description,
		stats.UnitBytes,
	)

	r.requestsTotal = stats.Int64(
		metric.HTTPRequestsTotal().Name,
		metric.HTTPRequestsTotal().Description,
		stats.UnitDimensionless,
	)

	r.requestsInflight = stats.Int64(
		metric.HTTPRequestsInflight().Name,
		metric.HTTPRequestsInflight().Description,
		stats.UnitDimensionless,
	)
}

func (r *httpRecorder) registerViews(cfg metric.HTTPRecorderConfig) error {
	requestTagKeys := []tag.Key{r.hostKey, r.endpointKey, r.methodKey, r.statusKey}
	inflightTagKeys := []tag.Key{r.hostKey, r.endpointKey}

	requestDurationView := &view.View{
		Name:        metric.HTTPRequestDuration().Name,
		Description: metric.HTTPRequestDuration().Description,
		TagKeys:     requestTagKeys,
		Measure:     r.requestDuration,
		Aggregation: view.Distribution(cfg.DurationBuckets...),
	}

	requestSizeView := &view.View{
		Name:        metric.HTTPRequestSize().Name,
		Description: metric.HTTPRequestSize().Description,
		TagKeys:     requestTagKeys,
		Measure:     r.requestSize,
		Aggregation: view.Distribution(cfg.SizeBuckets...),
	}

	responseSizeView := &view.View{
		Name:        metric.HTTPResponseSize().Name,
		Description: metric.HTTPResponseSize().Description,
		TagKeys:     requestTagKeys,
		Measure:     r.responseSize,
		Aggregation: view.Distribution(cfg.SizeBuckets...),
	}

	requestsTotalVIew := &view.View{
		Name:        metric.HTTPRequestsTotal().Name,
		Description: metric.HTTPRequestsTotal().Description,
		TagKeys:     requestTagKeys,
		Measure:     r.requestsTotal,
		Aggregation: view.Sum(),
	}

	requestsInflightVIew := &view.View{
		Name:        metric.HTTPRequestsInflight().Name,
		Description: metric.HTTPRequestsInflight().Description,
		TagKeys:     inflightTagKeys,
		Measure:     r.requestsInflight,
		Aggregation: view.Sum(),
	}

	return view.Register(
		requestDurationView,
		requestSizeView,
		responseSizeView,
		requestsTotalVIew,
		requestsInflightVIew,
	)
}

// ctxWithTagFromHTTPReqProperties generates new context that contains a tag map with values from
// the provided HTTP request property.
func (r *httpRecorder) ctxWithTagFromHTTPReqProperty(
	ctx context.Context, prop metric.HTTPRequestProperty,
) context.Context {
	newCtx, _ := tag.New(ctx,
		tag.Upsert(r.hostKey, prop.Host),
		tag.Upsert(r.endpointKey, prop.Endpoint),
		tag.Upsert(r.methodKey, prop.Method),
		tag.Upsert(r.statusKey, prop.Status),
	)
	return newCtx
}

// ctxWithTagFromHTTPInfProperty generates new context that contains a tag map with values from
// the provided HTTP inflight property.
func (r *httpRecorder) ctxWithTagFromHTTPInfProperty(
	ctx context.Context, prop metric.HTTPInflightProperty,
) context.Context {
	newCtx, _ := tag.New(ctx,
		tag.Upsert(r.hostKey, prop.Host),
		tag.Upsert(r.endpointKey, prop.Endpoint),
	)
	return newCtx
}

func (r *httpRecorder) RecordRequestDuration(
	ctx context.Context, prop metric.HTTPRequestProperty, duration time.Duration,
) {
	ctx = r.ctxWithTagFromHTTPReqProperty(ctx, prop)
	stats.Record(ctx, r.requestDuration.M(duration.Seconds()))
}

func (r *httpRecorder) RecordRequestSize(
	ctx context.Context, prop metric.HTTPRequestProperty, sizeBytes int64,
) {
	ctx = r.ctxWithTagFromHTTPReqProperty(ctx, prop)
	stats.Record(ctx, r.requestSize.M(sizeBytes))
}

func (r *httpRecorder) RecordResponseSize(
	ctx context.Context, prop metric.HTTPRequestProperty, sizeBytes int64,
) {
	ctx = r.ctxWithTagFromHTTPReqProperty(ctx, prop)
	stats.Record(ctx, r.responseSize.M(sizeBytes))
}

func (r *httpRecorder) AddTotalRequests(
	ctx context.Context, prop metric.HTTPRequestProperty, quantity int64,
) {
	ctx = r.ctxWithTagFromHTTPReqProperty(ctx, prop)
	stats.Record(ctx, r.requestsTotal.M(quantity))
}

func (r *httpRecorder) AddInflightRequests(
	ctx context.Context, prop metric.HTTPInflightProperty, quantity int64,
) {
	ctx = r.ctxWithTagFromHTTPInfProperty(ctx, prop)
	stats.Record(ctx, r.requestsInflight.M(quantity))
}
