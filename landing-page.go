package main

import (
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// LandingPageViewData is the view data passed to the HTML template
type LandingPageViewData struct {
	GoogleMapsKey     string
	DefaultSearchText string
	VersionInfo       *VersionInfo
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
		VersionInfo:   env.VersionInfo,

		// TODO: Maybe set this from env variable, or hard code when in Dev mode
		DefaultSearchText: "52.920279,-1.469559",
	}

	page := &LandingPage{
		ViewData: viewData,
	}

	// Dev workflow is easier if the template is parsed
	// on request, but Prod usage should cache the template.
	// Disable the preload via env variable when required
	if !env.DisableTemplatePreLoad {
		page.Template = page.ParseTemplate()
	}

	return page
}

func (page *LandingPage) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "text/html")
	response.WriteHeader(200)

	// Template may not have been created when the application started
	// so load it now if needed
	tmpl := page.Template
	if tmpl == nil {
		tmpl = page.ParseTemplate()
	}

	tmpl.Execute(response, page.ViewData)
}

// ParseTemplate parses the landing page template
func (page *LandingPage) ParseTemplate() *template.Template {
	tmpl, err := template.ParseFiles("templates/landing-page.html")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Panic("LandingPage template error")
	}

	log.WithFields(log.Fields{
		"tree": tmpl.Tree.Name,
	}).Debug("Parsed template")

	return tmpl
}
