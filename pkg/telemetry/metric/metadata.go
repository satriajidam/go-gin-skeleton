package metric

type metadata struct {
	Name        string
	Description string
}

// List of metrics metadata for the HTTP metrics backend.
var HTTPRequestDuration = metadata{
	Name:        "http_request_duration_milliseconds",
	Description: "The latency of the HTTP request in milliseconds.",
}
var HTTPRequestSize = metadata{
	Name:        "http_request_size_bytes",
	Description: "The size of the HTTP request in bytes.",
}
var HTTPResponseSize = metadata{
	Name:        "http_response_size_bytes",
	Description: "The size of the HTTP response in bytes.",
}
var HTTPRequestsTotal = metadata{
	Name:        "http_requests_total",
	Description: "The total number of completed HTTP requests.",
}
var HTTPRequestsInflight = metadata{
	Name:        "http_requests_inflight",
	Description: "The number of inflight requests being processed at the same time.",
}
