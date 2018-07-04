package main

import log "github.com/sirupsen/logrus"

// TestDistanceCalculator is a DistanceCalculator for testing without calling an external service
type TestDistanceCalculator struct {
}

// GetCoordinate always returns 52.920279, -1.469559
func (d *TestDistanceCalculator) GetCoordinate(originInput string) (Coordinate, error) {
	return Coordinate{
		Lat: 52.920279,
		Lng: -1.469559,
	}, nil
}

// GenerateDistances from interface DistanceCalculator returns hard coded data
func (d *TestDistanceCalculator) GenerateDistances(origin Coordinate) DistanceResponse {

	log.WithFields(log.Fields{
		"origin": origin,
	}).Debug("TestDistanceCalculator GenerateDistances")

	degrees := float64(2)
	// Number of degrees on lat and long for each point on the grid
	gridSize := 0.5

	gridMinLat := origin.Lat - degrees
	gridMaxLat := origin.Lat + degrees

	gridMinLng := origin.Lng - degrees
	gridMaxLng := origin.Lng + degrees

	var resultPoints []ResultPoint

	for lat := gridMinLat; lat < gridMaxLat; lat += gridSize {
		for lng := gridMinLng; lng < gridMaxLng; lng += gridSize {
			resultPoints = append(resultPoints, ResultPoint{
				Destination: Coordinate{
					Lat: lat,
					Lng: lng,
				},
				DurationGroup: 1,
			})
		}
	}

	return DistanceResponse{
		StatusCode:    200,
		StatusMessage: "Fake News!",
		Origin:        origin,
		Points:        resultPoints,
	}
}
