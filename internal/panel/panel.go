package panel

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Panel defines the interface for dashboard UI panels.
type Panel interface {
	Title() string
	Render(width, height int, active bool) string
}

// Container renders a panel border with header title and active focus indicator.
func RenderContainer(title, content string, width, height int, active bool) string {
	borderColor := lipgloss.Color("#31748f") // Dim teal
	if active {
		borderColor = lipgloss.Color("#ebbcba") // Vibrant rose gold highlight
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#e0def4"))

	if active {
		titleStyle = titleStyle.Foreground(lipgloss.Color("#ebbcba"))
		title = fmt.Sprintf("%s ▲▼", title)
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(width - 2).
		Height(height - 2)

	renderedTitle := titleStyle.Render(title)
	body := fmt.Sprintf("%s\n%s", renderedTitle, content)
	return style.Render(body)
}
