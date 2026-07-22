package graph

import (
	"fmt"
	"strings"
)

// CalculateBounds computes min and max for auto-scaling graph samples.
func CalculateBounds(samples []float64, minScale, maxScale float64) (float64, float64) {
	minV := minScale
	maxV := maxScale
	for _, v := range samples {
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
	return minV, maxV
}

// FormatAxisLabel formats numeric bounds into string labels with units.
func FormatAxisLabel(val float64, unit string) string {
	return fmt.Sprintf("%.1f%s", val, unit)
}

// AddYAxis prefixes the lines of a rendered graph with Y-axis tick marks.
func AddYAxis(graphStr string, height int, minVal, maxVal float64, suffix string) string {
	lines := strings.Split(graphStr, "\n")
	if len(lines) == 0 {
		return graphStr
	}
	result := make([]string, len(lines))
	for i, line := range lines {
		if i >= height {
			result[i] = "      " + line
			continue
		}
		val := maxVal - (float64(i)/float64(height-1))*(maxVal-minVal)
		var label string
		if i == 0 {
			label = fmt.Sprintf("%4.0f%s ┤ ", val, suffix)
		} else if i == height/2 {
			label = fmt.Sprintf("%4.0f%s ┤ ", val, suffix)
		} else if i == height-1 {
			label = fmt.Sprintf("%4.0f%s └ ", val, suffix)
		} else {
			label = "      │ "
		}
		result[i] = label + line
	}
	return strings.Join(result, "\n")
}
