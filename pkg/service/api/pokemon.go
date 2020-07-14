package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/satriajidam/go-gin-skeleton/pkg/service/domain"
)

// PokemonHTTPHandler provides methods for interacting with pokemon HTTP handler.
type PokemonHTTPHandler struct {
	service domain.PokemonService
}

// NewPokemonHTTPHandler creates new pokemon HTTP handler.
func NewPokemonHTTPHandler(service domain.PokemonService) *PokemonHTTPHandler {
	return &PokemonHTTPHandler{service}
}

func failedMsgPokemonNameNotExists(name string) string {
	return fmt.Sprintf("Pokemon named '%s' doesn't exist", name)
}

func failedMsgPokemonAction(action string) string {
	return fmt.Sprintf("Failed %s pokemon", action)
}

func successMsgPokemonAction(action string) string {
	return fmt.Sprintf("Success %s pokemon", action)
}

// GetPokemonByName gets a pokemon based on its name.
func (h *PokemonHTTPHandler) GetPokemonByName(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		responseFailed(ctx, http.StatusBadRequest, failedMsgMissingParam("name"), nil)
		return
	}

	p, err := h.service.GetPokemonByName(ctx, name)
	if err != nil {
		if err == domain.ErrNotFound {
			responseFailed(ctx, http.StatusNotFound, failedMsgPokemonNameNotExists(name), err)
			return
		}
		responseFailed(ctx, http.StatusInternalServerError, failedMsgPokemonAction(actionGet), err)
		return
	}

	responseSuccess(ctx, http.StatusOK, successMsgProviderAction(actionGet), p)
}
