package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/satriajidam/go-gin-skeleton/pkg/config"
)

var (
	// LogLevelPanic represents log level 5.
	LogLevelPanic = zerolog.PanicLevel.String()

	// LogLevelFatal represents log level 4.
	LogLevelFatal = zerolog.FatalLevel.String()

	// LogLevelError represents log level 3.
	LogLevelError = zerolog.ErrorLevel.String()

	// LogLevelWarn represents log level 2.
	LogLevelWarn = zerolog.WarnLevel.String()

	// LogLevelInfo represents log level 1.
	LogLevelInfo = zerolog.InfoLevel.String()

	// LogLevelDebug represents log level 0.
	LogLevelDebug = zerolog.DebugLevel.String()

	// LogLevelTrace represents log level -1.
	LogLevelTrace = zerolog.TraceLevel.String()

	stderrLogger zerolog.Logger
	stdoutLogger zerolog.Logger
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	switch config.Get().AppLogLevel {
	case LogLevelPanic:
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case LogLevelFatal:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case LogLevelError:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case LogLevelWarn:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case LogLevelInfo:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case LogLevelDebug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case LogLevelTrace:
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
	stderrLogger.Fatal().Err(err)
}
