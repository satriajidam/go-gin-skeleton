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
			}

			switch code {
			case http.StatusContinue:
			case http.StatusSwitchingProtocols:
			case http.StatusProcessing:
			case http.StatusEarlyHints:
			case http.StatusOK:
			case http.StatusCreated:
			case http.StatusAccepted:
			case http.StatusNonAuthoritativeInfo:
			case http.StatusNoContent:
			case http.StatusResetContent:
			case http.StatusPartialContent:
			case http.StatusMultiStatus:
			case http.StatusAlreadyReported:
			case http.StatusIMUsed:
			case http.StatusMultipleChoices:
			case http.StatusMovedPermanently:
			case http.StatusFound:
			case http.StatusSeeOther:
			case http.StatusNotModified:
			case http.StatusUseProxy:
			case http.StatusTemporaryRedirect:
			case http.StatusPermanentRedirect:
			case http.StatusBadRequest:
			case http.StatusUnauthorized:
			case http.StatusPaymentRequired:
			case http.StatusForbidden:
			case http.StatusNotFound:
			case http.StatusMethodNotAllowed:
			case http.StatusNotAcceptable:
			case http.StatusProxyAuthRequired:
			case http.StatusRequestTimeout:
			case http.StatusConflict:
			case http.StatusGone:
			case http.StatusLengthRequired:
			case http.StatusPreconditionFailed:
			case http.StatusRequestEntityTooLarge:
			case http.StatusRequestURITooLong:
			case http.StatusUnsupportedMediaType:
			case http.StatusRequestedRangeNotSatisfiable:
			case http.StatusExpectationFailed:
			case http.StatusTeapot:
			case http.StatusMisdirectedRequest:
			case http.StatusUnprocessableEntity:
			case http.StatusLocked:
			case http.StatusFailedDependency:
			case http.StatusTooEarly:
			case http.StatusUpgradeRequired:
			case http.StatusPreconditionRequired:
			case http.StatusTooManyRequests:
			case http.StatusRequestHeaderFieldsTooLarge:
			case http.StatusUnavailableForLegalReasons:
			case http.StatusInternalServerError:
			case http.StatusNotImplemented:
			case http.StatusBadGateway:
			case http.StatusServiceUnavailable:
			case http.StatusGatewayTimeout:
			case http.StatusHTTPVersionNotSupported:
			case http.StatusVariantAlsoNegotiates:
			case http.StatusInsufficientStorage:
			case http.StatusLoopDetected:
			case http.StatusNotExtended:
			case http.StatusNetworkAuthenticationRequired:
				ctx.String(code, http.StatusText(code))
			default:
				ctx.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			}
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
