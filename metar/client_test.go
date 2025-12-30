package metar

import (
	"strings"
	"testing"
)

func TestFetchValidation(t *testing.T) {
	tests := []struct {
		name        string
		icao        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid ICAO code",
			icao:        "KJFK",
			expectError: false,
		},
		{
			name:        "lowercase converts to uppercase",
			icao:        "kjfk",
			expectError: false,
		},
		{
			name:        "too short",
			icao:        "JFK",
			expectError: true,
			errorMsg:    "must be 4 characters",
		},
		{
			name:        "too long",
			icao:        "KJFKX",
			expectError: true,
			errorMsg:    "must be 4 characters",
		},
		{
			name:        "empty string",
			icao:        "",
			expectError: true,
			errorMsg:    "must be 4 characters",
		},
		{
			name:        "contains special characters",
			icao:        "KJ@K",
			expectError: true,
			errorMsg:    "must contain only letters and numbers",
		},
		{
			name:        "contains spaces",
			icao:        "KJ K",
			expectError: true,
			errorMsg:    "must contain only letters and numbers",
		},
		{
			name:        "valid with numbers",
			icao:        "K1FK",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can only test validation without hitting the network
			// For invalid inputs, Fetch should return an error before making any request
			if tt.expectError {
				_, err := Fetch(tt.icao)
				if err == nil {
					t.Errorf("Fetch(%q) expected error, got nil", tt.icao)
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Fetch(%q) error = %q, want error containing %q",
						tt.icao, err.Error(), tt.errorMsg)
				}
			}
			// Skip valid ICAO tests as they would hit the network
			// Those are covered by integration tests
		})
	}
}

// TestFetchIntegration tests the actual API call.
// Run with: go test -run TestFetchIntegration -integration
func TestFetchIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Test with a well-known airport that should always have METAR data
	metar, err := Fetch("KJFK")
	if err != nil {
		t.Fatalf("Fetch(KJFK) unexpected error: %v", err)
	}

	// Verify the response has expected fields
	if metar.StationID != "KJFK" {
		t.Errorf("StationID = %q, want KJFK", metar.StationID)
	}

	if metar.Raw == "" {
		t.Error("Raw METAR string is empty")
	}

	// Flight rules should be one of the valid values
	validFlightRules := map[string]bool{
		"VFR": true, "MVFR": true, "IFR": true, "LIFR": true,
	}
	if !validFlightRules[metar.FlightRules] {
		t.Errorf("FlightRules = %q, want one of VFR/MVFR/IFR/LIFR", metar.FlightRules)
	}
}

// TestFetchInvalidStation tests that an invalid station returns an error.
func TestFetchInvalidStation(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, err := Fetch("ZZZZ")
	if err == nil {
		t.Error("Fetch(ZZZZ) expected error for invalid station, got nil")
	}
}

func TestValidateICAO(t *testing.T) {
	tests := []struct {
		name        string
		icao        string
		want        string
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid uppercase",
			icao: "KJFK",
			want: "KJFK",
		},
		{
			name: "converts lowercase to uppercase",
			icao: "kjfk",
			want: "KJFK",
		},
		{
			name: "mixed case",
			icao: "KjFk",
			want: "KJFK",
		},
		{
			name: "valid with numbers",
			icao: "K1FK",
			want: "K1FK",
		},
		{
			name:        "too short",
			icao:        "JFK",
			expectError: true,
			errorMsg:    "must be 4 characters",
		},
		{
			name:        "too long",
			icao:        "KJFKX",
			expectError: true,
			errorMsg:    "must be 4 characters",
		},
		{
			name:        "empty string",
			icao:        "",
			expectError: true,
			errorMsg:    "must be 4 characters",
		},
		{
			name:        "contains special characters",
			icao:        "KJ@K",
			expectError: true,
			errorMsg:    "must contain only letters and numbers",
		},
		{
			name:        "contains spaces",
			icao:        "KJ K",
			expectError: true,
			errorMsg:    "must contain only letters and numbers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateICAO(tt.icao)

			if tt.expectError {
				if err == nil {
					t.Errorf("ValidateICAO(%q) expected error, got nil", tt.icao)
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("ValidateICAO(%q) error = %q, want error containing %q",
						tt.icao, err.Error(), tt.errorMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("ValidateICAO(%q) unexpected error: %v", tt.icao, err)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateICAO(%q) = %q, want %q", tt.icao, got, tt.want)
			}
		})
	}
}

func TestFetchMultipleValidation(t *testing.T) {
	tests := []struct {
		name        string
		icaos       []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "empty slice",
			icaos:       []string{},
			expectError: true,
			errorMsg:    "no ICAO codes provided",
		},
		{
			name:        "nil slice",
			icaos:       nil,
			expectError: true,
			errorMsg:    "no ICAO codes provided",
		},
		{
			name:        "one invalid ICAO",
			icaos:       []string{"KJFK", "BAD", "KLAX"},
			expectError: true,
			errorMsg:    "must be 4 characters",
		},
		{
			name:        "invalid characters in one",
			icaos:       []string{"KJFK", "KL@X"},
			expectError: true,
			errorMsg:    "must contain only letters and numbers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FetchMultiple(tt.icaos)
			if !tt.expectError {
				return // Skip valid cases, they would hit the network
			}
			if err == nil {
				t.Errorf("FetchMultiple(%v) expected error, got nil", tt.icaos)
				return
			}
			if !strings.Contains(err.Error(), tt.errorMsg) {
				t.Errorf("FetchMultiple(%v) error = %q, want error containing %q",
					tt.icaos, err.Error(), tt.errorMsg)
			}
		})
	}
}

