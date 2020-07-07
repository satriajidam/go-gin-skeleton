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
	http        *http.Server
	Router      *gin.Engine
	Port        string
	middlewares []gin.HandlerFunc
}

// NewServer creates new HTTP server.
func NewServer(port string, disallowUnknownJSONFields bool) *Server {
	if disallowUnknownJSONFields {
		gin.EnableJsonDecoderDisallowUnknownFields()
	}

	router := gin.New()

	loadPredefinedRoutes(router)

	return &Server{
		Router: router,
		Port:   port,
		middlewares: []gin.HandlerFunc{
			// Default gin middlewares.
			gin.Recovery(),
			requestid.New(),
			logger.New(port),
		},
	}
}

// AddMiddleware adds a gin middleware the HTTP server.
func (s *Server) AddMiddleware(h gin.HandlerFunc) {
	s.middlewares = append(s.middlewares, h)
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	log.Info(fmt.Sprintf("Start HTTP server on port %s", s.Port))
	s.Router.Use(s.middlewares...)
	s.http = &http.Server{
		Addr:    fmt.Sprintf(":%s", s.Port),
		Handler: s.Router,
	}
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
