package main

import (
	"github.com/satriajidam/go-gin-skeleton/pkg/config"
	"github.com/satriajidam/go-gin-skeleton/pkg/database/sql"
	"github.com/satriajidam/go-gin-skeleton/pkg/database/sql/sqlite"
	"github.com/satriajidam/go-gin-skeleton/pkg/server"
	"github.com/satriajidam/go-gin-skeleton/pkg/server/http"
	"github.com/satriajidam/go-gin-skeleton/pkg/server/prometheus"
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

	httpServer := http.NewServer(
		cfg.HTTPServerPort,
		cfg.HTTPServerDisallowUnknownJSONFields,
	)

	promServer := prometheus.NewServer(
		cfg.PrometheusServerPort,
		cfg.PrometheusServerMetricsSubsystem,
		[]string{},
	)

	promServer.Monitor(httpServer.Router, cfg.PrometheusServerMetricsPath)

	server.RunServersGracefully(cfg.GracefulTimeout, httpServer, promServer)
}
