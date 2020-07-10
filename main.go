package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
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

	httpServer := http.NewServer(cfg.HTTPServerPort, true)

	httpServer.GET("/provider/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")
		format, _ := ctx.GetQuery("format")
		switch format {
		case "yaml":
			ctx.YAML(200, map[string]string{"cloudProvider": name})
		case "json":
			ctx.JSON(200, map[string]string{"cloudProvider": name})
		default:
			ctx.String(200, fmt.Sprintf("provider: %s", name))
		}
	})

	httpServer.GET("/highlatency", func(ctx *gin.Context) {
		min, max := 0, 10
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Duration(rand.Intn(max-min+1)+5) * time.Second)
		ctx.String(200, fmt.Sprint("OK"))
	})

	promServer := prometheus.NewServer(
		cfg.PrometheusServerPort,
		cfg.PrometheusServerMetricsPath,
	)

	promServer.Monitor(
		&prometheus.Target{
			HTTPServer:         httpServer,
			GroupedStatus:      true,
			DisableMeasureSize: true,
		},
	)

	server.RunServersGracefully(cfg.GracefulTimeout, promServer, httpServer)
}
