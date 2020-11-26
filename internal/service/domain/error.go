package domain

import "errors"

var (
	// ErrNotFound occurs when the requested entity doesn't exist.
	ErrNotFound = errors.New("Entity not found")
	// ErrConflict occurs when an action tries to create entity that already exists.
	ErrConflict = errors.New("Entity already exists")
)
