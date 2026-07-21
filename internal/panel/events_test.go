package panel

import (
	"strings"
	"testing"
	"time"

	"theWatcher/internal/model"
)

func TestEventsPanel(t *testing.T) {
	p := NewEventsPanel()
	if p.Title() != "Events" {
		t.Errorf("expected title Events, got %s", p.Title())
	}

	p.AddEvent(model.SystemEvent{
		Timestamp: time.Now(),
		Level:     model.LevelWarning,
		Component: "CPU",
		Message:   "High CPU usage detected > 80%",
	})

	rendered := p.Render(80, 10, true)
	if !strings.Contains(rendered, "High CPU usage detected") {
		t.Errorf("expected warning event message in rendered output: %s", rendered)
	}
}
