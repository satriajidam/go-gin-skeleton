package main

import (
	"github.com/satriajidam/go-gin-skeleton/pkg/config"
	"github.com/satriajidam/go-gin-skeleton/pkg/database/mysql"
	"github.com/satriajidam/go-gin-skeleton/pkg/server"
	"github.com/satriajidam/go-gin-skeleton/pkg/server/http"
)

func main() {
	appConfig := config.Init()

	dbconn, err := mysql.Connect(appConfig)
	if err != nil {
		panic(err)
	}

	defer dbconn.Close()

	httpServer := http.New(appConfig)

	err = server.InitServers(httpServer)
	if err != nil {
		panic(err)
	}
}
