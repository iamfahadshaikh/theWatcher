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
	trend    string
	prevRate float64
}

func NewDiskPanel() *DiskPanel {
	return &DiskPanel{
		history:  history.NewRingBuffer[float64](300),
		renderer: graph.NewBrailleRenderer(),
		trend:    "~",
	}
}

func (p *DiskPanel) UpdateMetric(m metrics.DiskMetric) {
	rate := m.ReadRateKB + m.WriteRateKB
	if rate > p.prevRate {
		p.trend = "▲"
	} else if rate < p.prevRate {
		p.trend = "▼"
	} else {
		p.trend = "~"
	}
	p.prevRate = rate

	p.latest = m
	p.history.Push(rate)
}

func (p *DiskPanel) Title() string {
	return "Disk I/O"
}

func (p *DiskPanel) Render(width, height int, active bool, expanded bool) string {
	graphH := height - 4
	if expanded {
		graphH = height - 11
	}
	if graphH < 1 {
		graphH = 1
	}
	graphW := width - 4
	if expanded {
		graphW = width - 14
	}
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

	if expanded {
		brailleStr = graph.AddYAxis(brailleStr, graphH, 0, maxV, "K")
	}

	var statsStr string
	if expanded {
		statsStr = fmt.Sprintf("Read Rate    : %.1f KB/s\nWrite Rate   : %.1f KB/s\nRead IOPS    : %d\nWrite IOPS   : %d\nI/O Latency  : %.2f ms\nTotal Read   : %.1f MB\nTotal Written: %.1f MB",
			p.latest.ReadRateKB, p.latest.WriteRateKB, p.latest.ReadIOPS, p.latest.WriteIOPS, p.latest.LatencyMs, float64(p.latest.ReadBytes)/(1024*1024), float64(p.latest.WriteBytes)/(1024*1024))
	} else {
		statsStr = fmt.Sprintf("R: %s %.1fK/s W: %.1fK/s IOPS: %d Lat: %.1fms", p.trend, p.latest.ReadRateKB, p.latest.WriteRateKB, p.latest.ReadIOPS+p.latest.WriteIOPS, p.latest.LatencyMs)
	}

	content := fmt.Sprintf("%s\n%s", brailleStr, statsStr)
	return RenderContainer(p.Title(), content, width, height, active)
}
