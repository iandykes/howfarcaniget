package main

import (
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/gorilla/handlers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Service contains the elements involved in running a service
type Service struct {
	HTTPServer *http.Server
	Mux        *http.ServeMux
	Distance   *distance
	Env        *Environment
}

// NewService creates a new Service with the defined environment settings
func NewService(env *Environment) *Service {
	service := &Service{
		Mux:      createMux(env),
		Distance: newDistance(env),
		Env:      env,
	}

	setupUIRoutes(service)
	setupAPIRoutes(service)

	service.HTTPServer = &http.Server{
		Addr: ":" + env.Port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	if env.HTTPLoggingEnabled {
		httpLogger := &HTTPLogWriter{}

		service.HTTPServer.Handler = handlers.LoggingHandler(httpLogger, service.Mux)
	} else {
		service.HTTPServer.Handler = service.Mux
	}

	return service
}

func createMux(env *Environment) *http.ServeMux {
	mux := http.NewServeMux()

	// TODO: Figure out how to record http_request metrics
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", getCurrentHealth)

	if env.IncDebugHandlers {
		// https://www.robustperception.io/analysing-prometheus-memory-usage/
		// Install http://www.graphviz.org/
		// go tool pprof -svg http://localhost:5555/debug/pprof/heap > heap.svg
		// Open http://localhost:5555/debug/pprof/ in browser
		// See https://golang.org/pkg/net/http/pprof/
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	return mux
}

func getCurrentHealth(response http.ResponseWriter, request *http.Request) {

	// TODO. Define here what bad health looks like
	response.Header().Set("Content-Type", "text/plain")
	response.WriteHeader(200)
	response.Write([]byte("OK"))
}
