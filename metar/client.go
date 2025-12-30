// Package metar handles fetching and parsing METAR weather data.
// In Go, all files in the same folder share the same package name.
package metar

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"
)

// httpClient is reused across requests to avoid creating a new client each time.
// This is more efficient and follows HTTP best practices.
var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// METAR represents the weather data returned by the API.
// In Go, structs are like classes in other languages.
// The `json:"..."` tags tell Go how to map JSON fields to struct fields.
type METAR struct {
	Raw         string  `json:"rawOb"`         // Raw METAR string
	StationID   string  `json:"icaoId"`        // Airport ICAO code
	Name        string  `json:"name"`          // Airport name
	Temp        float64 `json:"temp"`          // Temperature in Celsius
	Dewpoint    float64 `json:"dewp"`          // Dewpoint in Celsius
	Wind        any     `json:"wdir"`          // Wind direction - can be "VRB" (string) or degrees (number)
	WindSpeed   int     `json:"wspd"`          // Wind speed in knots
	WindGust    int     `json:"wgst"`          // Wind gust in knots (0 if none)
	Visibility  any     `json:"visib"`         // Visibility - can be number or string like "10+"
	Altimeter   float64 `json:"altim"`         // Altimeter in millibars
	FlightRules string  `json:"fltcat"`        // VFR, MVFR, IFR, or LIFR
	Clouds      []Cloud `json:"clouds"`        // Cloud layers
	ObsTime     int64   `json:"obsTime"`       // Observation time (Unix timestamp)
}

// Cloud represents a cloud layer.
type Cloud struct {
	Cover string `json:"cover"` // SKC, FEW, SCT, BKN, OVC
	Base  int    `json:"base"`  // Cloud base in feet AGL
}

// apiResponse wraps the API response which is an array of METARs.
// We only request one, so we'll take the first element.
type apiResponse []METAR

// TAF represents Terminal Aerodrome Forecast data.
type TAF struct {
	StationID     string        `json:"icaoId"`        // Airport ICAO code
	Name          string        `json:"name"`          // Airport name
	RawTAF        string        `json:"rawTAF"`        // Raw TAF string
	IssueTime     string        `json:"issueTime"`     // When the TAF was issued
	ValidTimeFrom int64         `json:"validTimeFrom"` // Start of validity (Unix timestamp)
	ValidTimeTo   int64         `json:"validTimeTo"`   // End of validity (Unix timestamp)
	Forecasts     []TAFForecast `json:"fcsts"`         // Individual forecast periods
}

// TAFForecast represents a single forecast period within a TAF.
type TAFForecast struct {
	TimeFrom     int64   `json:"timeFrom"`     // Period start (Unix timestamp)
	TimeTo       int64   `json:"timeTo"`       // Period end (Unix timestamp)
	FcstChange   string  `json:"fcstChange"`   // Change indicator: FM, TEMPO, BECMG, PROB
	Probability  *int    `json:"probability"`  // Probability percentage (for PROB)
	WindDir      any     `json:"wdir"`         // Wind direction
	WindSpeed    int     `json:"wspd"`         // Wind speed in knots
	WindGust     *int    `json:"wgst"`         // Wind gust in knots
	Visibility   any     `json:"visib"`        // Visibility
	Weather      string  `json:"wxString"`     // Weather phenomena
	Clouds       []Cloud `json:"clouds"`       // Cloud layers
}

// tafAPIResponse wraps the TAF API response.
type tafAPIResponse []TAF

