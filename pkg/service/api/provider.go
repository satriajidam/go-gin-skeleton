package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/satriajidam/go-gin-skeleton/pkg/service/domain"
)

// ProviderHTTPHandler provides methods for interacting with provider HTTP handler.
type ProviderHTTPHandler struct {
	service domain.ProviderService
}

// NewHTTPHandler creates new provider HTTP handler.
func NewProviderHTTPHandler(service domain.ProviderService) *ProviderHTTPHandler {
	return &ProviderHTTPHandler{service}
}

func failedMsgProviderShortNameExists(shortName string) string {
	return fmt.Sprintf("Provider with '%s' short name already exists", shortName)
}

func failedMsgProviderUUIDNotFound(uuid string) string {
	return fmt.Sprintf("Provider with '%s' UUID doesn't exist", uuid)
}

func failedMsgProviderAction(action string) string {
	return fmt.Sprintf("Failed %s provider", action)
}

func successMsgProviderAction(action string) string {
	return fmt.Sprintf("Success %s provider", action)
}

// CreateProviderReq represents JSON request for creating new provider.
type CreateProviderReq struct {
	ShortName string `json:"shortName" binding:"required"`
	LongName  string `json:"longName" binding:"required"`
}

// CreateProvider creates new provider.
func (h *ProviderHTTPHandler) CreateProvider(ctx *gin.Context) {
	var req CreateProviderReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responseFailed(ctx, http.StatusBadRequest, failedMsgInvalidBody(), err)
		return
	}

	p, err := h.service.CreateProvider(ctx, req.ShortName, req.LongName)
	if err != nil {
		if err == domain.ErrConflict {
			responseFailed(ctx, http.StatusBadRequest, failedMsgProviderShortNameExists(req.ShortName), err)
			return
		}
		responseFailed(ctx, http.StatusInternalServerError, failedMsgProviderAction(actionCreate), err)
		return
	}

	responseSuccess(ctx, http.StatusCreated, successMsgProviderAction(actionCreate), p)
}

// UpdateProviderReq represents JSON request for updating existing provider.
type UpdateProviderReq struct {
	ShortName string `json:"shortName"`
	LongName  string `json:"longName"`
}

// UpdateProvider updates existing provider.
func (h *ProviderHTTPHandler) UpdateProvider(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		responseFailed(ctx, http.StatusBadRequest, failedMsgMissingParam("uuid"), nil)
		return
	}

	var req UpdateProviderReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responseFailed(ctx, http.StatusBadRequest, failedMsgInvalidBody(), err)
		return
	}

	if req.ShortName == "" && req.LongName == "" {
		responseFailed(ctx, http.StatusBadRequest, failedMsgEmptyPayload(), nil)
		return
	}

	p, err := h.service.UpdateProvider(ctx, uuid, req.ShortName, req.LongName)
	if err != nil {
		if err == domain.ErrConflict {
			responseFailed(ctx, http.StatusBadRequest, failedMsgProviderShortNameExists(req.ShortName), err)
			return
		}
		if err == domain.ErrNotFound {
			responseFailed(ctx, http.StatusNotFound, failedMsgProviderUUIDNotFound(uuid), err)
			return
		}
		responseFailed(ctx, http.StatusInternalServerError, failedMsgProviderAction(actionUpdate), err)
		return
	}

	responseSuccess(ctx, http.StatusOK, successMsgProviderAction(actionUpdate), p)
}

// GetProviderByUUID retrieves a provider based on its UUID.
func (h *ProviderHTTPHandler) GetProviderByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		responseFailed(ctx, http.StatusBadRequest, failedMsgMissingParam("uuid"), nil)
		return
	}

	p, err := h.service.GetProviderByUUID(ctx, uuid)
	if err != nil {
		if err == domain.ErrNotFound {
			responseFailed(ctx, http.StatusNotFound, failedMsgProviderUUIDNotFound(uuid), err)
			return
		}
		responseFailed(ctx, http.StatusInternalServerError, failedMsgProviderAction(actionGet), err)
		return
	}

	responseSuccess(ctx, http.StatusOK, successMsgProviderAction(actionGet), p)
}

// ListProviders retrieves all providers.
func (h *ProviderHTTPHandler) ListProviders(ctx *gin.Context) {
	limitStr, ok := ctx.GetQuery("limit")
	if !ok {
		limitStr = "0"
	}

	limitInt, err := strconv.Atoi(limitStr)
	if err != nil {
		responseFailed(ctx, http.StatusBadRequest, failedMsgInvalidQuery("limit"), err)
		return
	}

	ps, err := h.service.GetProviders(ctx, limitInt)
	if err != nil {
		responseFailed(ctx, http.StatusInternalServerError, failedMsgProviderAction(actionGet), err)
		return
	}

	responseSuccess(ctx, http.StatusOK, successMsgProviderAction(actionGet), ps)
}

// DeleteProviderByUUID deletes existing provider based on its UUID.
func (h *ProviderHTTPHandler) DeleteProviderByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		responseFailed(ctx, http.StatusBadRequest, failedMsgMissingParam("uuid"), nil)
		return
	}

	if err := h.service.DeleteProviderByUUID(ctx, uuid); err != nil {
		if err == domain.ErrNotFound {
			responseFailed(ctx, http.StatusNotFound, failedMsgProviderUUIDNotFound(uuid), err)
			return
		}
		responseFailed(ctx, http.StatusInternalServerError, failedMsgProviderAction(actionDelete), err)
		return
	}

	responseSuccess(ctx, http.StatusOK, successMsgProviderAction(actionDelete), nil)
}
