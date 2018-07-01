package main

import (
	"net/http"
)

func setupUIRoutes(service *Service) {
	service.Mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	service.Mux.HandleFunc("/favicon.ico", serveFavIcon)

	service.Mux.Handle("/", NewLandingPage(service.Env))
}

func serveFavIcon(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "static/favicon.ico")
}
