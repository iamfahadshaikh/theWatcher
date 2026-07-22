package panel

import (
	"fmt"

	"theWatcher/internal/graph"
	"theWatcher/internal/history"
	"theWatcher/internal/metrics"
	"theWatcher/internal/theme"
)

type CPUPanel struct {
	history  *history.RingBuffer[float64]
	latest   metrics.CPUMetric
	renderer graph.Renderer
	trend    string
	prevPct  float64
}

func NewCPUPanel() *CPUPanel {
	return &CPUPanel{
		history:  history.NewRingBuffer[float64](300),
		renderer: graph.NewBrailleRenderer(),
		trend:    "~",
	}
}

func (p *CPUPanel) UpdateMetric(m metrics.CPUMetric) {
	if m.UsagePct > p.prevPct {
		p.trend = "▲"
	} else if m.UsagePct < p.prevPct {
		p.trend = "▼"
	} else {
		p.trend = "~"
	}
	p.prevPct = m.UsagePct
	p.latest = m
	p.history.Push(m.UsagePct)
}

func (p *CPUPanel) Title() string {
	return "CPU Usage (%)"
}

func (p *CPUPanel) Render(width, height int, active bool) string {
	graphH := height - 4
	if graphH < 1 {
		graphH = 1
	}
	graphW := width - 4
	if graphW < 1 {
		graphW = 1
	}

	brailleStr := graph.RenderHistory(p.history, p.renderer, graph.GraphData{
		Width:  graphW,
		Height: graphH,
		Min:    0,
		Max:    100,
		Color:  string(theme.ColorLove),
	})

	statsStr := fmt.Sprintf("Cur: %s %.1f%%  Avg: %.1f%%  Peak: %.1f%%", p.trend, p.latest.UsagePct, p.latest.AveragePct, p.latest.PeakPct)

	content := fmt.Sprintf("%s\n%s", brailleStr, statsStr)
	return RenderContainer(p.Title(), content, width, height, active)
}
