package provider

import (
	"context"

	"github.com/satriajidam/go-gin-skeleton/pkg/service/domain"
)

type service struct {
	repo domain.ProviderRepository
}

// NewService creates new provider service.
func NewService(repo domain.ProviderRepository) domain.ProviderService {
	return &service{repo}
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
func (s *service) CreateProvider(ctx context.Context, shortName, longName string) error {
	conflicting, err := s.getProviderByShortName(ctx, shortName)
	if err != nil {
		return err
	}

	if conflicting != nil {
		return domain.ErrConflict
	}

	if err := s.repo.CreateOrUpdateProvider(ctx, domain.Provider{
		ShortName: shortName,
		LongName:  longName,
	}); err != nil {
		return err
	}

	return nil
}

// UpdateProvider updates existing provider.
func (s *service) UpdateProvider(ctx context.Context, uuid, shortName, longName string) error {
	conflicting, err := s.getProviderByShortName(ctx, shortName)
	if err != nil {
		return err
	}

	existing, err := s.GetProviderByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	if conflicting != nil && conflicting.UUID != existing.UUID {
		return domain.ErrConflict
	}

	if shortName != "" {
		existing.ShortName = shortName
	}

	if longName != "" {
		existing.LongName = longName
	}

	if err := s.repo.CreateOrUpdateProvider(ctx, *existing); err != nil {
		return err
	}

	return nil
}

// GetProviderByUUID gets a provider based on its UUID.
func (s *service) GetProviderByUUID(ctx context.Context, uuid string) (*domain.Provider, error) {
	result, err := s.repo.GetProviderByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ListProviders lists all providers.
func (s *service) ListProviders(ctx context.Context, limit int) ([]domain.Provider, error) {
	results, err := s.repo.ListProviders(ctx, limit)
	if err != nil {
		return nil, err
	}
	return results, nil
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
	return nil
}
