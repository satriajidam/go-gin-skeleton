package requestid

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const HeaderXRequestID = "X-Request-ID"

// New initializes the RequestID middleware.
func New() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get id from request
		rid := ctx.GetHeader(HeaderXRequestID)

		if rid == "" {
			rid = uuid.New().String()
			ctx.Request.Header.Set(HeaderXRequestID, rid)
		}

		ctx.Next()
	}
}
