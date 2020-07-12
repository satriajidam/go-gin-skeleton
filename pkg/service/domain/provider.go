package domain

import "context"

type Provider struct {
	UUID      string `json:"uuid"`
	ShortName string `json:"shortName"`
	LongName  string `json:"longName"`
}

// ProviderService provides methods for interacting with Provider service.
type ProviderService interface {
	CreateProvider(ctx context.Context, shortName, longName string) error
	UpdateProvider(ctx context.Context, shortName, longName string) error
	DeleteProviderByUUID(ctx context.Context, uuid string) error
	GetProviderByUUID(ctx context.Context, uuid string) (*Provider, error)
	ListProviders(ctx context.Context, limit int) ([]Provider, error)
}

// ProviderRepository provides methods for interacting with Provider repository.
type ProviderRepository interface {
	CreateOrUpdateProvider(ctx context.Context, p Provider) error
	DeleteProviderByUUID(ctx context.Context, uuid string) error
	GetProviderByUUID(ctx context.Context, uuid string) (*Provider, error)
	GetProviderByShortName(ctx context.Context, shortName string) (*Provider, error)
	ListProviders(ctx context.Context, limit int) ([]Provider, error)
}
