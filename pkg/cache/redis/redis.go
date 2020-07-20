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
	DebugMode      bool
}

// NewConnection creates new Redis connection.
func NewConnection(
	host, port, username, password, namespace string, dbnumber int,
	localCacheSize int, localCacheTTL time.Duration,
	mustAvailable, debugMode bool,
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
		log.Error(err, msgErrConnection(address))
		if mustAvailable {
			panic(err)
		}
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
		DebugMode:      debugMode,
	}
}

func (rc *RedisConnection) namespacedKey(key string) string {
	return fmt.Sprintf("%s:%s", rc.Namespace, key)
}

// SetCache caches an object using the specified key.
func (rc *RedisConnection) SetCache(
	ctx context.Context, key string, value interface{}, ttl time.Duration,
) error {
	err := rc.cache.Once(&cachev8.Item{
		Ctx:   ctx,
		Key:   rc.namespacedKey(key),
		Value: &value,
		TTL:   ttl,
	})
	if err != nil {
		if !rc.MustAvailable {
			log.Error(err, msgErrConnection(rc.address))
			return nil
		}
		return err
	}

	return nil
}

// GetCache gets cache for the specified key and assign the result to value.
func (rc *RedisConnection) GetCache(
	ctx context.Context, key string, value interface{},
) error {
	err := rc.cache.Get(ctx, rc.namespacedKey(key), &value)
	if err != nil {
		if err == redisv8.Nil {
			if rc.DebugMode {
				log.Warn(msgErrNoCache(rc.namespacedKey(key)))
			}
			return nil
		}
		if !rc.MustAvailable {
			log.Error(err, msgErrConnection(rc.address))
			return nil
		}
		return err
	}

	return nil
}
