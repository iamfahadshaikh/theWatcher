package panel

import (
	"fmt"
	"strings"

	"theWatcher/internal/metrics"
)

type PortsPanel struct {
	ports        []metrics.PortInfo
	scrollOffset int
}

func NewPortsPanel() *PortsPanel {
	return &PortsPanel{}
}

func (p *PortsPanel) ScrollUp() {
	p.scrollOffset--
	if p.scrollOffset < 0 {
		p.scrollOffset = 0
	}
}

func (p *PortsPanel) ScrollDown() {
	p.scrollOffset++
}

func (p *PortsPanel) UpdatePorts(ports []metrics.PortInfo) {
	p.ports = ports
}

func (p *PortsPanel) Title() string {
	return "Listening Ports"
}

func (p *PortsPanel) Render(width, height int, active bool, expanded bool) string {
	header := fmt.Sprintf("%-6s %-8s %s", "PROTO", "PORT", "PROCESS")
	var lines []string
	lines = append(lines, header)

	maxLines := height - 4
	if maxLines < 1 {
		maxLines = 1
	}

	if p.scrollOffset > len(p.ports)-maxLines {
		p.scrollOffset = len(p.ports) - maxLines
	}
	if p.scrollOffset < 0 {
		p.scrollOffset = 0
	}

	for i := p.scrollOffset; i < len(p.ports); i++ {
		if len(lines)-1 >= maxLines {
			break
		}
		port := p.ports[i]
		row := fmt.Sprintf("%-6s %-8s %s", port.Protocol, port.Port, port.Process)
		lines = append(lines, row)
	}

	content := strings.Join(lines, "\n")
	return RenderContainer(p.Title(), content, width, height, active)
}
