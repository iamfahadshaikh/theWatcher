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

	helpBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#f6c177")).
			Padding(1, 2)
)

// View renders the full dashboard grid layout.
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

	// Render Footer
	footerText := "q: Quit | 1-5: Timeframe | Tab: Switch Panel | h: Help"
	footerBar := footerStyle.Width(layout.Width).Render(footerText)

	if a.showHelp {
		helpText := "=== theWatcher Help ===\n\n" +
			"q / Ctrl+C : Quit application\n" +
			"1 - 5      : Change timeframe (30s, 1m, 5m, 15m, 1h)\n" +
			"Tab        : Focus next panel\n" +
			"Shift+Tab  : Focus previous panel\n" +
			"h          : Toggle this help screen"
		bodyView := helpBoxStyle.Width(layout.Width - 4).Height(layout.BodyHeight - 2).Render(helpText)
		return lipgloss.JoinVertical(lipgloss.Left, headerBar, bodyView, footerBar)
	}

	// Calculate widths
	w3 := layout.Width / 3
	w3Rem := layout.Width - (w3 * 2)
	w2 := layout.Width / 2
	w2Rem := layout.Width - w2
	wFull := layout.Width

	// Calculate row heights (total = layout.BodyHeight)
	// Allocate heights proportionally for 7 rows
	rowH := layout.BodyHeight / 7
	if rowH < 3 {
		rowH = 3
	}

	// Render Row 1: CPU, Memory, Disk
	r1p0 := a.panels[0].Render(w3, rowH, a.activeIdx == 0)
	r1p1 := a.panels[1].Render(w3, rowH, a.activeIdx == 1)
	r1p2 := a.panels[2].Render(w3Rem, rowH, a.activeIdx == 2)
	row1 := lipgloss.JoinHorizontal(lipgloss.Top, r1p0, r1p1, r1p2)

	// Render Row 2: Net RX, Net TX, Process Birth
	r2p3 := a.panels[3].Render(w3, rowH, a.activeIdx == 3)
	r2p4 := a.panels[4].Render(w3, rowH, a.activeIdx == 4)
	r2p5 := a.panels[5].Render(w3Rem, rowH, a.activeIdx == 5)
	row2 := lipgloss.JoinHorizontal(lipgloss.Top, r2p3, r2p4, r2p5)

	// Render Row 3: Listening Ports, Active Connections
	r3p6 := a.panels[6].Render(w2, rowH, a.activeIdx == 6)
	r3p7 := a.panels[7].Render(w2Rem, rowH, a.activeIdx == 7)
	row3 := lipgloss.JoinHorizontal(lipgloss.Top, r3p6, r3p7)

	// Render Row 4: Top Processes
	row4 := a.panels[8].Render(wFull, rowH, a.activeIdx == 8)

	// Render Row 5: Docker Containers
	row5 := a.panels[9].Render(wFull, rowH, a.activeIdx == 9)

	// Render Row 6: Login Timeline
	row6 := a.panels[10].Render(wFull, rowH, a.activeIdx == 10)

	// Render Row 7: Events
	remH := layout.BodyHeight - (rowH * 6)
	if remH < 3 {
		remH = 3
	}
	row7 := a.panels[11].Render(wFull, remH, a.activeIdx == 11)

	dashboard := lipgloss.JoinVertical(lipgloss.Left, row1, row2, row3, row4, row5, row6, row7)
	return lipgloss.JoinVertical(lipgloss.Left, headerBar, dashboard, footerBar)
}
