package provider

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler provides methods for interacting with provider handler.
type HTTPHandler interface {
	CreateOrUpdateProvider(*gin.Context)
	GetProviderByUUID(*gin.Context)
	ListProviders(*gin.Context)
}

type httpHandler struct {
	service Service
}

func NewHandler(service Service) HTTPHandler {
	return &httpHandler{service}
}

func (h *httpHandler) CreateOrUpdateProvider(ctx *gin.Context) {
	var req Provider

	if err := ctx.ShouldBind(&req); err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Invalid request body",
		})
		return
	}

	err := h.service.CreateOrUpdateProvider(req)
	if err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Failed creating/updating provider",
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Success creating/updating provider",
	})
}

func (h *httpHandler) GetProviderByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		errMsg := "Empty 'uuid' path parameter"
		_ = ctx.Error(fmt.Errorf(errMsg))
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": errMsg,
		})
		return
	}

	p, err := h.service.GetProviderByUUID(uuid)
	if err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Failed retrieving provider",
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Success retrieving provider",
		"data":    p,
	})
}

func (h *httpHandler) ListProviders(ctx *gin.Context) {
	limitStr, ok := ctx.GetQuery("limit")
	if !ok {
		limitStr = "0"
	}

	limitInt, err := strconv.Atoi(limitStr)
	if err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Invalid 'limit' query parameter",
		})
		return
	}

	ps, err := h.service.ListProviders(limitInt)
	if err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Failed retrieving providers",
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Success retrieving providers",
		"data":    ps,
	})
}
