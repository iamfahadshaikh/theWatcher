package panel

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"theWatcher/internal/theme"
)

type Panel interface {
	Title() string
	Render(width, height int, active bool, expanded bool) string
}

func RenderContainer(title, content string, width, height int, active bool) string {
	var borderColor lipgloss.Color
	var titleStyle lipgloss.Style
	var borderType lipgloss.Border

	if active {
		borderColor = theme.ColorGold
		borderType = theme.MaximizedBorder
		titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.ColorBase).
			Background(theme.ColorGold)
		title = fmt.Sprintf(" %s %s ", title, theme.IconFocus)
	} else {
		borderColor = theme.ColorSubtle
		borderType = theme.DefaultBorder
		titleStyle = theme.PanelTitleStyle
		title = fmt.Sprintf(" %s ", title)
	}

	style := lipgloss.NewStyle().
		Border(borderType).
		BorderForeground(borderColor).
		Width(width - 2).
		Height(height - 2)

	renderedTitle := titleStyle.Render(title)
	body := fmt.Sprintf("%s\n%s", renderedTitle, content)
	return style.Render(body)
}
