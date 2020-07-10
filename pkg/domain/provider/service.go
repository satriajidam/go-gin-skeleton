package provider

type Provider struct {
	UUID      string `json:"uuid"`
	ShortName string `json:"shortName"`
	LongName  string `json:"longName"`
}

func (p *Provider) toProviderModel() ProviderModel {
	return ProviderModel{
		UUID:      p.UUID,
		ShortName: p.ShortName,
		LongName:  p.LongName,
	}
}

// Service provides methods for interacting with provider service.
type Service interface {
	CreateOrUpdateProvider(p Provider) error
	GetProviderByUUID(uuid string) (*Provider, error)
	ListProviders(limit int) ([]Provider, error)
}

type service struct {
	repo Repository
}

// NewService creates new provider service.
func NewService(repo Repository) Service {
	return &service{repo}
}

// CreateProvider creates new provider or updates if it's already exist.
func (s *service) CreateOrUpdateProvider(p Provider) error {
	if err := s.repo.CreateOrUpdateProvider(p.toProviderModel()); err != nil {
		return err
	}
	return nil
}

// GetProviderByUUID gets a provider based on its UUID.
func (s *service) GetProviderByUUID(uuid string) (*Provider, error) {
	pm, err := s.repo.GetProviderByUUID(uuid)
	if err != nil {
		return nil, err
	}
	return &Provider{
		UUID:      pm.UUID,
		ShortName: pm.ShortName,
		LongName:  pm.LongName,
	}, nil
}

// ListProviders lists all providers.
func (s *service) ListProviders(limit int) ([]Provider, error) {
	pms, err := s.repo.ListProviders(limit)
	if err != nil {
		return nil, err
	}
	result := []Provider{}
	for _, pm := range pms {
		result = append(result, pm.toProvider())
	}
	return result, nil
}
