package logging

import (
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

/*
Logs can have either JSONFormat or TextFormat
*/
const (
	JSONFormat = "json"
	TextFormat = "text"

	Info  = "info"
	Debug = "debug"
)

//New instantiates a new logger with default options
func New() *log.Logger {
	logger := log.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}
	logger.SetLevel(logrus.InfoLevel)
	return logger
}

//NewLoggerLevel instantiates a new logger with level defined by user
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

//NewLoggerFormat instantiates a new logger with format defined by user
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

//NewLoggerLevelFormat instantiates a new logger with both level and format defined by user
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
