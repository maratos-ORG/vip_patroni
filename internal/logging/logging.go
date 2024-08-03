package logging

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger
)

func init() {
	format := "TEXT"
	colors := false
	log = logrus.New()
	switch format {
	case "JSON":
		log.Formatter = new(logrus.JSONFormatter)
	default:
		log.Formatter = new(logrus.TextFormatter)
		if colors {
			log.Formatter.(*logrus.TextFormatter).DisableColors = false
		} else {
			log.Formatter.(*logrus.TextFormatter).DisableColors = true
		}
		log.Formatter.(*logrus.TextFormatter).DisableTimestamp = false // remove timestamp from test output
	}
	log.Out = os.Stdout
	log.Level = logrus.DebugLevel
}

func WithFields(fields map[string]interface{}) *logrus.Entry {
	return log.WithFields(fields)
}

// Debug ...
func Debug(format string, v ...interface{}) {
	log.Debugf(format, v...)
}

// Info ...
func Info(format string, v ...interface{}) {
	log.Infof(format, v...)
}

// Warn ...
func Warn(format string, v ...interface{}) {
	log.Warnf(format, v...)
}

// Error ...
func Error(format string, v ...interface{}) {
	log.Errorf(format, v...)
}

// Fatal ...
func Fatal(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

// Warning ...
func Warning(format string, v ...interface{}) {
	log.Warningf(format, v...)
}

//Level
func Level(level string) bool {

	switch strings.ToLower(level) {
	case "trace":
		log.Level = logrus.TraceLevel
	case "debug":
		log.Level = logrus.DebugLevel
	case "info":
		log.Level = logrus.InfoLevel
	case "warning":
		log.Level = logrus.WarnLevel
	case "error":
		log.Level = logrus.ErrorLevel
	case "fatal":
		log.Level = logrus.FatalLevel
	default:
		log.Level = logrus.InfoLevel
	}
	return true
}

var (
// // ConfigError ...
// ConfigError = "%v type=config.error"

// // HTTPError ...
// HTTPError = "%v type=http.error"

// // HTTPWarn ...
// HTTPWarn = "%v type=http.warn"

// HTTPInfo ...
// HTTPInfo = "%v type=http.info"
)
