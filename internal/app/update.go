package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"theWatcher/internal/model"
)

// Update handles all BubbleTea messages, input events, and metric updates.
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
		return a, tea.Batch(
			tickCmd(a.tickPeriod),
			FetchMetricsCmd(a.collectors),
		)

	case MetricsMsg:
		a.cpuPanel.UpdateMetric(msg.CPU)
		a.memPanel.UpdateMetric(msg.Memory)
		a.diskPanel.UpdateMetric(msg.Disk)
		a.netRxPanel.UpdateMetric(msg.Network)
		a.netTxPanel.UpdateMetric(msg.Network)
		a.procBPanel.UpdateProcesses(msg.Processes, msg.BirthRate)
		a.portsPanel.UpdatePorts(msg.Ports)
		a.connPanel.UpdateConnections(msg.Conn)
		a.procsPanel.UpdateProcesses(msg.Processes, msg.BirthRate)
		a.dockerPanel.UpdateDockerStats(msg.Docker, msg.DockerErr)
		a.loginPanel.UpdateLoginEvents(msg.Logins)

		// Threshold Alert Generation
		if msg.CPU.UsagePct > 80.0 {
			a.eventsPanel.AddEvent(model.SystemEvent{
				Timestamp: msg.CollectedAt,
				Level:     model.LevelWarning,
				Component: "CPU",
				Message:   fmt.Sprintf("High CPU usage: %.1f%%", msg.CPU.UsagePct),
			})
		}
		if msg.Memory.UsedPct > 85.0 {
			a.eventsPanel.AddEvent(model.SystemEvent{
				Timestamp: msg.CollectedAt,
				Level:     model.LevelWarning,
				Component: "Memory",
				Message:   fmt.Sprintf("High Memory usage: %.1f%%", msg.Memory.UsedPct),
			})
		}
	}

	return a, nil
}
