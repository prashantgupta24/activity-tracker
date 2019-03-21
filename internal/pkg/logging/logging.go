package logging

import (
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

const (
	JSONFormat = "json"
	TextFormat = "text"

	Info  = "info"
	Debug = "debug"
)

func New() *log.Logger {
	logger := log.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}
	logger.SetLevel(logrus.InfoLevel)
	return logger
}

func NewLoggerLevel(logLevel string) *log.Logger {
	logger := log.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}

	switch {
	case logLevel == "debug":
		logger.SetLevel(logrus.DebugLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}
	return logger
}

func NewLoggerFormat(format string) *log.Logger {
	logger := log.New()

	switch {
	case format == "json":
		logger.Formatter = &logrus.JSONFormatter{}
	default:
		logger.Formatter = &logrus.TextFormatter{
			FullTimestamp: true,
		}
	}

	logger.SetLevel(logrus.InfoLevel)
	return logger
}

func NewLoggerLevelFormat(logLevel string, format string) *log.Logger {
	logger := log.New()
	switch {
	case logLevel == "debug":
		logger.SetLevel(logrus.DebugLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	switch {
	case format == "json":
		logger.Formatter = &logrus.JSONFormatter{}
	default:
		logger.Formatter = &logrus.TextFormatter{
			FullTimestamp: true,
		}
	}

	return logger
}
