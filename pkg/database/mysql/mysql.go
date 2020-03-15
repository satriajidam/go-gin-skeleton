package mysql

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/satriajidam/go-gin-skeleton/pkg/config"

	// Import MySQL driver.
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Connect initiates connection to a MysQL database.
func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?%s",
		config.Get().MySQLUsername,
		config.Get().MySQLPassword,
		config.Get().MySQLHost,
		config.Get().MySQLPort,
		config.Get().MySQLDatabase,
		config.Get().MySQLParams,
	)

	dbconn, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	dbconn.DB().SetMaxIdleConns(config.Get().GormMaxIdleConns)
	dbconn.DB().SetMaxOpenConns(config.Get().GormMaxOpenConns)

	dbconn.SingularTable(config.Get().GormSingularTable)
	dbconn.LogMode(config.IsDebugMode())

	return dbconn, nil
}
