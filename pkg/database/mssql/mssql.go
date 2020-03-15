package mssql

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/satriajidam/go-gin-skeleton/pkg/config"

	// Import Microsoft SQL Server driver.
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

// Connect initiates connection to a Microsoft SQL Server database.
func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s?database=%s&%s",
		config.Get().MSSQLUsername,
		config.Get().MSSQLPassword,
		config.Get().MSSQLHost,
		config.Get().MSSQLPort,
		config.Get().MSSQLDatabase,
		config.Get().MSSQLParams,
	)

	dbconn, err := gorm.Open("mssql", dsn)
	if err != nil {
		return nil, err
	}

	dbconn.DB().SetMaxIdleConns(config.Get().GormMaxIdleConns)
	dbconn.DB().SetMaxOpenConns(config.Get().GormMaxOpenConns)

	dbconn.SingularTable(config.Get().GormSingularTable)
	dbconn.LogMode(config.IsDebugMode())

	return dbconn, nil
}
