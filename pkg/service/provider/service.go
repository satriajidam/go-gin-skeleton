package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/satriajidam/go-gin-skeleton/pkg/cache/redis"
	"github.com/satriajidam/go-gin-skeleton/pkg/service/domain"
	"github.com/segmentio/ksuid"
)

type service struct {
	repo  domain.ProviderRepository
	cache *redis.Connection
}

// NewService creates new provider service.
func NewService(repo domain.ProviderRepository, cache *redis.Connection) domain.ProviderService {
	return &service{repo, cache}
}

func (s *service) getProviderByShortName(ctx context.Context, shortName string) (*domain.Provider, error) {
	p, err := s.repo.GetProviderByShortName(ctx, shortName)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
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

	_ = s.cache.SetCache(ctx, p.UUID, p, redis.DefaultCacheTTL)

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

	_ = s.cache.SetCache(ctx, existing.UUID, existing, redis.DefaultCacheTTL)

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
		_ = s.cache.SetCache(ctx, uuid, p, redis.DefaultCacheTTL)
	}

	return p, nil
}

// GetProviders gets all providers.
func (s *service) GetProviders(ctx context.Context, offset, limit int) ([]domain.Provider, error) {
	var (
		pms []domain.Provider
		err error
	)

	if offset < 0 {
		offset = 0
	}
	if limit < 1 {
		limit = 1
	}

	cacheKey := fmt.Sprintf("bulk_providers:%d:%d", offset, limit)
	_ = s.cache.GetCache(ctx, cacheKey, pms)

	if pms == nil {
		pms, err = s.repo.GetProviders(ctx, offset, limit)
		if err != nil {
			return nil, err
		}
		_ = s.cache.SetCache(ctx, cacheKey, pms, 5*time.Minute)
	}

	return pms, nil
}

// DeleteProviderByUUID deletes existing provider based on its UUID.
func (s *service) DeleteProviderByUUID(ctx context.Context, uuid string) error {
	_, err := s.GetProviderByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	if err := s.repo.DeleteProviderByUUID(ctx, uuid); err != nil {
		return err
	}

	_ = s.cache.DeleteCache(ctx, uuid)

	return nil
}
