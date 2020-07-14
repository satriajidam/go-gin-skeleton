package main

import (
	"github.com/satriajidam/go-gin-skeleton/pkg/config"
	"github.com/satriajidam/go-gin-skeleton/pkg/database/sql"
	"github.com/satriajidam/go-gin-skeleton/pkg/database/sql/sqlite"
	"github.com/satriajidam/go-gin-skeleton/pkg/server"
	"github.com/satriajidam/go-gin-skeleton/pkg/server/http"
	"github.com/satriajidam/go-gin-skeleton/pkg/server/prometheus"
	"github.com/satriajidam/go-gin-skeleton/pkg/service/provider"
)

func main() {
	cfg := config.Get()

	dbconn, err := sqlite.NewConnection(sql.DBConfig{
		Database:      cfg.SQLiteDatabase,
		MaxIdleConns:  cfg.SQLiteMaxIdleConns,
		MaxOpenConns:  cfg.SQLiteMaxOpenConns,
		SingularTable: cfg.SQLiteSingularTable,
		DebugMode:     cfg.SQLiteDebugMode,
	})
	defer func() {
		err := dbconn.Close()
		if err != nil {
			panic(err)
		}
	}()
	if err != nil {
		panic(err)
	}

	excludedMonitoringPaths := []string{"/_/health"}

	httpServer := http.NewServer(cfg.HTTPServerPort, true, excludedMonitoringPaths...)

	providerRepository := provider.NewRepository(dbconn, true)
	providerService := provider.NewService(providerRepository)
	providerHTTPHandler := provider.NewHTTPHandler(providerService)

	v1 := httpServer.Group("/v1")
	v1.POST("/provider", providerHTTPHandler.CreateProvider)
	v1.PUT("/provider/:uuid", providerHTTPHandler.UpdateProvider)
	v1.DELETE("/provider/:uuid", providerHTTPHandler.DeleteProviderByUUID)
	v1.GET("/provider/:uuid", providerHTTPHandler.GetProviderByUUID)
	v1.GET("/providers", providerHTTPHandler.ListProviders)

	promServer := prometheus.NewServer(
		cfg.PrometheusServerPort,
		cfg.PrometheusServerMetricsPath,
	)

	promServer.Monitor(
		&prometheus.Target{
			HTTPServer:   httpServer,
			ExcludePaths: excludedMonitoringPaths,
		},
	)

	server.RunServersGracefully(cfg.GracefulTimeout, promServer, httpServer)
}
