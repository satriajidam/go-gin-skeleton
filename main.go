package main

import (
	"fmt"

	"github.com/satriajidam/go-gin-skeleton/pkg/database"
	"github.com/satriajidam/go-gin-skeleton/pkg/database/postgres"
	"github.com/satriajidam/go-gin-skeleton/pkg/log"
)

func main() {
	dbengine := postgres.Init()
	dbconn, err := database.Connect(dbengine)
	if err != nil {
		log.Panic(err, fmt.Sprintf("failed connecting to %s database", dbengine.GetName()))
	}

	defer dbconn.Close()
}
