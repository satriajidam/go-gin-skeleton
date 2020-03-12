package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/satriajidam/go-gin-skeleton/pkg/config"
	"github.com/satriajidam/go-gin-skeleton/pkg/database"

	// Import PostgreSQL driver.
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// postgres stores PostgreSQL database configurations.
type postgres struct {
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

// Init initializes PostgreSQL database engine.
func Init() database.DBEngine {
	return &postgres{
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

// GetName returns PostgreSQL database engine name.
func (db *postgres) GetName() string {
	return "postgres"
}

// Connect initiates connection to a PostgreSQL database.
func (db *postgres) Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?%s",
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
