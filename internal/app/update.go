package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"theWatcher/internal/layout"
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
			if a.focusManager.Maximized {
				lay := layout.ComputeLayout(a.width, a.height)
				a.focusManager.NextFocus(lay.Rows)
			}
		case "shift+tab":
			if a.focusManager.Maximized {
				lay := layout.ComputeLayout(a.width, a.height)
				a.focusManager.PrevFocus(lay.Rows)
			}
		case " ":
			a.focusManager.ToggleMaximize()
		case "esc":
			a.focusManager.ResetMaximize()
		case "h", "?":
			a.showHelp = !a.showHelp
		case "left", "right", "up", "down":
			if !a.focusManager.Maximized {
				lay := layout.ComputeLayout(a.width, a.height)
				a.focusManager.Navigate(msg.String(), lay.Rows)
			} else {
				if msg.String() == "up" {
					switch a.focusManager.ActiveID {
					case layout.PanelProcesses:
						a.procsPanel.ScrollUp()
					case layout.PanelDocker:
						a.dockerPanel.ScrollUp()
					case layout.PanelLogins:
						a.loginPanel.ScrollUp()
					case layout.PanelEvents:
						a.eventsPanel.ScrollUp()
					case layout.PanelPorts:
						a.portsPanel.ScrollUp()
					}
				} else if msg.String() == "down" {
					switch a.focusManager.ActiveID {
					case layout.PanelProcesses:
						a.procsPanel.ScrollDown()
					case layout.PanelDocker:
						a.dockerPanel.ScrollDown()
					case layout.PanelLogins:
						a.loginPanel.ScrollDown()
					case layout.PanelEvents:
						a.eventsPanel.ScrollDown()
					case layout.PanelPorts:
						a.portsPanel.ScrollDown()
					}
				}
			}
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

		// 1. New Listening Port Detection
		currPorts := make(map[string]string)
		for _, p := range msg.Ports {
			currPorts[p.Port] = p.Process
			if len(a.prevPorts) > 0 {
				if _, ok := a.prevPorts[p.Port]; !ok {
					a.eventsPanel.AddEvent(metrics.SystemEvent{
						Timestamp: msg.CollectedAt,
						Level:     metrics.LevelWarning,
						Component: "Ports",
						Message:   fmt.Sprintf("New listening port: %s (%s)", p.Port, p.Process),
					})
				}
			}
		}
		a.prevPorts = currPorts

		// 2. Docker Container Event Detection
		currContainers := make(map[string]string)
		for _, d := range msg.Docker {
			currContainers[d.ID] = d.Name
			if len(a.prevContainers) > 0 {
				if _, ok := a.prevContainers[d.ID]; !ok {
					a.eventsPanel.AddEvent(metrics.SystemEvent{
						Timestamp: msg.CollectedAt,
						Level:     metrics.LevelInfo,
						Component: "Docker",
						Message:   fmt.Sprintf("Docker container started: %s", d.Name),
					})
				}
			}
		}
		for id, name := range a.prevContainers {
			found := false
			for _, d := range msg.Docker {
				if d.ID == id {
					found = true
					break
				}
			}
			if !found {
				a.eventsPanel.AddEvent(metrics.SystemEvent{
					Timestamp: msg.CollectedAt,
					Level:     metrics.LevelWarning,
					Component: "Docker",
					Message:   fmt.Sprintf("Docker container stopped: %s", name),
				})
			}
		}
		a.prevContainers = currContainers

		// 3. PID Burst Detection
		if msg.BirthRate > 15.0 {
			a.eventsPanel.AddEvent(metrics.SystemEvent{
				Timestamp: msg.CollectedAt,
				Level:     metrics.LevelWarning,
				Component: "Process",
				Message:   fmt.Sprintf("PID burst detected: %.1f spawns/sec", msg.BirthRate),
			})
		}

		// 4. Memory Reclaim Detection
		if a.prevMemPct > 0 {
			diff := a.prevMemPct - msg.Memory.UsedPct
			if diff > 10.0 {
				a.eventsPanel.AddEvent(metrics.SystemEvent{
					Timestamp: msg.CollectedAt,
					Level:     metrics.LevelInfo,
					Component: "Memory",
					Message:   fmt.Sprintf("Memory reclaimed: %.1f%% decrease", diff),
				})
			}
		}
		a.prevMemPct = msg.Memory.UsedPct
	}

	return a, nil
}
