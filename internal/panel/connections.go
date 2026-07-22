package panel

import (
	"fmt"

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
	total := p.summary.Established + p.summary.Listen + p.summary.TimeWait + p.summary.CloseWait + p.summary.SynSent + p.summary.SynRecv
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
		b5 := graph.RenderBar(fmt.Sprintf("SYN_S  %-3d", p.summary.SynSent), float64(p.summary.SynSent), maxVal, width-6, barStyle, lblStyle)
		b6 := graph.RenderBar(fmt.Sprintf("SYN_R  %-3d", p.summary.SynRecv), float64(p.summary.SynRecv), maxVal, width-6, barStyle, lblStyle)
		content = fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s\n\nTotal Active Sockets: %d", b1, b2, b3, b4, b5, b6, total)
	} else {
		content = fmt.Sprintf("%s\n%s\n%s\n%s", b1, b2, b3, b4)
	}

	return RenderContainer(p.Title(), content, width, height, active)
}
