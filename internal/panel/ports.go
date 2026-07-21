package panel

import (
	"fmt"
	"strings"

	"theWatcher/internal/metrics"
)

type PortsPanel struct {
	ports []metrics.PortInfo
}

func NewPortsPanel() *PortsPanel {
	return &PortsPanel{}
}

func (p *PortsPanel) UpdatePorts(ports []metrics.PortInfo) {
	p.ports = ports
}

func (p *PortsPanel) Title() string {
	return "Listening Ports"
}

func (p *PortsPanel) Render(width, height int, active bool) string {
	header := fmt.Sprintf("%-6s %-8s %s", "PROTO", "PORT", "PROCESS")
	var lines []string
	lines = append(lines, header)

	maxLines := height - 4
	if maxLines < 1 {
		maxLines = 1
	}

	for i, port := range p.ports {
		if i >= maxLines {
			break
		}
		row := fmt.Sprintf("%-6s %-8s %s", port.Protocol, port.Port, port.Process)
		lines = append(lines, row)
	}

	content := strings.Join(lines, "\n")
	return RenderContainer(p.Title(), content, width, height, active)
}
