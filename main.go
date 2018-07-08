package main

import (
	log "github.com/sirupsen/logrus"
)

var (
	// Set at build time

	// Value of BUILD_BUILDNUMBER
	version = "dev-version"

	// Date the build was done
	buildDate = "no-date"

	// Value of BUILD_SOURCEVERSION
	// Git commit hash for the current repo
	commitHash = "no-commit"
)

func main() {

	env := NewEnvironment(&VersionInfo{
		version,
		buildDate,
		commitHash,
	})
	InitialiseLogger(env.LogLevel)

	env.LogVariables()

	service := NewService(env)

	log.WithFields(log.Fields{
		"address": service.HTTPServer.Addr,
	}).Info("Starting HTTP server")

	if err := service.HTTPServer.ListenAndServe(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("ListenAndServe error")
	}
}
