package opentelemetry

import (
	"context"
	"time"

	"github.com/satriajidam/go-gin-skeleton/pkg/telemetry/metric/http"
)

type recorder struct {
}

// NewRecorder returns a new Recorder that uses OpenTelemetry as the backend.
func NewRecorder() (http.Recorder, error) {
	r := &recorder{}

	return r, nil
}

func (r recorder) RecordHTTPRequestDuration(
	ctx context.Context, props http.ReqProperties, duration time.Duration,
) {
}

func (r recorder) RecordHTTPRequestSize(
	ctx context.Context, props http.ReqProperties, sizeBytes int64,
) {
}

func (r recorder) RecordHTTPResponseSize(
	ctx context.Context, props http.ReqProperties, sizeBytes int64,
) {
}

func (r recorder) AddCompletedRequests(
	ctx context.Context, props http.ReqProperties, quantity int,
) {
}

func (r recorder) AddInflightRequests(
	ctx context.Context, props http.ReqProperties, quantity int,
) {
}
