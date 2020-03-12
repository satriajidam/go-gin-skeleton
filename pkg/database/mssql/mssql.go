package mssql

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/satriajidam/go-gin-skeleton/pkg/config"
	"github.com/satriajidam/go-gin-skeleton/pkg/database"

	// Import Microsoft SQL Server driver.
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

// mssql stores Microsoft SQL Server database configurations.
type mssql struct {
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

// Init initializes Microsoft SQL Server database engine.
func Init() database.DBEngine {
	return &mssql{
		Host:          config.Get().PostgresHost,
		Port:          config.Get().PostgresPort,
		Username:      config.Get().PostgresUsername,
		Password:      config.Get().PostgresPassword,
		Database:      config.Get().PostgresDatabase,
		Params:        config.Get().PostgresParams,
		MaxIdleConns:  config.Get().GormMaxIdleConns,
		MaxOpenConns:  config.Get().GormMaxOpenConns,
		SingularTable: config.Get().GormSingularTable,
	}
}

// GetName returns Microsoft SQL Server database engine name.
func (db *mssql) GetName() string {
	return "mssql"
}

// Connect initiates connection to a Microsoft SQL Server database.
func (db *mssql) Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s?database=%s&%s",
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
