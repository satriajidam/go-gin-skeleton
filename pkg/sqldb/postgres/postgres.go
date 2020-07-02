package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/satriajidam/go-gin-skeleton/pkg/sqldb"

	// Import PostgreSQL driver.
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// NewConnection initiates new connection to a PostgreSQL database using provided
// connection configs.
func NewConnection(conf sqldb.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?%s",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database,
		conf.Params,
	)

	dbconn, err := gorm.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	dbconn.DB().SetMaxIdleConns(conf.MaxIdleConns)
	dbconn.DB().SetMaxOpenConns(conf.MaxOpenConns)

	dbconn.SingularTable(conf.SingularTable)
	dbconn.LogMode(conf.DebugMode)

	return dbconn, nil
}
