package opentelemetry

import (
	"fmt"

	prom "github.com/prometheus/client_golang/prometheus"

	otelprom "go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/sdk/metric/controller/pull"
)

// NewPrometheusExporter creates a new global Prometheus OpenTelemetry exporter.
func NewPrometheusExporter(cfg otelprom.Config, opts ...pull.Option) (*otelprom.Exporter, error) {
	exporter, err := otelprom.InstallNewPipeline(cfg, opts...)
	if err != nil {
		panic(fmt.Errorf("failed to initialize prometheus exporter: %v", err))
	}
	return exporter, nil
}

// DefaultPrometheusExporter creates a new global Prometheus OpenTelemetry exporter
// using default Prometheus reporter config.
func DefaultPrometheusExporter(options ...pull.Option) (*otelprom.Exporter, error) {
	return NewPrometheusExporter(otelprom.Config{
		Registry:                   prom.NewRegistry(),
		Registerer:                 prom.DefaultRegisterer,
		Gatherer:                   prom.DefaultGatherer,
		DefaultSummaryQuantiles:    nil,
		DefaultHistogramBoundaries: nil,
	})
}
