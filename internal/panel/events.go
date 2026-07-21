package panel

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"theWatcher/internal/history"
	"theWatcher/internal/model"
)

type EventsPanel struct {
	history *history.RingBuffer[model.SystemEvent]
}

func NewEventsPanel() *EventsPanel {
	buf := history.NewRingBuffer[model.SystemEvent](100)
	buf.Push(model.SystemEvent{
		Timestamp: time.Now(),
		Level:     model.LevelInfo,
		Component: "System",
		Message:   "theWatcher Observability Engine started.",
	})
	return &EventsPanel{
		history: buf,
	}
}

func (p *EventsPanel) AddEvent(evt model.SystemEvent) {
	p.history.Push(evt)
}

func (p *EventsPanel) Title() string {
	return "Events"
}

func (p *EventsPanel) Render(width, height int, active bool) string {
	events := p.history.Slice()
	maxLines := height - 4
	if maxLines < 1 {
		maxLines = 1
	}

	if len(events) > maxLines {
		events = events[len(events)-maxLines:]
	}

	var lines []string
	for _, evt := range events {
		tStr := evt.Timestamp.Format("15:04:05")
		lvlStr := fmt.Sprintf("[%s]", evt.Level)

		lvlStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#9ccfd8")) // Info cyan
		if evt.Level == model.LevelWarning {
			lvlStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f6c177")).Bold(true)
		} else if evt.Level == model.LevelAlert {
			lvlStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#eb6f92")).Bold(true)
		}

		line := fmt.Sprintf("%s %-7s %-10s : %s", tStr, lvlStyle.Render(lvlStr), evt.Component, evt.Message)
		lines = append(lines, line)
	}

	content := strings.Join(lines, "\n")
	return RenderContainer(p.Title(), content, width, height, active)
}
