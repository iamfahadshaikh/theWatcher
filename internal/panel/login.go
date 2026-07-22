package panel

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"theWatcher/internal/metrics"
	"theWatcher/internal/theme"
)

type LoginPanel struct {
	events []metrics.LoginEvent
}

func NewLoginPanel() *LoginPanel {
	return &LoginPanel{}
}

func (p *LoginPanel) UpdateLoginEvents(events []metrics.LoginEvent) {
	p.events = events
}

func (p *LoginPanel) Title() string {
	return "Login Timeline"
}

func (p *LoginPanel) Render(width, height int, active bool) string {
	if len(p.events) == 0 {
		return RenderContainer(p.Title(), "No active or recent login events recorded.", width, height, active)
	}

	header := fmt.Sprintf("%-20s %-12s %-12s %s", "TIME", "USER", "TYPE", "FROM")
	var lines []string
	lines = append(lines, header)

	maxLines := height - 4
	if maxLines < 1 {
		maxLines = 1
	}

	for i, evt := range p.events {
		if i >= maxLines {
			break
		}
		tStr := evt.Time
		if len(tStr) > 18 {
			tStr = tStr[:18]
		}

		userStyle := lipgloss.NewStyle().Foreground(theme.ColorText)
		if evt.User == "root" {
			userStyle = lipgloss.NewStyle().Foreground(theme.ColorLove).Bold(true)
		} else {
			userStyle = lipgloss.NewStyle().Foreground(theme.ColorSuccess)
		}

		typeStyle := lipgloss.NewStyle().Foreground(theme.ColorFoam)
		switch evt.Type {
		case "root":
			typeStyle = lipgloss.NewStyle().Foreground(theme.ColorLove).Bold(true)
		case "SSH Login":
			typeStyle = lipgloss.NewStyle().Foreground(theme.ColorIris).Bold(true)
		case "tty":
			typeStyle = lipgloss.NewStyle().Foreground(theme.ColorGold)
		case "system":
			typeStyle = lipgloss.NewStyle().Foreground(theme.ColorPine)
		case "local":
			typeStyle = lipgloss.NewStyle().Foreground(theme.ColorFoam)
		}

		uStr := userStyle.Render(evt.User)
		tType := typeStyle.Render(evt.Type)
		row := fmt.Sprintf("%-20s %-12s %-12s %s", tStr, uStr, tType, evt.From)
		lines = append(lines, row)
	}

	content := strings.Join(lines, "\n")
	return RenderContainer(p.Title(), content, width, height, active)
}
