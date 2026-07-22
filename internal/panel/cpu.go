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

func (p *CPUPanel) Render(width, height int, active bool, expanded bool) string {
	graphH := height - 4
	if expanded {
		graphH = height - 9
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
		Color:  string(theme.ColorLove),
	})

	if expanded {
		brailleStr = graph.AddYAxis(brailleStr, graphH, 0, 100, "%")
	}

	var statsStr string
	if expanded {
		minV := 100.0
		slice := p.history.Slice()
		for _, v := range slice {
			if v < minV {
				minV = v
			}
		}
		if len(slice) == 0 {
			minV = 0.0
		}
		statsStr = fmt.Sprintf("Current : %.1f%%\nAverage : %.1f%%\nPeak    : %.1f%%\nMin     : %.1f%%\nSamples : %d",
			p.latest.UsagePct, p.latest.AveragePct, p.latest.PeakPct, minV, len(slice))
	} else {
		statsStr = fmt.Sprintf("Cur: %s %.1f%%  Avg: %.1f%%  Peak: %.1f%%", p.trend, p.latest.UsagePct, p.latest.AveragePct, p.latest.PeakPct)
	}

	content := fmt.Sprintf("%s\n%s", brailleStr, statsStr)
	return RenderContainer(p.Title(), content, width, height, active)
}
