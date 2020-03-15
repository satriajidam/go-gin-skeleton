package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/satriajidam/go-gin-skeleton/pkg/config"

	// Import PostgreSQL driver.
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Connect initiates connection to a PostgreSQL database.
func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?%s",
		config.Get().PostgresUsername,
		config.Get().PostgresPassword,
		config.Get().PostgresHost,
		config.Get().PostgresPort,
		config.Get().PostgresDatabase,
		config.Get().PostgresParams,
	)

	dbconn, err := gorm.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	dbconn.DB().SetMaxIdleConns(config.Get().GormMaxIdleConns)
	dbconn.DB().SetMaxOpenConns(config.Get().GormMaxOpenConns)

	dbconn.SingularTable(config.Get().GormSingularTable)
	dbconn.LogMode(config.IsDebugMode())

	return dbconn, nil
}
