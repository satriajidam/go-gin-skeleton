package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/satriajidam/go-gin-skeleton/pkg/config"
)

var (
	// LevelPanic represents log level 5.
	LevelPanic = zerolog.PanicLevel.String()

	// LevelFatal represents log level 4.
	LevelFatal = zerolog.FatalLevel.String()

	// LevelError represents log level 3.
	LevelError = zerolog.ErrorLevel.String()

	// LevelWarn represents log level 2.
	LevelWarn = zerolog.WarnLevel.String()

	// LevelInfo represents log level 1.
	LevelInfo = zerolog.InfoLevel.String()

	// LevelDebug represents log level 0.
	LevelDebug = zerolog.DebugLevel.String()

	// LevelTrace represents log level -1.
	LevelTrace = zerolog.TraceLevel.String()

	stderrLogger zerolog.Logger
	stdoutLogger zerolog.Logger
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	switch config.Get().AppLogLevel {
	case LevelPanic:
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case LevelFatal:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case LevelError:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case LevelWarn:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case LevelInfo:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case LevelDebug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case LevelTrace:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	stderrLogger = zerolog.New(formatConsoleWriter(os.Stderr)).With().Timestamp().Logger()
	stdoutLogger = zerolog.New(formatConsoleWriter(os.Stdout)).With().Timestamp().Logger()
}

func formatConsoleWriter(out *os.File) zerolog.ConsoleWriter {
	output := zerolog.ConsoleWriter{Out: out}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s=", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}

	return output
}

// Error prints error logs to Stderr.
func Error(err error) {
	stderrLogger.Error().Err(err)
}
