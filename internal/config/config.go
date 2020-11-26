package config

import (
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config stores application's configurations.
type Config struct {
	// The name of this application.
	AppName string `envconfig:"APP_NAME" default:"skeleton-server"`

	// Graceful shutdown timeout in seconds.
	GracefulTimeout time.Duration `envconfig:"GRACEFUL_TIMEOUT" default:"5s"`

	// HTTP Server configurations.
	HTTPServerPort                   string        `envconfig:"HTTP_SERVER_PORT" default:"80"`
	HTTPServerEnableCORS             bool          `envconfig:"HTTP_SERVER_ENABLE_CORS" default:"true"`
	HTTPServerEnablePredefinedRoutes bool          `envconfig:"HTTP_SERVER_ENABLE_PREDEFINED_ROUTES" default:"true"`
	HTTPServerAllowMethods           []string      `envconfig:"HTTP_SERVER_ALLOW_METHODS" default:""`
	HTTPServerAllowHeaders           []string      `envconfig:"HTTP_SERVER_ALLOW_HEADERS" default:""`
	HTTPServerAllowOrigins           []string      `envconfig:"HTTP_SERVER_ALLOW_ORIGINS" default:""`
	HTTPServerMaxAge                 time.Duration `envconfig:"HTTP_SERVER_MAX_AGE" default:""`
	HTTPServerMonitorGroupedStatus   bool          `envconfig:"HTTP_SERVER_MONITOR_GROUPED_STATUS" default:"false"`
	HTTPServerMonitorSkipPaths       []string      `envconfig:"HTTP_SERVER_MONITOR_SKIP_PATHS" default:"/_/health"`

	// Prometheus Server configurations.
	PrometheusServerPort          string `envconfig:"PROMETHEUS_SERVER_PORT" default:"9180"`
	PrometheusServerMetricsPath   string `envconfig:"PROMETHEUS_SERVER_METRICS_PATH" default:"/metrics"`
	PrometheusServerMetricsPrefix string `envconfig:"PROMETHEUS_SERVER_METRICS_PREFIX" default:"http"`

	// MySQL database configurations.
	MySQLHost     string `envconfig:"MYSQL_HOST" default:"127.0.0.1"`
	MySQLPort     string `envconfig:"MYSQL_PORT" default:"3306"`
	MySQLUsername string `envconfig:"MYSQL_USERNAME" default:""`
	MySQLPassword string `envconfig:"MYSQL_PASSWORD" default:""`
	MySQLDatabase string `envconfig:"MYSQL_DATABASE" default:""`
	// List of accepted MySQL parameters: https://github.com/go-sql-driver/mysql#parameters
	MySQLParams        string `envconfig:"MYSQL_PARAMS" default:"interpolateParams=true&charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=Local"`
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
	SQLiteDatabase      string `envconfig:"SQLITE_DATABASE" default:"file:database.db?mode=memory&cache=shared"`
	SQLiteDebugMode     bool   `envconfig:"SQLITE_DEBUG_MODE" default:"true"`
	SQLiteMaxIdleConns  int    `envconfig:"SQLITE_MAX_IDLE_CONNS" default:"1"`
	SQLiteMaxOpenConns  int    `envconfig:"SQLITE_MAX_OPEN_CONNS" default:"1"`
	SQLiteSingularTable bool   `envconfig:"SQLITE_SINGULAR_TABLE" default:"false"`

	// Redis configurations.
	RedisHost          string `envconfig:"REDIS_HOST" default:"127.0.0.1"`
	RedisPort          string `envconfig:"REDIS_PORT" default:"6379"`
	RedisUsername      string `envconfig:"REDIS_USERNAME" default:""`
	RedisPassword      string `envconfig:"REDIS_PASSWORD" default:""`
	RedisNamespace     string `envconfig:"REDIS_NAMESPACE" default:""`
	RedisDBNumber      int    `envconfig:"REDIS_DB_NUMBER" default:"0"`
	RedisMustAvailable bool   `envconfig:"REDIS_MUST_AVAILABLE" default:"false"`
	RedisDebugMode     bool   `envconfig:"REDIS_DEBUG_MODE" default:"true"`

	// External dependencies.
	PokeAPIAddressV2 string        `envconfig:"POKEAPI_ADDRESS" default:"https://pokeapi.co/api/v2"`
	PokeAPITimeout   time.Duration `envconfig:"POKEAPI_TIMEOUT" default:"15s"`
}

var (
	once      sync.Once
	singleton *Config
)

// Get retrieves singleton object of application configurations.
func Get() *Config {
	once.Do(func() {
		singleton = &Config{}
		envconfig.MustProcess("", singleton)
	})

	return singleton
}
