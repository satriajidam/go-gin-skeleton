// Package logger is based on: https://github.com/gin-contrib/logger
package logger

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/satriajidam/go-gin-skeleton/pkg/server/http/middleware/requestid"
)

// Config defines the config for logger middleware
type Config struct {
	Logger *zerolog.Logger
	// UTC a boolean stating whether to use UTC time zone or local.
	UTC            bool
	SkipPath       []string
	SkipPathRegexp *regexp.Regexp
}

// New initializes the logging middleware.
func New(config ...Config) gin.HandlerFunc {
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

	var sublog zerolog.Logger
	if newConfig.Logger == nil {
		sublog = log.Logger
	} else {
		sublog = *newConfig.Logger
	}

	return func(ctx *gin.Context) {
		start := time.Now()
		requestID := requestid.Get(ctx)
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		ctx.Next()
		track := true

		if _, ok := skip[path]; ok {
			track = false
		}

		if track &&
			newConfig.SkipPathRegexp != nil &&
			newConfig.SkipPathRegexp.MatchString(path) {
			track = false
		}

		if track {
			end := time.Now()
			latency := end.Sub(start)
			if newConfig.UTC {
				end = end.UTC()
			}

			var buf bytes.Buffer
			tee := io.TeeReader(ctx.Request.Body, &buf)
			body, _ := ioutil.ReadAll(tee)
			ctx.Request.Body = ioutil.NopCloser(&buf)

			errMsg := ""
			if len(ctx.Errors) > 0 {
				errMsg = ctx.Errors.String()
			}

			dumplogger := sublog.With().
				Str("request-id", requestID).
				Int("status", ctx.Writer.Status()).
				Str("method", ctx.Request.Method).
				Str("path", path).
				Str("client-ip", ctx.ClientIP()).
				Dur("latency", latency).
				Str("user-agent", ctx.Request.UserAgent()).
				Str("payload", string(body)).
				Logger()

			switch {
			case ctx.Writer.Status() >= http.StatusBadRequest && ctx.Writer.Status() < http.StatusInternalServerError:
				{
					dumplogger.Warn().Msg("")
				}
			case ctx.Writer.Status() >= http.StatusInternalServerError:
				{
					dumplogger.Error().Str("error", errMsg).Msg("")
				}
			default:
				dumplogger.Info().Msg("")
			}
		}
	}
}
