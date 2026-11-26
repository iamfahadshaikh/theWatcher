package panel

import (
	"fmt"
	"strings"

	"theWatcher/internal/metrics"
)

type DockerPanel struct {
	stats        []metrics.DockerStat
	filesystems  []metrics.FilesystemStat
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
	if err != nil || stats == nil || len(stats) == 0 {
		p.unavailable = true
		p.stats = nil
	} else {
		p.unavailable = false
		p.stats = stats
	}
}

func (p *DockerPanel) UpdateFilesystems(fs []metrics.FilesystemStat) {
	p.filesystems = fs
}

func (p *DockerPanel) Title() string {
	if p.unavailable {
		return "Filesystem Mounts (/proc/mounts)"
	}
	return "Docker Stats"
}

func (p *DockerPanel) Render(width, height int, active bool, expanded bool) string {
	if p.unavailable {
		if len(p.filesystems) == 0 {
			return RenderContainer(p.Title(), "Docker daemon not running | No filesystems available.", width, height, active)
		}
		var lines []string
		header := fmt.Sprintf("%-16s %-10s %-8s %-12s %-12s %-6s %s", "MOUNTPOINT", "DEVICE", "FSTYPE", "USED", "TOTAL", "USE%", "OPTIONS")
		lines = append(lines, header)
		maxLines := height - 4
		if maxLines < 1 {
			maxLines = 1
		}
		for i, fs := range p.filesystems {
			if i >= maxLines {
				break
			}
			usedGB := float64(fs.UsedBytes) / (1024 * 1024 * 1024)
			totalGB := float64(fs.TotalBytes) / (1024 * 1024 * 1024)
			row := fmt.Sprintf("%-16s %-10s %-8s %-10.1fGB %-10.1fGB %-5.1f%% %s",
				truncate(fs.MountPoint, 16), truncate(fs.Device, 10), truncate(fs.FSType, 8), usedGB, totalGB, fs.UsedPct, truncate(fs.MountOptions, 20))
			lines = append(lines, row)
		}
		return RenderContainer(p.Title(), strings.Join(lines, "\n"), width, height, active)
	}

	if len(p.stats) == 0 {
		return RenderContainer(p.Title(), "No active Docker containers.", width, height, active)
	}

	var header string
	if expanded {
		header = fmt.Sprintf("%-12s %-12s %-5s %-5s %-12s %-12s %-8s %-4s %-10s %-14s %-12s %s",
			"CONTAINER ID", "NAME", "CPU%", "MEM%", "NET RX/TX", "BLOCK R/W", "HEALTH", "REST", "UPTIME", "IP", "NETMODE", "IMAGE")
	} else {
		header = fmt.Sprintf("%-14s %-16s %-6s %-10s %-10s %-6s %-16s %-16s %-4s",
			"CONTAINER ID", "NAME", "CPU%", "MEM USAGE", "LIMIT", "MEM%", "NET RX / TX", "BLOCK R / W", "PIDS")
	}

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
		if expanded {
			netStr := fmt.Sprintf("%s/%s", stat.NetRx, stat.NetTx)
			blockStr := fmt.Sprintf("%s/%s", stat.BlockRead, stat.BlockWrite)
			row := fmt.Sprintf("%-12s %-12s %-5.1f %-5s %-12s %-12s %-8s %-4d %-10s %-14s %-12s %s",
				stat.ID, truncate(stat.Name, 12), stat.CPUPct, stat.MemPct, truncate(netStr, 12), truncate(blockStr, 12), stat.Health, stat.RestartCount, stat.Uptime, stat.IPAddress, stat.NetworkMode, stat.Image)
			lines = append(lines, row)
		} else {
			netStr := fmt.Sprintf("%s/%s", stat.NetRx, stat.NetTx)
			blockStr := fmt.Sprintf("%s/%s", stat.BlockRead, stat.BlockWrite)
			row := fmt.Sprintf("%-14s %-16s %-6.1f %-10s %-10s %-6s %-16s %-16s %-4d",
				stat.ID, stat.Name, stat.CPUPct, stat.MemUsage, stat.MemLimit, stat.MemPct, netStr, blockStr, stat.PIDs)
			lines = append(lines, row)
		}
	}

	content := strings.Join(lines, "\n")
	return RenderContainer(p.Title(), content, width, height, active)
}
