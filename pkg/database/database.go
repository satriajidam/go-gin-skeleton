package database

import "github.com/jinzhu/gorm"

// DBEngine represents a database engine.
type DBEngine interface {
	GetName() string
	Connect() (*gorm.DB, error)
}

// Connect initiates connection to the given database engine.
func Connect(db DBEngine) (*gorm.DB, error) {
	return db.Connect()
}
