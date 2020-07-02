package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/satriajidam/go-gin-skeleton/pkg/log"
)

// Server represents the HTTP server object.
type Server struct {
	router *gin.Engine
	port   string
}

// NewServer creates new HTTP server.
func NewServer(port, mode string, disallowUnknownJSONFields bool) *Server {
	gin.SetMode(mode)

	if disallowUnknownJSONFields {
		gin.EnableJsonDecoderDisallowUnknownFields()
	}

	return &Server{
		router: gin.New(),
		port:   port,
	}
}

// Start starts the HTTP server,
func (s *Server) Start() error {
	log.Info(fmt.Sprintf("Starting http server on port %s", s.port))
	return s.router.Run(fmt.Sprintf(":%s", s.port))
}
