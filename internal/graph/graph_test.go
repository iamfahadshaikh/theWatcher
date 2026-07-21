package graph

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestSparkline(t *testing.T) {
	data := []float64{0, 25, 50, 75, 100}
	res := RenderSparkline(data, 10, lipgloss.NewStyle())
	if res == "" {
		t.Errorf("expected non-empty sparkline rendering")
	}
}

func TestLineGraph(t *testing.T) {
	data := []float64{10, 20, 30, 40, 50}
	res := RenderLineGraph(data, 20, 5, 0, 100, lipgloss.NewStyle())
	if res == "" {
		t.Errorf("expected non-empty line graph rendering")
	}
	lines := strings.Split(res, "\n")
	if len(lines) != 5 {
		t.Errorf("expected 5 rows in line graph canvas, got %d", len(lines))
	}
}

func TestBarGraph(t *testing.T) {
	res := RenderBar("CPU 0", 75.0, 100.0, 40, lipgloss.NewStyle(), lipgloss.NewStyle())
	if !strings.Contains(res, "75.0%") {
		t.Errorf("expected bar graph output to contain 75.0%%, got: %s", res)
	}
}
