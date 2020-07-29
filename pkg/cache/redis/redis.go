package redis

import (
	"context"
	"fmt"
	"time"

	cachev8 "github.com/go-redis/cache/v8"
	redisv8 "github.com/go-redis/redis/v8"
	"github.com/satriajidam/go-gin-skeleton/pkg/log"
)

var (
	DefaultCacheTTL = 24 * time.Hour
)

// Connection stores Redis connection client & information.
type Connection struct {
	Client    *redisv8.Client
	cache     *cachev8.Cache
	namespace string
	DebugMode bool
}

// RedisConfig stores Redis common connection config.
type RedisConfig struct {
	Host      string
	Port      string
	Username  string
	Password  string
	Namespace string
	DBNumber  int
	DebugMode bool
}

// NewConnection creates new basic Redis connection.
func NewConnection(conf RedisConfig) (*Connection, error) {
	client := redisv8.NewClient(&redisv8.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		Username: conf.Username,
		Password: conf.Password,
		DB:       conf.DBNumber,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Error(err, msgErrFailedCommand(client.Options().Addr))
		return nil, err
	}

	return &Connection{
		Client: client,
		cache: cachev8.New(&cachev8.Options{
			Redis:      client,
			LocalCache: nil,
		}),
		namespace: conf.Namespace,
		DebugMode: conf.DebugMode,
	}, nil
}

func (c *Connection) namespacedKey(key string) string {
	return fmt.Sprintf("%s:%s", c.namespace, key)
}

// SetCache caches an object using the specified key.
func (c *Connection) SetCache(
	ctx context.Context, key string, value interface{}, ttl time.Duration,
) error {
	if err := c.cache.Set(&cachev8.Item{
		Ctx:            ctx,
		Key:            c.namespacedKey(key),
		Value:          value,
		TTL:            ttl,
		SkipLocalCache: true,
	}); err != nil {
		log.Error(err, msgErrFailedCommand(c.Client.Options().Addr))
		return err
	}
	return nil
}

// GetCache gets cache for the specified key and assign the result to value.
func (c *Connection) GetCache(ctx context.Context, key string, value interface{}) error {
	if err := c.cache.GetSkippingLocalCache(ctx, c.namespacedKey(key), value); err != nil {
		if err == cachev8.ErrCacheMiss {
			if c.DebugMode {
				log.Warn(msgErrNoCache(c.namespacedKey(key)))
			}
			return ErrNoCache
		}
		log.Error(err, msgErrFailedCommand(c.Client.Options().Addr))
		return err
	}
	return nil
}

// DeleteCache deletes cache in the specified key.
func (c *Connection) DeleteCache(ctx context.Context, key string) error {
	if err := c.cache.Delete(ctx, c.namespacedKey(key)); err != nil {
		log.Error(err, msgErrFailedCommand(c.Client.Options().Addr))
		return err
	}
	return nil
}

// Close closes the client, releasing any open resources.
func (c *Connection) Close() error {
	return c.Client.Close()
}
