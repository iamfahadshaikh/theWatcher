package panel

import (
	"fmt"

	"theWatcher/internal/graph"
	"theWatcher/internal/history"
	"theWatcher/internal/metrics"
	"theWatcher/internal/theme"
)

type NetworkPanel struct {
	isTX     bool
	history  *history.RingBuffer[float64]
	latest   metrics.NetworkMetric
	renderer graph.Renderer
	trend    string
	prevRate float64
}

func NewNetworkPanel(isTX bool) *NetworkPanel {
	return &NetworkPanel{
		isTX:     isTX,
		history:  history.NewRingBuffer[float64](300),
		renderer: graph.NewBrailleRenderer(),
		trend:    "~",
	}
}

func (p *NetworkPanel) UpdateMetric(m metrics.NetworkMetric) {
	rate := m.RxRateKB
	if p.isTX {
		rate = m.TxRateKB
	}

	if rate > p.prevRate {
		p.trend = "▲"
	} else if rate < p.prevRate {
		p.trend = "▼"
	} else {
		p.trend = "~"
	}
	p.prevRate = rate

	p.latest = m
	if p.isTX {
		p.history.Push(m.TxRateKB)
	} else {
		p.history.Push(m.RxRateKB)
	}
}

func (p *NetworkPanel) Title() string {
	if p.isTX {
		return "Outgoing Traffic"
	}
	return "Incoming Traffic"
}

func (p *NetworkPanel) Render(width, height int, active bool, expanded bool) string {
	graphH := height - 4
	if expanded {
		graphH = height - 10
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

	color := string(theme.ColorFoam)
	if p.isTX {
		color = string(theme.ColorIris)
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
		Color:  color,
	})

	if expanded {
		brailleStr = graph.AddYAxis(brailleStr, graphH, 0, maxV, "K")
	}

	var statsStr string
	if expanded {
		if p.isTX {
			statsStr = fmt.Sprintf("Current Rate : %.1f KB/s %s\nPeak Rate    : %.1f KB/s\nPackets      : %d\nErrors       : %d\nDropped      : %d\nInterface    : eth0",
				p.latest.TxRateKB, p.trend, maxV, p.latest.TxPackets, p.latest.TxErrors, p.latest.TxDrop)
		} else {
			statsStr = fmt.Sprintf("Current Rate : %.1f KB/s %s\nPeak Rate    : %.1f KB/s\nPackets      : %d\nErrors       : %d\nDropped      : %d\nInterface    : eth0",
				p.latest.RxRateKB, p.trend, maxV, p.latest.RxPackets, p.latest.RxErrors, p.latest.RxDrop)
		}
	} else {
		if p.isTX {
			statsStr = fmt.Sprintf("Rate: %s %.1fKB/s Pkts: %d Err: %d Drop: %d", p.trend, p.latest.TxRateKB, p.latest.TxPackets, p.latest.TxErrors, p.latest.TxDrop)
		} else {
			statsStr = fmt.Sprintf("Rate: %s %.1fKB/s Pkts: %d Err: %d Drop: %d", p.trend, p.latest.RxRateKB, p.latest.RxPackets, p.latest.RxErrors, p.latest.RxDrop)
		}
	}

	content := fmt.Sprintf("%s\n%s", brailleStr, statsStr)
	return RenderContainer(p.Title(), content, width, height, active)
}
