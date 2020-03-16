package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/satriajidam/go-gin-skeleton/pkg/config"
	"github.com/satriajidam/go-gin-skeleton/pkg/log"
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

	router := gin.New()

	return &Server{
		Router: router,
		Port:   config.Get().HTTPPort,
	}
}

// Start starts the HTTP server,
func (s *Server) Start() error {
	log.Info(fmt.Sprintf("starting http server on port %s", s.Port))
	return s.Router.Run(fmt.Sprintf(":%s", s.Port))
}
