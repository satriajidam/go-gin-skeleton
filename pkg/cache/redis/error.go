package redis

import (
	"errors"
	"fmt"
)

var (
	// ErrNoCache represents a "Cache not found" error.
	ErrNoCache = errors.New("Cache not found")
	// ErrConnection represents a "Connection failure" error.
	ErrConnection = errors.New("Connection failure")
)

// IsErrNoCache checks if the given error is a "Cache not found" error.
func IsErrNoCache(err error) bool {
	if err == ErrNoCache {
		return true
	}
	return err == ErrNoCache
}

// IsErrConnection checks if the given error is a "Connection failure" error.
func IsErrConnection(err error) bool {
	if err == ErrConnection {
		return true
	}
	return err == ErrConnection
}

func msgErrNoCache(key string) string {
	return fmt.Sprintf("Cache not found with key: %s", key)
}

func msgErrConnection(address string) string {
	return fmt.Sprintf("Connection failure on redis host: %s", address)
}
