package main

import (
	"fmt"
	"math"
	"sort"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"

	log "github.com/sirupsen/logrus"
)

// GoogleDistanceCalculator is a DistanceCalculator that uses Google's Distance Matrix API
// TODO: Come up with a better name and interface signature
// TODO: Define a fake "distance" implementation to use for testing without calling Distance Matrix API
type GoogleDistanceCalculator struct {
	Env *Environment
}

// GetCoordinate gets a Coordinate based on the string origin value
func (d *GoogleDistanceCalculator) GetCoordinate(originInput string) (Coordinate, error) {
	// TODO: Temp until work out what the user should enter
	origin, err := maps.ParseLatLng(originInput)
	if err != nil {

		log.WithFields(log.Fields{
			"originInput": originInput,
			"err":         err,
		}).Error("Error parsing input")

		return Coordinate{}, err
	}

	log.WithFields(log.Fields{
		"originInput": originInput,
		"LatLng":      origin,
	}).Debug("Parsed originInput")

	return toCoordinate(origin), nil
}

// GenerateDistances makes this type a DistanceCalculator implementation
func (d *GoogleDistanceCalculator) GenerateDistances(origin Coordinate) DistanceResponse {

	response := DistanceResponse{}

	response.Origin = origin

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

	client, err := maps.NewClient(maps.WithAPIKey(d.Env.GoogleAPIKey))
	if err != nil {
		log.WithFields(log.Fields{
			"origin": origin,
			"err":    err,
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

	r.Origins = []string{origin.String()}

	r.Destinations = destinations
	r.DepartureTime = "now"                     // TODO: Allow departure/arrival time
	r.TrafficModel = maps.TrafficModelBestGuess // TODO: Set this from caller

	log.WithFields(log.Fields{
		"APIKey":        d.Env.GoogleAPIKey,
		"DepartureTime": r.DepartureTime,
		"TrafficModel":  r.TrafficModel,
		"Origins":       origin,
		"Destinations":  len(destinations),
	}).Debug("Calling DistanceMatrix")

	resp, err := client.DistanceMatrix(context.Background(), r)
	if err != nil {
		log.WithFields(log.Fields{
			"origin": origin,
			"err":    err,
		}).Error("Error calling DistanceMatrix")

		response.StatusCode = 500
		response.StatusMessage = "Internal error calling Google API"

		return response
	}

	var results []ResultPoint

	for i := 0; i < len(gridPoints); i++ {
		elem := resp.Rows[0].Elements[i]

		results = append(results, ResultPoint{
			Destination:       toCoordinate(gridPoints[i]),
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
		"origin":  origin,
		"results": results,
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

func toCoordinate(latlng maps.LatLng) Coordinate {
	return Coordinate{
		Lat: latlng.Lat,
		Lng: latlng.Lng,
	}
}
