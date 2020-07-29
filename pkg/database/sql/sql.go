package sql

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/satriajidam/go-gin-skeleton/pkg/log"
)

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

// Connection stores SQL database connection client & information.
type Connection struct {
	DB      *gorm.DB
	address string
	dialect string
}

// NewConnection creates new SQL database connection.
func NewConnection(DB *gorm.DB, host, port, dialect string) *Connection {
	return &Connection{DB, fmt.Sprintf("%s:%s", host, port), dialect}
}

// LogError prints SQL database connection error log to stderr.
func (c *Connection) LogError(err error, msg string) {
	printMsg := fmt.Sprintf("%s error", c.dialect)
	if msg != "" {
		printMsg = fmt.Sprintf("%s: %s", printMsg, msg)
	}

	log.Stderr().Error().
		Timestamp().
		Str(fmt.Sprintf("%sHost", strings.ToLower(c.dialect)), c.address).
		Err(err).
		Msg(printMsg)
}

// LogWarn prints SQL database connection warning log to stdout.
func (c *Connection) LogWarn(err error, msg string) {
	printMsg := fmt.Sprintf("%s warning", c.dialect)
	if msg != "" {
		printMsg = fmt.Sprintf("%s: %s", printMsg, msg)
	}

	log.Stdout().Warn().
		Timestamp().
		Str(fmt.Sprintf("%sHost", strings.ToLower(c.dialect)), c.address).
		Err(err).
		Msg(printMsg)
}

// Close closes current db connection.
func (c *Connection) Close() error {
	return c.DB.Close()
}
