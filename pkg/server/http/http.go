package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server object.
type Server struct {
	router *gin.Engine
}

// New creates new HTTP server.
func New(appMode string, disallowUnknownJSONFields bool) *Server {
	if appMode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	if disallowUnknownJSONFields {
		gin.EnableJsonDecoderDisallowUnknownFields()
	}

	router := gin.Default()

	return &Server{router}
}

// Start starts the HTTP server.
func (s *Server) Start(port string) {
	s.router.Run(fmt.Sprintf(":%s"), port)
}
