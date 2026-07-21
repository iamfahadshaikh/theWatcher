package app

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#eb6f92")).
			Background(lipgloss.Color("#1f1d2e")).
			Padding(0, 1)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9ccfd8")).
			Background(lipgloss.Color("#1f1d2e")).
			Padding(0, 1)

	timeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f6c177")).
			Background(lipgloss.Color("#1f1d2e")).
			Padding(0, 1)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#908caa")).
			Background(lipgloss.Color("#191724")).
			Padding(0, 1)

	bodyStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#31748f")).
			Padding(1)
)

// View renders the terminal user interface.
func (a App) View() string {
	if a.width <= 0 || a.height <= 0 {
		return "Initializing theWatcher..."
	}

	layout := ComputeLayout(a.width, a.height)

	// Render Header
	title := headerStyle.Render("theWatcher")
	host := infoStyle.Render(fmt.Sprintf("HOST: %s", a.hostname))
	tf := infoStyle.Render(fmt.Sprintf("TF: %s", a.timeframe))
	currTime := timeStyle.Render(time.Now().Format("15:04:05 MST"))

	leftHeader := lipgloss.JoinHorizontal(lipgloss.Center, title, host, tf)
	spaceLen := layout.Width - lipgloss.Width(leftHeader) - lipgloss.Width(currTime)
	if spaceLen < 0 {
		spaceLen = 0
	}
	spacer := lipgloss.NewStyle().Background(lipgloss.Color("#1f1d2e")).Render(fmt.Sprintf("%*s", spaceLen, ""))
	headerBar := lipgloss.JoinHorizontal(lipgloss.Center, leftHeader, spacer, currTime)

	// Render Body placeholder
	content := fmt.Sprintf("theWatcher Dashboard Framework Ready\nWindow Size: %dx%d\nPress 'h' for help, 'q' to quit.", a.width, a.height)
	if a.showHelp {
		content = "=== HELP MENU ===\n'q' / Ctrl+C : Quit application\n'1' - '5'   : Change timeframe (30s, 1m, 5m, 15m, 1h)\n'Tab'       : Navigate next panel\n'h'         : Toggle help screen"
	}
	bodyView := bodyStyle.Width(layout.Width - 2).Height(layout.BodyHeight - 2).Render(content)

	// Render Footer
	footerText := "q: Quit | 1-5: Timeframe | Tab: Switch Panel | h: Help"
	footerBar := footerStyle.Width(layout.Width).Render(footerText)

	return lipgloss.JoinVertical(lipgloss.Left, headerBar, bodyView, footerBar)
}
