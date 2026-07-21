package panel

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"theWatcher/internal/graph"
	"theWatcher/internal/history"
	"theWatcher/internal/model"
)

type DiskPanel struct {
	history *history.RingBuffer[float64]
	latest  model.DiskMetric
}

func NewDiskPanel() *DiskPanel {
	return &DiskPanel{
		history: history.NewRingBuffer[float64](300),
	}
}

func (p *DiskPanel) UpdateMetric(m model.DiskMetric) {
	p.latest = m
	p.history.Push(m.ReadRateKB + m.WriteRateKB)
}

func (p *DiskPanel) Title() string {
	return "Disk"
}

func (p *DiskPanel) Render(width, height int, active bool) string {
	slice := p.history.Slice()
	graphH := height - 4
	if graphH < 1 {
		graphH = 1
	}
	graphW := width - 4
	if graphW < 1 {
		graphW = 1
	}

	lineGraph := graph.RenderLineGraph(slice, graphW, graphH, 0, 100, lipgloss.NewStyle().Foreground(lipgloss.Color("#a6e3a1")))
	statsStr := fmt.Sprintf("R: %.1fK/s W: %.1fK/s IOPS: %d Lat: %.1fms", p.latest.ReadRateKB, p.latest.WriteRateKB, p.latest.ReadIOPS+p.latest.WriteIOPS, p.latest.LatencyMs)

	content := fmt.Sprintf("%s\n%s", lineGraph, statsStr)
	return RenderContainer(p.Title(), content, width, height, active)
}
