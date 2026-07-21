package panel

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"theWatcher/internal/graph"
	"theWatcher/internal/history"
	"theWatcher/internal/model"
)

type NetworkPanel struct {
	isTX    bool
	history *history.RingBuffer[float64]
	latest  model.NetworkMetric
}

func NewNetworkPanel(isTX bool) *NetworkPanel {
	return &NetworkPanel{
		isTX:    isTX,
		history: history.NewRingBuffer[float64](300),
	}
}

func (p *NetworkPanel) UpdateMetric(m model.NetworkMetric) {
	p.latest = m
	if p.isTX {
		p.history.Push(m.TxRateKB)
	} else {
		p.history.Push(m.RxRateKB)
	}
}

func (p *NetworkPanel) Title() string {
	if p.isTX {
		return "Network TX"
	}
	return "Network RX"
}

func (p *NetworkPanel) Render(width, height int, active bool) string {
	slice := p.history.Slice()
	graphH := height - 4
	if graphH < 1 {
		graphH = 1
	}
	graphW := width - 4
	if graphW < 1 {
		graphW = 1
	}

	color := "#70a5fd" // Cyan/blue for RX
	if p.isTX {
		color = "#c4a7e7" // Purple/violet for TX
	}

	lineGraph := graph.RenderLineGraph(slice, graphW, graphH, 0, 100, lipgloss.NewStyle().Foreground(lipgloss.Color(color)))

	var statsStr string
	if p.isTX {
		statsStr = fmt.Sprintf("Rate: %.1fKB/s Pkts: %d Err: %d Drop: %d", p.latest.TxRateKB, p.latest.TxPackets, p.latest.TxErrors, p.latest.TxDrop)
	} else {
		statsStr = fmt.Sprintf("Rate: %.1fKB/s Pkts: %d Err: %d Drop: %d", p.latest.RxRateKB, p.latest.RxPackets, p.latest.RxErrors, p.latest.RxDrop)
	}

	content := fmt.Sprintf("%s\n%s", lineGraph, statsStr)
	return RenderContainer(p.Title(), content, width, height, active)
}
