package domain

import "context"

// Pokemon represents a pokemon entity.
type Pokemon struct {
	Name      string
	Height    int
	Weight    int
	Abilities []string
}

// PokemonService provides methods for interacting with Pokemon service.
type PokemonService interface {
	GetPokemonByName(ctx context.Context, name string) (*Pokemon, error)
}
