package main

import (
	"fmt"

	"github.com/satriajidam/go-gin-skeleton/pkg/config"
	"github.com/satriajidam/go-gin-skeleton/pkg/database"
	"github.com/satriajidam/go-gin-skeleton/pkg/database/mysql"
	"github.com/satriajidam/go-gin-skeleton/pkg/log"
	"github.com/satriajidam/go-gin-skeleton/pkg/server"
	"github.com/satriajidam/go-gin-skeleton/pkg/server/http"
)

func main() {
	appConfig := config.Init()

	dbengine := mysql.Init(appConfig)
	dbconn, err := database.Connect(dbengine)
	if err != nil {
		log.Panic(err, fmt.Sprintf("failed connecting to %s database", dbengine.GetName()))
	}

	defer dbconn.Close()

	httpServer := http.New(appConfig)

	err = server.InitServers(httpServer)
	if err != nil {
		log.Panic(err, fmt.Sprintf("failed to start server(s)"))
	}
}
