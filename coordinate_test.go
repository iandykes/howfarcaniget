package main

import "testing"

func TestString(t *testing.T) {

	testCases := []struct {
		lat      float64
		lng      float64
		expected string
	}{
		{0, 0, "0.000000,0.000000"},
		{51.234567, -1.654321, "51.234567,-1.654321"},
	}

	for _, testCase := range testCases {
		sut := Coordinate{testCase.lat, testCase.lng}
		actual := sut.String()

		if testCase.expected != actual {
			t.Errorf("%v != %v", testCase.expected, actual)
		}
	}

}
