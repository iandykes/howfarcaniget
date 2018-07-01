package main

import (
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type LandingPageViewData struct {
	GoogleMapsKey string
}

// LandingPage is the main UI page
type LandingPage struct {
	Template *template.Template
	ViewData *LandingPageViewData
}

// NewLandingPage creates a new landing page
func NewLandingPage(env *Environment) *LandingPage {
	viewData := &LandingPageViewData{
		GoogleMapsKey: env.GoogleMapsKey,
	}

	page := &LandingPage{
		ViewData: viewData,
	}

	tmpl, err := template.ParseFiles("templates/landing-page.html")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Panic("LandingPage template error")
	}

	log.WithFields(log.Fields{
		"tree": tmpl.Tree.Name,
	}).Debug("Parsed template")

	page.Template = tmpl

	return page
}

func (page *LandingPage) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "text/html")
	response.WriteHeader(200)

	page.Template.Execute(response, page.ViewData)
}
