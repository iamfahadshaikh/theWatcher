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
	borderColor := theme.ColorPine
	titleStyle := theme.PanelTitleStyle

	if active {
		borderColor = theme.ColorRose
		titleStyle = theme.PanelTitleFocusedStyle
		title = fmt.Sprintf("%s %s", title, theme.IconFocus)
	}

	borderType := theme.FocusedBorder
	style := lipgloss.NewStyle().
		Border(borderType).
		BorderForeground(borderColor).
		Width(width - 2).
		Height(height - 2)

	renderedTitle := titleStyle.Render(title)
	body := fmt.Sprintf("%s\n%s", renderedTitle, content)
	return style.Render(body)
}
