package app

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestAppNew(t *testing.T) {
	a := New()
	if a.hostname == "" {
		t.Errorf("expected non-empty hostname")
	}
	if a.timeframe != "30s" {
		t.Errorf("expected initial timeframe '30s', got '%s'", a.timeframe)
	}
}

func TestAppUpdateKeyEvents(t *testing.T) {
	a := New()

	// Test timeframe switching
	model, _ := a.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	updated := model.(App)
	if updated.timeframe != "1m" {
		t.Errorf("expected timeframe '1m', got '%s'", updated.timeframe)
	}

	// Test window size
	model, _ = updated.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	updated = model.(App)
	if updated.width != 100 || updated.height != 40 {
		t.Errorf("expected dimensions 100x40, got %dx%d", updated.width, updated.height)
	}

	// Test view rendering
	viewStr := updated.View()
	if viewStr == "" {
		t.Errorf("expected non-empty View string")
	}
}

func TestLayoutCalculation(t *testing.T) {
	layout := ComputeLayout(120, 40)
	if layout.Width != 120 || layout.Height != 40 {
		t.Errorf("expected 120x40 layout, got %dx%d", layout.Width, layout.Height)
	}
	if layout.BodyHeight != 38 {
		t.Errorf("expected BodyHeight 38, got %d", layout.BodyHeight)
	}
}
