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
	defer dbconn.Close()
	if err != nil {
		panic(err)
	}

	httpServer := http.NewServer(cfg.HTTPServerPort, true)

	promServer := prometheus.NewServer(
		cfg.PrometheusServerPort,
		cfg.PrometheusServerMetricsPath,
	)

	promServer.Monitor(
		&prometheus.Target{
			HTTPServer:    httpServer,
			GroupedStatus: true,
		},
	)

	providerRepository := provider.NewRepository(dbconn)
	providerService := provider.NewService(providerRepository)
	providerHTTPHandler := provider.NewHTTPHandler(providerService)

	httpServer.POST("/v1/provider", providerHTTPHandler.CreateProvider)
	httpServer.PUT("/v1/provider/:uuid", providerHTTPHandler.UpdateProvider)
	httpServer.DELETE("/v1/provider/:uuid", providerHTTPHandler.DeleteProviderByUUID)
	httpServer.GET("/v1/provider/:uuid", providerHTTPHandler.GetProviderByUUID)
	httpServer.GET("/v1/providers", providerHTTPHandler.ListProviders)

	server.RunServersGracefully(cfg.GracefulTimeout, promServer, httpServer)
}
