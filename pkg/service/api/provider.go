package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/satriajidam/go-gin-skeleton/pkg/service/domain"
)

const (
	providerEntity   = "provider"
	providerEntities = "providers"
)

// ProviderHTTPHandler provides methods for interacting with provider HTTP handler.
type ProviderHTTPHandler struct {
	service domain.ProviderService
}

// NewHTTPHandler creates new provider HTTP handler.
func NewProviderHTTPHandler(service domain.ProviderService) *ProviderHTTPHandler {
	return &ProviderHTTPHandler{service}
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
		ResponseFailed(ctx, FailedInvalidBody(), err)
		return
	}

	p, err := h.service.CreateProvider(ctx, req.ShortName, req.LongName)
	if err != nil {
		if err == domain.ErrConflict {
			ResponseFailed(ctx, FailedEntityConflict(providerEntity, "shortName", req.ShortName), err)
			return
		}
		ResponseFailed(ctx, FailedCreateEntity(providerEntity), err)
		return
	}

	ResponseSuccess(ctx, SuccessCreateEntity(providerEntity, p))
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
		ResponseFailed(ctx, FailedMissingParam("uuid"), nil)
		return
	}

	var req UpdateProviderReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ResponseFailed(ctx, FailedInvalidBody(), err)
		return
	}

	if req.ShortName == "" && req.LongName == "" {
		ResponseFailed(ctx, FailedEmptyPayload(), nil)
		return
	}

	p, err := h.service.UpdateProvider(ctx, uuid, req.ShortName, req.LongName)
	if err != nil {
		if err == domain.ErrConflict {
			ResponseFailed(ctx, FailedEntityConflict(providerEntity, "shortName", req.ShortName), err)
			return
		}
		if err == domain.ErrNotFound {
			ResponseFailed(ctx, FailedEntityNotFound(providerEntity, "uuid", uuid), err)
			return
		}
		ResponseFailed(ctx, FailedUpdateEntity(providerEntity), err)
		return
	}

	ResponseSuccess(ctx, SuccessUpdateEntity(providerEntity, p))
}

// GetProviderByUUID retrieves a provider based on its UUID.
func (h *ProviderHTTPHandler) GetProviderByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		ResponseFailed(ctx, FailedMissingParam("uuid"), nil)
		return
	}

	p, err := h.service.GetProviderByUUID(ctx, uuid)
	if err != nil {
		if err == domain.ErrNotFound {
			ResponseFailed(ctx, FailedEntityNotFound(providerEntity, "uuid", uuid), err)
			return
		}
		ResponseFailed(ctx, FailedGetEntity(providerEntity), err)
		return
	}

	ResponseSuccess(ctx, SuccessGetEntity(providerEntity, p))
}

// GetProviders gets all providers.
func (h *ProviderHTTPHandler) GetProviders(ctx *gin.Context) {
	offsetStr, ok := ctx.GetQuery("offset")
	if !ok {
		offsetStr = "0"
	}

	offsetInt, err := strconv.Atoi(offsetStr)
	if err != nil {
		ResponseFailed(ctx, FailedInvalidQuery("offset"), err)
		return
	}

	limitStr, ok := ctx.GetQuery("limit")
	if !ok {
		limitStr = "10"
	}

	limitInt, err := strconv.Atoi(limitStr)
	if err != nil {
		ResponseFailed(ctx, FailedInvalidQuery("limit"), err)
		return
	}

	ps, err := h.service.GetProviders(ctx, offsetInt, limitInt)
	if err != nil {
		ResponseFailed(ctx, FailedGetEntity(providerEntities), err)
		return
	}

	ResponseSuccess(ctx, SuccessGetEntity(providerEntities, ps))
}

// DeleteProviderByUUID deletes existing provider based on its UUID.
func (h *ProviderHTTPHandler) DeleteProviderByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		ResponseFailed(ctx, FailedMissingParam("uuid"), nil)
		return
	}

	if err := h.service.DeleteProviderByUUID(ctx, uuid); err != nil {
		if err == domain.ErrNotFound {
			ResponseFailed(ctx, FailedEntityNotFound(providerEntity, "uuid", uuid), err)
			return
		}
		ResponseFailed(ctx, FailedDeleteEntity(providerEntity), err)
		return
	}

	ResponseSuccess(ctx, SuccessDeleteEntity(providerEntity, nil))
}
