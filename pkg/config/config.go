package config

import (
	"sync"

	"github.com/kelseyhightower/envconfig"
)

// Config stores application's configurations.
type Config struct {
	// The name of this application.
	AppName string `envconfig:"APP_NAME" default:"gin"`

	// Available levels are based on https://github.com/rs/zerolog#leveled-logging.
	// zerolog allows for logging at the following levels (from highest to lowest):
	// - release
	// - debug
	AppMode string `envconfig:"APP_MODE" default:"debug"`

	// HTTP Server configurations.
	HTTPServerPort                      string `envconfig:"HTTP_SERVER_PORT" default:"80"`
	HTTPServerMode                      string `envconfig:"HTTP_SERVER_MODE" default:"debug"`
	HTTPServerDisallowUnknownJSONFields bool   `envconfig:"HTTP_SERVER_DISALLOW_UNKNOWN_JSON_FIELDS" default:"true"`

	// GRPC Server configurations.
	GRPCPort string `envconfig:"GRPC_PORT" default:"9090"`

	// Prometheus Server configurations.
	PrometheusPort string `envconfig:"PROMETHEUS_PORT" default:"9180"`

	// MySQL database configurations.
	MySQLHost     string `envconfig:"MYSQL_HOST" default:"127.0.0.1"`
	MySQLPort     string `envconfig:"MYSQL_PORT" default:"3306"`
	MySQLUsername string `envconfig:"MYSQL_USERNAME" default:""`
	MySQLPassword string `envconfig:"MYSQL_PASSWORD" default:""`
	MySQLDatabase string `envconfig:"MYSQL_DATABASE" default:""`
	// List of accepted MySQL parameters: https://github.com/go-sql-driver/mysql#parameters
	MySQLParams        string `envconfig:"MYSQL_PARAMS" default:"charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=Local"`
	MySQLDebugMode     bool   `envconfig:"MYSQL_DEBUG_MODE" default:"true"`
	MySQLMaxIdleConns  int    `envconfig:"MYSQL_MAX_IDLE_CONNS" default:"0"`
	MySQLMaxOpenConns  int    `envconfig:"MYSQL_MAX_OPEN_CONNS" default:"0"`
	MySQLSingularTable bool   `envconfig:"MYSQL_SINGULAR_TABLE" default:"false"`

	// PostgreSQL database configurations.
	PostgresHost     string `envconfig:"POSTGRES_HOST" default:"127.0.0.1"`
	PostgresPort     string `envconfig:"POSTGRES_PORT" default:"5432"`
	PostgresUsername string `envconfig:"POSTGRES_USERNAME" default:""`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD" default:""`
	PostgresDatabase string `envconfig:"POSTGRES_DATABASE" default:""`
	// List of accepted PostgreSQL parameters: https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters
	PostgresParams        string `envocnfig:"POSTGRES_PARAMS" default:"sslmode=require&fallback_application_name=gin"`
	PostgresDebugMode     bool   `envconfig:"POSTGRES_DEBUG_MODE" default:"true"`
	PostgresMaxIdleConns  int    `envconfig:"POSTGRES_MAX_IDLE_CONNS" default:"0"`
	PostgresMaxOpenConns  int    `envconfig:"POSTGRES_MAX_OPEN_CONNS" default:"0"`
	PostgresSingularTable bool   `envconfig:"POSTGRES_SINGULAR_TABLE" default:"false"`

	// Microsoft SQL Server database configurations.
	MSSQLHost     string `envconfig:"MSSQL_HOST" default:"127.0.0.1"`
	MSSQLPort     string `envconfig:"MSSQL_PORT" default:"1433"`
	MSSQLUsername string `envconfig:"MSSQL_USERNAME" default:""`
	MSSQLPassword string `envconfig:"MSSQL_PASSWORD" default:""`
	MSSQLDatabase string `envconfig:"MSSQL_DATABASE" default:""`
	// List of accepted Microsoft SQL Server parameters: https://github.com/denisenkom/go-mssqldb#connection-parameters-and-dsn
	MSSQLParams        string `envocnfig:"MSSQL_PARAMS" default:"encrypt=true&app+name=gin"`
	MSSQLDebugMode     bool   `envconfig:"MSSQL_DEBUG_MODE" default:"true"`
	MSSQLMaxIdleConns  int    `envconfig:"MSSQL_MAX_IDLE_CONNS" default:"0"`
	MSSQLMaxOpenConns  int    `envconfig:"MSSQL_MAX_OPEN_CONNS" default:"0"`
	MSSQLSingularTable bool   `envconfig:"MSSQL_SINGULAR_TABLE" default:"false"`

	// SQLite database configurations.
	SQLiteDatabase      string `envconfig:"SQLITE_DATABASE" default:":memory:"`
	SQLiteDebugMode     bool   `envconfig:"SQLITE_DEBUG_MODE" default:"true"`
	SQLiteMaxIdleConns  int    `envconfig:"SQLITE_MAX_IDLE_CONNS" default:"0"`
	SQLiteMaxOpenConns  int    `envconfig:"SQLITE_MAX_OPEN_CONNS" default:"0"`
	SQLiteSingularTable bool   `envconfig:"SQLITE_SINGULAR_TABLE" default:"false"`
}

var (
	once      sync.Once
	singleton *Config
)

const (
	releaseMode = "release"
	debugMode   = "debug"
)

// Get retrieves singleton object of application configurations.
func Get() *Config {
	once.Do(func() {
		singleton = &Config{}
		envconfig.MustProcess("", singleton)
	})

	return singleton
}

// IsReleaseMode checks if application is running in release mode.
func IsReleaseMode() bool {
	return Get().AppMode == releaseMode
}

// IsDebugMode checks if application is running in debug mode.
func IsDebugMode() bool {
	return Get().AppMode == debugMode
}
