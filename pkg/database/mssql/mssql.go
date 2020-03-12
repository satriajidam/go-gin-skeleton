package mssql

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/satriajidam/go-gin-skeleton/pkg/config"

	// Import Microsoft SQL Server driver.
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

// Connect initiates connection to a Microsoft SQL Server database.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s?database=%s&%s",
		cfg.MSSQLUsername,
		cfg.MSSQLPassword,
		cfg.MSSQLHost,
		cfg.MSSQLPort,
		cfg.MSSQLDatabase,
		cfg.MSSQLParams,
	)

	dbconn, err := gorm.Open("mssql", dsn)
	if err != nil {
		return nil, err
	}

	dbconn.DB().SetMaxIdleConns(cfg.GormMaxIdleConns)
	dbconn.DB().SetMaxOpenConns(cfg.GormMaxOpenConns)

	dbconn.SingularTable(cfg.GormSingularTable)
	dbconn.LogMode(cfg.IsDebugMode())

	return dbconn, nil
}
