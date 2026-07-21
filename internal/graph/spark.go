package graph

import (
	"math"

	"github.com/charmbracelet/lipgloss"
)

var sparkBlocks = []rune{' ', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

// RenderSparkline converts a slice of floats into a single-line block sparkline.
func RenderSparkline(data []float64, width int, style lipgloss.Style) string {
	if len(data) == 0 || width <= 0 {
		return ""
	}

	// Fit data to width
	points := data
	if len(data) > width {
		points = data[len(data)-width:]
	}

	minVal := 0.0
	maxVal := 100.0 // Default percent scale if points within 0..100
	for _, v := range points {
		if v > maxVal {
			maxVal = v
		}
		if v < minVal {
			minVal = v
		}
	}

	rangeVal := maxVal - minVal
	if rangeVal == 0 {
		rangeVal = 1.0
	}

	runes := make([]rune, len(points))
	for i, v := range points {
		idx := int(math.Floor(((v - minVal) / rangeVal) * float64(len(sparkBlocks)-1)))
		if idx < 0 {
			idx = 0
		}
		if idx >= len(sparkBlocks) {
			idx = len(sparkBlocks) - 1
		}
		runes[i] = sparkBlocks[idx]
	}

	return style.Render(string(runes))
}
