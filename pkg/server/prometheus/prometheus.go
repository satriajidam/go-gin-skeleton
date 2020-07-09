package prometheus

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/satriajidam/go-gin-skeleton/pkg/log"
	httpserver "github.com/satriajidam/go-gin-skeleton/pkg/server/http"
	ginmiddleware "github.com/satriajidam/go-gin-skeleton/pkg/server/prometheus/middleware/gin"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
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
	MetricsPrefix          string
	GroupedStatus          bool
	DisableMeasureSize     bool
	DisableMeasureInflight bool
}

// NewServer creates new Prometheus server.
func NewServer(port, path string) *Server {
	return &Server{
		Port: port,
		Path: path,
	}
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	log.Info(fmt.Sprintf("Start Prometheus server on port %s", s.Port))
	mux := http.NewServeMux()
	mux.Handle(s.Path, promhttp.Handler())
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
		mdlw := middleware.New(middleware.Config{
			Recorder: metrics.NewRecorder(metrics.Config{
				Prefix: t.MetricsPrefix,
			}),
			Service:                fmt.Sprintf("localhost:%s", t.HTTPServer.Port),
			GroupedStatus:          t.GroupedStatus,
			DisableMeasureSize:     t.DisableMeasureSize,
			DisableMeasureInflight: t.DisableMeasureInflight,
		})
		t.HTTPServer.AddMiddleware(ginmiddleware.Handler(t.HTTPServer.GetRoutePaths(), mdlw))
	}
}
