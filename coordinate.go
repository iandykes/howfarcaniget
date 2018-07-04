package main

import "fmt"

// Coordinate is a latitude/longitude pair
type Coordinate struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// String returns the string value of the Coordinate
func (c Coordinate) String() string {
	return fmt.Sprintf("%f,%f", c.Lat, c.Lng)
}
