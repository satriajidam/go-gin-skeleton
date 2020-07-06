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

// Color returns a colorized wrapper to fmt.Sprintf based on the given color string.
func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

var (
	// Black color.
	Black = Color("\033[1;30m%s\033[0m")
	// White color.
	White = Color("\033[1;37m%s\033[0m")
	// Red color.
	Red = Color("\033[1;31m%s\033[0m")
	// Green color.
	Green = Color("\033[1;32m%s\033[0m")
	// Yellow color.
	Yellow = Color("\033[1;33m%s\033[0m")
	// Purple color.
	Purple = Color("\033[1;34m%s\033[0m")
	// Magenta color.
	Magenta = Color("\033[1;35m%s\033[0m")
	// Teal color.
	Teal = Color("\033[1;36m%s\033[0m")
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

	output.FormatTimestamp = func(i interface{}) string {
		return Black(fmt.Sprintf("%s", i))
	}

	output.FormatLevel = func(i interface{}) string {
		level := strings.ToUpper(fmt.Sprintf("%v", i))
		logStr := fmt.Sprintf("| %-6s|", level)
		switch level {
		case "INFO":
			return Green(logStr)
		case "WARN":
			return Yellow(logStr)
		case "ERROR", "FATAL", "PANIC":
			return Red(logStr)
		default:
			return White(logStr)
		}
	}

	output.FormatMessage = func(i interface{}) string {
		return White(fmt.Sprintf("message=\"%s\"", i))
	}

	output.FormatFieldName = func(i interface{}) string {
		fieldName := strings.ToLower(fmt.Sprintf("%v", i))
		logStr := fmt.Sprintf("%s=", fieldName)
		switch fieldName {
		case "error":
			return Red(logStr)
		default:
			return Teal(logStr)
		}
	}

	output.FormatFieldValue = func(i interface{}) string {
		return White(fmt.Sprintf("\"%s\"", i))
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
