package prometheus

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/satriajidam/go-gin-skeleton/pkg/log"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	ginmiddleware "github.com/slok/go-http-metrics/middleware/gin"
)

// Server represents the implementation of Prometheus server object.
type Server struct {
	http *http.Server
	Port string
	Path string
}

// Target defines a target gin engine to monitor.
type Target struct {
	Engine        *gin.Engine
	MetricsPrefix string
}

// NewServer creates new Prometheus server.
func NewServer(port, path string) *Server {
	mux := http.NewServeMux()
	mux.Handle(path, promhttp.Handler())

	return &Server{
		http: &http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: mux,
		},
		Port: port,
	}
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	log.Info(fmt.Sprintf("Start Prometheus server on port %s", s.Port))
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
		})
		t.Engine.Use(ginmiddleware.Handler("", mdlw))
	}
}
