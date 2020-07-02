package main

import (
	"github.com/satriajidam/go-gin-skeleton/pkg/config"
	"github.com/satriajidam/go-gin-skeleton/pkg/server"
	"github.com/satriajidam/go-gin-skeleton/pkg/server/http"
	"github.com/satriajidam/go-gin-skeleton/pkg/sqldb"
	"github.com/satriajidam/go-gin-skeleton/pkg/sqldb/sqlite"
)

func main() {
	dbconn, err := sqlite.NewConnection(sqldb.Config{
		Database:      config.Get().SQLiteDatabase,
		MaxIdleConns:  config.Get().SQLiteMaxIdleConns,
		MaxOpenConns:  config.Get().SQLiteMaxOpenConns,
		SingularTable: config.Get().SQLiteSingularTable,
		DebugMode:     config.Get().SQLiteDebugMode,
	})
	if err != nil {
		panic(err)
	}

	defer dbconn.Close()

	httpServer := http.NewServer(
		config.Get().HTTPServerPort,
		config.Get().HTTPServerMode,
		config.Get().HTTPServerDisallowUnknownJSONFields,
	)

	if err := <-server.StartServers(httpServer); err != nil {
		panic(err)
	}
}
