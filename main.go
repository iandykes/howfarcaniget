package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {

	env := NewEnvironment()
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
