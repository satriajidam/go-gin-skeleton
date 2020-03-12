package config

import "github.com/kelseyhightower/envconfig"

const (
	// ReleaseMode is a mode for application in release state.
	ReleaseMode = "release"
	// DebugMode is a mode for application in debug state.
	DebugMode = "debug"
)

// Config ...
type Config struct {
	// The name of this application.
	AppName string `envconfig:"APP_NAME" default:"gin"`

	// Available levels are based on https://github.com/rs/zerolog#leveled-logging.
	// zerolog allows for logging at the following levels (from highest to lowest):
	// - release
	// - debug
	AppMode string `envconfig:"APP_MODE" default:"debug"`

	// Gin framework specific configs.
	GinMode                      string `envconfig:"GIN_MODE" default:"debug"`
	GinDisallowUnknownJSONFields bool   `envconfig:"GIN_DISALLOW_UNKNOWN_JSON_FIELDS" default:"true"`

	// Gorm ORM specific configs.
	GormEnableLog     bool `envconfig:"GORM_ENABLE_LOG" default:"false"`
	GormMaxIdleConns  int  `envconfig:"GORM_MAX_IDLE_CONNS" default:"0"`
	GormMaxOpenConns  int  `envconfig:"GORM_MAX_OPEN_CONNS" default:"0"`
	GormSingularTable bool `envconfig:"GORM_SINGULAR_TABLE" default:"false"`

	// Server port configurations.
	HTTPPort       string `envconfig:"HTTP_PORT" default:"80"`
	GRPCPort       string `envconfig:"GRPC_PORT" default:"9090"`
	PrometheusPort string `envconfig:"PROMETHEUS_PORT" default:"9180"`

	// MySQL database configurations.
	MySQLHost     string `envconfig:"MYSQL_HOST" default:"127.0.0.1"`
	MySQLPort     string `envconfig:"MYSQL_PORT" default:"3306"`
	MySQLUsername string `envconfig:"MYSQL_USERNAME" default:""`
	MySQLPassword string `envconfig:"MYSQL_PASSWORD" default:""`
	MySQLDatabase string `envconfig:"MYSQL_DATABASE" default:""`
	// List of accepted MySQL parameters: https://github.com/go-sql-driver/mysql#parameters
	MySQLParams string `envconfig:"MYSQL_PARAMS" default:"charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=Local"`

	// PostgreSQL database configurations.
	PostgresHost     string `envconfig:"POSTGRES_HOST" default:"127.0.0.1"`
	PostgresPort     string `envconfig:"POSTGRES_PORT" default:"5432"`
	PostgresUsername string `envconfig:"POSTGRES_USERNAME" default:""`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD" default:""`
	PostgresDatabase string `envconfig:"POSTGRES_DATABASE" default:""`
	// List of accepted PostgreSQL parameters: https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters
	PostgresParams string `envocnfig:"POSTGRES_PARAMS" default:"sslmode=require&fallback_application_name=gin"`

	// Microsoft SQL Server database configurations.
	MSSQLHost     string `envconfig:"MSSQL_HOST" default:"127.0.0.1"`
	MSSQLPort     string `envconfig:"MSSQL_PORT" default:"1433"`
	MSSQLUsername string `envconfig:"MSSQL_USERNAME" default:""`
	MSSQLPassword string `envconfig:"MSSQL_PASSWORD" default:""`
	MSSQLDatabase string `envconfig:"MSSQL_DATABASE" default:""`
	// List of accepted Microsoft SQL Server parameters: https://github.com/denisenkom/go-mssqldb#connection-parameters-and-dsn
	MSSQLParams string `envocnfig:"MSSQL_PARAMS" default:"encrypt=true&app+name=gin"`
}

// Get ...
func Get() Config {
	config := Config{}
	envconfig.MustProcess("", &config)
	return config
}
