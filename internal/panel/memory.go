package panel

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"theWatcher/internal/graph"
	"theWatcher/internal/history"
	"theWatcher/internal/model"
)

type MemoryPanel struct {
	history *history.RingBuffer[float64]
	latest  model.MemoryMetric
}

func NewMemoryPanel() *MemoryPanel {
	return &MemoryPanel{
		history: history.NewRingBuffer[float64](300),
	}
}

func (p *MemoryPanel) UpdateMetric(m model.MemoryMetric) {
	p.latest = m
	p.history.Push(m.UsedPct)
}

func (p *MemoryPanel) Title() string {
	return "Memory"
}

func (p *MemoryPanel) Render(width, height int, active bool) string {
	slice := p.history.Slice()
	graphH := height - 4
	if graphH < 1 {
		graphH = 1
	}
	graphW := width - 4
	if graphW < 1 {
		graphW = 1
	}

	lineGraph := graph.RenderLineGraph(slice, graphW, graphH, 0, 100, lipgloss.NewStyle().Foreground(lipgloss.Color("#9ccfd8")))
	statsStr := fmt.Sprintf("Used: %.0fM Free: %.0fM Buf: %.0fM Cache: %.0fM", p.latest.UsedMB, p.latest.FreeMB, p.latest.BufferMB, p.latest.CachedMB)

	content := fmt.Sprintf("%s\n%s", lineGraph, statsStr)
	return RenderContainer(p.Title(), content, width, height, active)
}
