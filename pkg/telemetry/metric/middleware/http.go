package middleware

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/satriajidam/go-gin-skeleton/pkg/telemetry/metric"
	"github.com/satriajidam/go-gin-skeleton/pkg/telemetry/metric/backend/opencensus"
)

// HTTPReporter abstracts the ways to report the required data for measuring HTTP metrics
// to the middleware object.
//
// The most important thing to pay attention to when implementing this interface
// is to make sure the `URLPath()` function groups all endpoints with path parameters
// as a single URL path.
// As an example, `/api/v1/providers/:name` endpoint has `:name` as a path parameter,
// when multiple requests hit the endpoint using the following URL path:
// - /api/v1/providers/aws
// - /api/v1/providers/gcp
// - /api/v1/providers/azure
// then the reporter should only produce 1 metric with `/api/v1/providers/:name` as
// its endpoint label and an aggregated value, instead of 3 metrics with 3 different
// endpoint labels matching the requested URL path and 3 non-aggregated values. Failing
// to do this will cause the produced metrics to have extremely high cardinality
// which is against the recommended best practice.
// Please refer to the following article to learn more about this topic:
// https://banzaicloud.com/blog/monitoring-gin-with-prometheus/
type HTTPReporter interface {
	Context() context.Context
	URLHost() string
	URLPath() string
	Method() string
	StatusCode() int
	RequestSize() int64
	ResponseSize() int64
}

// HTTPMiddlewareConfig stores configurations for the HTTP middleware factory.
type HTTPMiddlewareConfig struct {
	// Recorder is the way the HTTP metrics will be recorded in the different backends.
	// By default it will use OpenTelemetry.
	Recorder metric.HTTPRecorder
	// GroupedStatus will group the status label in the form of `\dxx`, for example,
	// 200, 201, and 203 will have the label `status="2xx"`. This impacts on the cardinality
	// of the metrics and also improves the performance of queries that are grouped by
	// status code because they are already aggregated in the metric.
	// By default it will be set to false.
	GroupedStatus bool
	// DisableMeasureReqSize will disable the recording metrics about the request size,
	// by default measuring request size is enabled (`DisableMeasureReqSize` is false).
	DisableMeasureReqSize bool
	// DisableMeasureRespSize will disable the recording metrics about the response size,
	// by default measuring response size is enabled (`DisableMeasureRespSize` is false).
	DisableMeasureRespSize bool
	// DisableMeasureInflight will disable the recording metrics about the inflight requests,
	// by default measuring inflights is enabled (`DisableMeasureInflight` is false).
	DisableMeasureInflight bool
}

func (c *HTTPMiddlewareConfig) defaults() {
	if c.Recorder == nil {
		c.Recorder = opencensus.NewHTTPRecorder(metric.HTTPRecorderConfig{})
	}
}

// HTTPMiddleware is an object that knows how to measure an HTTP handler by wrapping
// another handler.
type HTTPMiddleware struct {
	cfg *HTTPMiddlewareConfig
}

// NewHTTPMiddleware creates a new HTTP Middleware object.
func NewHTTPMiddleware(cfg HTTPMiddlewareConfig) HTTPMiddleware {
	cfg.defaults()
	return HTTPMiddleware{cfg: &cfg}
}

// Measure abstracts the HTTP handler implementation by only requesting an HTTP reporter,
// this reporter will return the required data to be measured.
// It accepts a next function that will be called as the wrapped logic before and after
// measurement actions.
func (m *HTTPMiddleware) Measure(reporter HTTPReporter, next func()) {
	ctx := reporter.Context()

	if !m.cfg.DisableMeasureInflight {
		prop := metric.HTTPInflightProperty{
			Host:     reporter.URLHost(),
			Endpoint: reporter.URLPath(),
			Method:   reporter.Method(),
		}
		m.cfg.Recorder.AddInflightRequests(ctx, prop, 1)
		defer m.cfg.Recorder.AddInflightRequests(ctx, prop, -1)
	}

	start := time.Now()
	defer func() {
		duration := time.Since(start)

		var status string
		if m.cfg.GroupedStatus {
			status = fmt.Sprintf("%dxx", reporter.StatusCode()/100)
		} else {
			status = strconv.Itoa(reporter.StatusCode())
		}

		prop := metric.HTTPRequestProperty{
			Host:     reporter.URLHost(),
			Endpoint: reporter.URLPath(),
			Method:   reporter.Method(),
			Status:   status,
		}

		m.cfg.Recorder.AddTotalRequests(ctx, prop, 1)
		m.cfg.Recorder.RecordRequestDuration(ctx, prop, duration)

		if !m.cfg.DisableMeasureReqSize {
			m.cfg.Recorder.RecordRequestSize(ctx, prop, reporter.RequestSize())
		}

		if !m.cfg.DisableMeasureRespSize {
			m.cfg.Recorder.RecordResponseSize(ctx, prop, reporter.ResponseSize())
		}
	}()

	next()
}
