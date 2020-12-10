package middleware

import (
	"context"
	"os"

	"github.com/satriajidam/go-gin-skeleton/pkg/telemetry/metric"
	"github.com/satriajidam/go-gin-skeleton/pkg/telemetry/metric/opentelemetry"
)

// HTTPReporter knows how to report the data to the Middleware object so it can
// measure the different framework/libraries.
type HTTPReporter interface {
	Method() string
	Context() context.Context
	URLPath() string
	StatusCode() int
	BytesWritten() int64
}

// HTTPMiddlewareConfig is the configuration for the HTTP middleware factory.
type HTTPMiddlewareConfig struct {
	// Recorder is the way the HTTP metrics will be recorded in the different backends.
	// By default it will use OpenTelemetry.
	Recorder *metric.HTTPMetricRecorder
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
	// DisableMeasureInflight will disable the recording metrics about the inflight requests number,
	// by default measuring inflights is enabled (`DisableMeasureInflight` is false).
	DisableMeasureInflight bool
}

func (c *HTTPMiddlewareConfig) defaults() {
	if c.Recorder == nil {
		recorder := opentelemetry.NewHTTPRecorder(opentelemetry.HTTPRecorderConfig{})
		c.Recorder = &recorder
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
}
