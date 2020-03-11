package main

import (
	"fmt"

	"github.com/satriajidam/go-gin-skeleton/pkg/database/mysql"
	"github.com/satriajidam/go-gin-skeleton/pkg/log"
)

func main() {
	db, err := mysql.DB()
	if err != nil {
		log.Panic(err, fmt.Sprintf("failed connecting to %s database", mysql.Engine))
	}

	defer db.Close()
}
