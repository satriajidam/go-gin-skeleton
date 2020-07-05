package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/satriajidam/go-gin-skeleton/pkg/config"
	"github.com/satriajidam/go-gin-skeleton/pkg/database/sql"
	"github.com/satriajidam/go-gin-skeleton/pkg/database/sql/sqlite"
	"github.com/satriajidam/go-gin-skeleton/pkg/server"
	"github.com/satriajidam/go-gin-skeleton/pkg/server/http"
)

func main() {
	dbconn, err := sqlite.NewConnection(sql.DBConfig{
		Database:      config.Get().SQLiteDatabase,
		MaxIdleConns:  config.Get().SQLiteMaxIdleConns,
		MaxOpenConns:  config.Get().SQLiteMaxOpenConns,
		SingularTable: config.Get().SQLiteSingularTable,
		DebugMode:     config.Get().SQLiteDebugMode,
	})
	defer dbconn.Close()
	if err != nil {
		panic(err)
	}

	httpServer := http.NewServer(
		config.Get().HTTPServerPort,
		config.Get().HTTPServerMode,
		config.Get().HTTPServerDisallowUnknownJSONFields,
	)

	if err := <-server.StartServers(httpServer); err != nil {
		panic(err)
	}

	// Graceful shutdown:
	// - https://chenyitian.gitbooks.io/gin-web-framework/docs/38.html
	// - https://medium.com/honestbee-tw-engineer/gracefully-shutdown-in-go-http-server-5f5e6b83da5a
	// Wait for interrupt signal to gracefully shutdown the server.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit

	// Set graceful shutdown timeout to configured seconds.
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(config.Get().GracefulTimeout)*time.Second,
	)
	defer cancel()

	server.StopServers(ctx, httpServer)
}
