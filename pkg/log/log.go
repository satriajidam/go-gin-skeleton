package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/logrusorgru/aurora"
	"github.com/rs/zerolog"
)

type logLevel string

const (
	debugLevel logLevel = "debug"
	infoLevel  logLevel = "info"
	warnLevel  logLevel = "warn"
	errorLevel logLevel = "error"
)

type config struct {
	logAsJSON bool     `envconfig:"LOG_AS_JSON" default:"false"`
	logLevel  logLevel `envconfig:"LOG_LEVEL" default:"info"`
}

type logger struct {
	stderr zerolog.Logger
	stdout zerolog.Logger
}

var (
	once      sync.Once
	singleton *logger
)

const (
	// LogFieldError for error log field.
	LogFieldError = "error"

	// LevelInfo for info log level.
	LevelInfo = "INFO"
	// LevelWarn for warn log level.
	LevelWarn = "WARN"
	// LevelError for error log level.
	LevelError = "ERROR"
	// LevelFatal for fatal log level.
	LevelFatal = "FATAL"
	// LevelPanic for panic log level.
	LevelPanic = "PANIC"
	// LevelDebug for debug log level.
	LevelDebug = "DEBUG"
)

func init() {
	once.Do(func() {
		logConfig := &config{}
		envconfig.MustProcess("", logConfig)

		var logLevel zerolog.Level

		switch logConfig.logLevel {
		case debugLevel:
			logLevel = zerolog.DebugLevel
		case infoLevel:
			logLevel = zerolog.InfoLevel
		case warnLevel:
			logLevel = zerolog.WarnLevel
		case errorLevel:
			logLevel = zerolog.ErrorLevel
		default:
			panic(fmt.Errorf("unsupported log level: %s", logConfig.logLevel))
		}

		var stdErrWriter io.Writer
		var stdOutWriter io.Writer

		if logConfig.logAsJSON {
			stdErrWriter = os.Stderr
			stdOutWriter = os.Stdout
		} else {
			stdErrWriter = formatConsoleWriter(os.Stderr)
			stdOutWriter = formatConsoleWriter(os.Stdout)
		}

		singleton = &logger{
			stderr: zerolog.New(stdErrWriter).
				Level(logLevel).
				With().
				Logger(),
			stdout: zerolog.New(stdOutWriter).
				Level(logLevel).
				With().
				Logger(),
		}
	})
}

func formatConsoleWriter(out *os.File) zerolog.ConsoleWriter {
	output := zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC3339}

	output.FormatTimestamp = func(i interface{}) string {
		return fmt.Sprint(aurora.BrightBlack(i))
	}

	output.FormatLevel = func(i interface{}) string {
		level := strings.ToUpper(fmt.Sprintf("%v", i))
		logStr := fmt.Sprintf("| %s |", level)
		switch level {
		case LevelInfo:
			return fmt.Sprint(aurora.BrightGreen(logStr))
		case LevelWarn:
			return fmt.Sprint(aurora.BrightYellow(logStr))
		case LevelError, LevelFatal, LevelPanic:
			return fmt.Sprint(aurora.BrightRed(logStr))
		default:
			return fmt.Sprint(aurora.BrightWhite(logStr))
		}
	}

	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s=\"%s\"", aurora.BrightCyan("message"), aurora.BrightWhite(i))
	}

	output.FormatFieldName = func(i interface{}) string {
		fieldName := fmt.Sprintf("%v", i)
		logStr := fmt.Sprintf("%s=", fieldName)
		switch fieldName {
		case LogFieldError:
			return fmt.Sprint(aurora.BrightRed(logStr))
		default:
			return fmt.Sprint(aurora.BrightCyan(logStr))
		}
	}

	output.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprint(aurora.BrightWhite(fmt.Sprintf("\"%s\"", i)))
	}

	output.FormatErrFieldName = func(i interface{}) string {
		return fmt.Sprint(aurora.BrightRed(fmt.Sprintf("%s=", i)))
	}

	output.FormatErrFieldValue = func(i interface{}) string {
		return fmt.Sprint(aurora.BrightRed(fmt.Sprintf("%s", i)))
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
