package app

import (
	"fmt"
	"strings"
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
	}, a.focusManager.Maximized, lay.Width)

	var statusLine string
	if a.searchMode {
		statusLine = lipgloss.NewStyle().
			Foreground(theme.ColorGold).
			Background(theme.ColorSurface).
			Bold(true).
			Render(fmt.Sprintf(" SEARCH: %s▋", a.searchBuffer))
	} else if a.statusMessage != "" && time.Since(a.statusMsgTime) < 5*time.Second {
		statusLine = lipgloss.NewStyle().
			Foreground(theme.ColorFoam).
			Bold(true).
			Render(fmt.Sprintf(" %s", a.statusMessage))
	} else {
		alerts := a.alertEngine.RenderFooterAlerts()
		if alerts != "" {
			statusLine = lipgloss.NewStyle().
				Foreground(theme.ColorLove).
				Bold(true).
				Render(alerts)
		}
	}

	var finalView string

	if a.showHelp {
		helpText := fmt.Sprintf("=== theWatcher Help ===\n\n" +
			"q / Ctrl+C : Quit application\n" +
			"/          : Open search filter across processes\n" +
			"c / m / p  : Sort processes by CPU, MEM, or PID\n" +
			"t / n      : Sort processes by Runtime or Name\n" +
			"R          : Reset session metrics (min, max, averages)\n" +
			"e          : Export JSON metrics snapshot\n" +
			"Arrow Keys : Navigate focus in the dashboard grid\n" +
			"Tab        : In Maximized view, focus next panel\n" +
			"Shift+Tab  : In Maximized view, focus previous panel\n" +
			"Space      : Toggle full-screen zoom on focused panel\n" +
			"Esc        : Restore standard multi-panel grid view\n" +
			"1 - 5      : Change timeframe (30s, 1m, 5m, 15m, 1h)\n" +
			"h / ?      : Toggle this help overlay")
		bodyView := helpBoxStyle.Width(lay.Width - 4).Height(lay.BodyHeight - 2).Render(helpText)
		finalView = lipgloss.JoinVertical(lipgloss.Left, headerBar, bodyView, statusLine, footerBar)
	} else {
		activeIdx := int(a.focusManager.ActiveID)

		// Full-screen Maximized Mode
		if a.focusManager.Maximized && activeIdx >= 0 && activeIdx < len(a.panels) {
			maxBody := a.panels[activeIdx].Render(lay.Width, lay.BodyHeight, true, true)
			finalView = lipgloss.JoinVertical(lipgloss.Left, headerBar, maxBody, statusLine, footerBar)
		} else {
			// Multi-panel Grid Mode
			w3 := lay.Width / 3
			w3Rem := lay.Width - (w3 * 2)
			w2 := lay.Width / 2
			w2Rem := lay.Width - w2
			wFull := lay.Width

			var dashboardRows []string

			if lay.Rows.ShowRow1 {
				r1p0 := a.panels[0].Render(w3, lay.Rows.Row1Height, activeIdx == 0, false)
				r1p1 := a.panels[1].Render(w3, lay.Rows.Row1Height, activeIdx == 1, false)
				r1p2 := a.panels[2].Render(w3Rem, lay.Rows.Row1Height, activeIdx == 2, false)
				dashboardRows = append(dashboardRows, lipgloss.JoinHorizontal(lipgloss.Top, r1p0, r1p1, r1p2))
			}

			if lay.Rows.ShowRow2 {
				r2p3 := a.panels[3].Render(w3, lay.Rows.Row2Height, activeIdx == 3, false)
				r2p4 := a.panels[4].Render(w3, lay.Rows.Row2Height, activeIdx == 4, false)
				r2p5 := a.panels[5].Render(w3Rem, lay.Rows.Row2Height, activeIdx == 5, false)
				dashboardRows = append(dashboardRows, lipgloss.JoinHorizontal(lipgloss.Top, r2p3, r2p4, r2p5))
			}

			if lay.Rows.ShowRow3 {
				r3p6 := a.panels[6].Render(w2, lay.Rows.Row3Height, activeIdx == 6, false)
				r3p7 := a.panels[7].Render(w2Rem, lay.Rows.Row3Height, activeIdx == 7, false)
				dashboardRows = append(dashboardRows, lipgloss.JoinHorizontal(lipgloss.Top, r3p6, r3p7))
			}

			if lay.Rows.ShowRow4 {
				dashboardRows = append(dashboardRows, a.panels[8].Render(wFull, lay.Rows.Row4Height, activeIdx == 8, false))
			}

			if lay.Rows.ShowRow5 {
				dashboardRows = append(dashboardRows, a.panels[9].Render(wFull, lay.Rows.Row5Height, activeIdx == 9, false))
			}

			if lay.Rows.ShowRow6 {
				dashboardRows = append(dashboardRows, a.panels[10].Render(wFull, lay.Rows.Row6Height, activeIdx == 10, false))
			}

			if lay.Rows.ShowRow7 {
				dashboardRows = append(dashboardRows, a.panels[11].Render(wFull, lay.Rows.Row7Height, activeIdx == 11, false))
			}

			dashboard := lipgloss.JoinVertical(lipgloss.Left, dashboardRows...)
			finalView = lipgloss.JoinVertical(lipgloss.Left, headerBar, dashboard, statusLine, footerBar)
		}
	}

	// Safety Clamp: Guarantee the rendered output is never taller than terminal height
	lines := strings.Split(finalView, "\n")
	if len(lines) > a.height {
		lines = lines[:a.height]
	} else if len(lines) < a.height {
		for len(lines) < a.height {
			lines = append(lines, "")
		}
	}
	return strings.Join(lines, "\n")
}
