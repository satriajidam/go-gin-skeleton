package config

import "github.com/kelseyhightower/envconfig"

// Config ...
type Config struct {
	// Available levels are based on https://github.com/rs/zerolog#leveled-logging.
	// zerolog allows for logging at the following levels (from highest to lowest):
	// - panic (zerolog.PanicLevel, 5)
	// - fatal (zerolog.FatalLevel, 4)
	// - error (zerolog.ErrorLevel, 3)
	// - warn (zerolog.WarnLevel, 2)
	// - info (zerolog.InfoLevel, 1)
	// - debug (zerolog.DebugLevel, 0)
	// - trace (zerolog.TraceLevel, -1)
	LogLevel string `envconfig:"log_level" default:"debug"`

	// Frameworks specific logging config. Turn off by default.
	GinMode       string `envconfig:"gin_mode" default:"release"`
	GormEnableLog bool   `envconfig:"gorm_enable_log" default:"false"`

	// Server port configurations.
	HTTPPort       string `envconfig:"http_port" default:"80"`
	GRPCPort       string `envconfig:"grpc_port" default:"9090"`
	PrometheusPort string `envconfig:"prometheus_port" default:"9180"`

	// Gorm SQL database configurations.
	// Available database modes for Gorm are:
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
}

// Get ...
func Get() Config {
	config := Config{}
	envconfig.MustProcess("", &config)
	return config
}
