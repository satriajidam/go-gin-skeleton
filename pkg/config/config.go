package config

import "github.com/kelseyhightower/envconfig"

// Config ...
type Config struct {
	// The name of this application.
	AppName string `envconfig:"APP_NAME" default:"gin"`

	// Available levels are based on https://github.com/rs/zerolog#leveled-logging.
	// zerolog allows for logging at the following levels (from highest to lowest):
	// - panic (zerolog.PanicLevel, 5)
	// - fatal (zerolog.FatalLevel, 4)
	// - error (zerolog.ErrorLevel, 3)
	// - warn (zerolog.WarnLevel, 2)
	// - info (zerolog.InfoLevel, 1)
	// - debug (zerolog.DebugLevel, 0)
	// - trace (zerolog.TraceLevel, -1)
	LogLevel string `envconfig:"LOG_LEVEL" default:"debug"`

	// Gin framework specific configs.
	GinMode string `envconfig:"GIN_MODE" default:"release"`

	// Gorm ORM specific configs.
	GormEnableLog    bool `envconfig:"GORM_ENABLE_LOG" default:"false"`
	GormMaxIdleConns int  `envconfig:"GORM_MAX_IDLE_CONNS" default:"0"`
	GormMaxOpenConns int  `envconfig:"GORM_MAX_OPEN_CONNS" default:"0"`

	// Server port configurations.
	HTTPPort       string `envconfig:"HTTP_PORT" default:"80"`
	GRPCPort       string `envconfig:"GRPC_PORT" default:"9090"`
	PrometheusPort string `envconfig:"PROMETHEUS_PORT" default:"9180"`

	// Gorm SQL database configurations.
	// Available database modes for Gorm are:
	// - sqlite3
	// - mysql
	// - postgres
	// - mssql
	DBMode   string `envconfig:"db_mode" default:"mysql"`
	DBHost   string `envconfig:"db_host" default:"127.0.0.1"`
	DBPort   string `envconfig:"db_port" default:"3306"`
	DBUser   string `envconfig:"db_user" default:""`
	DBPass   string `envconfig:"db_pass" default:""`
	DBName   string `envconfig:"db_name" default:""`
	DBParams string `envconfig:"db_params" default:""`

	// MySQL database configurations.
	MySQLEnable   bool   `envconfig:"MYSQL_ENABLE" default:"false"`
	MySQLHost     string `envconfig:"MYSQL_HOST" default:"127.0.0.1"`
	MySQLPort     string `envconfig:"MYSQL_PORT" default:"3306"`
	MySQLUsername string `envconfig:"MYSQL_USERNAME" default:""`
	MySQLPassword string `envconfig:"MYSQL_PASSWORD" default:""`
	MySQLDatabase string `envconfig:"MYSQL_DATABASE" default:""`
	// List of accepted MySQL parameters: https://github.com/go-sql-driver/mysql#parameters
	MySQLParams string `envconfig:"MYSQL_PARAMS" default:"charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=Local"`

	// PostgreSQL database configurations.
	PostgresEnable   bool   `envconfig:"POSTGRES_ENABLE" default:"false"`
	PostgresHost     string `envconfig:"POSTGRES_HOST" default:"127.0.0.1"`
	PostgresPort     string `envconfig:"POSTGRES_PORT" default:"5432"`
	PostgresUsername string `envconfig:"POSTGRES_USERNAME" default:""`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD" default:""`
	PostgresDatabase string `envconfig:"POSTGRES_DATABASE" default:""`
	// List of accepted PostgreSQL parameters: https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters
	PostgresParams string `envocnfig:"POSTGRES_PARAMS" default:"sslmode=require&fallback_application_name=gin"`
}

// Get ...
func Get() Config {
	config := Config{}
	envconfig.MustProcess("", &config)
	return config
}
