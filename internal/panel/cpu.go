package panel

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"theWatcher/internal/graph"
	"theWatcher/internal/history"
	"theWatcher/internal/model"
)

type CPUPanel struct {
	history *history.RingBuffer[float64]
	latest  model.CPUMetric
}

func NewCPUPanel() *CPUPanel {
	return &CPUPanel{
		history: history.NewRingBuffer[float64](300),
	}
}

func (p *CPUPanel) UpdateMetric(m model.CPUMetric) {
	p.latest = m
	p.history.Push(m.UsagePct)
}

func (p *CPUPanel) Title() string {
	return "CPU"
}

func (p *CPUPanel) Render(width, height int, active bool) string {
	slice := p.history.Slice()
	graphH := height - 4
	if graphH < 1 {
		graphH = 1
	}
	graphW := width - 4
	if graphW < 1 {
		graphW = 1
	}

	lineGraph := graph.RenderLineGraph(slice, graphW, graphH, 0, 100, lipgloss.NewStyle().Foreground(lipgloss.Color("#eb6f92")))
	statsStr := fmt.Sprintf("Cur: %.1f%%  Avg: %.1f%%  Peak: %.1f%%", p.latest.UsagePct, p.latest.AveragePct, p.latest.PeakPct)

	content := fmt.Sprintf("%s\n%s", lineGraph, statsStr)
	return RenderContainer(p.Title(), content, width, height, active)
}
