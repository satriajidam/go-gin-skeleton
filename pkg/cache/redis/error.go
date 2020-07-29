package redis

import (
	"errors"
)

var (
	// ErrNoCache represents a "Cache not found" error.
	ErrNoCache = errors.New("Cache not found")
	// ErrFailedCommand represents a "Failed command" error.
	ErrFailedCommand = errors.New("Failed command")
)

// IsErrNoCache checks if the given error is a "Cache not found" error.
func IsErrNoCache(err error) bool {
	if err == ErrNoCache {
		return true
	}
	return err == ErrNoCache
}

// IsErrFailedCommand checks if the given error is a "Failed command" error.
func IsErrFailedCommand(err error) bool {
	if err == ErrFailedCommand {
		return true
	}
	return err == ErrFailedCommand
}
