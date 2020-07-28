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
	Client        *redisv8.Client
	cache         *cachev8.Cache
	Namespace     string
	MustAvailable bool
	DebugMode     bool
}

// RedisConfig stores Redis common connection config.
type RedisConfig struct {
	Host          string
	Port          string
	Username      string
	Password      string
	Namespace     string
	DBNumber      int
	MustAvailable bool
	DebugMode     bool
}

// NewConnection creates new basic Redis connection.
func NewConnection(conf RedisConfig) *Connection {
	client := redisv8.NewClient(&redisv8.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		Username: conf.Username,
		Password: conf.Password,
		DB:       conf.DBNumber,
	})

	ctx := context.Background()

	_, err := client.Ping(ctx).Result()
	if err != nil && err != redisv8.Nil {
		log.Error(err, msgErrFailedCommand(client.Options().Addr))
		if conf.MustAvailable {
			panic(err)
		}
	}

	cacheOpts := &cachev8.Options{
		Redis:      client,
		LocalCache: nil,
	}

	return &Connection{
		Client:        client,
		cache:         cachev8.New(cacheOpts),
		Namespace:     conf.Namespace,
		MustAvailable: conf.MustAvailable,
		DebugMode:     conf.DebugMode,
	}
}

func (c *Connection) namespacedKey(key string) string {
	return fmt.Sprintf("%s:%s", c.Namespace, key)
}

// SetCache caches an object using the specified key.
func (c *Connection) SetCache(
	ctx context.Context, key string, value interface{}, ttl time.Duration,
) error {
	if err := c.cache.Once(&cachev8.Item{
		Ctx:            ctx,
		Key:            c.namespacedKey(key),
		Value:          value,
		TTL:            ttl,
		SkipLocalCache: true,
	}); err != nil {
		if !c.MustAvailable {
			log.Error(err, msgErrFailedCommand(c.Client.Options().Addr))
			return nil
		}
		return err
	}

	return nil
}

// GetCache gets cache for the specified key and assign the result to value.
func (c *Connection) GetCache(
	ctx context.Context, key string, value interface{},
) error {
	if err := c.cache.GetSkippingLocalCache(ctx, c.namespacedKey(key), value); err != nil {
		if err == cachev8.ErrCacheMiss {
			if c.DebugMode {
				log.Warn(msgErrNoCache(c.namespacedKey(key)))
			}
			return nil
		}
		if !c.MustAvailable {
			log.Error(err, msgErrFailedCommand(c.Client.Options().Addr))
			return nil
		}
		return err
	}

	return nil
}

// DeleteCache deletes cache in the specified key.
func (c *Connection) DeleteCache(ctx context.Context, key string) error {
	if err := c.cache.Delete(ctx, c.namespacedKey(key)); err != nil {
		if !c.MustAvailable {
			log.Error(err, msgErrFailedCommand(c.Client.Options().Addr))
			return nil
		}
		return err
	}

	return nil
}

// Close closes the client, releasing any open resources.
func (c *Connection) Close() error {
	return c.Client.Close()
}
