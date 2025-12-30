package metar

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Color definitions for flight rules
var (
	vfrColor  = lipgloss.Color("#22c55e") // Green (good flying conditions)
	mvfrColor = lipgloss.Color("#eab308") // Yellow/Orange (marginal conditions)
	ifrColor  = lipgloss.Color("#ef4444") // Red (instrument required)
	lifrColor = lipgloss.Color("#d946ef") // Magenta (very poor conditions)

	// UI colors
	headerColor  = lipgloss.Color("#60a5fa") // Light blue
	labelColor   = lipgloss.Color("#9ca3af") // Gray
	valueColor   = lipgloss.Color("#f3f4f6") // White-ish
	borderColor  = lipgloss.Color("#4b5563") // Dark gray
	stationColor = lipgloss.Color("#f3f4f6") // White
)

// Styles
var (
	boxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(0, 1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(headerColor)

	stationStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(stationColor)

	labelStyle = lipgloss.NewStyle().
			Foreground(labelColor)

	valueStyle = lipgloss.NewStyle().
			Foreground(valueColor)

	// Flight rules styles - pre-defined for reuse
	vfrStyle   = lipgloss.NewStyle().Foreground(vfrColor).Bold(true)
	mvfrStyle  = lipgloss.NewStyle().Foreground(mvfrColor).Bold(true)
	ifrStyle   = lipgloss.NewStyle().Foreground(ifrColor).Bold(true)
	lifrStyle  = lipgloss.NewStyle().Foreground(lifrColor).Bold(true)
)

// coverMap maps cloud cover abbreviations to full descriptions.
// Defined at package level for efficiency (avoids recreating on each call).
var coverMap = map[string]string{
	"SKC": "Clear",
	"CLR": "Clear",
	"FEW": "Few",
	"SCT": "Scattered",
	"BKN": "Broken",
	"OVC": "Overcast",
	"OVX": "Obscured",
}

// weatherMap maps weather phenomenon codes to human-readable descriptions.
var weatherMap = map[string]string{
	// Intensity
	"-":  "Light",
	"+":  "Heavy",
	"VC": "Vicinity",

	// Descriptor
	"MI": "Shallow",
	"PR": "Partial",
	"BC": "Patches",
	"DR": "Drifting",
	"BL": "Blowing",
	"SH": "Showers",
	"TS": "Thunderstorm",
	"FZ": "Freezing",

	// Precipitation
	"DZ": "Drizzle",
	"RA": "Rain",
	"SN": "Snow",
	"SG": "Snow Grains",
	"IC": "Ice Crystals",
	"PL": "Ice Pellets",
	"GR": "Hail",
	"GS": "Small Hail",
	"UP": "Unknown Precip",

	// Obscuration
	"BR": "Mist",
	"FG": "Fog",
	"FU": "Smoke",
	"VA": "Volcanic Ash",
	"DU": "Dust",
	"SA": "Sand",
	"HZ": "Haze",
	"PY": "Spray",

	// Other
	"PO": "Dust Whirls",
	"SQ": "Squalls",
	"FC": "Funnel Cloud",
	"SS": "Sandstorm",
	"DS": "Duststorm",
}

// Decode converts a METAR struct into a styled, human-readable string.
func Decode(m *METAR) string {
	var sb strings.Builder

	// Station header
	stationText := stationStyle.Render(m.StationID)
	if m.Name != "" {
		stationText += labelStyle.Render(" · ") + valueStyle.Render(m.Name)
	}
	sb.WriteString(stationText + "\n")

	// Observation time
	if m.ObsTime > 0 {
		obsTime := time.Unix(m.ObsTime, 0).UTC()
		sb.WriteString(formatLine("Time", obsTime.Format("02 Jan 2006 15:04")+" UTC"))
	}

	// Flight category with color
	sb.WriteString(formatFlightLine(m.FlightRules))

	// Weather data
	sb.WriteString(formatLine("Wind", formatWind(m.Wind, m.WindSpeed, m.WindGust)))
	sb.WriteString(formatLine("Visibility", formatVisibility(m.Visibility)))
	sb.WriteString(formatLine("Temp", fmt.Sprintf("%.0f°C (Dewpoint: %.0f°C)", m.Temp, m.Dewpoint)))

	// Altimeter
	altInHg := m.Altimeter * 0.02953
	sb.WriteString(formatLine("Altimeter", fmt.Sprintf("%.2f inHg / %.0f hPa", altInHg, m.Altimeter)))

	// Clouds (last line, no trailing newline)
	cloudsLabel := labelStyle.Render(fmt.Sprintf("%-11s", "Clouds"))
	if len(m.Clouds) > 0 {
		sb.WriteString(cloudsLabel + valueStyle.Render(formatClouds(m.Clouds)))
	} else {
		sb.WriteString(cloudsLabel + valueStyle.Render("Clear"))
	}

	// Wrap in box
	return boxStyle.Render(sb.String())
}

// formatLine creates a styled label: value line
func formatLine(label, value string) string {
	paddedLabel := fmt.Sprintf("%-11s", label)
	return labelStyle.Render(paddedLabel) + valueStyle.Render(value) + "\n"
}

// formatTAFLine creates a styled indented line for TAF forecast details
func formatTAFLine(label, value string) string {
	paddedLabel := fmt.Sprintf("  %-9s", label)
	return labelStyle.Render(paddedLabel) + valueStyle.Render(value) + "\n"
}

// formatFlightLine creates a color-coded flight rules line
func formatFlightLine(fr string) string {
	var style lipgloss.Style

	switch fr {
	case "VFR":
		style = vfrStyle
	case "MVFR":
		style = mvfrStyle
	case "IFR":
		style = ifrStyle
	case "LIFR":
		style = lifrStyle
	default:
		style = valueStyle
	}

	paddedLabel := fmt.Sprintf("%-11s", "Flight")
	return labelStyle.Render(paddedLabel) + style.Render(fr) + "\n"
}

// formatWind converts wind data to a readable string.
func formatWind(dir any, speed, gust int) string {
	if speed == 0 {
		return "Calm"
	}

	var result string

	switch d := dir.(type) {
	case string:
		if d == "VRB" {
			result = fmt.Sprintf("Variable at %d kt", speed)
		} else {
			result = fmt.Sprintf("%s° at %d kt", d, speed)
		}
	case float64:
		result = fmt.Sprintf("%.0f° at %d kt", d, speed)
	default:
		result = fmt.Sprintf("%d kt", speed)
	}

	if gust > 0 {
		result += fmt.Sprintf(", gusting %d kt", gust)
	}

	return result
}

// formatVisibility makes visibility human-readable.
func formatVisibility(vis any) string {
	if s, ok := vis.(string); ok {
		return s + " SM"
	}

	v, ok := vis.(float64)
	if !ok {
		return "Unknown"
	}

	if v >= 10 {
		return "10+ SM"
	}
	return fmt.Sprintf("%.0f SM", v)
}

// formatClouds converts cloud layers to readable text.
func formatClouds(clouds []Cloud) string {
	descriptions := make([]string, 0, len(clouds))

	for _, c := range clouds {
		cover := expandCloudCover(c.Cover)
		if c.Base > 0 {
			descriptions = append(descriptions, fmt.Sprintf("%s @ %d ft", cover, c.Base))
		} else {
			descriptions = append(descriptions, cover)
		}
	}

	return strings.Join(descriptions, ", ")
}

// expandCloudCover converts abbreviations to full words.
func expandCloudCover(cover string) string {
	if expanded, ok := coverMap[cover]; ok {
		return expanded
	}
	return cover
}

// decodeWeather converts weather codes like "-RA BR" to "Light Rain, Mist".
func decodeWeather(wxString string) string {
	if wxString == "" {
		return ""
	}

	// Split by spaces to get individual weather groups
	groups := strings.Fields(wxString)
	decoded := make([]string, 0, len(groups))

	for _, group := range groups {
		decoded = append(decoded, decodeWeatherGroup(group))
	}

	return strings.Join(decoded, ", ")
}

// decodeWeatherGroup decodes a single weather group like "-RA" or "TSRA".
func decodeWeatherGroup(group string) string {
	if group == "" {
		return ""
	}

	var parts []string
	remaining := group

	// Check for intensity prefix (- or +)
	if len(remaining) > 0 && (remaining[0] == '-' || remaining[0] == '+') {
		if desc, ok := weatherMap[string(remaining[0])]; ok {
			parts = append(parts, desc)
		}
		remaining = remaining[1:]
	}

	// Check for VC (vicinity) prefix
	if strings.HasPrefix(remaining, "VC") {
		parts = append(parts, weatherMap["VC"])
		remaining = remaining[2:]
	}

	// Process remaining codes in 2-character chunks
	for len(remaining) >= 2 {
		code := remaining[:2]
		if desc, ok := weatherMap[code]; ok {
			parts = append(parts, desc)
		} else {
			parts = append(parts, code) // Keep unknown codes as-is
		}
		remaining = remaining[2:]
	}

	// Handle any leftover single character
	if len(remaining) > 0 {
		parts = append(parts, remaining)
	}

	return strings.Join(parts, " ")
}

// TAF header style
var tafHeaderStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#a78bfa")). // Purple for TAF header
	Bold(true)

// DecodeTAF converts a TAF struct into a styled, human-readable string.
func DecodeTAF(t *TAF) string {
	var sb strings.Builder

	// Station header
	stationText := stationStyle.Render(t.StationID)
	if t.Name != "" {
		stationText += labelStyle.Render(" · ") + valueStyle.Render(t.Name)
	}
	sb.WriteString(stationText + "\n")

	// TAF label
	sb.WriteString(tafHeaderStyle.Render("TAF FORECAST") + "\n")

	// Valid period
	if t.ValidTimeFrom > 0 && t.ValidTimeTo > 0 {
		from := time.Unix(t.ValidTimeFrom, 0).UTC()
		to := time.Unix(t.ValidTimeTo, 0).UTC()
		sb.WriteString(formatLine("Valid", fmt.Sprintf("%s to %s UTC",
			from.Format("02 Jan 15:04"), to.Format("02 Jan 15:04"))))
	}

	// Forecast periods
	for i, f := range t.Forecasts {
		sb.WriteString(formatTAFForecast(f, i == 0, i == len(t.Forecasts)-1))
	}

	return boxStyle.Render(sb.String())
}

// Separator style for TAF periods
var separatorStyle = lipgloss.NewStyle().Foreground(borderColor)

// formatTAFForecast formats a single TAF forecast period.
func formatTAFForecast(f TAFForecast, isFirst, isLast bool) string {
	var sb strings.Builder

	// Add separator before non-first forecast periods
	if !isFirst {
		sb.WriteString(separatorStyle.Render("────────────────────────────") + "\n")
	}

	// Time period with change indicator
	fromTime := time.Unix(f.TimeFrom, 0).UTC()
	toTime := time.Unix(f.TimeTo, 0).UTC()

	var prefix string
	switch f.FcstChange {
	case "FM":
		prefix = "From  "
	case "TEMPO":
		prefix = "Tempo "
	case "BECMG":
		prefix = "Becmg "
	case "PROB":
		if f.Probability != nil {
			prefix = fmt.Sprintf("Prob%-2d", *f.Probability)
		} else {
			prefix = "Prob  "
		}
	default:
		prefix = "Init  "
	}

	// Format time with day name (e.g., "Sun 18:00 - Mon 00:00")
	timeStr := fmt.Sprintf("%s%s %s - %s %s",
		prefix,
		fromTime.Format("Mon"),
		fromTime.Format("15:04"),
		toTime.Format("Mon"),
		toTime.Format("15:04"))
	sb.WriteString(headerStyle.Render(timeStr) + "\n")

	// Wind
	if f.WindSpeed > 0 {
		var gust int
		if f.WindGust != nil {
			gust = *f.WindGust
		}
		sb.WriteString(formatTAFLine("Wind", formatWind(f.WindDir, f.WindSpeed, gust)))
	}

	// Visibility
	if f.Visibility != nil && f.Visibility != "" {
		sb.WriteString(formatTAFLine("Visib", formatVisibility(f.Visibility)))
	}

	// Weather (decoded)
	if f.Weather != "" {
		decoded := decodeWeather(f.Weather)
		sb.WriteString(formatTAFLine("Weather", decoded))
	}

	// Clouds
	if len(f.Clouds) > 0 {
		cloudsLine := formatTAFLine("Clouds", formatClouds(f.Clouds))
		if isLast {
			// Remove trailing newline for last item
			cloudsLine = strings.TrimSuffix(cloudsLine, "\n")
		}
		sb.WriteString(cloudsLine)
	} else if isLast {
		// Remove trailing newline if clouds was the last item added
		str := sb.String()
		if strings.HasSuffix(str, "\n") {
			return strings.TrimSuffix(str, "\n")
		}
	}

	return sb.String()
}
