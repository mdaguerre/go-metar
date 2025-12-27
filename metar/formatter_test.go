package metar

import (
	"strings"
	"testing"
)

func TestFormatWind(t *testing.T) {
	tests := []struct {
		name     string
		dir      any
		speed    int
		gust     int
		expected string
	}{
		{
			name:     "calm winds",
			dir:      float64(0),
			speed:    0,
			gust:     0,
			expected: "Calm",
		},
		{
			name:     "numeric direction",
			dir:      float64(270),
			speed:    10,
			gust:     0,
			expected: "270° at 10 kt",
		},
		{
			name:     "numeric direction with gust",
			dir:      float64(180),
			speed:    15,
			gust:     25,
			expected: "180° at 15 kt, gusting 25 kt",
		},
		{
			name:     "variable winds",
			dir:      "VRB",
			speed:    5,
			gust:     0,
			expected: "Variable at 5 kt",
		},
		{
			name:     "variable winds with gust",
			dir:      "VRB",
			speed:    8,
			gust:     15,
			expected: "Variable at 8 kt, gusting 15 kt",
		},
		{
			name:     "string direction",
			dir:      "360",
			speed:    12,
			gust:     0,
			expected: "360° at 12 kt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatWind(tt.dir, tt.speed, tt.gust)
			if result != tt.expected {
				t.Errorf("formatWind(%v, %d, %d) = %q, want %q",
					tt.dir, tt.speed, tt.gust, result, tt.expected)
			}
		})
	}
}

func TestFormatVisibility(t *testing.T) {
	tests := []struct {
		name     string
		vis      any
		expected string
	}{
		{
			name:     "10+ statute miles",
			vis:      float64(10),
			expected: "10+ SM",
		},
		{
			name:     "greater than 10",
			vis:      float64(15),
			expected: "10+ SM",
		},
		{
			name:     "limited visibility",
			vis:      float64(3),
			expected: "3 SM",
		},
		{
			name:     "string visibility",
			vis:      "10+",
			expected: "10+ SM",
		},
		{
			name:     "unknown type",
			vis:      nil,
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatVisibility(tt.vis)
			if result != tt.expected {
				t.Errorf("formatVisibility(%v) = %q, want %q",
					tt.vis, result, tt.expected)
			}
		})
	}
}

func TestExpandCloudCover(t *testing.T) {
	tests := []struct {
		abbrev   string
		expected string
	}{
		{"SKC", "Clear"},
		{"CLR", "Clear"},
		{"FEW", "Few"},
		{"SCT", "Scattered"},
		{"BKN", "Broken"},
		{"OVC", "Overcast"},
		{"OVX", "Obscured"},
		{"UNKNOWN", "UNKNOWN"}, // unknown codes returned as-is
	}

	for _, tt := range tests {
		t.Run(tt.abbrev, func(t *testing.T) {
			result := expandCloudCover(tt.abbrev)
			if result != tt.expected {
				t.Errorf("expandCloudCover(%q) = %q, want %q",
					tt.abbrev, result, tt.expected)
			}
		})
	}
}

func TestFormatClouds(t *testing.T) {
	tests := []struct {
		name     string
		clouds   []Cloud
		expected string
	}{
		{
			name:     "single layer",
			clouds:   []Cloud{{Cover: "FEW", Base: 3000}},
			expected: "Few @ 3000 ft",
		},
		{
			name:     "multiple layers",
			clouds:   []Cloud{{Cover: "SCT", Base: 2500}, {Cover: "BKN", Base: 5000}},
			expected: "Scattered @ 2500 ft, Broken @ 5000 ft",
		},
		{
			name:     "no base reported",
			clouds:   []Cloud{{Cover: "OVC", Base: 0}},
			expected: "Overcast",
		},
		{
			name:     "empty clouds",
			clouds:   []Cloud{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatClouds(tt.clouds)
			if result != tt.expected {
				t.Errorf("formatClouds(%v) = %q, want %q",
					tt.clouds, result, tt.expected)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	metar := &METAR{
		StationID:   "KJFK",
		Name:        "John F Kennedy International",
		Temp:        15,
		Dewpoint:    10,
		Wind:        float64(270),
		WindSpeed:   10,
		WindGust:    0,
		Visibility:  float64(10),
		Altimeter:   1013.25,
		FlightRules: "VFR",
		Clouds:      []Cloud{{Cover: "FEW", Base: 5000}},
		ObsTime:     1704200000,
	}

	result := Decode(metar)

	// Check that key elements are present in the output
	checks := []string{
		"KJFK",
		"John F Kennedy International",
		"VFR",
		"270° at 10 kt",
		"10+ SM",
		"15°C",
		"10°C",
		"Few @ 5000 ft",
	}

	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf("Decode() output missing %q", check)
		}
	}
}
