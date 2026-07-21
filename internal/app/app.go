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
	panels     []panel.Panel
}

// New constructs a new App model with all dashboard panels.
func New() App {
	host, err := os.Hostname()
	if err != nil || host == "" {
		host = "localhost"
	}

	panels := []panel.Panel{
		panel.NewCPUPanel(),
		panel.NewMemoryPanel(),
		panel.NewDiskPanel(),
		panel.NewNetworkPanel(false), // RX
		panel.NewNetworkPanel(true),  // TX
		panel.NewProcessPanel(true),  // Process Birth
		panel.NewPortsPanel(),
		panel.NewConnectionsPanel(),
		panel.NewProcessPanel(false), // Top Processes
		panel.NewDockerPanel(),
		panel.NewLoginPanel(),
		panel.NewEventsPanel(),
	}

	return App{
		hostname:   host,
		timeframe:  "30s",
		activeIdx:  0,
		tickPeriod: 500 * time.Millisecond,
		panels:     panels,
	}
}

// Init initializes the BubbleTea event loop with a periodic tick.
func (a App) Init() tea.Cmd {
	return tickCmd(a.tickPeriod)
}

func tickCmd(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
