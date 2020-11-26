package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/satriajidam/go-gin-skeleton/internal/service/domain"
	"github.com/satriajidam/go-gin-skeleton/pkg/cache/redis"
)

var (
	singleCacheTTL = 12 * time.Hour
	pagedCacheTTL  = 1 * time.Hour
)

type cache struct {
	rc     *redis.Connection
	prefix string
}

// NewService creates new provider cache.
func NewCache(rc *redis.Connection) domain.ProviderCache {
	return &cache{rc, "provider"}
}

func (c *cache) prefixedKey(key string) string {
	return fmt.Sprintf("%s:%s", c.prefix, key)
}

// GetCacheByUUID gets a cached provider based on its UUID.
func (c *cache) GetCacheByUUID(ctx context.Context, uuid string) (*domain.Provider, error) {
	var p domain.Provider

	if err := c.rc.GetCache(ctx, c.prefixedKey(uuid), &p); err != nil {
		if err == redis.ErrNoCache {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

// SetCacheByUUID caches a provider using its UUID as the cache key.
func (c *cache) SetCacheByUUID(ctx context.Context, p domain.Provider) error {
	return c.rc.SetCache(ctx, c.prefixedKey(p.UUID), p, singleCacheTTL)
}

// DeleteCacheByUUID removes a cached provider based on its UUID.
func (c *cache) DeleteCacheByUUID(ctx context.Context, uuid string) error {
	return c.rc.DeleteCache(ctx, c.prefixedKey(uuid))
}

// GetCacheByShortName gets a cached provider based on its short name.
func (c *cache) GetCacheByShortName(ctx context.Context, shortName string) (*domain.Provider, error) {
	var uuid string

	if err := c.rc.GetCache(ctx, c.prefixedKey(shortName), &uuid); err != nil {
		if err == redis.ErrNoCache {
			return nil, nil
		}
		return nil, err
	}

	if uuid != "" {
		return c.GetCacheByUUID(ctx, uuid)
	}

	return nil, nil
}

// SetCacheByShortName caches a provider UUID using its short name as the cache key.
func (c *cache) SetCacheByShortName(ctx context.Context, shortName, uuid string) error {
	return c.rc.SetCache(ctx, c.prefixedKey(shortName), uuid, singleCacheTTL)
}

// DeleteCacheByShortName removes a cached provider based on its short name.
func (c *cache) DeleteCacheByShortName(ctx context.Context, shortName string) error {
	return c.rc.DeleteCache(ctx, c.prefixedKey(shortName))
}

// SetCache caches a provider.
func (c *cache) SetCache(ctx context.Context, p domain.Provider) error {
	if err := c.SetCacheByUUID(ctx, p); err != nil {
		return err
	}

	if err := c.SetCacheByShortName(ctx, p.ShortName, p.UUID); err != nil {
		return err
	}

	return nil
}

func (c *cache) pagedCacheKey(offset, limit int) string {
	return c.prefixedKey(fmt.Sprintf("paged:%d:%d", offset, limit))
}

func (c *cache) getPagedCache(ctx context.Context, uuids []string) ([]domain.Provider, error) {
	ps := []domain.Provider{}

	for _, uuid := range uuids {
		p, err := c.GetCacheByUUID(ctx, uuid)
		if err != nil {
			return nil, err
		}
		if p == nil {
			return nil, nil
		}
		ps = append(ps, *p)
	}

	return ps, nil
}

// GetPagedCache gets paged providers based on the offset & limit.
func (c *cache) GetPagedCache(ctx context.Context, offset, limit int) ([]domain.Provider, error) {
	var uuids []string

	err := c.rc.GetCache(ctx, c.pagedCacheKey(offset, limit), &uuids)
	if err != nil {
		if err == redis.ErrNoCache {
			return nil, nil
		}
		return nil, err
	}

	if uuids == nil {
		return nil, nil
	}

	return c.getPagedCache(ctx, uuids)
}

// SetPagedCache caches paged providers using the offset & limit as the cache key.
func (c *cache) SetPagedCache(ctx context.Context, offset, limit int, ps []domain.Provider) error {
	uuids := []string{}
	for _, p := range ps {
		if err := c.SetCacheByUUID(ctx, p); err != nil {
			return err
		}
		uuids = append(uuids, p.UUID)
	}
	return c.rc.SetCache(ctx, c.pagedCacheKey(offset, limit), uuids, pagedCacheTTL)
}

// DeleteCache removes all paged caches.
func (c *cache) DeleteAllPagedCache(ctx context.Context) error {
	if err := c.rc.DeleteCacheByPrefix(ctx, c.prefixedKey("paged:*")); err != nil {
		return err
	}
	return nil
}

// DeleteCache removes all caches that store value of the given provider.
func (c *cache) DeleteCache(ctx context.Context, p domain.Provider) error {
	if err := c.DeleteCacheByUUID(ctx, p.UUID); err != nil {
		return err
	}

	if err := c.DeleteCacheByShortName(ctx, p.ShortName); err != nil {
		return err
	}

	if err := c.DeleteAllPagedCache(ctx); err != nil {
		return err
	}

	return nil
}
