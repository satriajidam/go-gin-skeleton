package http

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// List of predefined routes.
// They can be overwritten by re-declaring the same relative path but with different handler function
// on the HTTP server object's router.
var predefinedRoutes = []route{
	{
		method:       http.MethodGet,
		relativePath: "/_/health",
		logPayload:   false,
		handlers:     []gin.HandlerFunc{healthCheck},
	},
	{
		method:       http.MethodGet,
		relativePath: "/_/status/:code",
		logPayload:   false,
		handlers:     []gin.HandlerFunc{simulateStatusCode},
	},
	{
		method:       http.MethodPost,
		relativePath: "/_/status/:code",
		logPayload:   false,
		handlers:     []gin.HandlerFunc{simulateStatusCode},
	},
	{
		method:       http.MethodGet,
		relativePath: "/_/latency/:seconds",
		logPayload:   false,
		handlers:     []gin.HandlerFunc{simulateLatency},
	},
	{
		method:       http.MethodPost,
		relativePath: "/_/latency/:seconds",
		logPayload:   false,
		handlers:     []gin.HandlerFunc{simulateLatency},
	},
}

func getStatusCodeAndText(code int) (int, string) {
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
		return code, http.StatusText(code)
	}

	return http.StatusBadRequest, http.StatusText(http.StatusBadRequest)
}

// healthCheck is a simple endpoint for checking HTTP server health.
func healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{"status": "healthy"})
}

// simulateStatusCode simulates response based on the given status code.
func simulateStatusCode(ctx *gin.Context) {
	code, err := strconv.Atoi(ctx.Param("code"))
	if err != nil {
		ctx.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	statusCode, statusText := getStatusCodeAndText(code)

	if statusCode >= http.StatusInternalServerError &&
		statusCode <= http.StatusNetworkAuthenticationRequired {
		_ = ctx.Error(fmt.Errorf(statusText))
	}

	ctx.String(statusCode, statusText)
}

// simulateLatency simulates request with N seconds latency.
func simulateLatency(ctx *gin.Context) {
	seconds, err := strconv.Atoi(ctx.Param("seconds"))
	if err != nil {
		ctx.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	min, max := 1, 5
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(max-min+1)+seconds) * time.Second)

	ctx.String(http.StatusOK, http.StatusText(http.StatusOK))
}
