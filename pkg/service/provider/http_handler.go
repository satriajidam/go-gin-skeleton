package provider

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/satriajidam/go-gin-skeleton/pkg/service/domain"
)

const (
	actionCreate  = "creating"
	actionUpdate  = "updating"
	actionDelete  = "deleting"
	statusSuccess = "success"
	statusFailed  = "failed"
)

// Handler provides methods for interacting with provider handler.
type HTTPHandler struct {
	service domain.ProviderService
}

func NewHTTPHandler(service domain.ProviderService) *HTTPHandler {
	return &HTTPHandler{service}
}

func failedMsgInvalidBody() string {
	return "Invalid request body"
}

func failedMsgEmptyUUID() string {
	return "Empty 'uuid' path parameter"
}

func failedMsgShortNameExists(shortName string) string {
	return fmt.Sprintf("Provider with '%s' short name already exists", shortName)
}

func failedMsgUUIDNotFound(uuid string) string {
	return fmt.Sprintf("Provider with '%s' UUID doesn't exist", uuid)
}

func (h *HTTPHandler) responseFailed(ctx *gin.Context, code int, msg string) {
	ctx.JSON(code, map[string]interface{}{
		"status":  statusFailed,
		"message": msg,
	})
}

func (h *HTTPHandler) responseSuccess(ctx *gin.Context, code int, msg string, data interface{}) {
	respMap := map[string]interface{}{
		"status":  statusSuccess,
		"message": msg,
	}

	if data != nil {
		respMap["data"] = data
	}

	ctx.JSON(code, respMap)
}

type CreateProviderReq struct {
	ShortName string `json:"shortName" binding:"required"`
	LongName  string `json:"longName" binding:"required"`
}

// CreateProvider creates new provider.
func (h *HTTPHandler) CreateProvider(ctx *gin.Context) {
	var req CreateProviderReq
	if err := ctx.ShouldBind(&req); err != nil {
		_ = ctx.Error(err)
		h.responseFailed(ctx, http.StatusBadRequest, failedMsgInvalidBody())
		return
	}

	if err := h.service.CreateProvider(ctx, req.ShortName, req.LongName); err != nil {
		if err == domain.ErrConflict {
			_ = ctx.Error(err)
			h.responseFailed(ctx, http.StatusBadRequest, failedMsgShortNameExists(req.ShortName))
			return
		}
		_ = ctx.Error(err)
		h.responseFailed(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed %s new provider", actionCreate))
		return
	}

	h.responseSuccess(ctx, http.StatusCreated, fmt.Sprintf("Success %s new provider", actionCreate), nil)
}

type UpdateProviderReq struct {
	ShortName string `json:"shortName"`
	LongName  string `json:"longName"`
}

// UpdateProvider updates existing provider.
func (h *HTTPHandler) UpdateProvider(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		_ = ctx.Error(fmt.Errorf(failedMsgEmptyUUID()))
		h.responseFailed(ctx, http.StatusBadRequest, failedMsgEmptyUUID())
		return
	}

	var req UpdateProviderReq
	if err := ctx.ShouldBind(&req); err != nil {
		_ = ctx.Error(err)
		h.responseFailed(ctx, http.StatusBadRequest, failedMsgInvalidBody())
		return
	}

	if req.ShortName == "" && req.LongName == "" {
		errMsg := "Empty payload"
		_ = ctx.Error(fmt.Errorf(errMsg))
		h.responseFailed(ctx, http.StatusBadRequest, errMsg)
		return
	}

	if err := h.service.UpdateProvider(ctx, uuid, req.ShortName, req.LongName); err != nil {
		if err == domain.ErrConflict {
			_ = ctx.Error(err)
			h.responseFailed(ctx, http.StatusBadRequest, failedMsgShortNameExists(req.ShortName))
			return
		}
		if err == domain.ErrNotFound {
			_ = ctx.Error(err)
			h.responseFailed(ctx, http.StatusNotFound, failedMsgUUIDNotFound(uuid))
			return
		}
		_ = ctx.Error(err)
		h.responseFailed(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed %s existing provider", actionUpdate))
		return
	}

	h.responseSuccess(ctx, http.StatusOK, fmt.Sprintf("Success %s existing provider", actionUpdate), nil)
}

// GetProviderByUUID gets a provider based on its UUID.
func (h *HTTPHandler) GetProviderByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		_ = ctx.Error(fmt.Errorf(failedMsgEmptyUUID()))
		h.responseFailed(ctx, http.StatusBadRequest, failedMsgEmptyUUID())
		return
	}

	p, err := h.service.GetProviderByUUID(ctx, uuid)
	if err != nil {
		if err == domain.ErrNotFound {
			_ = ctx.Error(err)
			h.responseFailed(ctx, http.StatusNotFound, failedMsgUUIDNotFound(uuid))
			return
		}
		_ = ctx.Error(err)
		h.responseFailed(ctx, http.StatusInternalServerError, "Failed retrieving provider")
		return
	}

	h.responseSuccess(ctx, http.StatusOK, "Success retrieving provider", p)
}

// ListProviders lists all providers.
func (h *HTTPHandler) ListProviders(ctx *gin.Context) {
	limitStr, ok := ctx.GetQuery("limit")
	if !ok {
		limitStr = "0"
	}

	limitInt, err := strconv.Atoi(limitStr)
	if err != nil {
		_ = ctx.Error(err)
		h.responseFailed(ctx, http.StatusBadRequest, "Invalid 'limit' query parameter")
		return
	}

	ps, err := h.service.ListProviders(ctx, limitInt)
	if err != nil {
		_ = ctx.Error(err)
		h.responseFailed(ctx, http.StatusInternalServerError, "Failed retrieving providers")
		return
	}

	h.responseSuccess(ctx, http.StatusOK, "Success retrieving providers", ps)
}

// DeleteProviderByUUID deletes existing provider based on its UUID.
func (h *HTTPHandler) DeleteProviderByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		_ = ctx.Error(fmt.Errorf(failedMsgEmptyUUID()))
		h.responseFailed(ctx, http.StatusBadRequest, failedMsgEmptyUUID())
		return
	}

	if err := h.service.DeleteProviderByUUID(ctx, uuid); err != nil {
		if err == domain.ErrNotFound {
			_ = ctx.Error(err)
			h.responseFailed(ctx, http.StatusNotFound, failedMsgUUIDNotFound(uuid))
			return
		}
		_ = ctx.Error(err)
		h.responseFailed(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed %s existing provider", actionDelete))
		return
	}

	h.responseSuccess(ctx, http.StatusOK, fmt.Sprintf("Success %s existing provider", actionDelete), nil)
}
