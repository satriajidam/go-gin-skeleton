package sqlite

import (
	"github.com/jinzhu/gorm"
	"github.com/satriajidam/go-gin-skeleton/pkg/database/sql"

	// Import SQLite driver.
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// NewConnection creates a new connection to an SQLite database using provided
// connection configs.
func NewConnection(conf sql.DBConfig) (*sql.Connection, error) {
	db, err := gorm.Open("sqlite3", conf.Database)
	if err != nil {
		return nil, err
	}

	db.DB().SetMaxIdleConns(conf.MaxIdleConns)
	db.DB().SetMaxOpenConns(conf.MaxOpenConns)

	db.SingularTable(conf.SingularTable)
	db.LogMode(conf.DebugMode)

	return sql.NewConnection(db, conf.Host, conf.Port, "SQLite"), nil
}
