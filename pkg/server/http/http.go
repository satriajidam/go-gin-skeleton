package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/satriajidam/go-gin-skeleton/pkg/log"
)

// Server represents the HTTP server object.
type Server struct {
	http *http.Server
	port string
}

// NewServer creates new HTTP server.
func NewServer(port, mode string, disallowUnknownJSONFields bool) *Server {
	gin.SetMode(mode)

	if disallowUnknownJSONFields {
		gin.EnableJsonDecoderDisallowUnknownFields()
	}

	return &Server{
		http: &http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: gin.Default(),
		},
		port: port,
	}
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	log.Info(fmt.Sprintf("Start HTTP server on port %s", s.port))
	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop stops the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	log.Info(fmt.Sprintf("Stop HTTP server on port %s", s.port))
	if err := s.http.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
