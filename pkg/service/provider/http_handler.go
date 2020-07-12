package provider

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/satriajidam/go-gin-skeleton/pkg/service/domain"
)

const (
	actionGet     = "retrieving"
	actionCreate  = "creating"
	actionUpdate  = "updating"
	actionDelete  = "deleting"
	statusSuccess = "success"
	statusFailed  = "failed"
)

// Handler provides methods for interacting with provider HTTP handler.
type HTTPHandler struct {
	service domain.ProviderService
}

// NewHTTPHandler creates new provider HTTP handler.
func NewHTTPHandler(service domain.ProviderService) *HTTPHandler {
	return &HTTPHandler{service}
}

func failedMsgInvalidBody() string {
	return "Invalid request body"
}

func failedMsgMissingUUID() string {
	return "Missing 'uuid' path parameter"
}

func failedMsgEmptyPayload() string {
	return "Empty payload"
}

func failedMsgInvalidLimit() string {
	return "Invalid 'limit' query parameter"
}

func failedMsgShortNameExists(shortName string) string {
	return fmt.Sprintf("Provider with '%s' short name already exists", shortName)
}

func failedMsgUUIDNotFound(uuid string) string {
	return fmt.Sprintf("Provider with '%s' UUID doesn't exist", uuid)
}

func failedMsgProviderAction(action string) string {
	return fmt.Sprintf("Failed %s provider", action)
}

func successMsgProviderAction(action string) string {
	return fmt.Sprintf("Success %s provider", action)
}

func responseFailed(ctx *gin.Context, code int, msg string, err error) {
	if err != nil {
		_ = ctx.Error(err)
	}
	ctx.JSON(code, map[string]interface{}{
		"status":  statusFailed,
		"message": msg,
	})
}

func responseSuccess(ctx *gin.Context, code int, msg string, data interface{}) {
	respMap := map[string]interface{}{
		"status":  statusSuccess,
		"message": msg,
	}

	if data != nil {
		respMap["data"] = data
	}

	ctx.JSON(code, respMap)
}

// CreateProviderReq represents JSON request for creating new provider.
type CreateProviderReq struct {
	ShortName string `json:"shortName" binding:"required"`
	LongName  string `json:"longName" binding:"required"`
}

// CreateProvider creates new provider.
func (h *HTTPHandler) CreateProvider(ctx *gin.Context) {
	var req CreateProviderReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responseFailed(ctx, http.StatusBadRequest, failedMsgInvalidBody(), err)
		return
	}

	if err := h.service.CreateProvider(ctx, req.ShortName, req.LongName); err != nil {
		if err == domain.ErrConflict {
			responseFailed(ctx, http.StatusBadRequest, failedMsgShortNameExists(req.ShortName), err)
			return
		}
		responseFailed(ctx, http.StatusInternalServerError, failedMsgProviderAction(actionCreate), err)
		return
	}

	responseSuccess(ctx, http.StatusCreated, successMsgProviderAction(actionCreate), nil)
}

// UpdateProviderReq represents JSON request for updating existing provider.
type UpdateProviderReq struct {
	ShortName string `json:"shortName"`
	LongName  string `json:"longName"`
}

// UpdateProvider updates existing provider.
func (h *HTTPHandler) UpdateProvider(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		responseFailed(ctx, http.StatusBadRequest, failedMsgMissingUUID(), nil)
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

	if err := h.service.UpdateProvider(ctx, uuid, req.ShortName, req.LongName); err != nil {
		if err == domain.ErrConflict {
			responseFailed(ctx, http.StatusBadRequest, failedMsgShortNameExists(req.ShortName), err)
			return
		}
		if err == domain.ErrNotFound {
			responseFailed(ctx, http.StatusNotFound, failedMsgUUIDNotFound(uuid), err)
			return
		}
		responseFailed(ctx, http.StatusInternalServerError, failedMsgProviderAction(actionUpdate), err)
		return
	}

	responseSuccess(ctx, http.StatusOK, successMsgProviderAction(actionUpdate), nil)
}

// GetProviderByUUID retrieves a provider based on its UUID.
func (h *HTTPHandler) GetProviderByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		responseFailed(ctx, http.StatusBadRequest, failedMsgMissingUUID(), nil)
		return
	}

	p, err := h.service.GetProviderByUUID(ctx, uuid)
	if err != nil {
		if err == domain.ErrNotFound {
			responseFailed(ctx, http.StatusNotFound, failedMsgUUIDNotFound(uuid), err)
			return
		}
		responseFailed(ctx, http.StatusInternalServerError, failedMsgProviderAction(actionGet), err)
		return
	}

	responseSuccess(ctx, http.StatusOK, successMsgProviderAction(actionGet), p)
}

// ListProviders retrieves all providers.
func (h *HTTPHandler) ListProviders(ctx *gin.Context) {
	limitStr, ok := ctx.GetQuery("limit")
	if !ok {
		limitStr = "0"
	}

	limitInt, err := strconv.Atoi(limitStr)
	if err != nil {
		responseFailed(ctx, http.StatusBadRequest, failedMsgInvalidLimit(), err)
		return
	}

	ps, err := h.service.ListProviders(ctx, limitInt)
	if err != nil {
		responseFailed(ctx, http.StatusInternalServerError, failedMsgProviderAction(actionGet), err)
		return
	}

	responseSuccess(ctx, http.StatusOK, successMsgProviderAction(actionGet), ps)
}

// DeleteProviderByUUID deletes existing provider based on its UUID.
func (h *HTTPHandler) DeleteProviderByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		responseFailed(ctx, http.StatusBadRequest, failedMsgMissingUUID(), nil)
		return
	}

	if err := h.service.DeleteProviderByUUID(ctx, uuid); err != nil {
		if err == domain.ErrNotFound {
			responseFailed(ctx, http.StatusNotFound, failedMsgUUIDNotFound(uuid), err)
			return
		}
		responseFailed(ctx, http.StatusInternalServerError, failedMsgProviderAction(actionDelete), err)
		return
	}

	responseSuccess(ctx, http.StatusOK, successMsgProviderAction(actionDelete), nil)
}
