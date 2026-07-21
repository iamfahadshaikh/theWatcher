package collector

import (
	"fmt"
	"testing"
)

func TestLoginCollectorMock(t *testing.T) {
	mockWho := "root     pts/0        2026-07-21 12:00 (:0)\n" +
		"alice    pts/1        2026-07-21 12:05 (192.168.1.50)\n"

	mockLast := "root     pts/0        :0               Tue Jul 21 12:00   still logged in\n"

	col := NewLoginCollectorWithRunner(func(name string, args ...string) ([]byte, error) {
		if name == "who" {
			return []byte(mockWho), nil
		}
		if name == "last" {
			return []byte(mockLast), nil
		}
		return nil, fmt.Errorf("unknown command")
	})

	events, err := col.Collect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(events) < 2 {
		t.Fatalf("expected at least 2 login events, got %d", len(events))
	}

	if events[0].User != "root" || events[1].User != "alice" {
		t.Errorf("unexpected event parsing: %+v", events)
	}
}