// isAlphanumeric checks if all characters in the string are alphanumeric.
func isAlphanumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// Fetch retrieves METAR data for the given ICAO airport code.
// In Go, function names starting with uppercase are "exported" (public).
// Lowercase names are private to the package.
func Fetch(icao string) (*METAR, error) {
	// Convert to uppercase - ICAO codes are always uppercase
	icao = strings.ToUpper(icao)

	// Validate ICAO code format (4 alphanumeric characters)
	if len(icao) != 4 {
		return nil, fmt.Errorf("invalid ICAO code: must be 4 characters (e.g., KJFK)")
	}
	if !isAlphanumeric(icao) {
		return nil, fmt.Errorf("invalid ICAO code: must contain only letters and numbers")
	}

	// Build the API URL
	// aviationweather.gov provides free METAR data in JSON format
	url := fmt.Sprintf(
		"https://aviationweather.gov/api/data/metar?ids=%s&format=json",
		icao,
	)

	// Make the GET request using the shared HTTP client
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch METAR: %w", err)
	}
	// defer ensures this runs when the function exits, even if there's an error.
	// Always close response bodies to avoid resource leaks!
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// Parse the JSON response
	var data apiResponse
	// json.NewDecoder reads from the response body and decodes into our struct
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check if we got any results
	if len(data) == 0 {
		return nil, fmt.Errorf("no METAR found for %s - check the ICAO code", icao)
	}

	// Return a pointer to the first (and only) METAR
	// The & operator gets the memory address (creates a pointer)
	return &data[0], nil
}

// ValidateICAO checks if an ICAO code is valid (4 alphanumeric characters).
// Returns the uppercase ICAO code and an error if invalid.
func ValidateICAO(icao string) (string, error) {
	icao = strings.ToUpper(icao)

	if len(icao) != 4 {
		return "", fmt.Errorf("invalid ICAO code %q: must be 4 characters", icao)
	}
	if !isAlphanumeric(icao) {
		return "", fmt.Errorf("invalid ICAO code %q: must contain only letters and numbers", icao)
	}

	return icao, nil
}

// FetchMultiple retrieves METAR data for multiple ICAO airport codes in a single request.
// Returns a slice of METARs and any errors encountered during validation.
func FetchMultiple(icaos []string) ([]*METAR, error) {
	if len(icaos) == 0 {
		return nil, fmt.Errorf("no ICAO codes provided")
	}

	// Validate all ICAO codes first
	validICAOs := make([]string, 0, len(icaos))
	for _, icao := range icaos {
		validated, err := ValidateICAO(icao)
		if err != nil {
			return nil, err
		}
		validICAOs = append(validICAOs, validated)
	}

	// Build the API URL with comma-separated ICAOs
	url := fmt.Sprintf(
		"https://aviationweather.gov/api/data/metar?ids=%s&format=json",
		strings.Join(validICAOs, ","),
	)

	// Make the GET request
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch METAR: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// Parse the JSON response
	var data apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("no METAR data found for the requested airports")
	}

	// Convert to pointer slice
	result := make([]*METAR, len(data))
	for i := range data {
		result[i] = &data[i]
	}

	return result, nil
}

// FetchTAF retrieves TAF data for the given ICAO airport code.
func FetchTAF(icao string) (*TAF, error) {
	icao, err := ValidateICAO(icao)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(
		"https://aviationweather.gov/api/data/taf?ids=%s&format=json",
		icao,
	)

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch TAF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var data tafAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("no TAF found for %s - check the ICAO code", icao)
	}

	return &data[0], nil
}

// FetchMultipleTAF retrieves TAF data for multiple ICAO airport codes.
func FetchMultipleTAF(icaos []string) ([]*TAF, error) {
	if len(icaos) == 0 {
		return nil, fmt.Errorf("no ICAO codes provided")
	}

	validICAOs := make([]string, 0, len(icaos))
	for _, icao := range icaos {
		validated, err := ValidateICAO(icao)
		if err != nil {
			return nil, err
		}
		validICAOs = append(validICAOs, validated)
	}

	url := fmt.Sprintf(
		"https://aviationweather.gov/api/data/taf?ids=%s&format=json",
		strings.Join(validICAOs, ","),
	)

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch TAF: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var data tafAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("no TAF data found for the requested airports")
	}

	result := make([]*TAF, len(data))
	for i := range data {
		result[i] = &data[i]
	}

	return result, nil
}
