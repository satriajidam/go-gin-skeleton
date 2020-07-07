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
	router      *gin.Engine
	Port        string
	middlewares []gin.HandlerFunc
	routes      []route
}

type route struct {
	method       string
	relativePath string
	handlers     []gin.HandlerFunc
}

// NewServer creates new HTTP server.
func NewServer(port string) *Server {
	return &Server{
		router: gin.New(),
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

func loadRoutes(router *gin.Engine, routes []route) {
	for _, route := range routes {
		switch route.method {
		case http.MethodGet:
			router.GET(route.relativePath, route.handlers...)
		case http.MethodHead:
			router.HEAD(route.relativePath, route.handlers...)
		case http.MethodPost:
			router.POST(route.relativePath, route.handlers...)
		case http.MethodPut:
			router.PUT(route.relativePath, route.handlers...)
		case http.MethodPatch:
			router.PATCH(route.relativePath, route.handlers...)
		case http.MethodDelete:
			router.DELETE(route.relativePath, route.handlers...)
		case http.MethodOptions:
			router.OPTIONS(route.relativePath, route.handlers...)
		}
	}
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	log.Info(fmt.Sprintf("Start HTTP server on port %s", s.Port))
	s.router.Use(s.middlewares...)
	loadPredefinedRoutes(s.router)
	loadRoutes(s.router, s.routes)
	s.http = &http.Server{
		Addr:    fmt.Sprintf(":%s", s.Port),
		Handler: s.router,
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

// POST registers HTTP server endpoint with Post method .
func (s *Server) POST(relativePath string, handlers ...gin.HandlerFunc) {
	s.routes = append(s.routes, route{
		method:       http.MethodPost,
		relativePath: relativePath,
		handlers:     handlers,
	})
}

// GET registers HTTP server endpoint with Get method.
func (s *Server) GET(relativePath string, handlers ...gin.HandlerFunc) {
	s.routes = append(s.routes, route{
		method:       http.MethodGet,
		relativePath: relativePath,
		handlers:     handlers,
	})
}

// DELETE registers HTTP server endpoint with Delete method.
func (s *Server) DELETE(relativePath string, handlers ...gin.HandlerFunc) {
	s.routes = append(s.routes, route{
		method:       http.MethodDelete,
		relativePath: relativePath,
		handlers:     handlers,
	})
}

// PATCH registers HTTP server endpoint with Patch method.
func (s *Server) PATCH(relativePath string, handlers ...gin.HandlerFunc) {
	s.routes = append(s.routes, route{
		method:       http.MethodPatch,
		relativePath: relativePath,
		handlers:     handlers,
	})
}

// PUT registers HTTP server endpoint with Put method.
func (s *Server) PUT(relativePath string, handlers ...gin.HandlerFunc) {
	s.routes = append(s.routes, route{
		method:       http.MethodPut,
		relativePath: relativePath,
		handlers:     handlers,
	})
}

// OPTIONS registers HTTP server endpoint with Options method.
func (s *Server) OPTIONS(relativePath string, handlers ...gin.HandlerFunc) {
	s.routes = append(s.routes, route{
		method:       http.MethodOptions,
		relativePath: relativePath,
		handlers:     handlers,
	})
}

// HEAD registers HTTP server endpoint with Head method.
func (s *Server) HEAD(relativePath string, handlers ...gin.HandlerFunc) {
	s.routes = append(s.routes, route{
		method:       http.MethodHead,
		relativePath: relativePath,
		handlers:     handlers,
	})
}
