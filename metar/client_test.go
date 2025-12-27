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
