package sqlite

import (
	"github.com/jinzhu/gorm"
	"github.com/satriajidam/go-gin-skeleton/pkg/config"

	// Import SQLite driver.
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Connect initiates connection to a SQLite database.
func Connect() (*gorm.DB, error) {
	dbconn, err := gorm.Open("sqlite3", config.Get().SQLiteDatabase)
	if err != nil {
		return nil, err
	}

	dbconn.DB().SetMaxIdleConns(config.Get().GormMaxIdleConns)
	dbconn.DB().SetMaxOpenConns(config.Get().GormMaxOpenConns)

	dbconn.SingularTable(config.Get().GormSingularTable)
	dbconn.LogMode(config.IsDebugMode())

	return dbconn, nil
}
