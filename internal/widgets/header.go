package widgets

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"theWatcher/internal/config"
	"theWatcher/internal/metrics"
	"theWatcher/internal/theme"
)

func RenderHeader(info metrics.SystemInfo, currentTime string, timeframe string, width int) string {
	if width <= 0 {
		return ""
	}

	// Row 1
	title := theme.HeaderTitleStyle.Render(fmt.Sprintf("%s %s", config.AppName, config.Version))
	host := theme.HeaderItemStyle.Render(fmt.Sprintf("HOST: %s", info.Hostname))
	refresh := theme.HeaderItemStyle.Render("REFRESH: 500ms")
	tf := theme.HeaderItemStyle.Render(fmt.Sprintf("TF: %s", timeframe))
	tStr := theme.HeaderValueStyle.Render(fmt.Sprintf("TIME: %s", currentTime))

	r1Left := lipgloss.JoinHorizontal(lipgloss.Center, title, host, refresh, tf)
	sp1 := width - lipgloss.Width(r1Left) - lipgloss.Width(tStr)
	if sp1 < 0 {
		sp1 = 0
	}
	r1Spacer := lipgloss.NewStyle().Background(theme.ColorSurface).Render(fmt.Sprintf("%*s", sp1, ""))
	row1 := lipgloss.JoinHorizontal(lipgloss.Center, r1Left, r1Spacer, tStr)

	// Row 2
	kern := theme.HeaderItemStyle.Render(fmt.Sprintf("KERNEL: %s", info.Kernel))
	upt := theme.HeaderItemStyle.Render(fmt.Sprintf("UPTIME: %s", info.Uptime))
	load := theme.HeaderValueStyle.Render(fmt.Sprintf("LOAD: %s", info.LoadAvg))

	r2Left := lipgloss.JoinHorizontal(lipgloss.Center, kern, upt, load)
	sp2 := width - lipgloss.Width(r2Left)
	if sp2 < 0 {
		sp2 = 0
	}
	r2Spacer := lipgloss.NewStyle().Background(theme.ColorSurface).Render(fmt.Sprintf("%*s", sp2, ""))
	row2 := lipgloss.JoinHorizontal(lipgloss.Center, r2Left, r2Spacer)

	return lipgloss.JoinVertical(lipgloss.Left, row1, row2)
}
