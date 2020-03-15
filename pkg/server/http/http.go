package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/satriajidam/go-gin-skeleton/pkg/config"
)

// Server represents the HTTP server object.
type Server struct {
	Router *gin.Engine
	Port   string
}

// New creates new HTTP server.
func New() *Server {
	if config.IsReleaseMode() {
		gin.SetMode(gin.ReleaseMode)
	}

	if config.Get().GinDisallowUnknownJSONFields {
		gin.EnableJsonDecoderDisallowUnknownFields()
	}

	router := gin.Default()

	return &Server{
		Router: router,
		Port:   config.Get().HTTPPort,
	}
}

// Start starts the HTTP server,
func (s *Server) Start() error {
	return s.Router.Run(fmt.Sprintf(":%s", s.Port))
}
