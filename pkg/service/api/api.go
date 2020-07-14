package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	actionGet     = "retrieving"
	actionCreate  = "creating"
	actionUpdate  = "updating"
	actionDelete  = "deleting"
	statusSuccess = "success"
	statusFailed  = "failed"
)

// HTTPResponse represents JSON response for HTTP handler.
type HTTPResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func responseFailed(ctx *gin.Context, code int, msg string, err error) {
	if err != nil {
		// Attach error to current context to push it to the logger middleware.
		_ = ctx.Error(err)
	}
	ctx.JSON(code, HTTPResponse{
		Status:  statusFailed,
		Message: msg,
		Data:    nil,
	})
}

func responseSuccess(ctx *gin.Context, code int, msg string, data interface{}) {
	resp := HTTPResponse{
		Status:  statusSuccess,
		Message: msg,
		Data:    nil,
	}

	if data != nil {
		resp.Data = data
	}

	ctx.JSON(code, resp)
}

func failedMsgInvalidBody() string {
	return "Invalid request body"
}

func failedMsgEmptyPayload() string {
	return "Empty payload"
}

func failedMsgMissingParam(param string) string {
	return fmt.Sprintf("Missing '%s' path parameter", param)
}

func failedMsgMissingQuery(query string) string {
	return fmt.Sprintf("Missing '%s' query parameter", query)
}

func failedMsgInvalidParam(param string) string {
	return fmt.Sprintf("Invalid '%s' path parameter", param)
}

func failedMsgInvalidQuery(query string) string {
	return fmt.Sprintf("Invalid '%s' query parameter", query)
}
