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
	LogMode       bool
}

// Init initializes PostgreSQL database engine.
func Init(cfg *config.Config) database.DBEngine {
	return &postgres{
		Host:          cfg.PostgresHost,
		Port:          cfg.PostgresPort,
		Username:      cfg.PostgresUsername,
		Password:      cfg.PostgresPassword,
		Database:      cfg.PostgresDatabase,
		Params:        cfg.PostgresParams,
		MaxIdleConns:  cfg.GormMaxIdleConns,
		MaxOpenConns:  cfg.GormMaxOpenConns,
		SingularTable: cfg.GormSingularTable,
		LogMode:       cfg.IsDebugMode(),
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

	dbconn.DB().SetMaxIdleConns(db.MaxIdleConns)
	dbconn.DB().SetMaxOpenConns(db.MaxOpenConns)

	dbconn.SingularTable(db.SingularTable)
	dbconn.LogMode(db.LogMode)

	return dbconn, nil
}
