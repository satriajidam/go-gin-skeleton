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
		return fmt.Sprintf("\"%s\"", i)
	}

	return output
}

// Stdout returns logger which prints to console stdout.
func Stdout() *zerolog.Logger {
	logger := singleton.stdout
	return &logger
}

// Stderr returns logger which prints to console stderr.
func Stderr() *zerolog.Logger {
	logger := singleton.stderr
	return &logger
}

// Panic prints panic level logs to Stderr.
func Panic(err error, msg string) {
	Stderr().Panic().Timestamp().Err(err).Msg(msg)
}

// Fatal prints fatal level logs to Stderr.
func Fatal(err error, msg string) {
	Stderr().Fatal().Timestamp().Err(err).Msg(msg)
}

// Error prints error level logs to Stderr.
func Error(err error, msg string) {
	Stderr().Error().Timestamp().Err(err).Msg(msg)
}

// Warn prints warn level logs to Stdout.
func Warn(msg string) {
	Stdout().Warn().Timestamp().Msg(msg)
}

// Info prints info level logs to Stdout.
func Info(msg string) {
	Stdout().Info().Timestamp().Msg(msg)
}

// Debug prints debug level logs to Stdout.
func Debug(msg string) {
	Stdout().Debug().Timestamp().Msg(msg)
}
