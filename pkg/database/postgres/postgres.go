package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/satriajidam/go-gin-skeleton/pkg/config"
	"github.com/satriajidam/go-gin-skeleton/pkg/log"
)

// Config stores PostgreSQL database configurations.
type Config struct {
	Host          string
	Port          string
	Username      string
	Password      string
	Database      string
	Params        string
	MaxIdleConns  int
	MaxOpenConns  int
	SingularTable bool
}

var cfg *Config

func init() {
	cfg = &Config{
		Host:          config.Get().PostgresHost,
		Port:          config.Get().PostgresPort,
		Username:      config.Get().PostgresUsername,
		Password:      config.Get().PostgresPassword,
		Database:      config.Get().PostgresDatabase,
		Params:        config.Get().PostgresParams,
		MaxIdleConns:  config.Get().GormMaxIdleConns,
		MaxOpenConns:  config.Get().GormMaxOpenConns,
		SingularTable: config.Get().GormSingularTable,
	}
}

// DB returns PostgreSQL database connection.
func DB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Params,
	)

	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	logLevel := config.Get().AppLogLevel

	if logLevel == log.LevelDebug || logLevel == log.LevelTrace {
		db.LogMode(true)
	} else {
		db.LogMode(false)
	}

	db.DB().SetMaxIdleConns(cfg.MaxIdleConns)
	db.DB().SetMaxOpenConns(cfg.MaxOpenConns)

	db.SingularTable(cfg.SingularTable)

	return db, nil
}
