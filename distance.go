package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sort"
	"time"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"

	log "github.com/sirupsen/logrus"
)

// ResultPoint is contains the result for one point in the /distances call
type ResultPoint struct {
	Destination       maps.LatLng   `json:"destination,omitempty"`
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
	Origin        maps.LatLng   `json:"origin,omitempty"`
	Points        []ResultPoint `json:"points,omitempty"`
}

// TODO: Come up with a better name and interface signature
// TODO: Define a fake "distance" implementation to use for testing without calling Distance Matric API
type distance struct {
	env *Environment
}

func newDistance(env *Environment) *distance {
	distance := &distance{
		env: env,
	}

	return distance
}

func (d *distance) ServeHTTP(response http.ResponseWriter, request *http.Request) {

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

	// TODO: Somewhere needs to convert the text the user entered into a lat,long

	// Could use Places API to work that out? Might be better to do that as
	// a autocomplete drop down in UI in case there are options.

	result := d.generateDistances(s)

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

func (d *distance) generateDistances(originInput string) DistanceResponse {

	response := DistanceResponse{}

	origin, err := maps.ParseLatLng(originInput)
	if err != nil {

		log.WithFields(log.Fields{
			"originInput": originInput,
			"err":         err,
		}).Error("Error parsing input")

		response.StatusCode = 400
		response.StatusMessage = "Error parsing origin"

		return response
	}

	response.Origin = origin

	log.WithFields(log.Fields{
		"origin": origin,
	}).Debug("Parsed originInput")

	rounded := maps.LatLng{}
	rounded.Lat = float64(int32(origin.Lat*10)) / 10
	rounded.Lng = float64(int32(origin.Lng*10)) / 10

	log.WithFields(log.Fields{
		"origin":  origin,
		"rounded": rounded,
	}).Debug("Applied rounding on origin")

	// TODO: based on the value of x in "how far can I go in x hours", need to set degrees and gridSize to
	// produce total destinations not exceeding the maximum (maybe 100? 64 is fine, 156 is not)
	// degress will need to increase as x increases, but as degrees increases gridSize must increase too
	// in order to reduce the destinations enough

	// +/- total degrees around the truncated origin point
	degrees := float64(2)
	// Number of degrees on lat and long for each point on the grid
	gridSize := 0.5

	gridMinLat := rounded.Lat - degrees
	gridMaxLat := rounded.Lat + degrees

	gridMinLng := rounded.Lng - degrees
	gridMaxLng := rounded.Lng + degrees

	var gridPoints []maps.LatLng

	// TODO: Investigate a spiral of points instead of a grid
	for lat := gridMinLat; lat < gridMaxLat; lat += gridSize {
		for lng := gridMinLng; lng < gridMaxLng; lng += gridSize {
			gridPoints = append(gridPoints, maps.LatLng{
				Lat: lat,
				Lng: lng,
			})
		}
	}

	log.WithFields(log.Fields{
		"origin":          origin,
		"rounded":         rounded,
		"len(gridPoints)": len(gridPoints),
		"degrees":         degrees,
		"gridSize":        gridSize,
	}).Debug("Calculated grid")

	client, err := maps.NewClient(maps.WithAPIKey(d.env.GoogleAPIKey))
	if err != nil {
		log.WithFields(log.Fields{
			"originInput": originInput,
			"err":         err,
		}).Error("Error creating Google API client")

		response.StatusCode = 500
		response.StatusMessage = "Internal error creating Google API client"

		return response
	}

	var destinations []string

	for _, val := range gridPoints {
		destinations = append(destinations, fmt.Sprintf("%f,%f", val.Lat, val.Lng))
	}

	r := &maps.DistanceMatrixRequest{}

	r.Origins = []string{originInput}
	r.Destinations = destinations
	r.DepartureTime = "now"                     // TODO: Allow departure/arrival time
	r.TrafficModel = maps.TrafficModelBestGuess // TODO: Set this from caller

	log.WithFields(log.Fields{
		"APIKey":        d.env.GoogleAPIKey,
		"DepartureTime": r.DepartureTime,
		"TrafficModel":  r.TrafficModel,
		"Origins":       originInput,
		"Destinations":  len(destinations),
	}).Debug("Calling DistanceMatrix")

	resp, err := client.DistanceMatrix(context.Background(), r)
	if err != nil {
		log.WithFields(log.Fields{
			"originInput": originInput,
			"err":         err,
		}).Error("Error calling DistanceMatrix")

		response.StatusCode = 500
		response.StatusMessage = "Internal error calling Google API"

		return response
	}

	var results []ResultPoint

	for i := 0; i < len(gridPoints); i++ {
		elem := resp.Rows[0].Elements[i]

		results = append(results, ResultPoint{
			Destination:       gridPoints[i],
			Status:            elem.Status,
			Duration:          elem.Duration,
			DurationInTraffic: elem.DurationInTraffic,
			DistanceMetres:    elem.Distance.Meters,
			DurationGroup:     calcDurationGroup(elem),
		})
	}

	// sort by durationGroup
	sort.Slice(results, func(i, j int) bool { return results[i].DurationGroup < results[j].DurationGroup })

	log.WithFields(log.Fields{
		"originInput": originInput,
		"results":     results,
	}).Debug("Result points")

	response.StatusCode = 200
	response.StatusMessage = "OK"
	response.Points = results

	return response
}

func calcDurationGroup(elem *maps.DistanceMatrixElement) int {

	// Take the longest duration (traffic duration may not be available)
	duration := elem.Duration
	if elem.DurationInTraffic > duration {
		duration = elem.DurationInTraffic
	}

	return int(math.Ceil(duration.Hours()))
}
