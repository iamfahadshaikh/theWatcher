package panel

import (
	"fmt"
	"strings"
	"time"

	"theWatcher/internal/config"
	"theWatcher/internal/graph"
	"theWatcher/internal/history"
	"theWatcher/internal/metrics"
	"theWatcher/internal/theme"
)

type CPUPanel struct {
	history   *history.RingBuffer[history.Sample]
	latest    metrics.CPUMetric
	renderer  graph.Renderer
	trend     string
	prevPct   float64
	Timeframe string
}

func NewCPUPanel() *CPUPanel {
	return &CPUPanel{
		history:   history.NewRingBuffer[history.Sample](300),
		renderer:  graph.NewBrailleRenderer(),
		trend:     "~",
		Timeframe: "30s",
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
	p.history.Push(history.Sample{Time: time.Now(), Value: m.UsagePct})
}

func (p *CPUPanel) Title() string {
	return "CPU Usage (%)"
}

func (p *CPUPanel) Render(width, height int, active bool, expanded bool) string {
	graphH := height - 5
	if expanded {
		graphH = height - 10
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
		Color:  string(theme.ColorLove),
	})

	timestamps := make([]time.Time, len(slice))
	for i, s := range slice {
		timestamps[i] = s.Time
	}

	brailleStr = graph.AddYAxis(brailleStr, graphH, minVal, maxVal, "%")
	xAxisStr := "        " + graph.RenderXAxis(graphW, timestamps, p.Timeframe)

	var statsStr string
	if expanded {
		r1 := fmt.Sprintf("UTIL   | Cur: %.1f%% %s  Avg: %.1f%%  Min: %.1f%%  Max: %.1f%%",
			p.latest.UsagePct, p.trend, p.latest.AveragePct, p.latest.SessionMin, p.latest.SessionMax)
		r2 := fmt.Sprintf("SESS   | Avg: %.1f%%  Min: %.1f%%  Max: %.1f%%  Samples: %d",
			p.latest.SessionAvg, p.latest.SessionMin, p.latest.SessionMax, p.latest.SessionSamples)
		var raplStr string
		if p.latest.PackagePowerW > 0 {
			raplStr = fmt.Sprintf("Power: %.1fW  ", p.latest.PackagePowerW)
		}
		var tempStr string
		if p.latest.TempC > 0 {
			tempStr = fmt.Sprintf("Temp: %.1f°C  ", p.latest.TempC)
		}
		r3 := fmt.Sprintf("SENSOR | %s%sFreq: %.0fMHz (Range: %.0f-%.0f)",
			tempStr, raplStr, p.latest.FreqAvgMHz, p.latest.FreqMinMHz, p.latest.FreqMaxMHz)
		var spikeStr string
		if p.latest.LastSpike.PeakValue > 0 {
			age := time.Since(p.latest.LastSpike.EndTime).Round(time.Second)
			spikeStr = fmt.Sprintf("Peak: %.1f%%  Dur: %s  Age: %s ago",
				p.latest.LastSpike.PeakValue, p.latest.LastSpike.Duration.Round(time.Millisecond), age)
		} else {
			spikeStr = "No spikes recorded (>80%)"
		}
		r4 := fmt.Sprintf("SPIKE  | %s", spikeStr)
		var coreFreqs []string
		for idx, f := range p.latest.PerCoreFreqMHz {
			if idx >= 8 {
				coreFreqs = append(coreFreqs, "...")
				break
			}
			coreFreqs = append(coreFreqs, fmt.Sprintf("%.0f", f))
		}
		r5 := fmt.Sprintf("CORES  | %s MHz", strings.Join(coreFreqs, " "))
		statsStr = fmt.Sprintf("%s\n%s\n%s\n%s\n%s", r1, r2, r3, r4, r5)
	} else {
		statsStr = fmt.Sprintf("Cur: %s %.1f%%  Avg: %.1f%%  Peak: %.1f%%", p.trend, p.latest.UsagePct, p.latest.AveragePct, p.latest.PeakPct)
	}

	content := fmt.Sprintf("%s\n%s\n%s", brailleStr, xAxisStr, statsStr)
	return RenderContainer(p.Title(), content, width, height, active)
}
