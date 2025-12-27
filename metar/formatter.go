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
)

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

// formatFlightLine creates a color-coded flight rules line
func formatFlightLine(fr string) string {
	var style lipgloss.Style

	switch fr {
	case "VFR":
		style = lipgloss.NewStyle().Foreground(vfrColor).Bold(true)
	case "MVFR":
		style = lipgloss.NewStyle().Foreground(mvfrColor).Bold(true)
	case "IFR":
		style = lipgloss.NewStyle().Foreground(ifrColor).Bold(true)
	case "LIFR":
		style = lipgloss.NewStyle().Foreground(lifrColor).Bold(true)
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
	coverMap := map[string]string{
		"SKC": "Clear",
		"CLR": "Clear",
		"FEW": "Few",
		"SCT": "Scattered",
		"BKN": "Broken",
		"OVC": "Overcast",
		"OVX": "Obscured",
	}

	if expanded, ok := coverMap[cover]; ok {
		return expanded
	}
	return cover
}
