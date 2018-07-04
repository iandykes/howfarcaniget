package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	logLevelEnvVarName               = "LOG_LEVEL"
	portEnvVarName                   = "PORT"
	googleAPIKeyEnvVarName           = "GOOGLE_API_KEY"
	googleMapsKeyEnvVarName          = "GOOGLE_MAPS_KEY"
	incDebugHandlersEnvVarName       = "INCLUDE_DEBUG_HANDLERS"
	httpLoggingEnabledEnvVarName     = "HTTP_LOGGING_ENABLED"
	disableTemplatePreLoadEnvVarName = "DISABLE_TEMPLATE_PRELOAD"
	distanceCalculatorEnvVarName     = "DISTANCE_CALCULATOR"

	googleDistanceCalculatorName = "google"
	testDistanceCalculatorName   = "test"
)

// Environment contains the environment variables for the app
type Environment struct {
	// LogLevel is the logging level to output to console
	// Supported values: Debug, Info, Warn, Error.  Default is Info
	LogLevel string
	// Port is the HTTP port to start the server on.  Default is 80
	Port string
	// GoogleAPIKey must be set to a valid key for calling the Distance Matrix API
	GoogleAPIKey string
	// GoogleMapsKey is the HTTP referrer restricted API key for client side maps scripts
	GoogleMapsKey string

	// TODO: Define a type for bool env vars to reduce this repetition

	// IncDebugHandlers specifies whether to add pprof HTTP endpoints.
	// Use 0: false, 1: true. Default is false.
	IncDebugHandlersValue string
	IncDebugHandlers      bool

	// HttpLoggingEnabledValue - 0 or 1 for logging HTTP values
	HTTPLoggingEnabledValue string
	HTTPLoggingEnabled      bool

	// DisableTemplatePreLoadValue - 0 or 1. Use 1 to reload the template every request.
	// Useful for dev workflow when changing template contents without an app restart
	DisableTemplatePreLoadValue string
	DisableTemplatePreLoad      bool

	// Name of the DistanceCalculator to load. Either google or test. Default google
	DistanceCalculatorName string
}

// NewEnvironment creates an Environment pointer
func NewEnvironment() *Environment {
	env := &Environment{
		LogLevel:               "Info",
		Port:                   "80",
		IncDebugHandlers:       false,
		HTTPLoggingEnabled:     false,
		DistanceCalculatorName: "google",
	}

	if level, hasLevel := os.LookupEnv(logLevelEnvVarName); hasLevel {
		env.LogLevel = level
	}

	if port, hasPort := os.LookupEnv(portEnvVarName); hasPort {
		env.Port = port
	}

	if calc, hasCalc := os.LookupEnv(distanceCalculatorEnvVarName); hasCalc {
		env.DistanceCalculatorName = calc
	}

	if incDebugHandlers, hasIncDebug := os.LookupEnv(incDebugHandlersEnvVarName); hasIncDebug {
		env.IncDebugHandlersValue = incDebugHandlers
		env.IncDebugHandlers = incDebugHandlers == "1"
	}

	if httpLogging, hasHTTPLogging := os.LookupEnv(httpLoggingEnabledEnvVarName); hasHTTPLogging {
		env.HTTPLoggingEnabledValue = httpLogging
		env.HTTPLoggingEnabled = httpLogging == "1"
	}

	if disableTemplatePreLoad, hasdisableTemplate := os.LookupEnv(disableTemplatePreLoadEnvVarName); hasdisableTemplate {
		env.DisableTemplatePreLoadValue = disableTemplatePreLoad
		env.DisableTemplatePreLoad = disableTemplatePreLoad == "1"
	}

	env.GoogleMapsKey = os.Getenv(googleMapsKeyEnvVarName)
	env.GoogleAPIKey = os.Getenv(googleAPIKeyEnvVarName)

	return env
}

// LogVariables dumps the current variables to log
func (e *Environment) LogVariables() {
	log.WithFields(log.Fields{
		logLevelEnvVarName:               e.LogLevel,
		portEnvVarName:                   e.Port,
		googleAPIKeyEnvVarName:           e.GoogleAPIKey,
		googleMapsKeyEnvVarName:          e.GoogleMapsKey,
		incDebugHandlersEnvVarName:       e.IncDebugHandlersValue,
		httpLoggingEnabledEnvVarName:     e.HTTPLoggingEnabledValue,
		disableTemplatePreLoadEnvVarName: e.DisableTemplatePreLoadValue,
		distanceCalculatorEnvVarName:     e.DistanceCalculatorName,
	}).Info("Environment variables")
}
