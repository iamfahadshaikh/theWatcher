package panel

import (
	"strings"
	"testing"
)

func TestPanelRendering(t *testing.T) {
	cpu := NewCPUPanel()
	if cpu.Title() != "CPU" {
		t.Errorf("expected CPU title, got %s", cpu.Title())
	}

	rendered := cpu.Render(40, 10, true)
	if !strings.Contains(rendered, "CPU ▲▼") {
		t.Errorf("expected active title indicator in rendered output: %s", rendered)
	}

	mem := NewMemoryPanel()
	renderedMem := mem.Render(40, 10, false)
	if strings.Contains(renderedMem, "▲▼") {
		t.Errorf("did not expect active indicator for inactive panel")
	}
}
