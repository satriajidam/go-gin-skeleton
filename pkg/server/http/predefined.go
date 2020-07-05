package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type predefinedRoute struct {
	httpMethod    string
	relateivePath string
	handlerFunc   gin.HandlerFunc
}

// List of predefined routes.
// They can be overwritten after the HTTP server object is created.
var predefinedRoutes = []predefinedRoute{
	{
		httpMethod:    http.MethodGet,
		relateivePath: "/_/status/:code",
		handlerFunc:
		// Generates responses based on the given status code.
		func(ctx *gin.Context) {
			code, err := strconv.Atoi(ctx.Param("code"))
			if err != nil {
				ctx.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
				return
			}

			var respCode int
			var respMessage string

			switch code {
			case http.StatusContinue,
				http.StatusSwitchingProtocols,
				http.StatusProcessing,
				http.StatusEarlyHints,
				http.StatusOK,
				http.StatusCreated,
				http.StatusAccepted,
				http.StatusNonAuthoritativeInfo,
				http.StatusNoContent,
				http.StatusResetContent,
				http.StatusPartialContent,
				http.StatusMultiStatus,
				http.StatusAlreadyReported,
				http.StatusIMUsed,
				http.StatusMultipleChoices,
				http.StatusMovedPermanently,
				http.StatusFound,
				http.StatusSeeOther,
				http.StatusNotModified,
				http.StatusUseProxy,
				http.StatusTemporaryRedirect,
				http.StatusPermanentRedirect,
				http.StatusBadRequest,
				http.StatusUnauthorized,
				http.StatusPaymentRequired,
				http.StatusForbidden,
				http.StatusNotFound,
				http.StatusMethodNotAllowed,
				http.StatusNotAcceptable,
				http.StatusProxyAuthRequired,
				http.StatusRequestTimeout,
				http.StatusConflict,
				http.StatusGone,
				http.StatusLengthRequired,
				http.StatusPreconditionFailed,
				http.StatusRequestEntityTooLarge,
				http.StatusRequestURITooLong,
				http.StatusUnsupportedMediaType,
				http.StatusRequestedRangeNotSatisfiable,
				http.StatusExpectationFailed,
				http.StatusTeapot,
				http.StatusMisdirectedRequest,
				http.StatusUnprocessableEntity,
				http.StatusLocked,
				http.StatusFailedDependency,
				http.StatusTooEarly,
				http.StatusUpgradeRequired,
				http.StatusPreconditionRequired,
				http.StatusTooManyRequests,
				http.StatusRequestHeaderFieldsTooLarge,
				http.StatusUnavailableForLegalReasons,
				http.StatusInternalServerError,
				http.StatusNotImplemented,
				http.StatusBadGateway,
				http.StatusServiceUnavailable,
				http.StatusGatewayTimeout,
				http.StatusHTTPVersionNotSupported,
				http.StatusVariantAlsoNegotiates,
				http.StatusInsufficientStorage,
				http.StatusLoopDetected,
				http.StatusNotExtended,
				http.StatusNetworkAuthenticationRequired:
				respCode = code
				respMessage = http.StatusText(code)
			default:
				respCode = http.StatusBadRequest
				respMessage = http.StatusText(http.StatusBadRequest)
			}

			ctx.String(respCode, respMessage)
			return
		},
	},
}

func loadPredefinedRoutes(router *gin.Engine) {
	for _, route := range predefinedRoutes {
		switch route.httpMethod {
		case http.MethodGet:
			router.GET(route.relateivePath, route.handlerFunc)
		case http.MethodHead:
			router.HEAD(route.relateivePath, route.handlerFunc)
		case http.MethodPost:
			router.POST(route.relateivePath, route.handlerFunc)
		case http.MethodPut:
			router.PUT(route.relateivePath, route.handlerFunc)
		case http.MethodPatch:
			router.PATCH(route.relateivePath, route.handlerFunc)
		case http.MethodDelete:
			router.DELETE(route.relateivePath, route.handlerFunc)
		case http.MethodOptions:
			router.OPTIONS(route.relateivePath, route.handlerFunc)
		}
	}
}
