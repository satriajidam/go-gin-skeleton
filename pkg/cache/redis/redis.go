package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/VictoriaMetrics/fastcache"
	cachev8 "github.com/go-redis/cache/v8"
	redisv8 "github.com/go-redis/redis/v8"
	"github.com/satriajidam/go-gin-skeleton/pkg/log"
)

// RedisConnection stores Redis connection client & information.
type RedisConnection struct {
	address        string
	Client         *redisv8.Client
	cache          *cachev8.Cache
	skipLocalCache bool
	Namespace      string
	MustAvailable  bool
}

// NewConnection creates new Redis connection.
func NewConnection(
	host, port, username, password, namespace string,
	dbnumber, localCacheSize int,
	localCacheTTL time.Duration,
	mustAvailable bool,
) *RedisConnection {
	address := fmt.Sprintf("%s:%s", host, port)
	client := redisv8.NewClient(&redisv8.Options{
		Addr:     address,
		Username: username,
		Password: password,
		DB:       dbnumber,
	})
	ctx := context.Background()

	_, err := client.Ping(ctx).Result()
	if err != nil && err != redisv8.Nil {
		if mustAvailable {
			panic(err)
		}
		log.Error(err, msgConnErr(address))
	}

	cacheOpts := &cachev8.Options{
		Redis: client,
	}

	skipLocalCache := true
	if localCacheSize > 0 {
		cacheOpts.LocalCache = fastcache.New(localCacheSize << 20)
		cacheOpts.LocalCacheTTL = localCacheTTL
		skipLocalCache = false
	}

	return &RedisConnection{
		address:        address,
		Client:         client,
		cache:          cachev8.New(cacheOpts),
		skipLocalCache: skipLocalCache,
		Namespace:      namespace,
		MustAvailable:  mustAvailable,
	}
}
