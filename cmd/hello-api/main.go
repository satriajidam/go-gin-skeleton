package main

import (
	"github.com/satriajidam/go-gin-skeleton/internal/config"
	"github.com/satriajidam/go-gin-skeleton/pkg/server"
	"github.com/satriajidam/go-gin-skeleton/pkg/server/http"
	"github.com/satriajidam/go-gin-skeleton/pkg/server/prometheus"
)

func main() {
	cfg := config.Get()

	httpServer := http.NewServer(
		cfg.HTTPServerPort,
		cfg.HTTPServerEnableCORS,
		cfg.HTTPServerEnablePredefinedRoutes,
	)
	httpServer.CORS.AllowOrigins = cfg.HTTPServerAllowOrigins
	httpServer.CORS.AllowMethods = cfg.HTTPServerAllowMethods
	httpServer.CORS.AllowHeaders = cfg.HTTPServerAllowHeaders
	httpServer.CORS.MaxAge = cfg.HTTPServerMaxAge

	promServer := prometheus.NewServer(
		cfg.PrometheusServerPort,
		cfg.PrometheusServerMetricsPath,
	)

	promServer.Monitor(
		&prometheus.Target{
			HTTPServer:    httpServer,
			ExcludePaths:  cfg.HTTPServerMonitorSkipPaths,
			GroupedStatus: cfg.HTTPServerMonitorGroupedStatus,
		},
	)

	server.RunServersGracefully(cfg.GracefulTimeout, promServer, httpServer)
}
