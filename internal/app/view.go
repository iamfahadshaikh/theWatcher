package app

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"theWatcher/internal/layout"
	"theWatcher/internal/theme"
	"theWatcher/internal/widgets"
)

var helpBoxStyle = lipgloss.NewStyle().
	Border(theme.MaximizedBorder).
	BorderForeground(theme.ColorGold).
	Padding(1, 2)

func (a App) View() string {
	if a.width <= 0 || a.height <= 0 {
		return "Initializing theWatcher..."
	}

	lay := layout.ComputeLayout(a.width, a.height)
	currTime := time.Now().Format("15:04:05 MST")

	headerBar := widgets.RenderHeader(a.latestSystem, currTime, a.timeframe, lay.Width)

	footerBar := widgets.RenderFooter(widgets.StatusSummary{
		CPU:        a.latestCPU.UsagePct,
		RAM:        a.latestMem.UsedPct,
		RxRateKB:   a.latestNet.RxRateKB,
		TxRateKB:   a.latestNet.TxRateKB,
		Processes:  a.procCount,
		Containers: a.dockerCount,
		Timeframe:  a.timeframe,
	}, lay.Width)

	if a.showHelp {
		helpText := fmt.Sprintf("=== theWatcher Help ===\n\n" +
			"q / Ctrl+C : Quit application\n" +
			"Tab        : Move focus to next panel\n" +
			"Shift+Tab  : Move focus to previous panel\n" +
			"Space      : Toggle full-screen zoom on focused panel\n" +
			"Esc        : Restore standard multi-panel grid view\n" +
			"1 - 5      : Change timeframe (30s, 1m, 5m, 15m, 1h)\n" +
			"h / ?      : Toggle this help overlay")
		bodyView := helpBoxStyle.Width(lay.Width - 4).Height(lay.BodyHeight - 2).Render(helpText)
		return lipgloss.JoinVertical(lipgloss.Left, headerBar, bodyView, footerBar)
	}

	activeIdx := int(a.focusManager.ActiveID)

	// Full-screen Maximized Mode
	if a.focusManager.Maximized && activeIdx >= 0 && activeIdx < len(a.panels) {
		maxBody := a.panels[activeIdx].Render(lay.Width, lay.BodyHeight, true)
		return lipgloss.JoinVertical(lipgloss.Left, headerBar, maxBody, footerBar)
	}

	// Multi-panel Grid Mode
	w3 := lay.Width / 3
	w3Rem := lay.Width - (w3 * 2)
	w2 := lay.Width / 2
	w2Rem := lay.Width - w2
	wFull := lay.Width

	rowH := lay.BodyHeight / 7
	if rowH < 3 {
		rowH = 3
	}

	r1p0 := a.panels[0].Render(w3, rowH, activeIdx == 0)
	r1p1 := a.panels[1].Render(w3, rowH, activeIdx == 1)
	r1p2 := a.panels[2].Render(w3Rem, rowH, activeIdx == 2)
	row1 := lipgloss.JoinHorizontal(lipgloss.Top, r1p0, r1p1, r1p2)

	r2p3 := a.panels[3].Render(w3, rowH, activeIdx == 3)
	r2p4 := a.panels[4].Render(w3, rowH, activeIdx == 4)
	r2p5 := a.panels[5].Render(w3Rem, rowH, activeIdx == 5)
	row2 := lipgloss.JoinHorizontal(lipgloss.Top, r2p3, r2p4, r2p5)

	r3p6 := a.panels[6].Render(w2, rowH, activeIdx == 6)
	r3p7 := a.panels[7].Render(w2Rem, rowH, activeIdx == 7)
	row3 := lipgloss.JoinHorizontal(lipgloss.Top, r3p6, r3p7)

	row4 := a.panels[8].Render(wFull, rowH, activeIdx == 8)
	row5 := a.panels[9].Render(wFull, rowH, activeIdx == 9)
	row6 := a.panels[10].Render(wFull, rowH, activeIdx == 10)

	remH := lay.BodyHeight - (rowH * 6)
	if remH < 3 {
		remH = 3
	}
	row7 := a.panels[11].Render(wFull, remH, activeIdx == 11)

	dashboard := lipgloss.JoinVertical(lipgloss.Left, row1, row2, row3, row4, row5, row6, row7)
	return lipgloss.JoinVertical(lipgloss.Left, headerBar, dashboard, footerBar)
}
