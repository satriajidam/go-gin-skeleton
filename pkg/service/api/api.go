package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	actionGet     = "getting"
	actionCreate  = "creating"
	actionUpdate  = "updating"
	actionDelete  = "deleting"
	statusSuccess = "success"
	statusFailed  = "failed"
)

// HTTPResponse represents JSON response for HTTP handler.
type HTTPResponse struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ResponseFailed proccess an HTTP response for failed request.
func ResponseFailed(ctx *gin.Context, resp HTTPResponse, err error) {
	if err != nil {
		// Attach error to current context to push it to the logger middleware.
		_ = ctx.Error(err)
	}
	ctx.JSON(resp.Code, resp)
}

// ResponseSuccess proccess an HTTP response for successful request.
func ResponseSuccess(ctx *gin.Context, resp HTTPResponse) {
	ctx.JSON(resp.Code, resp)
}

func SuccessGetEntity(entityName string, data interface{}) HTTPResponse {
	return HTTPResponse{
		Status:  statusSuccess,
		Code:    http.StatusOK,
		Message: fmt.Sprintf("Success %s %s", actionGet, entityName),
		Data:    data,
	}
}

func SuccessCreateEntity(entityName string, data interface{}) HTTPResponse {
	return HTTPResponse{
		Status:  statusSuccess,
		Code:    http.StatusCreated,
		Message: fmt.Sprintf("Success %s %s", actionCreate, entityName),
		Data:    data,
	}
}

func SuccessUpdateEntity(entityName string, data interface{}) HTTPResponse {
	return HTTPResponse{
		Status:  statusSuccess,
		Code:    http.StatusOK,
		Message: fmt.Sprintf("Success %s %s", actionUpdate, entityName),
		Data:    data,
	}
}

func SuccessDeleteEntity(entityName string, data interface{}) HTTPResponse {
	return HTTPResponse{
		Status:  statusSuccess,
		Code:    http.StatusOK,
		Message: fmt.Sprintf("Success %s %s", actionDelete, entityName),
		Data:    data,
	}
}

func FailedInvalidBody() HTTPResponse {
	return HTTPResponse{
		Status:  statusFailed,
		Code:    http.StatusBadRequest,
		Message: "Invalid request body",
		Data:    nil,
	}
}

func FailedEmptyPayload() HTTPResponse {
	return HTTPResponse{
		Status:  statusFailed,
		Code:    http.StatusBadRequest,
		Message: "Empty payload",
		Data:    nil,
	}
}

func FailedMissingParam(param string) HTTPResponse {
	return HTTPResponse{
		Status:  statusFailed,
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf("Missing '%s' path parameter", param),
		Data:    nil,
	}
}

func FailedInvalidParam(param string) HTTPResponse {
	return HTTPResponse{
		Status:  statusFailed,
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf("Invalid '%s' path parameter", param),
		Data:    nil,
	}
}

func FailedMissingQuery(query string) HTTPResponse {
	return HTTPResponse{
		Status:  statusFailed,
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf("Missing '%s' query parameter", query),
		Data:    nil,
	}
}

func FailedInvalidQuery(query string) HTTPResponse {
	return HTTPResponse{
		Status:  statusFailed,
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf("Invalid '%s' query parameter", query),
		Data:    nil,
	}
}

func FailedEntityNotFound(entityName, fieldName, fieldValue string) HTTPResponse {
	return HTTPResponse{
		Status:  statusFailed,
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("No %s was found with %s: %s", entityName, fieldName, fieldValue),
		Data:    nil,
	}
}

func FailedEntityConflict(entityName, fieldName, fieldValue string) HTTPResponse {
	return HTTPResponse{
		Status:  statusFailed,
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf("Duplicate entry for %s with %s: %s", entityName, fieldName, fieldValue),
		Data:    nil,
	}
}

func FailedGetEntity(entityName string) HTTPResponse {
	return HTTPResponse{
		Status: statusFailed,
		Code:   http.StatusInternalServerError,
		Message: fmt.Sprintf(
			"Failed %s %s: %s",
			entityName, actionGet, http.StatusText(http.StatusInternalServerError),
		),
		Data: nil,
	}
}

func FailedCreateEntity(entityName string) HTTPResponse {
	return HTTPResponse{
		Status: statusFailed,
		Code:   http.StatusInternalServerError,
		Message: fmt.Sprintf(
			"Failed %s %s: %s",
			actionCreate, entityName, http.StatusText(http.StatusInternalServerError),
		),
		Data: nil,
	}
}

func FailedUpdateEntity(entityName string) HTTPResponse {
	return HTTPResponse{
		Status: statusFailed,
		Code:   http.StatusInternalServerError,
		Message: fmt.Sprintf(
			"Failed %s %s: %s",
			actionUpdate, entityName, http.StatusText(http.StatusInternalServerError),
		),
		Data: nil,
	}
}

func FailedDeleteEntity(entityName string) HTTPResponse {
	return HTTPResponse{
		Status: statusFailed,
		Code:   http.StatusInternalServerError,
		Message: fmt.Sprintf(
			"Failed %s %s: %s",
			actionDelete, entityName, http.StatusText(http.StatusInternalServerError),
		),
		Data: nil,
	}
}
