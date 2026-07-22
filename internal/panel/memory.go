package panel

import (
	"fmt"

	"theWatcher/internal/graph"
	"theWatcher/internal/history"
	"theWatcher/internal/metrics"
	"theWatcher/internal/theme"
)

type MemoryPanel struct {
	history  *history.RingBuffer[float64]
	latest   metrics.MemoryMetric
	renderer graph.Renderer
	trend    string
	prevPct  float64
}

func NewMemoryPanel() *MemoryPanel {
	return &MemoryPanel{
		history:  history.NewRingBuffer[float64](300),
		renderer: graph.NewBrailleRenderer(),
		trend:    "~",
	}
}

func (p *MemoryPanel) UpdateMetric(m metrics.MemoryMetric) {
	if m.UsedPct > p.prevPct {
		p.trend = "▲"
	} else if m.UsedPct < p.prevPct {
		p.trend = "▼"
	} else {
		p.trend = "~"
	}
	p.prevPct = m.UsedPct
	p.latest = m
	p.history.Push(m.UsedPct)
}

func (p *MemoryPanel) Title() string {
	return "Memory Usage"
}

func (p *MemoryPanel) Render(width, height int, active bool, expanded bool) string {
	graphH := height - 4
	if expanded {
		graphH = height - 10
	}
	if graphH < 1 {
		graphH = 1
	}
	graphW := width - 4
	if expanded {
		graphW = width - 12
	}
	if graphW < 1 {
		graphW = 1
	}

	brailleStr := graph.RenderHistory(p.history, p.renderer, graph.GraphData{
		Width:  graphW,
		Height: graphH,
		Min:    0,
		Max:    100,
		Color:  string(theme.ColorFoam),
	})

	if expanded {
		brailleStr = graph.AddYAxis(brailleStr, graphH, 0, 100, "%")
	}

	var statsStr string
	if expanded {
		statsStr = fmt.Sprintf("Total Memory : %.0f MB\nUsed Memory  : %.0f MB (%.1f%%) %s\nFree Memory  : %.0f MB\nBuffers      : %.0f MB\nCached       : %.0f MB",
			p.latest.TotalMB, p.latest.UsedMB, p.latest.UsedPct, p.trend, p.latest.FreeMB, p.latest.BufferMB, p.latest.CachedMB)
	} else {
		statsStr = fmt.Sprintf("Used: %s %.0fM Free: %.0fM Buf: %.0fM Cache: %.0fM", p.trend, p.latest.UsedMB, p.latest.FreeMB, p.latest.BufferMB, p.latest.CachedMB)
	}

	content := fmt.Sprintf("%s\n%s", brailleStr, statsStr)
	return RenderContainer(p.Title(), content, width, height, active)
}
