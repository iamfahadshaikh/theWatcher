package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles all BubbleTea messages and input events.
func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return a, tea.Quit
		case "1":
			a.timeframe = "30s"
		case "2":
			a.timeframe = "1m"
		case "3":
			a.timeframe = "5m"
		case "4":
			a.timeframe = "15m"
		case "5":
			a.timeframe = "1h"
		case "tab":
			if len(a.panels) > 0 {
				a.activeIdx = (a.activeIdx + 1) % len(a.panels)
			}
		case "shift+tab":
			if len(a.panels) > 0 {
				a.activeIdx = (a.activeIdx - 1 + len(a.panels)) % len(a.panels)
			}
		case "h":
			a.showHelp = !a.showHelp
		}

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height

	case TickMsg:
		return a, tickCmd(a.tickPeriod)
	}

	return a, nil
}
