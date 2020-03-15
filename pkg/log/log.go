package log

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/satriajidam/go-gin-skeleton/pkg/config"
)

type logger struct {
	stderr zerolog.Logger
	stdout zerolog.Logger
}

var (
	once      sync.Once
	singleton *logger
)

func init() {
	once.Do(func() {
		logLevel := zerolog.DebugLevel

		if config.IsReleaseMode() {
			logLevel = zerolog.InfoLevel
		}

		singleton = &logger{
			stderr: zerolog.New(formatConsoleWriter(os.Stderr)).
				Level(logLevel).
				With().
				Logger(),
			stdout: zerolog.New(formatConsoleWriter(os.Stdout)).
				Level(logLevel).
				With().
				Logger(),
		}
	})
}

func formatConsoleWriter(out *os.File) zerolog.ConsoleWriter {
	output := zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC3339}

	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}

	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("message=\"%s\"", i)
	}

	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s=", i)
	}

	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("\"%s\"", i))
	}

	return output
}

// Panic prints panic level logs to Stderr.
func Panic(err error, msg string) {
	singleton.stderr.Panic().Timestamp().Err(err).Msg(msg)
}

// Fatal prints fatal level logs to Stderr.
func Fatal(err error, msg string) {
	singleton.stderr.Fatal().Timestamp().Err(err).Msg(msg)
}

// Error prints error level logs to Stderr.
func Error(err error, msg string) {
	singleton.stderr.Error().Timestamp().Err(err).Msg(msg)
}

// Warn prints warn level logs to Stdout.
func Warn(msg string) {
	singleton.stdout.Warn().Timestamp().Msg(msg)
}

// Info prints info level logs to Stdout.
func Info(msg string) {
	singleton.stdout.Info().Timestamp().Msg(msg)
}

// Debug prints debug level logs to Stdout.
func Debug(msg string) {
	singleton.stdout.Debug().Timestamp().Msg(msg)
}
