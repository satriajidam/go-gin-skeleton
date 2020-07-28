package provider

import (
	"context"

	"github.com/satriajidam/go-gin-skeleton/pkg/service/domain"
	"github.com/segmentio/ksuid"
)

type service struct {
	repo  domain.ProviderRepository
	cache domain.ProviderCache
}

// NewService creates new provider service.
func NewService(repo domain.ProviderRepository, cache domain.ProviderCache) domain.ProviderService {
	return &service{repo, cache}
}

func (s *service) getProviderByShortName(ctx context.Context, shortName string) (*domain.Provider, error) {
	var result *domain.Provider

	if cached, _ := s.cache.GetCacheByShortName(ctx, shortName); cached == nil {
		p, err := s.repo.GetProviderByShortName(ctx, shortName)
		if err != nil {
			return nil, err
		}

		if p != nil {
			go func() {
				_ = s.cache.SetCache(ctx, *p)
			}()
			result = p
		}
	} else {
		result = cached
	}

	return result, nil
}

// CreateProvider creates new provider.
func (s *service) CreateProvider(ctx context.Context, shortName, longName string) (*domain.Provider, error) {
	conflicting, err := s.getProviderByShortName(ctx, shortName)
	if err != nil && err != domain.ErrNotFound {
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
		_ = s.cache.SetCache(ctx, p)
	}()

	return &p, nil
}

// UpdateProvider updates existing provider.
func (s *service) UpdateProvider(
	ctx context.Context, uuid, shortName, longName string,
) (*domain.Provider, error) {
	conflicting, err := s.getProviderByShortName(ctx, shortName)
	if err != nil && err != domain.ErrNotFound {
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
		_ = s.cache.SetCache(ctx, *existing)
	}()

	return existing, nil
}

// GetProviderByUUID gets a provider based on its UUID.
func (s *service) GetProviderByUUID(ctx context.Context, uuid string) (*domain.Provider, error) {
	var result *domain.Provider

	if cached, _ := s.cache.GetCacheByUUID(ctx, uuid); cached == nil {
		p, err := s.repo.GetProviderByUUID(ctx, uuid)
		if err != nil {
			return nil, err
		}

		if p != nil {
			go func() {
				_ = s.cache.SetCache(ctx, *p)
			}()
			result = p
		}
	} else {
		result = cached
	}

	return result, nil
}

// GetProviders gets all providers.
func (s *service) GetProviders(ctx context.Context, offset, limit int) ([]domain.Provider, error) {
	var result []domain.Provider

	if offset < 0 {
		offset = 0
	}

	if limit < 1 {
		limit = 1
	}

	if cached, _ := s.cache.GetPagedCache(ctx, offset, limit); (cached == nil) || (len(cached) < limit) {
		ps, err := s.repo.GetProviders(ctx, offset, limit)
		if err != nil {
			return nil, err
		}

		if len(ps) > 0 {
			go func() {
				_ = s.cache.SetPagedCache(ctx, offset, limit, ps)
			}()
			result = ps
		}
	} else {
		result = cached
	}

	return result, nil
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
		_ = s.cache.DeleteCache(ctx, *p)
	}()

	return nil
}
