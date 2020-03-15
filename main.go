package main

import (
	"github.com/satriajidam/go-gin-skeleton/pkg/database/mysql"
	"github.com/satriajidam/go-gin-skeleton/pkg/log"
	"github.com/satriajidam/go-gin-skeleton/pkg/server"
	"github.com/satriajidam/go-gin-skeleton/pkg/server/http"
)

func main() {
	dbconn, err := mysql.Connect()
	if err != nil {
		log.Panic(err, "failed connecting to mysql database")
		// panic(err)
	}

	defer dbconn.Close()

	httpServer := http.New()

	err = server.InitServers(httpServer)
	if err != nil {
		panic(err)
	}
}
