package panel

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"theWatcher/internal/graph"
	"theWatcher/internal/metrics"
	"theWatcher/internal/theme"
)

type ConnectionsPanel struct {
	summary metrics.ConnSummary
}

func NewConnectionsPanel() *ConnectionsPanel {
	return &ConnectionsPanel{}
}

func (p *ConnectionsPanel) UpdateConnections(summary metrics.ConnSummary) {
	p.summary = summary
}

func (p *ConnectionsPanel) Title() string {
	return "Active Connections"
}

func (p *ConnectionsPanel) Render(width, height int, active bool, expanded bool) string {
	total := p.summary.Established + p.summary.Listen + p.summary.TimeWait + p.summary.CloseWait + p.summary.SynSent + p.summary.SynRecv + p.summary.FinWait
	maxVal := float64(total)
	if maxVal < 10 {
		maxVal = 10
	}

	barStyle := lipgloss.NewStyle().Foreground(theme.ColorFoam)
	lblStyle := lipgloss.NewStyle().Bold(true).Foreground(theme.ColorText)

	b1 := graph.RenderBar(fmt.Sprintf("ESTAB  %-3d", p.summary.Established), float64(p.summary.Established), maxVal, width-6, barStyle, lblStyle)
	b2 := graph.RenderBar(fmt.Sprintf("LISTEN %-3d", p.summary.Listen), float64(p.summary.Listen), maxVal, width-6, barStyle, lblStyle)
	b3 := graph.RenderBar(fmt.Sprintf("TIME_W %-3d", p.summary.TimeWait), float64(p.summary.TimeWait), maxVal, width-6, barStyle, lblStyle)
	b4 := graph.RenderBar(fmt.Sprintf("CLOS_W %-3d", p.summary.CloseWait), float64(p.summary.CloseWait), maxVal, width-6, barStyle, lblStyle)

	var content string
	if expanded {
		r1 := fmt.Sprintf("STATE | ESTAB: %d  LISTEN: %d  TIME_W: %d  CLOS_W: %d  FIN_W: %d  SYN_S: %d  SYN_R: %d",
			p.summary.Established, p.summary.Listen, p.summary.TimeWait, p.summary.CloseWait, p.summary.FinWait, p.summary.SynSent, p.summary.SynRecv)
		r2 := fmt.Sprintf("PROTO | UDP Sockets: %d  UNIX Sockets: %d  RAW Sockets: %d  Conn Rate: %.1f/s",
			p.summary.UDP, p.summary.UNIX, p.summary.RAW, p.summary.ConnRate)
		var ipFields []string
		for _, ip := range p.summary.TopIPs {
			ipFields = append(ipFields, fmt.Sprintf("%s (%d)", ip.IP, ip.Count))
		}
		r3 := fmt.Sprintf("TOP IP| %s", strings.Join(ipFields, "  "))
		var portFields []string
		for _, pt := range p.summary.TopPorts {
			portFields = append(portFields, fmt.Sprintf("%d (%d)", pt.Port, pt.Count))
		}
		r4 := fmt.Sprintf("PORT  | %s", strings.Join(portFields, "  "))
		content = fmt.Sprintf("%s\n%s\n%s\n%s", r1, r2, r3, r4)
	} else {
		if height < 6 {
			content = fmt.Sprintf("ESTAB:%d  LISTEN:%d  TIME_W:%d  CLOS_W:%d", p.summary.Established, p.summary.Listen, p.summary.TimeWait, p.summary.CloseWait)
		} else {
			content = fmt.Sprintf("%s\n%s\n%s\n%s", b1, b2, b3, b4)
		}
	}

	return RenderContainer(p.Title(), content, width, height, active)
}
