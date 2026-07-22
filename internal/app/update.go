package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"theWatcher/internal/metrics"
)

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
			a.focusManager.NextFocus()
		case "shift+tab":
			a.focusManager.PrevFocus()
		case " ":
			a.focusManager.ToggleMaximize()
		case "esc":
			a.focusManager.ResetMaximize()
		case "h", "?":
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
		a.latestSystem = msg.System
		a.latestCPU = msg.CPU
		a.latestMem = msg.Memory
		a.latestNet = msg.Network

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

		a.procCount = len(msg.Processes)
		a.dockerCount = len(msg.Docker)

		if msg.CPU.UsagePct > 80.0 {
			a.eventsPanel.AddEvent(metrics.SystemEvent{
				Timestamp: msg.CollectedAt,
				Level:     metrics.LevelWarning,
				Component: "CPU",
				Message:   fmt.Sprintf("High CPU usage: %.1f%%", msg.CPU.UsagePct),
			})
		}
		if msg.Memory.UsedPct > 85.0 {
			a.eventsPanel.AddEvent(metrics.SystemEvent{
				Timestamp: msg.CollectedAt,
				Level:     metrics.LevelWarning,
				Component: "Memory",
				Message:   fmt.Sprintf("High Memory usage: %.1f%%", msg.Memory.UsedPct),
			})
		}
	}

	return a, nil
}
