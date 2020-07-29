package mysql

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/satriajidam/go-gin-skeleton/pkg/database/sql"

	// Import MySQL driver.
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// NewConnection creates a new connection to a MySQL database using provided
// connection configs.
func NewConnection(conf sql.DBConfig) (*sql.Connection, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?%s",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database,
		conf.Params,
	)

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.DB().SetMaxIdleConns(conf.MaxIdleConns)
	db.DB().SetMaxOpenConns(conf.MaxOpenConns)

	db.SingularTable(conf.SingularTable)
	db.LogMode(conf.DebugMode)

	return sql.NewConnection(db, conf.Host, conf.Port, "MySQL"), nil
}
