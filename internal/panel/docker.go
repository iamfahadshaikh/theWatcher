package panel

import (
	"fmt"
	"strings"

	"theWatcher/internal/metrics"
)

type DockerPanel struct {
	stats        []metrics.DockerStat
	unavailable  bool
	scrollOffset int
}

func NewDockerPanel() *DockerPanel {
	return &DockerPanel{}
}

func (p *DockerPanel) ScrollUp() {
	p.scrollOffset--
	if p.scrollOffset < 0 {
		p.scrollOffset = 0
	}
}

func (p *DockerPanel) ScrollDown() {
	p.scrollOffset++
}

func (p *DockerPanel) UpdateDockerStats(stats []metrics.DockerStat, err error) {
	if err != nil || stats == nil {
		p.unavailable = true
		p.stats = nil
	} else {
		p.unavailable = false
		p.stats = stats
	}
}

func (p *DockerPanel) Title() string {
	return "Docker Stats"
}

func (p *DockerPanel) Render(width, height int, active bool, expanded bool) string {
	if p.unavailable {
		return RenderContainer(p.Title(), "Docker daemon not running or not installed.", width, height, active)
	}

	if len(p.stats) == 0 {
		return RenderContainer(p.Title(), "No active Docker containers.", width, height, active)
	}

	header := fmt.Sprintf("%-14s %-16s %-6s %-10s %-10s %-6s %-16s %-16s %-4s",
		"CONTAINER ID", "NAME", "CPU%", "MEM USAGE", "LIMIT", "MEM%", "NET RX / TX", "BLOCK R / W", "PIDS")
	var lines []string
	lines = append(lines, header)

	maxLines := height - 4
	if maxLines < 1 {
		maxLines = 1
	}

	if p.scrollOffset > len(p.stats)-maxLines {
		p.scrollOffset = len(p.stats) - maxLines
	}
	if p.scrollOffset < 0 {
		p.scrollOffset = 0
	}

	for i := p.scrollOffset; i < len(p.stats); i++ {
		if len(lines)-1 >= maxLines {
			break
		}
		stat := p.stats[i]
		netStr := fmt.Sprintf("%s/%s", stat.NetRx, stat.NetTx)
		blockStr := fmt.Sprintf("%s/%s", stat.BlockRead, stat.BlockWrite)
		row := fmt.Sprintf("%-14s %-16s %-6.1f %-10s %-10s %-6s %-16s %-16s %-4d",
			stat.ID, stat.Name, stat.CPUPct, stat.MemUsage, stat.MemLimit, stat.MemPct, netStr, blockStr, stat.PIDs)
		lines = append(lines, row)
	}

	content := strings.Join(lines, "\n")
	return RenderContainer(p.Title(), content, width, height, active)
}
