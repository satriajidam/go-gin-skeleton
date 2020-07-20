package redis

import (
	"errors"
	"fmt"
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

func msgErrNoCache(key string) string {
	return fmt.Sprintf("Cache not found with key: %s", key)
}

func msgErrFailedCommand(address string) string {
	return fmt.Sprintf("Failed command on redis host: %s", address)
}
