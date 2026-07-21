package graph

import (
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"theWatcher/internal/theme"
)

type BrailleRenderer struct{}

func NewBrailleRenderer() *BrailleRenderer {
	return &BrailleRenderer{}
}

func (b *BrailleRenderer) Render(data GraphData) string {
	if data.Width <= 0 || data.Height <= 0 {
		return ""
	}

	if len(data.Samples) == 0 {
		emptyRow := strings.Repeat(" ", data.Width)
		rows := make([]string, data.Height)
		for i := range rows {
			rows[i] = emptyRow
		}
		return lipgloss.NewStyle().Foreground(lipgloss.Color(data.Color)).Render(strings.Join(rows, "\n"))
	}

	dotsWidth := data.Width * 2
	dotsHeight := data.Height * 4

	// Sample or truncate to dotsWidth
	samples := data.Samples
	if len(samples) > dotsWidth {
		samples = samples[len(samples)-dotsWidth:]
	}

	minV := data.Min
	maxV := data.Max
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

	// Matrix of Braille codepoints offset from 0x2800
	grid := make([][]rune, data.Height)
	for r := 0; r < data.Height; r++ {
		grid[r] = make([]rune, data.Width)
		for c := 0; c < data.Width; c++ {
			grid[r][c] = 0x2800
		}
	}

	startDotX := dotsWidth - len(samples)
	if startDotX < 0 {
		startDotX = 0
	}

	dotBitMap := [4][2]rune{
		{0x01, 0x08},
		{0x02, 0x10},
		{0x04, 0x20},
		{0x40, 0x80},
	}

	for i, v := range samples {
		dotX := startDotX + i
		if dotX < 0 || dotX >= dotsWidth {
			continue
		}

		norm := (v - minV) / (maxV - minV)
		if norm < 0 {
			norm = 0
		}
		if norm > 1 {
			norm = 1
		}

		dotY := int(math.Round((1.0 - norm) * float64(dotsHeight-1)))
		if dotY < 0 {
			dotY = 0
		}
		if dotY >= dotsHeight {
			dotY = dotsHeight - 1
		}

		cellX := dotX / 2
		cellY := dotY / 4
		localX := dotX % 2
		localY := dotY % 4

		if cellY >= 0 && cellY < data.Height && cellX >= 0 && cellX < data.Width {
			grid[cellY][cellX] |= dotBitMap[localY][localX]
		}
	}

	rows := make([]string, data.Height)
	for r := 0; r < data.Height; r++ {
		rows[r] = string(grid[r])
	}

	color := data.Color
	if color == "" {
		color = string(theme.ColorLove)
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(strings.Join(rows, "\n"))
}
