package opentelemetry

import (
	"fmt"

	otelprom "go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/sdk/metric/controller/pull"
)

// NewPrometheusExporter creates a new global OpenTelemetry Prometheus exporter.
func NewPrometheusExporter(cfg otelprom.Config, opts ...pull.Option) (*otelprom.Exporter, error) {
	exporter, err := otelprom.InstallNewPipeline(cfg, opts...)
	if err != nil {
		panic(fmt.Errorf("failed creating opentelemetry prometheus exporter: %v", err))
	}
	return exporter, nil
}

// DefaultPrometheusExporter creates a new global OpenTelemetry Prometheus exporter
// using default Prometheus reporter config.
func DefaultPrometheusExporter() (*otelprom.Exporter, error) {
	return NewPrometheusExporter(otelprom.Config{})
}
