package redis

import "fmt"

type ErrNoCache struct {
	Key string
}

func (e ErrNoCache) Error() string {
	return fmt.Sprintf("Cache key not found: %s", e.Key)
}

func msgConnErr(address string) string {
	return fmt.Sprintf("Connection error on %s redis host", address)
}
