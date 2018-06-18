package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// InitialiseLogger sets up the logger
func InitialiseLogger(level string) {
	log.SetOutput(os.Stdout)

	formatter := new(log.TextFormatter)
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02 15:04:05.000"

	log.SetFormatter(formatter)

	log.SetLevel(getLogLevel(level))

	log.Debug("Configured logger")
}

func getLogLevel(level string) log.Level {
	switch level {
	case "Debug":
		return log.DebugLevel
	case "Info":
		return log.InfoLevel
	case "Warn":
		return log.WarnLevel
	}

	return log.ErrorLevel
}

// HTTPLogWriter provides the Writer interface to a logrus logger
type HTTPLogWriter struct{}

// Write logs the value at Debug as a HTTP Request
func (l *HTTPLogWriter) Write(p []byte) (n int, err error) {

	log.WithFields(log.Fields{
		"apacheLogFormat": string(p[:]),
	}).Debug("HTTP Request")

	return len(p), nil
}
