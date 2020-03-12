package mysql

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/satriajidam/go-gin-skeleton/pkg/config"

	// Import MySQL driver.
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Config stores MySQL database configurations.
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

// Engine sets to MySQL.
const Engine = "mysql"

var cfg *Config

func init() {
	cfg = &Config{
		Host:          config.Get().MySQLHost,
		Port:          config.Get().MySQLPort,
		Username:      config.Get().MySQLUsername,
		Password:      config.Get().MySQLPassword,
		Database:      config.Get().MySQLPassword,
		Params:        config.Get().MySQLParams,
		MaxIdleConns:  config.Get().GormMaxIdleConns,
		MaxOpenConns:  config.Get().GormMaxOpenConns,
		SingularTable: config.Get().GormSingularTable,
	}
}

// DB returns MysQL database connection.
func DB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?%s",
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
