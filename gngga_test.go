package nmea

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGNGGAGoodSentence(t *testing.T) {
	goodMsg := "$GNGGA,203415.000,6325.6138,N,01021.4290,E,1,8,2.42,72.5,M,41.5,M,,*7C"
	sentence, err := Parse(goodMsg)

	assert.NoError(t, err, "Unexpected error parsing good sentence")

	lat, _ := NewLatLong("6325.6138 N")
	lon, _ := NewLatLong("01021.4290 E")
	// Attributes of the parsed sentence, and their expected values.
	expected := GNGGA{
		Sentence: Sentence{
			Type:     "GNGGA",
			Fields:   []string{"203415.000", "6325.6138", "N", "01021.4290", "E", "1", "8", "2.42", "72.5", "M", "41.5", "M", "", ""},
			Checksum: "7C",
			Raw:      "$GNGGA,203415.000,6325.6138,N,01021.4290,E,1,8,2.42,72.5,M,41.5,M,,*7C",
		},
		Time:          "203415.000",
		Latitude:      lat,
		Longitude:     lon,
		FixQuality:    GPS,
		NumSatellites: "8",
		HDOP:          "2.42",
		Altitude:      "72.5",
		Separation:    "41.5",
		DGPSAge:       "",
		DGPSId:        "",
	}

	assert.EqualValues(t, expected, sentence, "Sentence values do not match")
}

func TestGNGGABadType(t *testing.T) {
	badType := "$GPRMC,220516,A,5133.82,N,00042.24,W,173.8,231.8,130694,004.2,W*70"
	s, err := Parse(badType)

	assert.NoError(t, err, "Unexpected error parsing sentence")
	assert.NotEqual(t, "GNGGA", s.GetSentence().Type, "Unexpected sentence type")
}

func TestGNGGABadLatitude(t *testing.T) {
	badLat := "$GNGGA,034225.077,A,S,15124.5567,E,1,03,9.7,-25.0,M,21.0,M,,0000*24"
	_, err := Parse(badLat)

	assert.Error(t, err, "Parse error not returned")
	assert.Equal(t, "GNGGA decode error: cannot parse [A S], unknown format", err.Error(), "Error message does not match")
}

func TestGNGGABadLongitude(t *testing.T) {
	badLon := "$GNGGA,034225.077,3356.4650,S,A,E,1,03,9.7,-25.0,M,21.0,M,,0000*12"
	_, err := Parse(badLon)

	assert.Error(t, err, "Parse error not returned")
	assert.Equal(t, "GNGGA decode error: cannot parse [A E], unknown format", err.Error(), "Error message does not match")
}

func TestGNGGABadFixQuality(t *testing.T) {
	// Make sure bad fix mode is detected.
	badMode := "$GNGGA,034225.077,3356.4650,S,15124.5567,E,10,03,9.7,-25.0,M,21.0,M,,0000*7F"
	_, err := Parse(badMode)

	assert.Error(t, err, "Parse error not returned")
	assert.Equal(t, "Invalid fix quality [10]", err.Error(), "Error message not as expected")
}
