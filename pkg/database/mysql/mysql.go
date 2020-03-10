package mysql

import "github.com/satriajidam/go-gin-skeleton/pkg/config"

// Config ...
type Config struct {
	Host         string
	Port         string
	Username     string
	Password     string
	Database     string
	Params       string
	MaxIdleConns int
	MaxOpenConns int
}

var cfg *Config

func init() {
	if config.Get().MySQLEnable {
		cfg = Config{
			Host:         config.Get().MySQLHost,
			Port:         config.Get().MySQLPort,
			Username:     config.Get().MySQLUsername,
			Password:     config.Get().MySQLPassword,
			Database:     config.Get().MySQLPassword,
			Params:       config.Get().MySQLParams,
			MaxIdleConns: config.Get().GormMaxIdleConns,
			MaxOpenConns: config.Get().GormMaxOpenConns,
		}
	}
}
