package widgets

import (
	"fmt"

	"theWatcher/internal/theme"
)

type StatusSummary struct {
	CPU        float64
	RAM        float64
	RxRateKB   float64
	TxRateKB   float64
	Processes  int
	Containers int
	Timeframe  string
}

func RenderFooter(s StatusSummary, width int) string {
	if width <= 0 {
		return ""
	}

	statusBar := fmt.Sprintf("STATUS | CPU: %.1f%% | RAM: %.1f%% | RX: %.1fKB/s | TX: %.1fKB/s | PROC: %d | DOCKER: %d | TF: %s",
		s.CPU, s.RAM, s.RxRateKB, s.TxRateKB, s.Processes, s.Containers, s.Timeframe)

	keyHints := "q: Quit | Tab: Focus | Space: Maximize | Esc: Restore | 1-5: Timeframe | h: Help"
	fullStr := fmt.Sprintf("%s  --  %s", statusBar, keyHints)

	return theme.FooterStyle.Width(width).Render(fullStr)
}
