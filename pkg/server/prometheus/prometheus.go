package prometheus

import (
	"context"
	"fmt"
	"net/http"

	"github.com/satriajidam/go-gin-skeleton/pkg/log"
	httpserver "github.com/satriajidam/go-gin-skeleton/pkg/server/http"
	"github.com/satriajidam/go-gin-skeleton/pkg/telemetry/metric"
	metricbackend "github.com/satriajidam/go-gin-skeleton/pkg/telemetry/metric/backend/opencensus"
	"github.com/satriajidam/go-gin-skeleton/pkg/telemetry/metric/middleware"
	ginmiddleware "github.com/satriajidam/go-gin-skeleton/pkg/telemetry/metric/middleware/gin"
	"github.com/satriajidam/go-gin-skeleton/pkg/util"
)

// Server represents the implementation of Prometheus server object.
type Server struct {
	http *http.Server
	Port string
	Path string
}

// Target defines a target gin engine to monitor.
type Target struct {
	HTTPServer             *httpserver.Server
	ExcludePaths           []string
	MetricsPrefix          string
	GroupedStatus          bool
	DisableMeasureReqSize  bool
	DisableMeasureRespSize bool
	DisableMeasureInflight bool
}

func (t *Target) filterMonitoredPaths(paths []string) []string {
	included := []string{}

	for _, p := range paths {
		if t.isExcluded(p) {
			continue
		}
		included = append(included, p)
	}

	return included
}

func (t *Target) isExcluded(path string) bool {
	for _, e := range t.ExcludePaths {
		if path == e {
			return true
		}
	}
	return false
}

// NewServer creates new Prometheus server.
func NewServer(port, path string) *Server {
	if port == "" {
		port = "9180"
	}

	if path == "" {
		path = "/metrics"
	}

	return &Server{
		Port: port,
		Path: path,
	}
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	log.Info(fmt.Sprintf("Start Prometheus server on port %s", s.Port))

	mux := http.NewServeMux()
	handler, err := metricbackend.DefaultPrometheusExporter()
	if err != nil {
		return err
	}

	mux.Handle(s.Path, handler)
	s.http = &http.Server{
		Addr:    fmt.Sprintf(":%s", s.Port),
		Handler: mux,
	}

	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Stop stops the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	log.Info(fmt.Sprintf("Stop Prometheus server on port %s", s.Port))
	if err := s.http.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

// Monitor registers gin engine(s) to monitor.
func (s *Server) Monitor(targets ...*Target) {
	for _, t := range targets {
		mdlw := middleware.NewHTTPMiddleware(middleware.HTTPMiddlewareConfig{
			Recorder:               metricbackend.NewHTTPRecorder(metric.HTTPRecorderConfig{}),
			Host:                   fmt.Sprintf("%s:%s", util.GetHostname(), t.HTTPServer.Port),
			GroupedStatus:          t.GroupedStatus,
			DisableMeasureReqSize:  t.DisableMeasureReqSize,
			DisableMeasureRespSize: t.DisableMeasureRespSize,
			DisableMeasureInflight: t.DisableMeasureInflight,
		})

		t.HTTPServer.AddMiddleware(ginmiddleware.HTTPHandler(mdlw))
	}
}
