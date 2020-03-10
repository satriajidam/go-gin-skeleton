package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/satriajidam/go-gin-skeleton/pkg/config"
)

func initLogger(config config.Config) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	switch config.LogLevel {
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
}

func initDBConnection(config config.Config) (*gorm.DB, error) {
	mode := config.DBMode
	host := config.DBHost
	port := config.DBPort
	user := config.DBUser
	pass := config.DBPass
	name := config.DBName
	params := config.DBParams
	dsn := ""

	switch mode {
	case "sqlite3":
		// TODO
	case "mysql":
		if params == "" {
			params = "charset=utf8&parseTime=True&loc=Local"
		}

		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", user, pass, host, port, name, params)
	case "postgres":
		if params == "" {
			params = "sslmode=disable"
		}

		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s", user, pass, host, port, name, params)
	case "mssql":
		if params == "" {
			params = "connection+timeout=0&dial+timeout=15&encrypt=false"
		}

		dsn = fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s&%s", user, pass, host, port, name, params)
	default:
		return nil, fmt.Errorf("unknown database mode: %s", mode)
	}

	db, err := gorm.Open(mode, dsn)
	if err != nil {
		return db, err
	}

	db.LogMode(config.GormEnableLog)

	return db, nil
}

func main() {
	config := config.Get()

	initLogger(config)

	db, err := initDBConnection(config)
	if err != nil {
		db.Close()
		log.Print(err.Error())
		panic(err)
	}

}
