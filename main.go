package main

import (
	"github.com/satriajidam/go-gin-skeleton/pkg/database/mysql"
	"github.com/satriajidam/go-gin-skeleton/pkg/log"
)

func main() {
	db, err := mysql.DB()
	if err != nil {
		db.Close()
		log.Error(err)
		panic(err)
	}
}
