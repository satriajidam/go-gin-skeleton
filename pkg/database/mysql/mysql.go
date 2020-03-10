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
			Host: config.Get().MySQLHost,
			Port: config.Get().MySQLPort,
		}
	}
}
