package gin

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/satriajidam/go-gin-skeleton/pkg/telemetry/metric/middleware"
)

// HTTPHandler returns a gin handler for reporting HTTP metrics.
func HTTPHandler(m middleware.HTTPMiddleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := &httpReporter{c: c}
		m.Measure(r, func() { c.Next() })
	}
}

type httpReporter struct {
	c *gin.Context
}

func (r *httpReporter) Context() context.Context {
	return r.c.Request.Context()
}

func (r *httpReporter) URLHost() string {
	return r.c.Request.Host
}

func (r *httpReporter) URLPath() string {
	return r.c.FullPath()
}

func (r *httpReporter) Method() string {
	return r.c.Request.Method
}

func (r *httpReporter) StatusCode() int {
	return r.c.Writer.Status()
}

func (r *httpReporter) RequestSize() int64 {
	size := 0

	if r.c.Request.URL != nil {
		size = len(r.c.Request.URL.String())
	}

	size += len(r.c.Request.Proto)
	size += len(r.c.Request.Method)

	for name, values := range r.c.Request.Header {
		size += len(name)
		size += len(strings.Join(values, ""))
	}

	size += len(r.c.Request.Host)

	if r.c.Request.ContentLength != -1 {
		size += int(r.c.Request.ContentLength)
	}

	return int64(size)
}

func (r *httpReporter) ResponseSize() int64 {
	return int64(r.c.Writer.Size())
}
