// Package logger is based on: https://github.com/gin-contrib/logger
package logger

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/satriajidam/go-gin-skeleton/pkg/log"
	"github.com/satriajidam/go-gin-skeleton/pkg/server/http/middleware/requestid"
)

// Config defines the config for logger middleware
type Config struct {
	Stdout *zerolog.Logger
	Stderr *zerolog.Logger
	// UTC a boolean stating whether to use UTC time zone or local.
	UTC      bool
	Routes   []Route
	SkipPath []string
}

type Route struct {
	Method       string
	RelativePath string
	LogPayload   bool
}

type logFields struct {
	requestID string
	status    int
	method    string
	path      string
	clientIP  string
	latency   time.Duration
	userAgent string
	payload   string
}

func createDumplogger(logger *zerolog.Logger, fields logFields) zerolog.Logger {
	return logger.With().
		Str("requestID", fields.requestID).
		Int("status", fields.status).
		Str("method", fields.method).
		Str("path", fields.path).
		Str("clientIP", fields.clientIP).
		Dur("latency", fields.latency).
		Str("userAgent", fields.userAgent).
		Str("payload", fields.payload).
		Logger()
}

func pathKey(method, path string) string {
	return fmt.Sprintf("%s %s", method, path)
}

// New initializes the logging middleware.
func New(port string, config ...Config) gin.HandlerFunc {
	var newConfig Config
	if len(config) > 0 {
		newConfig = config[0]
	}

	var skipped map[string]struct{}
	if length := len(newConfig.SkipPath); length > 0 {
		skipped = make(map[string]struct{}, length)
		for _, p := range newConfig.SkipPath {
			skipped[p] = struct{}{}
		}
	}

	var logged map[string]bool
	if length := len(newConfig.Routes); length > 0 {
		logged = make(map[string]bool, length)
		for _, p := range newConfig.Routes {
			logged[pathKey(p.Method, p.RelativePath)] = p.LogPayload
		}
	}

	var stdout *zerolog.Logger
	if newConfig.Stdout == nil {
		stdout = log.Stdout()
	} else {
		stdout = newConfig.Stdout
	}

	var stderr *zerolog.Logger
	if newConfig.Stderr == nil {
		stderr = log.Stderr()
	} else {
		stderr = newConfig.Stderr
	}

	return func(ctx *gin.Context) {
		start := time.Now()
		requestID := ctx.GetHeader(requestid.HeaderXRequestID)
		method := ctx.Request.Method
		routePath := ctx.FullPath()
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		payload := ""
		if yes, ok := logged[pathKey(method, routePath)]; ok && yes {
			var buf bytes.Buffer
			tee := io.TeeReader(ctx.Request.Body, &buf)
			body, _ := ioutil.ReadAll(tee)
			ctx.Request.Body = ioutil.NopCloser(&buf)
			payload = string(body)
		}

		ctx.Next()

		track := true
		if _, ok := skipped[routePath]; ok {
			track = false
		}

		if track {
			end := time.Now()
			latency := end.Sub(start)
			if newConfig.UTC {
				end = end.UTC()
			}

			errMsg := ""
			switch {
			case len(ctx.Errors) > 0:
				errMsg = ctx.Errors.String()
			case len(ctx.Errors) == 0 &&
				ctx.Writer.Status() >= http.StatusBadRequest &&
				ctx.Writer.Status() <= http.StatusNetworkAuthenticationRequired:
				errMsg = http.StatusText(ctx.Writer.Status())
			}

			msg := fmt.Sprintf("HTTP request to port %s", port)

			fields := logFields{
				requestID: requestID,
				status:    ctx.Writer.Status(),
				method:    ctx.Request.Method,
				path:      path,
				clientIP:  ctx.ClientIP(),
				latency:   latency,
				userAgent: ctx.Request.UserAgent(),
				payload:   payload,
			}

			dumpStdout := createDumplogger(stdout, fields)
			dumpStderr := createDumplogger(stderr, fields)

			switch {
			case ctx.Writer.Status() >= http.StatusBadRequest && ctx.Writer.Status() < http.StatusInternalServerError:
				{
					dumpStdout.Warn().Timestamp().Str(log.LogFieldError, errMsg).Msg(msg)
				}
			case ctx.Writer.Status() >= http.StatusInternalServerError:
				{
					dumpStderr.Error().Timestamp().Str(log.LogFieldError, errMsg).Msg(msg)
				}
			default:
				dumpStdout.Info().Timestamp().Msg(msg)
			}
		}
	}
}
