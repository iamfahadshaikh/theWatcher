package graph

import (
	"fmt"
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
