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
	SkipPath []string
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

// New initializes the logging middleware.
func New(port string, config ...Config) gin.HandlerFunc {
	var newConfig Config
	if len(config) > 0 {
		newConfig = config[0]
	}

	var skip map[string]struct{}
	if length := len(newConfig.SkipPath); length > 0 {
		skip = make(map[string]struct{}, length)
		for _, path := range newConfig.SkipPath {
			skip[path] = struct{}{}
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
		requestID := requestid.Get(ctx)
		routePath := ctx.FullPath()
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		var buf bytes.Buffer
		tee := io.TeeReader(ctx.Request.Body, &buf)
		body, _ := ioutil.ReadAll(tee)
		ctx.Request.Body = ioutil.NopCloser(&buf)

		ctx.Next()

		track := true
		if _, ok := skip[routePath]; ok {
			track = false
		}

		if track {
			end := time.Now()
			latency := end.Sub(start)
			if newConfig.UTC {
				end = end.UTC()
			}

			errMsg := ""
			if len(ctx.Errors) > 0 {
				errMsg = ctx.Errors.String()
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
				payload:   string(body),
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
