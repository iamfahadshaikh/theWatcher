package panel

import (
	"fmt"
	"strings"
	"time"

	"theWatcher/internal/graph"
	"theWatcher/internal/history"
	"theWatcher/internal/metrics"
	"theWatcher/internal/theme"
)

type NetworkPanel struct {
	isTX      bool
	history   *history.RingBuffer[history.Sample]
	latest    metrics.NetworkMetric
	renderer  graph.Renderer
	trend     string
	prevRate  float64
	Timeframe string
}

func NewNetworkPanel(isTX bool) *NetworkPanel {
	return &NetworkPanel{
		isTX:      isTX,
		history:   history.NewRingBuffer[history.Sample](300),
		renderer:  graph.NewBrailleRenderer(),
		trend:     "~",
		Timeframe: "30s",
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
		p.history.Push(history.Sample{Time: time.Now(), Value: m.TxRateKB})
	} else {
		p.history.Push(history.Sample{Time: time.Now(), Value: m.RxRateKB})
	}
}

func (p *NetworkPanel) Title() string {
	if p.isTX {
		return "Outgoing Traffic"
	}
	return "Incoming Traffic"
}

func (p *NetworkPanel) Render(width, height int, active bool, expanded bool) string {
	graphH := height - 5
	if expanded {
		graphH = height - 11
	}
	if graphH < 1 {
		graphH = 1
	}
	graphW := width - 12
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
		if v.Value > maxV {
			maxV = v.Value
		}
	}

	var visibleSamples []float64
	if len(slice) > graphW {
		visibleSamples = make([]float64, graphW)
		for i := 0; i < graphW; i++ {
			visibleSamples[i] = slice[len(slice)-graphW+i].Value
		}
	} else {
		visibleSamples = make([]float64, len(slice))
		for i, s := range slice {
			visibleSamples[i] = s.Value
		}
	}

	var avg, vMin, vMax float64
	if len(visibleSamples) > 0 {
		vMin = visibleSamples[0]
		vMax = visibleSamples[0]
		sum := 0.0
		for _, v := range visibleSamples {
			sum += v
			if v < vMin {
				vMin = v
			}
			if v > vMax {
				vMax = v
			}
		}
		avg = sum / float64(len(visibleSamples))
	}

	brailleStr := graph.RenderHistory(p.history, p.renderer, graph.GraphData{
		Width:        graphW,
		Height:       graphH,
		Min:          0,
		Max:          maxV,
		Color:        color,
		ShowRefLines: true,
		AvgVal:       avg,
		MinVal:       vMin,
		MaxVal:       vMax,
		AvgColor:     string(theme.ColorGold),
		MinColor:     string(theme.ColorIris),
		MaxColor:     string(theme.ColorRose),
	})

	timestamps := make([]time.Time, len(slice))
	for i, s := range slice {
		timestamps[i] = s.Time
	}

	brailleStr = graph.AddYAxis(brailleStr, graphH, 0, maxV, "K")
	xAxisStr := "        " + graph.RenderXAxis(graphW, timestamps, p.Timeframe)

	var statsStr string
	if expanded {
		var r1 string
		if p.isTX {
			r1 = fmt.Sprintf("NET   | Out Rate: %.1fK/s %s  Peak: %.1fK/s  Pkts: %d  Err: %d  Drop: %d",
				p.latest.TxRateKB, p.trend, maxV, p.latest.TxPackets, p.latest.TxErrors, p.latest.TxDrop)
		} else {
			r1 = fmt.Sprintf("NET   | In Rate: %.1fK/s %s  Peak: %.1fK/s  Pkts: %d  Err: %d  Drop: %d",
				p.latest.RxRateKB, p.trend, maxV, p.latest.RxPackets, p.latest.RxErrors, p.latest.RxDrop)
		}
		r2 := fmt.Sprintf("%-10s %-8s %-8s %-6s %-8s %-10s %-10s", "INTERFACE", "CARRIER", "SPEED", "MTU", "DUPLEX", "RX RATE", "TX RATE")

		var rows []string
		rows = append(rows, r1, r2)
		for idx, iface := range p.latest.Interfaces {
			if idx >= 2 {
				break
			}
			state := "down"
			if iface.Carrier == 1 {
				state = "up"
			}
			speedStr := fmt.Sprintf("%dM", iface.Speed)
			if iface.Speed <= 0 {
				speedStr = "N/A"
			}
			row := fmt.Sprintf("%-10s %-8s %-8s %-6d %-8s %-10.1f %-10.1f",
				truncate(iface.Name, 10), state, speedStr, iface.MTU, iface.Duplex, iface.RxRateKB, iface.TxRateKB)
			rows = append(rows, row)
		}
		for len(rows) < 4 {
			rows = append(rows, "")
		}
		pr := p.latest.Proto
		r5 := fmt.Sprintf("PROTO | TCP Packets: %d  UDP Packets: %d  ICMP Packets: %d", pr.TCPPackets, pr.UDPPackets, pr.ICMPPackets)
		r6 := fmt.Sprintf("PROTO | IPv4 Packets: %d  IPv6 Packets: %d  Retransmissions: %d", pr.IPv4Packets, pr.IPv6Packets, pr.Retransmits)
		rows = append(rows, r5, r6)
		statsStr = strings.Join(rows, "\n")
	} else {
		if p.isTX {
			statsStr = fmt.Sprintf("Rate: %s %.1fKB/s Pkts: %d Err: %d Drop: %d", p.trend, p.latest.TxRateKB, p.latest.TxPackets, p.latest.TxErrors, p.latest.TxDrop)
		} else {
			statsStr = fmt.Sprintf("Rate: %s %.1fKB/s Pkts: %d Err: %d Drop: %d", p.trend, p.latest.RxRateKB, p.latest.RxPackets, p.latest.RxErrors, p.latest.RxDrop)
		}
	}

	content := fmt.Sprintf("%s\n%s\n%s", brailleStr, xAxisStr, statsStr)
	return RenderContainer(p.Title(), content, width, height, active)
}
