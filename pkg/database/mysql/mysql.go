package mysql

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/satriajidam/go-gin-skeleton/pkg/config"

	// Import MySQL driver.
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Connect initiates connection to a MysQL database.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?%s",
		cfg.MySQLUsername,
		cfg.MySQLPassword,
		cfg.MySQLHost,
		cfg.MySQLPort,
		cfg.MySQLDatabase,
		cfg.MySQLParams,
	)

	dbconn, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	dbconn.DB().SetMaxIdleConns(cfg.GormMaxIdleConns)
	dbconn.DB().SetMaxOpenConns(cfg.GormMaxOpenConns)

	dbconn.SingularTable(cfg.GormSingularTable)
	dbconn.LogMode(cfg.IsDebugMode())

	return dbconn, nil
}
