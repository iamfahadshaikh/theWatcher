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

type DiskPanel struct {
	history   *history.RingBuffer[history.Sample]
	latest    metrics.DiskMetric
	renderer  graph.Renderer
	trend     string
	prevRate  float64
	Timeframe string
}

func NewDiskPanel() *DiskPanel {
	return &DiskPanel{
		history:   history.NewRingBuffer[history.Sample](300),
		renderer:  graph.NewBrailleRenderer(),
		trend:     "~",
		Timeframe: "30s",
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
	p.history.Push(history.Sample{Time: time.Now(), Value: rate})
}

func (p *DiskPanel) Title() string {
	return "Disk I/O"
}

func (p *DiskPanel) Render(width, height int, active bool, expanded bool) string {
	graphH := height - 5
	if expanded {
		graphH = height - 12
	}
	if graphH < 1 {
		graphH = 1
	}
	graphW := width - 12
	if graphW < 1 {
		graphW = 1
	}

	slice := p.history.Slice()
	maxV := 10.0
	for _, v := range slice {
		if v.Value > maxV {
			maxV = v.Value
		}
	}

	brailleStr := graph.RenderHistory(p.history, p.renderer, graph.GraphData{
		Width:  graphW,
		Height: graphH,
		Min:    0,
		Max:    maxV,
		Color:  string(theme.ColorSuccess),
	})

	timestamps := make([]time.Time, len(slice))
	for i, s := range slice {
		timestamps[i] = s.Time
	}

	brailleStr = graph.AddYAxis(brailleStr, graphH, 0, maxV, "K")
	xAxisStr := "        " + graph.RenderXAxis(graphW, timestamps, p.Timeframe)

	var statsStr string
	if expanded {
		r1 := fmt.Sprintf("I/O   | Read: %.1fK/s (%d IOPS)  Write: %.1fK/s (%d IOPS)  Lat: %.1fms  QD: %d  Util: %.1f%%  Temp: %.1f°C",
			p.latest.ReadRateKB, p.latest.ReadIOPS, p.latest.WriteRateKB, p.latest.WriteIOPS, p.latest.LatencyMs, p.latest.QueueDepth, p.latest.DiskUtilPct, p.latest.DiskTempC)
		r2 := fmt.Sprintf("%-20s %-8s %-10s %-10s %-10s %-6s %-6s %-16s", "MOUNTPOINT", "TYPE", "SIZE", "USED", "FREE", "USE%", "INODE%", "OPTIONS")
		
		var fsRows []string
		fsRows = append(fsRows, r1, r2)
		for _, fs := range p.latest.Filesystems {
			if len(fsRows)-2 >= 6 {
				break
			}
			totStr := formatBytes(fs.TotalBytes)
			usedStr := formatBytes(fs.UsedBytes)
			freeStr := formatBytes(fs.AvailBytes)
			roStr := "rw"
			if fs.ReadOnly {
				roStr = "ro"
			}
			optStr := fs.MountOptions
			if len(optStr) > 16 {
				optStr = optStr[:13] + "..."
			}
			row := fmt.Sprintf("%-20s %-8s %-10s %-10s %-10s %-6.1f %-6.1f %-16s",
				truncate(fs.MountPoint, 20), truncate(fs.FSType, 8), totStr, usedStr, freeStr, fs.UsedPct, fs.InodesPct, fmt.Sprintf("%s,%s", roStr, optStr))
			fsRows = append(fsRows, row)
		}
		for len(fsRows) < 8 {
			fsRows = append(fsRows, "")
		}
		statsStr = strings.Join(fsRows, "\n")
	} else {
		statsStr = fmt.Sprintf("R: %s %.1fK/s W: %.1fK/s IOPS: %d Lat: %.1fms", p.trend, p.latest.ReadRateKB, p.latest.WriteRateKB, p.latest.ReadIOPS+p.latest.WriteIOPS, p.latest.LatencyMs)
	}

	content := fmt.Sprintf("%s\n%s\n%s", brailleStr, xAxisStr, statsStr)
	return RenderContainer(p.Title(), content, width, height, active)
}

func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cI", float64(b)/float64(div), "KMGTPE"[exp])
}
