package app

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestAppNew(t *testing.T) {
	a := New()
	if a.timeframe != "30s" {
		t.Errorf("expected initial timeframe '30s', got '%s'", a.timeframe)
	}
}

func TestAppUpdateKeyEvents(t *testing.T) {
	a := New()

	model, _ := a.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	updated := model.(App)
	if updated.timeframe != "1m" {
		t.Errorf("expected timeframe '1m', got '%s'", updated.timeframe)
	}

	model, _ = updated.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	updated = model.(App)
	if updated.width != 100 || updated.height != 40 {
		t.Errorf("expected dimensions 100x40, got %dx%d", updated.width, updated.height)
	}

	viewStr := updated.View()
	if viewStr == "" {
		t.Errorf("expected non-empty View string")
	}
}

func TestAppMaximizeSpaceKey(t *testing.T) {
	a := New()
	if a.focusManager.Maximized {
		t.Errorf("expected initial Maximized to be false")
	}

	model, _ := a.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	updated := model.(App)
	if !updated.focusManager.Maximized {
		t.Errorf("expected Space key to maximize active panel")
	}

	model, _ = updated.Update(tea.KeyMsg{Type: tea.KeyEsc})
	updated = model.(App)
	if updated.focusManager.Maximized {
		t.Errorf("expected Esc key to restore grid view")
	}
}
