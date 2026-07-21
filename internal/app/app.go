package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"theWatcher/internal/config"
	"theWatcher/internal/layout"
	"theWatcher/internal/metrics"
	"theWatcher/internal/panel"
)

type TickMsg time.Time

type App struct {
	width        int
	height       int
	timeframe    string
	showHelp     bool
	tickPeriod   time.Duration
	focusManager *layout.FocusManager
	collectors   *Collectors
	latestSystem metrics.SystemInfo
	latestCPU    metrics.CPUMetric
	latestMem    metrics.MemoryMetric
	latestNet    metrics.NetworkMetric

	cpuPanel    *panel.CPUPanel
	memPanel    *panel.MemoryPanel
	diskPanel   *panel.DiskPanel
	netRxPanel  *panel.NetworkPanel
	netTxPanel  *panel.NetworkPanel
	procBPanel  *panel.ProcessPanel
	portsPanel  *panel.PortsPanel
	connPanel   *panel.ConnectionsPanel
	procsPanel  *panel.ProcessPanel
	dockerPanel *panel.DockerPanel
	loginPanel  *panel.LoginPanel
	eventsPanel *panel.EventsPanel

	panels []panel.Panel
}

func New() App {
	cpuP := panel.NewCPUPanel()
	memP := panel.NewMemoryPanel()
	diskP := panel.NewDiskPanel()
	netRxP := panel.NewNetworkPanel(false)
	netTxP := panel.NewNetworkPanel(true)
	procBP := panel.NewProcessPanel(true)
	portsP := panel.NewPortsPanel()
	connP := panel.NewConnectionsPanel()
	procsP := panel.NewProcessPanel(false)
	dockerP := panel.NewDockerPanel()
	loginP := panel.NewLoginPanel()
	eventsP := panel.NewEventsPanel()

	panels := []panel.Panel{
		cpuP, memP, diskP, netRxP, netTxP, procBP,
		portsP, connP, procsP, dockerP, loginP, eventsP,
	}

	return App{
		timeframe:    "30s",
		tickPeriod:   config.RefreshRate,
		focusManager: layout.NewFocusManager(),
		collectors:   NewCollectors(),
		cpuPanel:     cpuP,
		memPanel:     memP,
		diskPanel:    diskP,
		netRxPanel:   netRxP,
		netTxPanel:   netTxP,
		procBPanel:   procBP,
		portsPanel:   portsP,
		connPanel:    connP,
		procsPanel:   procsP,
		dockerPanel:  dockerP,
		loginPanel:   loginP,
		eventsPanel:  eventsP,
		panels:       panels,
	}
}

func (a App) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(a.tickPeriod),
		FetchMetricsCmd(a.collectors),
	)
}

func tickCmd(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
