package middleware

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/satriajidam/go-gin-skeleton/pkg/telemetry/metric"
	"github.com/satriajidam/go-gin-skeleton/pkg/telemetry/metric/backend/opentelemetry"
)

// HTTPReporter knows how to report the data to the Middleware object so it can
// measure the different framework/libraries.
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
	// Host is an optional identifier for the metrics host, this can be useful if
	// the same app has multiple servers (e.g API, metrics and healthchecks).
	// By default it will be set to the current hostname.
	Host string
	// GroupedStatus will group the status label in the form of `\dxx`, for example,
	// 200, 201, and 203 will have the label `status="2xx"`. This impacts on the cardinality
	// of the metrics and also improves the performance of queries that are grouped by
	// status code because there are already aggregated in the metric.
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
		c.Recorder = opentelemetry.NewHTTPRecorder(opentelemetry.HTTPRecorderConfig{})
	}

	if c.Host == "" {
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "localhost"
		}
		c.Host = hostname
	}
}

// HTTPMiddleware is an object that knows how to measure an HTTP handler by wrapping
// another handler.
//
// Depending on the framework/library we want to measure, this can change a lot,
// to abstract the way how we measure on the different libraries, Middleware will
// recieve an `HTTPReporter` that knows how to get the data the Middleware object
// needs to measure.
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
			Host:     m.cfg.Host,
			Endpoint: reporter.URLPath(),
		}
		m.cfg.Recorder.AddInflightRequests(ctx, prop, 1)
		defer m.cfg.Recorder.AddInflightRequests(ctx, prop, -1)
	}

	start := time.Now()
	defer func() {
		duration := time.Since(start)

		// If we need to group the status code, it uses the
		// first number of the status code because it is the
		// least required identification way.
		var status string
		if m.cfg.GroupedStatus {
			status = fmt.Sprintf("%dxx", reporter.StatusCode()/100)
		} else {
			status = strconv.Itoa(reporter.StatusCode())
		}

		prop := metric.HTTPRequestProperty{
			Host:     m.cfg.Host,
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

	// Call the wrapped logic.
	next()
}
