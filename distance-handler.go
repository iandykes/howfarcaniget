package main

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// ResultPoint is contains the result for one point in the /distances call
type ResultPoint struct {
	Destination       Coordinate    `json:"destination,omitempty"`
	DurationGroup     int           `json:"durationGroup,omitempty"`
	Status            string        `json:"status,omitempty"`
	Duration          time.Duration `json:"duration,omitempty"`
	DurationInTraffic time.Duration `json:"durationInTraffic,omitempty"`
	DistanceMetres    int           `json:"distanceMetres,omitempty"`
}

// DistanceResponse is the API response from the /distances call
type DistanceResponse struct {
	StatusCode    int           `json:"statusCode,omitempty"`
	StatusMessage string        `json:"statusMessage,omitempty"`
	Origin        Coordinate    `json:"origin,omitempty"`
	Points        []ResultPoint `json:"points,omitempty"`
}

// DistanceCalculator generates the distance response for an input origin
// TODO: Need to change this when I decide what people can use for originInput
type DistanceCalculator interface {
	GenerateDistances(origin Coordinate) DistanceResponse
	GetCoordinate(originInput string) (Coordinate, error)
}

// DistanceHandler is the API handler for the /api/distances call
type DistanceHandler struct {
	Env  *Environment
	Calc DistanceCalculator
}

// NewDistanceHandler creates a new DistanceHandler
func NewDistanceHandler(env *Environment) *DistanceHandler {
	handler := &DistanceHandler{
		Env: env,
	}

	log.WithFields(log.Fields{
		"name": env.DistanceCalculatorName,
	}).Debug("Selecting DistanceCalculator")

	if env.DistanceCalculatorName == googleDistanceCalculatorName {
		handler.Calc = &GoogleDistanceCalculator{Env: env}
	} else {
		// no other choices now, so use test is the env var is not google
		handler.Calc = &TestDistanceCalculator{}
	}

	return handler
}

func (d *DistanceHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	queryStringValues := request.URL.Query()

	s := queryStringValues.Get("s")
	if s == "" {
		log.WithFields(log.Fields{
			"queryStringValues": queryStringValues,
		}).Warn("Missing 's' in query string")

		response.Header().Set("Content-Type", "text/plain")
		response.WriteHeader(400)
		response.Write([]byte("Missing 's' query string value [origin]"))
		return
	}

	coordinate, err := d.Calc.GetCoordinate(s)
	if err != nil {
		log.WithFields(log.Fields{
			"s":   s,
			"err": err,
		}).Error("Error parsing input")

		response.Header().Set("Content-Type", "text/plain")
		response.WriteHeader(400)
		response.Write([]byte("Error parsing origin"))
		return
	}

	log.WithFields(log.Fields{
		"coordinate": coordinate,
	}).Debug("Got coordinate from string input")

	result := d.Calc.GenerateDistances(coordinate)

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		response.Header().Set("Content-Type", "text/plain")
		response.WriteHeader(500)
		response.Write([]byte("Internal error marshalling JSON response"))
	} else {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(result.StatusCode)
		response.Write(jsonBytes)
	}
}
