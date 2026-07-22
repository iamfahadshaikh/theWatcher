package panel

import (
	"fmt"

	"theWatcher/internal/graph"
	"theWatcher/internal/history"
	"theWatcher/internal/metrics"
	"theWatcher/internal/theme"
)

type DiskPanel struct {
	history  *history.RingBuffer[float64]
	latest   metrics.DiskMetric
	renderer graph.Renderer
}

func NewDiskPanel() *DiskPanel {
	return &DiskPanel{
		history:  history.NewRingBuffer[float64](300),
		renderer: graph.NewBrailleRenderer(),
	}
}

func (p *DiskPanel) UpdateMetric(m metrics.DiskMetric) {
	p.latest = m
	p.history.Push(m.ReadRateKB + m.WriteRateKB)
}

func (p *DiskPanel) Title() string {
	return "Disk I/O"
}

func (p *DiskPanel) Render(width, height int, active bool) string {
	graphH := height - 4
	if graphH < 1 {
		graphH = 1
	}
	graphW := width - 4
	if graphW < 1 {
		graphW = 1
	}

	slice := p.history.Slice()
	maxV := 10.0
	for _, v := range slice {
		if v > maxV {
			maxV = v
		}
	}

	brailleStr := graph.RenderHistory(p.history, p.renderer, graph.GraphData{
		Width:  graphW,
		Height: graphH,
		Min:    0,
		Max:    maxV,
		Color:  string(theme.ColorSuccess),
	})

	statsStr := fmt.Sprintf("R: %.1fK/s W: %.1fK/s IOPS: %d Lat: %.1fms", p.latest.ReadRateKB, p.latest.WriteRateKB, p.latest.ReadIOPS+p.latest.WriteIOPS, p.latest.LatencyMs)

	content := fmt.Sprintf("%s\n%s", brailleStr, statsStr)
	return RenderContainer(p.Title(), content, width, height, active)
}
