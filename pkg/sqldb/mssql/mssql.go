package mssql

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/satriajidam/go-gin-skeleton/pkg/sqldb"

	// Import Microsoft SQL Server driver.
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

// NewConnection initiates new connection to a Microsoft SQL Server database using provided
// connection configs.
func NewConnection(conf sqldb.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s?database=%s&%s",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database,
		conf.Params,
	)

	dbconn, err := gorm.Open("mssql", dsn)
	if err != nil {
		return nil, err
	}

	dbconn.DB().SetMaxIdleConns(conf.MaxIdleConns)
	dbconn.DB().SetMaxOpenConns(conf.MaxOpenConns)

	dbconn.SingularTable(conf.SingularTable)
	dbconn.LogMode(conf.DebugMode)

	return dbconn, nil
}
