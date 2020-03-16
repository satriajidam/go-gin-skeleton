package main

import (
	"github.com/satriajidam/go-gin-skeleton/pkg/database/sqlite"
	"github.com/satriajidam/go-gin-skeleton/pkg/server"
	"github.com/satriajidam/go-gin-skeleton/pkg/server/http"
)

func main() {
	dbconn, err := sqlite.Connect()
	if err != nil {
		panic(err)
	}

	defer dbconn.Close()

	httpServer := http.New()

	if err := <-server.StartServers(httpServer); err != nil {
		panic(err)
	}
}
