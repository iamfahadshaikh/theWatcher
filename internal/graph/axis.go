package graph

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// ComputeAdaptiveBounds calculates scaleMin and scaleMax from sample data.
// It adds ~10% headroom above and below the visible range, enforces a minimum
// scale range of minRange, and clamps the final bounds to [floorMin, ceilMax].
func ComputeAdaptiveBounds(samples []float64, floorMin, ceilMax, minRange float64) (float64, float64) {
	if len(samples) == 0 {
		return floorMin, floorMin + minRange
	}

	visMin := samples[0]
	visMax := samples[0]
	for _, v := range samples[1:] {
		if v < visMin {
			visMin = v
		}
		if v > visMax {
			visMax = v
		}
	}

	rawRange := visMax - visMin
	headroom := rawRange * 0.10
	if headroom < 0.5 {
		headroom = 0.5
	}

	scaleMin := visMin - headroom
	scaleMax := visMax + headroom

	// Enforce minimum visible range
	if scaleMax-scaleMin < minRange {
		diff := minRange - (scaleMax - scaleMin)
		scaleMin -= diff / 2
		scaleMax += diff / 2
	}

	// Clamp to [floorMin, ceilMax]
	if scaleMin < floorMin {
		// Shift range up if we hit the floor
		scaleMax += (floorMin - scaleMin)
		scaleMin = floorMin
	}
	if scaleMax > ceilMax {
		// Shift range down if we hit the ceiling
		scaleMin -= (scaleMax - ceilMax)
		scaleMax = ceilMax
	}
	// Re-clamp floor just in case shifting down pushed min below floor
	if scaleMin < floorMin {
		scaleMin = floorMin
	}

	// Round to clean integer boundaries to keep axis labels nice
	scaleMin = math.Floor(scaleMin)
	scaleMax = math.Ceil(scaleMax)

	if scaleMax <= scaleMin {
		scaleMax = scaleMin + minRange
	}

	return scaleMin, scaleMax
}

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

	useM := false
	if suffix == "K" && maxVal >= 1024.0 {
		useM = true
		maxVal = maxVal / 1024.0
		minVal = minVal / 1024.0
		suffix = "M"
	}

	result := make([]string, len(lines))
	for i, line := range lines {
		if i >= height {
			result[i] = "        " + line
			continue
		}

		val := maxVal
		if height > 1 {
			val = maxVal - (float64(i)/float64(height-1))*(maxVal-minVal)
		}
		var label string

		isTick := false
		isBottom := i == height-1
		tickChar := "┤ "
		if isBottom {
			tickChar = "└ "
		}

		if height >= 5 {
			h4 := height / 4
			h2 := height / 2
			h34 := (height * 3) / 4
			if i == 0 || i == h4 || i == h2 || i == h34 || isBottom {
				isTick = true
			}
		} else {
			if i == 0 || i == height/2 || isBottom {
				isTick = true
			}
		}

		if isTick {
			valStr := ""
			if useM {
				if val >= 10.0 {
					valStr = fmt.Sprintf("%3.0f", val)
				} else {
					valStr = fmt.Sprintf("%3.1f", val)
				}
			} else {
				valStr = fmt.Sprintf("%3.0f", val)
			}
			label = fmt.Sprintf("%4s%s %s", valStr, suffix, tickChar)
		} else {
			label = "      │ "
		}

		if len(label) < 8 {
			label = label + strings.Repeat(" ", 8-len(label))
		} else if len(label) > 8 {
			label = label[:8]
		}

		result[i] = label + line
	}
	return strings.Join(result, "\n")
}

// RenderXAxis generates a line of timestamps evenly spaced across the graph width.
func RenderXAxis(width int, timestamps []time.Time, timeframe string) string {
	if width < 15 {
		return ""
	}

	activeTimes := timestamps
	if len(activeTimes) > width {
		activeTimes = activeTimes[len(activeTimes)-width:]
	}
	if len(activeTimes) == 0 {
		return strings.Repeat(" ", width)
	}

	startCol := width - len(activeTimes)
	if startCol < 0 {
		startCol = 0
	}

	idx0 := 0
	idx1 := (len(activeTimes) - 1) / 3
	idx2 := 2 * (len(activeTimes) - 1) / 3
	idx3 := len(activeTimes) - 1

	col0 := startCol + idx0
	col1 := startCol + idx1
	col2 := startCol + idx2
	col3 := startCol + idx3

	format := "15:04:05"
	if timeframe == "15m" || timeframe == "1h" {
		format = "15:04"
	}

	l0 := activeTimes[idx0].Format(format)
	l1 := activeTimes[idx1].Format(format)
	l2 := activeTimes[idx2].Format(format)
	l3 := activeTimes[idx3].Format(format)

	buf := make([]byte, width)
	for i := range buf {
		buf[i] = ' '
	}

	writeAt := func(s string, start int) {
		if start < 0 {
			start = 0
		}
		if start+len(s) > width {
			start = width - len(s)
		}
		for i := 0; i < len(s); i++ {
			if start+i >= 0 && start+i < width {
				buf[start+i] = s[i]
			}
		}
	}

	writeAt(l0, col0)
	if len(activeTimes) >= 3 {
		writeAt(l1, col1-len(l1)/2)
		writeAt(l2, col2-len(l2)/2)
	}
	writeAt(l3, col3-len(l3)+1)

	return string(buf)
}
