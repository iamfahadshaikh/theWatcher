package panel

import (
	"strings"
	"testing"
	"time"

	"theWatcher/internal/metrics"
)

func TestEventsPanel(t *testing.T) {
	p := NewEventsPanel()
	if p.Title() != "Events & Alerts" {
		t.Errorf("expected title Events & Alerts, got %s", p.Title())
	}

	p.AddEvent(metrics.SystemEvent{
		Timestamp: time.Now(),
		Level:     metrics.LevelWarning,
		Component: "CPU",
		Message:   "High CPU usage detected > 80%",
	})

	rendered := p.Render(80, 10, true)
	if !strings.Contains(rendered, "High CPU usage detected") {
		t.Errorf("expected warning event message in rendered output: %s", rendered)
	}
}
