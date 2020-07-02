package sql

// DBConfig stores SQL database common connection config.
type DBConfig struct {
	Host          string
	Port          string
	Database      string
	Username      string
	Password      string
	Params        string
	MaxIdleConns  int
	MaxOpenConns  int
	SingularTable bool
	DebugMode     bool
}
