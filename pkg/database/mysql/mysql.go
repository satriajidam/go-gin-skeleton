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
	LogMode       bool
}

// Init initializes MySQL database engine.
func Init(cfg *config.Config) database.DBEngine {
	return &mysql{
		Host:          cfg.MySQLHost,
		Port:          cfg.MySQLPort,
		Username:      cfg.MySQLUsername,
		Password:      cfg.MySQLPassword,
		Database:      cfg.MySQLPassword,
		Params:        cfg.MySQLParams,
		MaxIdleConns:  cfg.GormMaxIdleConns,
		MaxOpenConns:  cfg.GormMaxOpenConns,
		SingularTable: cfg.GormSingularTable,
		LogMode:       cfg.IsDebugMode(),
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

	dbconn.DB().SetMaxIdleConns(db.MaxIdleConns)
	dbconn.DB().SetMaxOpenConns(db.MaxOpenConns)

	dbconn.SingularTable(db.SingularTable)
	dbconn.LogMode(db.LogMode)

	return dbconn, nil
}
