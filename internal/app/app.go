package app

import (
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"theWatcher/internal/panel"
)

// TickMsg represents a periodic timer event.
type TickMsg time.Time

// App holds the primary application state.
type App struct {
	width      int
	height     int
	hostname   string
	timeframe  string
	activeIdx  int
	showHelp   bool
	err        error
	tickPeriod time.Duration
	collectors *Collectors

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

// New constructs a new App model with all dashboard panels and collectors.
func New() App {
	host, err := os.Hostname()
	if err != nil || host == "" {
		host = "localhost"
	}

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
		hostname:    host,
		timeframe:   "30s",
		activeIdx:   0,
		tickPeriod:  500 * time.Millisecond,
		collectors:  NewCollectors(),
		cpuPanel:    cpuP,
		memPanel:    memP,
		diskPanel:   diskP,
		netRxPanel:  netRxP,
		netTxPanel:  netTxP,
		procBPanel:  procBP,
		portsPanel:  portsP,
		connPanel:   connP,
		procsPanel:  procsP,
		dockerPanel: dockerP,
		loginPanel:  loginP,
		eventsPanel: eventsP,
		panels:      panels,
	}
}

// Init initializes the BubbleTea event loop with a periodic tick and initial metric fetch.
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
