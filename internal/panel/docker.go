package panel

import (
	"fmt"
	"strings"

	"theWatcher/internal/model"
)

type DockerPanel struct {
	stats       []model.DockerStat
	unavailable bool
}

func NewDockerPanel() *DockerPanel {
	return &DockerPanel{}
}

func (p *DockerPanel) UpdateDockerStats(stats []model.DockerStat, err error) {
	if err != nil || stats == nil {
		p.unavailable = true
		p.stats = nil
	} else {
		p.unavailable = false
		p.stats = stats
	}
}

func (p *DockerPanel) Title() string {
	return "Docker Containers"
}

func (p *DockerPanel) Render(width, height int, active bool) string {
	if p.unavailable {
		return RenderContainer(p.Title(), "Docker daemon not running or not installed.", width, height, active)
	}

	if len(p.stats) == 0 {
		return RenderContainer(p.Title(), "No active Docker containers.", width, height, active)
	}

	header := fmt.Sprintf("%-14s %-16s %-8s %-20s %s", "CONTAINER ID", "NAME", "CPU %", "MEM USAGE / LIMIT", "NET I/O")
	var lines []string
	lines = append(lines, header)

	maxLines := height - 4
	if maxLines < 1 {
		maxLines = 1
	}

	for i, stat := range p.stats {
		if i >= maxLines {
			break
		}
		row := fmt.Sprintf("%-14s %-16s %-8.1f %-20s %s", stat.ID, stat.Name, stat.CPUPct, stat.MemUsage, stat.NetIO)
		lines = append(lines, row)
	}

	content := strings.Join(lines, "\n")
	return RenderContainer(p.Title(), content, width, height, active)
}
