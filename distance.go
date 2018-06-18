package main

import (
	"fmt"
	"math"
	"net/http"
	"sort"
	"time"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"

	log "github.com/sirupsen/logrus"
)

type resultPoint struct {
	destination       string
	durationGroup     int
	status            string
	duration          time.Duration
	durationInTraffic time.Duration
	distanceMetres    int
}

// TODO: Come up with a better name and interface signature
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
	// TODO: Take originInput from querystring
	d.generateDistances("52.920279,-1.469559")

	response.Header().Set("Content-Type", "text/plain")
	response.WriteHeader(200)
	response.Write([]byte("TODO"))
}

// TODO: needs to return a value!
func (d *distance) generateDistances(originInput string) {
	origin, err := maps.ParseLatLng(originInput)
	if err != nil {

		log.WithFields(log.Fields{
			"originInput": originInput,
			"err":         err,
		}).Error("Error parsing input")

		return // TODO: Return something so the API can yield a 400 response
	}

	log.WithFields(log.Fields{
		"origin": origin,
	}).Debug("Parsed originInput")

	rounded := maps.LatLng{}
	rounded.Lat = float64(int32(origin.Lat*10)) / 10
	rounded.Lng = float64(int32(origin.Lng*10)) / 10

	log.WithFields(log.Fields{
		"origin":  origin,
		"rounded": rounded,
	}).Debug("Applied rouding on origin")

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

	var destinations []string

	for lat := gridMinLat; lat < gridMaxLat; lat += gridSize {
		for lng := gridMinLng; lng < gridMaxLng; lng += gridSize {
			destinations = append(destinations, fmt.Sprintf("%f,%f", lat, lng))
		}
	}

	log.WithFields(log.Fields{
		"origin":            origin,
		"rounded":           rounded,
		"len(destinations)": len(destinations),
		"degrees":           degrees,
		"gridSize":          gridSize,
	}).Debug("Calculated grid")

	client, err := maps.NewClient(maps.WithAPIKey(d.env.GoogleAPIKey))
	if err != nil {
		log.WithFields(log.Fields{
			"originInput": originInput,
			"err":         err,
		}).Error("Error creating Google API client")

		return // TODO: Return something so the API can yield a 500 response
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

		return // TODO: Return something so the API can yield a 500 response
	}

	var results []resultPoint

	for i := 0; i < len(destinations); i++ {
		elem := resp.Rows[0].Elements[i]

		results = append(results, resultPoint{
			destination:       destinations[i],
			status:            elem.Status,
			duration:          elem.Duration,
			durationInTraffic: elem.DurationInTraffic,
			distanceMetres:    elem.Distance.Meters,
			durationGroup:     calcDurationGroup(elem),
		})
	}

	// sort by durationGroup
	sort.Slice(results, func(i, j int) bool { return results[i].durationGroup < results[j].durationGroup })

	log.WithFields(log.Fields{
		"originInput": originInput,
		"results":     results,
	}).Debug("Result points")

	// TODO: return data that the API can pass back to client
}

func calcDurationGroup(elem *maps.DistanceMatrixElement) int {

	// Take the longest duration (traffic duration may not be available)
	duration := elem.Duration
	if elem.DurationInTraffic > duration {
		duration = elem.DurationInTraffic
	}

	return int(math.Ceil(duration.Hours()))
}
