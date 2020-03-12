package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/satriajidam/go-gin-skeleton/pkg/config"

	// Import PostgreSQL driver.
	_ "github.com/jinzhu/gorm/dialects/postgres"
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

// Engine sets to PostgreSQL.
const Engine = "postgres"

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

	db, err := gorm.Open(Engine, dsn)
	if err != nil {
		return nil, err
	}

	switch config.Get().AppMode {
	case config.ReleaseMode:
		db.LogMode(false)
	case config.DebugMode:
	default:
		db.LogMode(true)
	}

	db.DB().SetMaxIdleConns(cfg.MaxIdleConns)
	db.DB().SetMaxOpenConns(cfg.MaxOpenConns)

	db.SingularTable(cfg.SingularTable)

	return db, nil
}
