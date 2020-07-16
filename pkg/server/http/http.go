package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/satriajidam/go-gin-skeleton/pkg/log"
	"github.com/satriajidam/go-gin-skeleton/pkg/server/http/middleware/logger"
	"github.com/satriajidam/go-gin-skeleton/pkg/server/http/middleware/requestid"
)

var (
	CORSDefaultAllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"}
	// CORS safelisted-request-header: https://fetch.spec.whatwg.org/#cors-safelisted-request-header
	// CORS forbidden-header-name: https://fetch.spec.whatwg.org/#forbidden-header-name
	CORSDefaultAllowHeaders = []string{
		"Accept",
		"Accept-Charset",
		"Accept-Encoding",
		"Accept-Language",
		"Content-Language",
		"Content-Length",
		"Content-Type",
		"Host",
		"Origin",
	}
	CORSDefaultAllowCredentials = true
	CORSDefaultMaxAge           = 12 * time.Hour
)

// Server represents the implementation of HTTP server object.
type Server struct {
	RouterGroup
	http         *http.Server
	router       *gin.Engine
	loggerConfig *logger.Config
	middlewares  []gin.HandlerFunc
	routes       []route
	enableCORS   bool
	CORS         *cors.Config
	Port         string
}

type route struct {
	method       string
	relativePath string
	handlers     []gin.HandlerFunc
}

// NewServer creates new HTTP server.
func NewServer(port string, enableCORS bool, enablePredefinedRoutes bool) *Server {
	routes := []route{}

	if enablePredefinedRoutes {
		routes = append(routes, predefinedRoutes...)
	}

	server := &Server{
		router: gin.New(),
		middlewares: []gin.HandlerFunc{
			// Default gin middlewares.
			gin.Recovery(),
			requestid.New(),
		},
		loggerConfig: &logger.Config{
			Stdout:   log.Stdout(),
			Stderr:   log.Stderr(),
			SkipPath: []string{},
		},
		enableCORS: enableCORS,
		CORS: &cors.Config{
			AllowCredentials: CORSDefaultAllowCredentials,
		},
		routes: routes,
		Port:   port,
	}

	server.RouterGroup = RouterGroup{
		server: server,
	}

	return server
}

// AddMiddleware adds a gin middleware the HTTP server.
func (s *Server) AddMiddleware(h gin.HandlerFunc) {
	s.middlewares = append(s.middlewares, h)
}

// LoggerSkipPaths registers endpoint paths that you want to skip from being logged
// by the logger middleware.
func (s *Server) LoggerSkipPaths(paths ...string) {
	s.loggerConfig.SkipPath = append(s.loggerConfig.SkipPath, paths...)
}

// GetRoutePaths retrieves all route paths registerd to this HTTP server.
func (s *Server) GetRoutePaths() []string {
	paths := []string{}
	for _, r := range s.routes {
		paths = append(paths, r.relativePath)
	}
	return paths
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

func (s *Server) setupCORS() {
	if s.enableCORS {
		if len(s.CORS.AllowOrigins) <= 0 {
			s.CORS.AllowAllOrigins = true
		} else {
			s.CORS.AllowAllOrigins = false
		}
		if len(s.CORS.AllowMethods) <= 0 {
			s.CORS.AllowMethods = CORSDefaultAllowMethods
		}
		if len(s.CORS.AllowHeaders) <= 0 {
			s.CORS.AllowHeaders = CORSDefaultAllowHeaders
		}
		if s.CORS.MaxAge <= 0 {
			s.CORS.MaxAge = CORSDefaultMaxAge
		}
		fmt.Printf("%+v\n", s.CORS)
		s.AddMiddleware(cors.New(*s.CORS))
	}
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	log.Info(fmt.Sprintf("Start HTTP server on port %s", s.Port))
	s.setupCORS()
	s.AddMiddleware(logger.New(s.Port, *s.loggerConfig))
	s.router.Use(s.middlewares...)
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

// POST registers HTTP server endpoint with Post method.
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

// RouterGroup groups path under one path prefix.
type RouterGroup struct {
	prefix string
	server *Server
}

// Group creates new RouterGroup with the given path prefix.
func (rg *RouterGroup) Group(prefix string) *RouterGroup {
	return &RouterGroup{
		prefix: prefix,
		server: rg.server,
	}
}

// POST registers HTTP server endpoint with Post method.
func (rg *RouterGroup) POST(relativePath string, handlers ...gin.HandlerFunc) {
	rg.server.POST(rg.prefix+relativePath, handlers...)
}

// GET registers HTTP server endpoint with Get method.
func (rg *RouterGroup) GET(relativePath string, handlers ...gin.HandlerFunc) {
	rg.server.GET(rg.prefix+relativePath, handlers...)
}

// DELETE registers HTTP server endpoint with Delete method.
func (rg *RouterGroup) DELETE(relativePath string, handlers ...gin.HandlerFunc) {
	rg.server.DELETE(rg.prefix+relativePath, handlers...)
}

// PATCH registers HTTP server endpoint with Patch method.
func (rg *RouterGroup) PATCH(relativePath string, handlers ...gin.HandlerFunc) {
	rg.server.PATCH(rg.prefix+relativePath, handlers...)
}

// PUT registers HTTP server endpoint with Put method.
func (rg *RouterGroup) PUT(relativePath string, handlers ...gin.HandlerFunc) {
	rg.server.PUT(rg.prefix+relativePath, handlers...)
}

// OPTIONS registers HTTP server endpoint with Options method.
func (rg *RouterGroup) OPTIONS(relativePath string, handlers ...gin.HandlerFunc) {
	rg.server.OPTIONS(rg.prefix+relativePath, handlers...)
}

// HEAD registers HTTP server endpoint with Head method.
func (rg *RouterGroup) HEAD(relativePath string, handlers ...gin.HandlerFunc) {
	rg.server.HEAD(rg.prefix+relativePath, handlers...)
}
