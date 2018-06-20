package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	logLevelEnvVarName         = "LOG_LEVEL"
	portEnvVarName             = "HTTP_PORT"
	googleAPIKeyEnvVarName     = "GOOGLE_API_KEY"
	incDebugHandlersEnvVarName = "INCLUDE_DEBUG_HANDLERS"
)

// Environment contains the environment variables for the app
type Environment struct {
	// LogLevel is the logging level to output to console
	// Supported values: Debug, Info, Warn, Error.  Default is Info
	LogLevel string
	// Port is the HTTP port to start the server on.  Default is 8080
	Port string
	// GoogleAPIKey must be set to a valid key for calling the Distance Matrix API
	GoogleAPIKey string

	// IncDebugHandlers specifies whether to add pprof HTTP endpoints.
	// Use 0: false, 1: true. Default is false.
	IncDebugHandlersValue string
	IncDebugHandlers      bool
}

// NewEnvironment creates an Environment pointer
func NewEnvironment() *Environment {
	env := &Environment{
		LogLevel:         "Info",
		Port:             "8080",
		IncDebugHandlers: false,
	}

	if level, hasLevel := os.LookupEnv(logLevelEnvVarName); hasLevel {
		env.LogLevel = level
	}

	if port, hasPort := os.LookupEnv(portEnvVarName); hasPort {
		env.Port = port
	}

	if incDebugHandlers, hasIncDebug := os.LookupEnv(incDebugHandlersEnvVarName); hasIncDebug {
		env.IncDebugHandlersValue = incDebugHandlers
		env.IncDebugHandlers = incDebugHandlers == "1"
	}

	env.GoogleAPIKey = os.Getenv(googleAPIKeyEnvVarName)

	return env
}

// LogVariables dumps the current variables to log
func (e *Environment) LogVariables() {
	log.WithFields(log.Fields{
		logLevelEnvVarName:         e.LogLevel,
		portEnvVarName:             e.Port,
		googleAPIKeyEnvVarName:     e.GoogleAPIKey,
		incDebugHandlersEnvVarName: e.IncDebugHandlersValue,
	}).Info("Environment variables")
}
