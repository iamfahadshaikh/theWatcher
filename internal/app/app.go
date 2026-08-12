package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"theWatcher/internal/config"
	"theWatcher/internal/layout"
	"theWatcher/internal/metrics"
	"theWatcher/internal/panel"
	"theWatcher/internal/theme"
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
	store        *MetricsStore
	scheduler    *Scheduler
	latestSystem  metrics.SystemInfo
	latestCPU     metrics.CPUMetric
	latestMem     metrics.MemoryMetric
	latestNet     metrics.NetworkMetric
	latestDisk    metrics.DiskMetric
	latestConn    metrics.ConnSummary
	latestDocker  []metrics.DockerStat
	latestDockerErr error

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

	procCount   int
	dockerCount int

	prevPorts      map[string]string
	prevContainers map[string]string
	prevMemPct     float64
	alertEngine    *AlertEngine
	searchMode     bool
	searchBuffer   string
	statusMessage  string
	statusMsgTime  time.Time
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

	// Load user config and set theme
	cfg := config.LoadConfig()
	theme.SetTheme(cfg.Theme)

	cols := NewCollectors()
	store := NewMetricsStore()
	sched := NewScheduler(cols, store)

	return App{
		timeframe:      "30s",
		tickPeriod:     config.GetRefreshPeriod(),
		focusManager:   layout.NewFocusManager(),
		collectors:     cols,
		store:          store,
		scheduler:      sched,
		cpuPanel:       cpuP,
		memPanel:       memP,
		diskPanel:      diskP,
		netRxPanel:     netRxP,
		netTxPanel:     netTxP,
		procBPanel:     procBP,
		portsPanel:     portsP,
		connPanel:      connP,
		procsPanel:     procsP,
		dockerPanel:    dockerP,
		loginPanel:     loginP,
		eventsPanel:    eventsP,
		panels:         panels,
		prevPorts:      make(map[string]string),
		prevContainers: make(map[string]string),
		alertEngine:    NewAlertEngine(),
	}
}

func (a App) Init() tea.Cmd {
	a.scheduler.Start()
	return tea.Batch(
		tickCmd(a.tickPeriod),
		FetchMetricsCmd(a.store),
	)
}

func (a App) Stop() {
	if a.scheduler != nil {
		a.scheduler.Stop()
	}
}

func tickCmd(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
