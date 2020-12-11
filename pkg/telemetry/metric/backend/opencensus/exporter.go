package opencensus

import (
	"fmt"

	ocenprom "contrib.go.opencensus.io/exporter/prometheus"
)

// NewPrometheusExporter creates a new global OpenCensus Prometheus exporter.
func NewPrometheusExporter(opts ocenprom.Options) (*ocenprom.Exporter, error) {
	exporter, err := ocenprom.NewExporter(opts)
	if err != nil {
		panic(fmt.Errorf("failed creating opencensus prometheus exporter: %v", err))
	}
	return exporter, nil
}

// DefaultPrometheusExporter creates a new global OpenCensus Prometheus exporter
// using default Prometheus reporter config.
func DefaultPrometheusExporter() (*ocenprom.Exporter, error) {
	return NewPrometheusExporter(ocenprom.Options{})
}
