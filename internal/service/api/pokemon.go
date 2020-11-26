package api

import (
	"github.com/gin-gonic/gin"
	"github.com/satriajidam/go-gin-skeleton/internal/service/domain"
)

const pokemonEntity = "pokemon"

// PokemonHTTPHandler provides methods for interacting with pokemon HTTP handler.
type PokemonHTTPHandler struct {
	service domain.PokemonService
}

// NewPokemonHTTPHandler creates new pokemon HTTP handler.
func NewPokemonHTTPHandler(service domain.PokemonService) *PokemonHTTPHandler {
	return &PokemonHTTPHandler{service}
}

// GetPokemonByName gets a pokemon based on its name.
func (h *PokemonHTTPHandler) GetPokemonByName(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		ResponseFailed(ctx, FailedMissingParam("name"), nil)
		return
	}

	p, err := h.service.GetPokemonByName(ctx, name)
	if err != nil {
		if err == domain.ErrNotFound {
			ResponseFailed(ctx, FailedEntityNotFound(pokemonEntity, "name", name), err)
			return
		}
		ResponseFailed(ctx, FailedGetEntity(pokemonEntity), err)
		return
	}

	ResponseSuccess(ctx, SuccessGetEntity(pokemonEntity, p))
}
