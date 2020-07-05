// Package requestid is based on: https://github.com/gin-contrib/requestid
package requestid

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const headerXRequestID = "X-Request-ID"

// Config defines the config for RequestID middleware
type Config struct {
	Generator func() string
}

// New initializes the RequestID middleware.
func New(config ...Config) gin.HandlerFunc {
	var cfg Config
	if len(config) > 0 {
		cfg = config[0]
	}

	// Set config default values
	if cfg.Generator == nil {
		cfg.Generator = func() string {
			return uuid.New().String()
		}
	}

	return func(ctx *gin.Context) {
		// Get id from request
		rid := ctx.GetHeader(headerXRequestID)

		if rid == "" {
			rid = cfg.Generator()
			ctx.Header(headerXRequestID, rid)
		}

		ctx.Next()
	}
}

// Get returns the request identifier
func Get(ctx *gin.Context) string {
	return ctx.Writer.Header().Get(headerXRequestID)
}
