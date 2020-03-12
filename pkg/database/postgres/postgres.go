package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/satriajidam/go-gin-skeleton/pkg/config"

	// Import PostgreSQL driver.
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Connect initiates connection to a PostgreSQL database.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?%s",
		cfg.PostgresUsername,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDatabase,
		cfg.PostgresParams,
	)

	dbconn, err := gorm.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	dbconn.DB().SetMaxIdleConns(cfg.GormMaxIdleConns)
	dbconn.DB().SetMaxOpenConns(cfg.GormMaxOpenConns)

	dbconn.SingularTable(cfg.GormSingularTable)
	dbconn.LogMode(cfg.IsDebugMode())

	return dbconn, nil
}
