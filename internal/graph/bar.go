package graph

import (
	"fmt"
	"math"

	"github.com/charmbracelet/lipgloss"
)

// RenderBar renders a horizontal progress bar with label and value percentage.
func RenderBar(label string, val, maxVal float64, totalWidth int, barStyle, labelStyle lipgloss.Style) string {
	if totalWidth <= 0 {
		return ""
	}

	percent := 0.0
	if maxVal > 0 {
		percent = (val / maxVal) * 100.0
	}
	if percent > 100.0 {
		percent = 100.0
	}
	if percent < 0.0 {
		percent = 0.0
	}

	valStr := fmt.Sprintf(" %.1f%%", percent)
	lbl := labelStyle.Render(label)

	// Available width for bar brackets and fill
	barCap := totalWidth - lipgloss.Width(lbl) - len(valStr) - 2
	if barCap < 3 {
		barCap = 3
	}

	fillLen := int(math.Round((percent / 100.0) * float64(barCap)))
	if fillLen < 0 {
		fillLen = 0
	}
	if fillLen > barCap {
		fillLen = barCap
	}

	unfillLen := barCap - fillLen
	filled := fmt.Sprintf("%*s", fillLen, "")
	unfilled := fmt.Sprintf("%*s", unfillLen, "")

	filledStr := barStyle.Render(fmt.Sprintf("%s", stringsRepeat('█', fillLen)))
	unfilledStr := lipgloss.NewStyle().Foreground(lipgloss.Color("#444444")).Render(fmt.Sprintf("%s", stringsRepeat('░', unfillLen)))

	_ = filled
	_ = unfilled

	barStr := fmt.Sprintf("[%s%s]", filledStr, unfilledStr)
	return fmt.Sprintf("%s %s%s", lbl, barStr, valStr)
}

func stringsRepeat(r rune, count int) string {
	if count <= 0 {
		return ""
	}
	runes := make([]rune, count)
	for i := range runes {
		runes[i] = r
	}
	return string(runes)
}
