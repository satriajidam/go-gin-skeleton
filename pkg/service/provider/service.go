package provider

import (
	"context"
	"fmt"

	"github.com/satriajidam/go-gin-skeleton/pkg/cache/redis"
	"github.com/satriajidam/go-gin-skeleton/pkg/service/domain"
	"github.com/segmentio/ksuid"
)

type service struct {
	repo        domain.ProviderRepository
	cache       *redis.Connection
	cachePrefix string
}

// NewService creates new provider service.
func NewService(repo domain.ProviderRepository, cache *redis.Connection) domain.ProviderService {
	return &service{repo, cache, "provider"}
}

func (s *service) prefixedKey(key string) string {
	return fmt.Sprintf("%s:%s", s.cachePrefix, key)
}

func (s *service) getProviderByShortName(ctx context.Context, shortName string) (*domain.Provider, error) {
	var p *domain.Provider

	err := s.cache.GetCache(ctx, s.prefixedKey(shortName), p)
	if err != nil {
		return nil, err
	}

	p, err = s.repo.GetProviderByShortName(ctx, shortName)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	go func() {
		_ = s.cache.SetCache(ctx, s.prefixedKey(shortName), p, redis.DefaultCacheTTL)
	}()

	return p, nil
}

// CreateProvider creates new provider.
func (s *service) CreateProvider(ctx context.Context, shortName, longName string) (*domain.Provider, error) {
	conflicting, err := s.getProviderByShortName(ctx, shortName)
	if err != nil {
		return nil, err
	}

	if conflicting != nil {
		return nil, domain.ErrConflict
	}

	p := domain.Provider{
		UUID:      ksuid.New().String(),
		ShortName: shortName,
		LongName:  longName,
	}

	if err := s.repo.CreateProvider(ctx, p); err != nil {
		return nil, err
	}

	go func() {
		_ = s.cache.SetCache(ctx, s.prefixedKey(p.UUID), p, redis.DefaultCacheTTL)
	}()

	return &p, nil
}

// UpdateProvider updates existing provider.
func (s *service) UpdateProvider(
	ctx context.Context, uuid, shortName, longName string,
) (*domain.Provider, error) {
	conflicting, err := s.getProviderByShortName(ctx, shortName)
	if err != nil {
		return nil, err
	}

	existing, err := s.GetProviderByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if conflicting != nil && conflicting.UUID != existing.UUID {
		return nil, domain.ErrConflict
	}

	if shortName != "" {
		existing.ShortName = shortName
	}

	if longName != "" {
		existing.LongName = longName
	}

	if err := s.repo.UpdateProvider(ctx, *existing); err != nil {
		return nil, err
	}

	go func() {
		_ = s.cache.SetCache(ctx, s.prefixedKey(existing.UUID), existing, redis.DefaultCacheTTL)
	}()

	return existing, nil
}

// GetProviderByUUID gets a provider based on its UUID.
func (s *service) GetProviderByUUID(ctx context.Context, uuid string) (*domain.Provider, error) {
	var (
		p   *domain.Provider
		err error
	)

	_ = s.cache.GetCache(ctx, uuid, p)

	if p == nil {
		p, err = s.repo.GetProviderByUUID(ctx, uuid)
		if err != nil {
			return nil, err
		}
		go func() {
			_ = s.cache.SetCache(ctx, s.prefixedKey(uuid), p, redis.DefaultCacheTTL)
		}()
	}

	return p, nil
}

// GetProviders gets all providers.
func (s *service) GetProviders(ctx context.Context, offset, limit int) ([]domain.Provider, error) {
	var (
		ps  []domain.Provider
		err error
	)

	if offset < 0 {
		offset = 0
	}
	if limit < 1 {
		limit = 1
	}

	cacheKey := fmt.Sprintf("bulk:%d:%d", offset, limit)
	ps, err = s.getCacheBulkProviders(ctx, s.prefixedKey(cacheKey))
	if err != nil {
		return nil, err
	}

	if ps == nil {
		ps, err = s.repo.GetProviders(ctx, offset, limit)
		if err != nil {
			return nil, err
		}
	} else if len(ps) < limit {
		missingLimit := limit - len(ps)
		missingOffset := offset + missingLimit + 1
		missingPS, err := s.repo.GetProviders(ctx, missingOffset, missingLimit)
		if err != nil {
			return nil, err
		}
		ps = append(ps, missingPS...)
	}

	go func() {
		_ = s.setCacheBulkProviders(ctx, s.prefixedKey(cacheKey), ps)
	}()

	return ps, nil
}

func (s *service) getCacheBulkProviders(ctx context.Context, cacheKey string) ([]domain.Provider, error) {
	var uuids []string

	err := s.cache.GetCache(ctx, s.prefixedKey(cacheKey), uuids)
	if err != nil {
		return nil, err
	}

	if uuids == nil {
		return nil, nil
	}

	var ps []domain.Provider

	for _, uuid := range uuids {
		var p *domain.Provider
		err := s.cache.GetCache(ctx, s.prefixedKey(uuid), p)
		if err != nil {
			return nil, err
		}
		if p == nil {
			continue
		}
		ps = append(ps, *p)
	}

	return ps, nil
}

func (s *service) setCacheBulkProviders(ctx context.Context, cacheKey string, ps []domain.Provider) error {
	uuids := []string{}
	for _, p := range ps {
		uuids = append(uuids, p.UUID)
	}
	return s.cache.SetCache(ctx, s.prefixedKey(cacheKey), uuids, redis.DefaultCacheTTL)
}

// DeleteProviderByUUID deletes existing provider based on its UUID.
func (s *service) DeleteProviderByUUID(ctx context.Context, uuid string) error {
	p, err := s.GetProviderByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	if err := s.repo.DeleteProviderByUUID(ctx, uuid); err != nil {
		return err
	}

	go func() {
		_ = s.cache.DeleteCache(ctx, s.prefixedKey(p.UUID))
		_ = s.cache.DeleteCache(ctx, s.prefixedKey(p.ShortName))
	}()

	return nil
}