// TestFetchMultipleIntegration tests fetching multiple airports from the API.
func TestFetchMultipleIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	icaos := []string{"KJFK", "KLAX"}
	metars, err := FetchMultiple(icaos)
	if err != nil {
		t.Fatalf("FetchMultiple(%v) unexpected error: %v", icaos, err)
	}

	if len(metars) != 2 {
		t.Errorf("FetchMultiple(%v) returned %d results, want 2", icaos, len(metars))
	}

	// Verify both airports are present (order may vary)
	found := make(map[string]bool)
	for _, m := range metars {
		found[m.StationID] = true
		if m.Raw == "" {
			t.Errorf("METAR for %s has empty Raw string", m.StationID)
		}
	}

	for _, icao := range icaos {
		if !found[icao] {
			t.Errorf("FetchMultiple(%v) missing result for %s", icaos, icao)
		}
	}
}

// TestFetchMultipleSingleAirport verifies FetchMultiple works with a single airport.
func TestFetchMultipleSingleAirport(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	metars, err := FetchMultiple([]string{"KJFK"})
	if err != nil {
		t.Fatalf("FetchMultiple([KJFK]) unexpected error: %v", err)
	}

	if len(metars) != 1 {
		t.Errorf("FetchMultiple([KJFK]) returned %d results, want 1", len(metars))
	}

	if metars[0].StationID != "KJFK" {
		t.Errorf("StationID = %q, want KJFK", metars[0].StationID)
	}
}

// TestFetchTAFValidation tests TAF fetch validation.
func TestFetchTAFValidation(t *testing.T) {
	tests := []struct {
		name        string
		icao        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "too short",
			icao:        "JFK",
			expectError: true,
			errorMsg:    "must be 4 characters",
		},
		{
			name:        "invalid characters",
			icao:        "KJ@K",
			expectError: true,
			errorMsg:    "must contain only letters and numbers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FetchTAF(tt.icao)
			if !tt.expectError {
				return
			}
			if err == nil {
				t.Errorf("FetchTAF(%q) expected error, got nil", tt.icao)
				return
			}
			if !strings.Contains(err.Error(), tt.errorMsg) {
				t.Errorf("FetchTAF(%q) error = %q, want error containing %q",
					tt.icao, err.Error(), tt.errorMsg)
			}
		})
	}
}

// TestFetchTAFIntegration tests fetching TAF from the API.
func TestFetchTAFIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	taf, err := FetchTAF("KJFK")
	if err != nil {
		t.Fatalf("FetchTAF(KJFK) unexpected error: %v", err)
	}

	if taf.StationID != "KJFK" {
		t.Errorf("StationID = %q, want KJFK", taf.StationID)
	}

	if taf.RawTAF == "" {
		t.Error("RawTAF is empty")
	}

	if len(taf.Forecasts) == 0 {
		t.Error("Forecasts slice is empty")
	}
}

// TestFetchMultipleTAFValidation tests TAF multi-fetch validation.
func TestFetchMultipleTAFValidation(t *testing.T) {
	tests := []struct {
		name        string
		icaos       []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "empty slice",
			icaos:       []string{},
			expectError: true,
			errorMsg:    "no ICAO codes provided",
		},
		{
			name:        "one invalid ICAO",
			icaos:       []string{"KJFK", "BAD"},
			expectError: true,
			errorMsg:    "must be 4 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FetchMultipleTAF(tt.icaos)
			if !tt.expectError {
				return
			}
			if err == nil {
				t.Errorf("FetchMultipleTAF(%v) expected error, got nil", tt.icaos)
				return
			}
			if !strings.Contains(err.Error(), tt.errorMsg) {
				t.Errorf("FetchMultipleTAF(%v) error = %q, want error containing %q",
					tt.icaos, err.Error(), tt.errorMsg)
			}
		})
	}
}

// TestFetchMultipleTAFIntegration tests fetching multiple TAFs from the API.
func TestFetchMultipleTAFIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	icaos := []string{"KJFK", "KLAX"}
	tafs, err := FetchMultipleTAF(icaos)
	if err != nil {
		t.Fatalf("FetchMultipleTAF(%v) unexpected error: %v", icaos, err)
	}

	if len(tafs) != 2 {
		t.Errorf("FetchMultipleTAF(%v) returned %d results, want 2", icaos, len(tafs))
	}

	found := make(map[string]bool)
	for _, taf := range tafs {
		found[taf.StationID] = true
		if taf.RawTAF == "" {
			t.Errorf("TAF for %s has empty RawTAF", taf.StationID)
		}
	}

	for _, icao := range icaos {
		if !found[icao] {
			t.Errorf("FetchMultipleTAF(%v) missing result for %s", icaos, icao)
		}
	}
}
