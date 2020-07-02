package mysql

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/satriajidam/go-gin-skeleton/pkg/database/sql"

	// Import MySQL driver.
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// NewConnection initiates new connection to a MySQL database using provided
// connection configs.
func NewConnection(conf sql.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?%s",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database,
		conf.Params,
	)

	dbconn, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	dbconn.DB().SetMaxIdleConns(conf.MaxIdleConns)
	dbconn.DB().SetMaxOpenConns(conf.MaxOpenConns)

	dbconn.SingularTable(conf.SingularTable)
	dbconn.LogMode(conf.DebugMode)

	return dbconn, nil
}
