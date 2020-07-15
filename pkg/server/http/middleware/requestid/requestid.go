package requestid

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const HeaderXRequestID = "X-Request-ID"

// New initializes the request ID middleware.
func New() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get id from request
		rid := ctx.GetHeader(HeaderXRequestID)

		if rid == "" {
			rid = uuid.New().String()
			ctx.Request.Header.Set(HeaderXRequestID, rid)
		}

		// Attach the request ID to the response writer.
		ctx.Header(HeaderXRequestID, rid)

		ctx.Next()
	}
}

// Get gets the request ID from the response writer.
func Get(ctx *gin.Context) string {
	return ctx.Writer.Header().Get(HeaderXRequestID)
}
