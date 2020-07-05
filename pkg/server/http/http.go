package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/satriajidam/go-gin-skeleton/pkg/log"
	"github.com/satriajidam/go-gin-skeleton/pkg/server/http/middleware/logger"
	"github.com/satriajidam/go-gin-skeleton/pkg/server/http/middleware/requestid"
)

// Server represents the implementation of HTTP server object.
type Server struct {
	http   *http.Server
	Router *gin.Engine
	Port   string
}

// NewServer creates new HTTP server.
func NewServer(port, mode string, disallowUnknownJSONFields bool) *Server {
	gin.SetMode(mode)

	if disallowUnknownJSONFields {
		gin.EnableJsonDecoderDisallowUnknownFields()
	}

	router := gin.New()

	// Setup middlewares.
	router.Use(
		gin.Recovery(),
		requestid.New(),
		logger.New(port),
	)

	loadPredefinedRoutes(router)

	return &Server{
		http: &http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: router,
		},
		Router: router,
		Port:   port,
	}
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	log.Info(fmt.Sprintf("Start HTTP server on port %s", s.Port))
	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop stops the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	log.Info(fmt.Sprintf("Stop HTTP server on port %s", s.Port))
	if err := s.http.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
