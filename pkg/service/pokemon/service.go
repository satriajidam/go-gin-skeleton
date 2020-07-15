package pokemon

import (
	"context"

	"github.com/satriajidam/go-gin-skeleton/pkg/service/client/pokeapi"
	"github.com/satriajidam/go-gin-skeleton/pkg/service/domain"
)

type service struct {
	client *pokeapi.Client
}

// NewService creates new pokemon service.
func NewService(client *pokeapi.Client) domain.PokemonService {
	return &service{client}
}

// GetPokemonByName gets a pokemon based on its name.
func (s *service) GetPokemonByName(ctx context.Context, name string) (*domain.Pokemon, error) {
	pokemon, err := s.client.GetPokemonByName(name)
	if err != nil {
		if err == pokeapi.ErrNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	abilities := []string{}
	for _, a := range pokemon.Abilities {
		abilities = append(abilities, a.Ability.Name)
	}

	return &domain.Pokemon{
		Name:      pokemon.Name,
		Height:    pokemon.Height,
		Weight:    pokemon.Weight,
		Abilities: abilities,
	}, nil
}
