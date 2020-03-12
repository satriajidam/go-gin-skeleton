package mysql

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/satriajidam/go-gin-skeleton/pkg/config"
	"github.com/satriajidam/go-gin-skeleton/pkg/database"

	// Import MySQL driver.
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// mysql stores MySQL database configurations.
type mysql struct {
	Host          string
	Port          string
	Username      string
	Password      string
	Database      string
	Params        string
	MaxIdleConns  int
	MaxOpenConns  int
	SingularTable bool
}

// Init initializes MySQL database engine.
func Init() database.DBEngine {
	return &mysql{
		Host:          config.Get().MySQLHost,
		Port:          config.Get().MySQLPort,
		Username:      config.Get().MySQLUsername,
		Password:      config.Get().MySQLPassword,
		Database:      config.Get().MySQLPassword,
		Params:        config.Get().MySQLParams,
		MaxIdleConns:  config.Get().GormMaxIdleConns,
		MaxOpenConns:  config.Get().GormMaxOpenConns,
		SingularTable: config.Get().GormSingularTable,
	}
}

// GetName returns MySQL database engine name.
func (db *mysql) GetName() string {
	return "mysql"
}

// Connect initiates connection to a MysQL database.
func (db *mysql) Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?%s",
		db.Username,
		db.Password,
		db.Host,
		db.Port,
		db.Database,
		db.Params,
	)

	dbconn, err := gorm.Open(db.GetName(), dsn)
	if err != nil {
		return nil, err
	}

	switch config.Get().AppMode {
	case config.ReleaseMode:
		dbconn.LogMode(false)
	case config.DebugMode:
	default:
		dbconn.LogMode(true)
	}

	dbconn.DB().SetMaxIdleConns(db.MaxIdleConns)
	dbconn.DB().SetMaxOpenConns(db.MaxOpenConns)

	dbconn.SingularTable(db.SingularTable)

	return dbconn, nil
}
