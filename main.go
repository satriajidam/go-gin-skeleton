package main

import (
	"fmt"

	"github.com/satriajidam/go-gin-skeleton/pkg/config"
	"github.com/satriajidam/go-gin-skeleton/pkg/database"
	"github.com/satriajidam/go-gin-skeleton/pkg/database/mysql"
	"github.com/satriajidam/go-gin-skeleton/pkg/log"
	"github.com/satriajidam/go-gin-skeleton/pkg/server/http"
)

func main() {
	dbengine := mysql.Init()
	dbconn, err := database.Connect(dbengine)
	if err != nil {
		log.Panic(err, fmt.Sprintf("failed connecting to %s database", dbengine.GetName()))
	}

	defer dbconn.Close()

	httpServer := http.New(config.Get().AppMode, config.Get().GinDisallowUnknownJSONFields)
	httpServer.Start(config.Get().HTTPPort)
}
