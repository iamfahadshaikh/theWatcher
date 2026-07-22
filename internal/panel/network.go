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
}

func NewNetworkPanel(isTX bool) *NetworkPanel {
	return &NetworkPanel{
		isTX:     isTX,
		history:  history.NewRingBuffer[float64](300),
		renderer: graph.NewBrailleRenderer(),
	}
}

func (p *NetworkPanel) UpdateMetric(m metrics.NetworkMetric) {
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

func (p *NetworkPanel) Render(width, height int, active bool) string {
	graphH := height - 4
	if graphH < 1 {
		graphH = 1
	}
	graphW := width - 4
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

	var statsStr string
	if p.isTX {
		statsStr = fmt.Sprintf("Rate: %.1fKB/s Pkts: %d Err: %d Drop: %d", p.latest.TxRateKB, p.latest.TxPackets, p.latest.TxErrors, p.latest.TxDrop)
	} else {
		statsStr = fmt.Sprintf("Rate: %.1fKB/s Pkts: %d Err: %d Drop: %d", p.latest.RxRateKB, p.latest.RxPackets, p.latest.RxErrors, p.latest.RxDrop)
	}

	content := fmt.Sprintf("%s\n%s", brailleStr, statsStr)
	return RenderContainer(p.Title(), content, width, height, active)
}
