package panel

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"theWatcher/internal/history"
	"theWatcher/internal/metrics"
	"theWatcher/internal/theme"
)

type EventsPanel struct {
	history      *history.RingBuffer[metrics.SystemEvent]
	scrollOffset int
}

func NewEventsPanel() *EventsPanel {
	buf := history.NewRingBuffer[metrics.SystemEvent](200)
	buf.Push(metrics.SystemEvent{
		Timestamp: time.Now(),
		Level:     metrics.LevelInfo,
		Component: "System",
		Message:   "theWatcher Observability Engine started.",
	})
	return &EventsPanel{
		history: buf,
	}
}

func (p *EventsPanel) ScrollUp() {
	p.scrollOffset--
	if p.scrollOffset < 0 {
		p.scrollOffset = 0
	}
}

func (p *EventsPanel) ScrollDown() {
	p.scrollOffset++
}

func (p *EventsPanel) AddEvent(evt metrics.SystemEvent) {
	p.history.Push(evt)
}

func (p *EventsPanel) Title() string {
	return "Events & Alerts"
}

func (p *EventsPanel) Render(width, height int, active bool, expanded bool) string {
	events := p.history.Slice()
	maxLines := height - 4
	if maxLines < 1 {
		maxLines = 1
	}

	if p.scrollOffset > len(events)-maxLines {
		p.scrollOffset = len(events) - maxLines
	}
	if p.scrollOffset < 0 {
		p.scrollOffset = 0
	}

	// In expanded view we show the scrolled slice, otherwise the latest maxLines items.
	var visibleEvents []metrics.SystemEvent
	if expanded {
		endIdx := p.scrollOffset + maxLines
		if endIdx > len(events) {
			endIdx = len(events)
		}
		if p.scrollOffset < len(events) {
			visibleEvents = events[p.scrollOffset:endIdx]
		}
	} else {
		if len(events) > maxLines {
			visibleEvents = events[len(events)-maxLines:]
		} else {
			visibleEvents = events
		}
	}

	var lines []string
	for _, evt := range visibleEvents {
		tStr := evt.Timestamp.Format("15:04:05")
		lvlStr := fmt.Sprintf("[%s]", evt.Level)

		lvlStyle := lipgloss.NewStyle().Foreground(theme.ColorFoam)
		if evt.Level == metrics.LevelWarning {
			lvlStyle = lipgloss.NewStyle().Foreground(theme.ColorGold).Bold(true)
		} else if evt.Level == metrics.LevelAlert {
			lvlStyle = lipgloss.NewStyle().Foreground(theme.ColorLove).Bold(true)
		}

		line := fmt.Sprintf("%s %-7s %-10s : %s", tStr, lvlStyle.Render(lvlStr), evt.Component, evt.Message)
		lines = append(lines, line)
	}

	content := strings.Join(lines, "\n")
	return RenderContainer(p.Title(), content, width, height, active)
}
