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

	color := data.Color
	if color == "" {
		color = string(theme.ColorLove)
	}

	multiSamples := data.MultiSamples
	multiColors := data.MultiColors

	if len(multiSamples) == 0 && len(data.Samples) > 0 {
		multiSamples = [][]float64{data.Samples}
		multiColors = []string{color}
	}

	if len(multiSamples) == 0 {
		emptyRow := strings.Repeat(" ", data.Width)
		rows := make([]string, data.Height)
		for i := range rows {
			rows[i] = emptyRow
		}
		return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(strings.Join(rows, "\n"))
	}

	minV := data.Min
	maxV := data.Max
	for sIdx, samples := range multiSamples {
		if len(samples) > data.Width {
			multiSamples[sIdx] = samples[len(samples)-data.Width:]
			samples = multiSamples[sIdx]
		}
		for _, v := range samples {
			if v > maxV {
				maxV = v
			}
			if v < minV {
				minV = v
			}
		}
	}
	if maxV <= minV {
		maxV = minV + 1.0
	}

	grid := make([][]rune, data.Height)
	colorGrid := make([][]string, data.Height)
	for r := 0; r < data.Height; r++ {
		grid[r] = make([]rune, data.Width)
		colorGrid[r] = make([]string, data.Width)
		for c := 0; c < data.Width; c++ {
			grid[r][c] = ' '
			colorGrid[r][c] = ""
		}
	}

	valRow := func(val float64) int {
		norm := (val - minV) / (maxV - minV)
		if norm < 0 {
			norm = 0
		}
		if norm > 1 {
			norm = 1
		}
		y := (1.0 - norm) * float64(data.Height-1)
		return int(math.Round(y))
	}

	// 1. Draw reference lines if enabled (background)
	if data.ShowRefLines {
		drawRefLine := func(val float64, char rune, col string) {
			if col == "" {
				col = string(theme.ColorSubtle)
			}
			rIdx := valRow(val)
			if rIdx >= 0 && rIdx < data.Height {
				for c := 0; c < data.Width; c++ {
					grid[rIdx][c] = char
					colorGrid[rIdx][c] = col
				}
			}
		}
		// Draw Min first, then Max, then Avg so Avg takes priority if rows overlap
		drawRefLine(data.MinVal, '╌', data.MinColor)
		drawRefLine(data.MaxVal, '╌', data.MaxColor)
		drawRefLine(data.AvgVal, '╌', data.AvgColor)
	}

	// 2. Draw trend lines (foreground - overwrites background)
	for sIdx, samples := range multiSamples {
		col := color
		if sIdx < len(multiColors) && multiColors[sIdx] != "" {
			col = multiColors[sIdx]
		}

		startCol := data.Width - len(samples)
		if startCol < 0 {
			startCol = 0
		}

		rows := make([]float64, len(samples))
		for i, v := range samples {
			norm := (v - minV) / (maxV - minV)
			if norm < 0 {
				norm = 0
			}
			if norm > 1 {
				norm = 1
			}
			rows[i] = (1.0 - norm) * float64(data.Height-1)
		}

		for i := 0; i < len(samples); i++ {
			c := startCol + i
			if c >= data.Width {
				break
			}

			yCurr := rows[i]
			rCurr := int(math.Round(yCurr))

			if i == len(samples)-1 {
				if rCurr >= 0 && rCurr < data.Height {
					grid[rCurr][c] = '─'
					colorGrid[rCurr][c] = col
				}
				continue
			}

			yNext := rows[i+1]
			rNext := int(math.Round(yNext))

			if rCurr == rNext {
				if rCurr >= 0 && rCurr < data.Height {
					grid[rCurr][c] = '─'
					colorGrid[rCurr][c] = col
				}
			} else if rCurr < rNext {
				// Going down
				if rCurr >= 0 && rCurr < data.Height {
					grid[rCurr][c] = '╮'
					colorGrid[rCurr][c] = col
				}
				for r := rCurr + 1; r < rNext; r++ {
					if r >= 0 && r < data.Height {
						grid[r][c] = '│'
						colorGrid[r][c] = col
					}
				}
				if rNext >= 0 && rNext < data.Height {
					grid[rNext][c] = '╰'
					colorGrid[rNext][c] = col
				}
			} else {
				// Going up
				if rCurr >= 0 && rCurr < data.Height {
					grid[rCurr][c] = '╯'
					colorGrid[rCurr][c] = col
				}
				for r := rNext + 1; r < rCurr; r++ {
					if r >= 0 && r < data.Height {
						grid[r][c] = '│'
						colorGrid[r][c] = col
					}
				}
				if rNext >= 0 && rNext < data.Height {
					grid[rNext][c] = '╭'
					colorGrid[rNext][c] = col
				}
			}
		}
	}

	// 3. Render grid grouping adjacent cells with the same color run
	var sb strings.Builder
	for r := 0; r < data.Height; r++ {
		if r > 0 {
			sb.WriteRune('\n')
		}
		runColor := ""
		var runBuilder strings.Builder
		for c := 0; c < data.Width; c++ {
			cellColor := colorGrid[r][c]
			if cellColor == "" {
				cellColor = color
			}
			if cellColor != runColor {
				if runBuilder.Len() > 0 {
					sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(runColor)).Render(runBuilder.String()))
					runBuilder.Reset()
				}
				runColor = cellColor
			}
			runBuilder.WriteRune(grid[r][c])
		}
		if runBuilder.Len() > 0 {
			sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(runColor)).Render(runBuilder.String()))
		}
	}

	return sb.String()
}
