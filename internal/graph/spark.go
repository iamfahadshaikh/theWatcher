package graph

import (
	"math"

	"github.com/charmbracelet/lipgloss"
)

var sparkBlocks = []rune{' ', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

type SparkRenderer struct{}

func NewSparkRenderer() *SparkRenderer {
	return &SparkRenderer{}
}

func (s *SparkRenderer) Render(data GraphData) string {
	if len(data.Samples) == 0 || data.Width <= 0 {
		return ""
	}

	points := data.Samples
	if len(points) > data.Width {
		points = points[len(points)-data.Width:]
	}

	minVal := data.Min
	maxVal := data.Max
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

	color := data.Color
	if color == "" {
		color = "#eb6f92"
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(string(runes))
}

// RenderSparkline helper function
func RenderSparkline(data []float64, width int, style lipgloss.Style) string {
	r := NewSparkRenderer()
	return r.Render(GraphData{
		Samples: data,
		Width:   width,
		Min:     0,
		Max:     100,
	})
}
