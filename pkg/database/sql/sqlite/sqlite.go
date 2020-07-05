package sqlite

import (
	"github.com/jinzhu/gorm"
	"github.com/satriajidam/go-gin-skeleton/pkg/database/sql"

	// Import SQLite driver.
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// NewConnection creates a new connection to an SQLite database using provided
// connection configs.
func NewConnection(conf sql.DBConfig) (*gorm.DB, error) {
	dbconn, err := gorm.Open("sqlite3", conf.Database)
	if err != nil {
		return nil, err
	}

	dbconn.DB().SetMaxIdleConns(conf.MaxIdleConns)
	dbconn.DB().SetMaxOpenConns(conf.MaxOpenConns)

	dbconn.SingularTable(conf.SingularTable)
	dbconn.LogMode(conf.DebugMode)

	return dbconn, nil
}
