// Package gin is based on https://github.com/slok/go-http-metrics/blob/master/middleware/gin/gin.go
// with a slight modifications to reduce metrics cardinality.
package gin

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/slok/go-http-metrics/middleware"
)

// Handler returns a Gin measuring middleware.
func Handler(paths []string, m middleware.Middleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		// In order to avoid high cardinality metrics, check each incoming request
		// path to a list of registered route paths.
		// This will make path with parameter like /provider/:name recorded as
		// /provider/:name instead of /provider/aws or /provider/gcp.
		// Ref: https://github.com/slok/go-http-metrics#custom-handler-id
		path := c.FullPath()
		if !contains(path, paths) {
			c.Next()
			return
		}
		r := &reporter{c: c}
		m.Measure(path, r, func() {
			c.Next()
		})
	}
}

func contains(path string, paths []string) bool {
	for _, p := range paths {
		if path == p {
			return true
		}
	}
	return false
}

type reporter struct {
	c *gin.Context
}

func (r *reporter) Method() string { return r.c.Request.Method }

func (r *reporter) Context() context.Context { return r.c.Request.Context() }

func (r *reporter) URLPath() string { return r.c.Request.URL.Path }

func (r *reporter) StatusCode() int { return r.c.Writer.Status() }

func (r *reporter) BytesWritten() int64 { return int64(r.c.Writer.Size()) }
