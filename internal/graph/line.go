package graph

import (
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderLineGraph plots a time series array onto a 2D grid canvas.
func RenderLineGraph(data []float64, width, height int, minScale, maxScale float64, style lipgloss.Style) string {
	if width <= 0 || height <= 0 {
		return ""
	}

	if len(data) == 0 {
		emptyRow := strings.Repeat(" ", width)
		rows := make([]string, height)
		for i := range rows {
			rows[i] = emptyRow
		}
		return style.Render(strings.Join(rows, "\n"))
	}

	// Truncate or sample data to fit width
	points := data
	if len(points) > width {
		points = points[len(points)-width:]
	}

	// Calculate bounds
	minV := minScale
	maxV := maxScale
	for _, v := range points {
		if v > maxV {
			maxV = v
		}
		if v < minV {
			minV = v
		}
	}
	if maxV <= minV {
		maxV = minV + 1.0
	}

	// Grid canvas: height rows by width columns
	canvas := make([][]rune, height)
	for r := 0; r < height; r++ {
		row := make([]rune, width)
		for c := 0; c < width; c++ {
			row[c] = ' '
		}
		canvas[r] = row
	}

	startCol := width - len(points)
	if startCol < 0 {
		startCol = 0
	}

	for i, v := range points {
		col := startCol + i
		if col < 0 || col >= width {
			continue
		}

		// Map v to row: row 0 is top, height-1 is bottom
		ratio := (v - minV) / (maxV - minV)
		if ratio < 0 {
			ratio = 0
		}
		if ratio > 1 {
			ratio = 1
		}

		row := int(math.Round((1.0 - ratio) * float64(height-1)))
		if row < 0 {
			row = 0
		}
		if row >= height {
			row = height - 1
		}

		canvas[row][col] = '█'
	}

	lines := make([]string, height)
	for r := 0; r < height; r++ {
		lines[r] = string(canvas[r])
	}

	return style.Render(strings.Join(lines, "\n"))
}
