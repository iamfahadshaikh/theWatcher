package panel

import (
	"fmt"
	"time"

	"theWatcher/internal/config"
	"theWatcher/internal/graph"
	"theWatcher/internal/history"
	"theWatcher/internal/metrics"
	"theWatcher/internal/theme"
)

type MemoryPanel struct {
	history   *history.RingBuffer[history.Sample]
	latest    metrics.MemoryMetric
	renderer  graph.Renderer
	trend     string
	prevPct   float64
	Timeframe string
}

func NewMemoryPanel() *MemoryPanel {
	return &MemoryPanel{
		history:   history.NewRingBuffer[history.Sample](300),
		renderer:  graph.NewBrailleRenderer(),
		trend:     "~",
		Timeframe: "30s",
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
	p.history.Push(history.Sample{Time: time.Now(), Value: m.UsedPct})
}

func (p *MemoryPanel) Title() string {
	return "Memory Usage"
}

func (p *MemoryPanel) Render(width, height int, active bool, expanded bool) string {
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

	minVal := 0.0
	maxVal := 100.0
	slice := p.history.Slice()

	if config.GlobalConfig.Graph.Scale == "auto" {
		// Only consider visible samples to compute the Y-axis range
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
		minVal, maxVal = graph.ComputeAdaptiveBounds(visibleSamples, 0.0, 100.0, 10.0)
	}

	brailleStr := graph.RenderHistory(p.history, p.renderer, graph.GraphData{
		Width:  graphW,
		Height: graphH,
		Min:    minVal,
		Max:    maxVal,
		Color:  string(theme.ColorFoam),
	})

	timestamps := make([]time.Time, len(slice))
	for i, s := range slice {
		timestamps[i] = s.Time
	}

	brailleStr = graph.AddYAxis(brailleStr, graphH, minVal, maxVal, "%")
	xAxisStr := "        " + graph.RenderXAxis(graphW, timestamps, p.Timeframe)

	var statsStr string
	if expanded {
		r1 := fmt.Sprintf("RAM  | Tot: %.0fM  Used: %.0fM (%.1f%%) %s  Free: %.0fM  Avail: %.0fM",
			p.latest.TotalMB, p.latest.UsedMB, p.latest.UsedPct, p.trend, p.latest.FreeMB, p.latest.AvailableMB)
		r2 := fmt.Sprintf("SESS | Avg: %.1f%%  Min: %.1f%%  Max: %.1f%%  Samples: %d",
			p.latest.SessionAvg, p.latest.SessionMin, p.latest.SessionMax, p.latest.SessionSamples)
		r3 := fmt.Sprintf("PAGES| Buf: %.0fM  Cache: %.0fM  SwapCached: %.0fM  Dirty: %.1fM  WB: %.1fM",
			p.latest.BufferMB, p.latest.CachedMB, p.latest.SwapCachedMB, p.latest.DirtyMB, p.latest.WritebackMB)
		var swapPct float64
		if p.latest.SwapTotalMB > 0 {
			swapPct = (p.latest.SwapUsedMB / p.latest.SwapTotalMB) * 100.0
		}
		r4 := fmt.Sprintf("SWAP | Tot: %.0fM  Used: %.0fM (%.1f%%)  Free: %.0fM",
			p.latest.SwapTotalMB, p.latest.SwapUsedMB, swapPct, p.latest.SwapFreeMB)
		r5 := fmt.Sprintf("HPG  | Huge Total: %d  Huge Free: %d  THP Status: %s",
			p.latest.HugePagesTotal, p.latest.HugePagesFree, p.latest.THPStatus)
		r6 := fmt.Sprintf("DIAG | Pressure: %s  OOM Kills: %d",
			p.latest.MemPressure, p.latest.OOMKills)
		statsStr = fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s", r1, r2, r3, r4, r5, r6)
	} else {
		statsStr = fmt.Sprintf("Used: %s %.0fM Free: %.0fM Buf: %.0fM Cache: %.0fM", p.trend, p.latest.UsedMB, p.latest.FreeMB, p.latest.BufferMB, p.latest.CachedMB)
	}

	content := fmt.Sprintf("%s\n%s\n%s", brailleStr, xAxisStr, statsStr)
	return RenderContainer(p.Title(), content, width, height, active)
}
