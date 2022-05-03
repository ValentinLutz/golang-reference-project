package internal

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"time"
)

type LogLevel string

const (
	TRACE LogLevel = "TRACE"
	DEBUG          = "DEBUG"
	INFO           = "INFO"
	WARN           = "WARN"
	ERROR          = "ERROR"
	FATAL          = "FATAL"
	PANIC          = "PANIC"
)

type LoggerConfig struct {
	Pretty bool     `yaml:"pretty"`
	Level  LogLevel `yaml:"level"`
}

func NewLogger() zerolog.Logger {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	logger := zerolog.New(consoleWriter).
		With().
		Timestamp().
		Logger()

	return logger
}

func SetLogLevel(logLevel LogLevel) {
	zerolog.SetGlobalLevel(logLevel.toZeroLogLevel())
}

func (ll *LogLevel) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var logLevelString string
	err := unmarshal(&logLevelString)
	if err != nil {
		return err
	}

	logLevel := LogLevel(logLevelString)

	switch logLevel {
	case TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC:
		*ll = logLevel
		return nil
	}
	return fmt.Errorf("log level is invalid: %s", logLevel)
}

func (ll LogLevel) toZeroLogLevel() zerolog.Level {
	switch ll {
	case TRACE:
		return zerolog.TraceLevel
	case DEBUG:
		return zerolog.DebugLevel
	case INFO:
		return zerolog.InfoLevel
	case WARN:
		return zerolog.WarnLevel
	case ERROR:
		return zerolog.ErrorLevel
	case FATAL:
		return zerolog.FatalLevel
	case PANIC:
		return zerolog.PanicLevel
	}
	return zerolog.NoLevel
}
