package metric

var (
	// Latency in seconds buckets for histogram metric (5ms to 10s).
	durationBuckets = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}
	// Size in bytes buckets for histogram metrics (100B to 1GB).
	sizeBuckets = []float64{100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000}
)

type metadata struct {
	Name        string
	Description string
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
