package http

import "github.com/satriajidam/go-gin-skeleton/pkg/telemetry/metric"

// HTTPMiddlewareConfig is the configuration for the HTTP middleware factory.
type HTTPMiddlewareConfig struct {
	// Recorder is the way the HTTP metrics will be recorded in the different backends.
	Recorder metric.HTTPMetricRecorder
	// Host is an optional identifier for the metrics host, this can be useful if
	// the same app has multiple servers (e.g API, metrics and healthchecks).
	Host string
	// GroupedStatus will group the status label in the form of `\dxx`, for example,
	// 200, 201, and 203 will have the label `status="2xx"`. This impacts on the cardinality
	// of the metrics and also improves the performance of queries that are grouped by
	// status code because there are already aggregated in the metric.
	// By default will be false.
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
