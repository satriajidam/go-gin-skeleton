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
	stderrLogger zerolog.Logger
	stdoutLogger zerolog.Logger
}

var (
	once      sync.Once
	singleton *logger
)

func init() {
	once.Do(func() {
		singleton = &logger{
			stderrLogger: zerolog.New(formatConsoleWriter(os.Stderr)).
				Level(getLogLevel(config.Get().AppMode)).
				With().
				Logger(),
			stdoutLogger: zerolog.New(formatConsoleWriter(os.Stdout)).
				Level(getLogLevel(config.Get().AppMode)).
				With().
				Logger(),
		}
	})
}

func getLogLevel(appMode string) zerolog.Level {
	var level zerolog.Level

	switch appMode {
	case config.ReleaseMode:
		level = zerolog.InfoLevel
	case config.DebugMode:
	default:
		level = zerolog.DebugLevel
	}

	return level
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
	singleton.stderrLogger.Panic().Timestamp().Err(err).Msg(msg)
}

// Fatal prints fatal level logs to Stderr.
func Fatal(err error, msg string) {
	singleton.stderrLogger.Fatal().Timestamp().Err(err).Msg(msg)
}

// Error prints error level logs to Stderr.
func Error(err error, msg string) {
	singleton.stderrLogger.Error().Timestamp().Err(err).Msg(msg)
}
