package collector

import (
	"context"
	"fmt"
	"testing"
	"theWatcher/internal/metrics"
)

func TestLoginCollectorMock(t *testing.T) {
	mockWho := "root     pts/0        2026-07-21 12:00 (:0)\n" +
		"alice    pts/1        2026-07-21 12:05 (192.168.1.50)\n"

	mockLast := "root     pts/0        :0               Tue Jul 21 12:00   still logged in\n"

	col := NewLoginCollectorWithRunner(func(ctx context.Context, name string, args ...string) ([]byte, error) {
		if name == "who" {
			return []byte(mockWho), nil
		}
		if name == "last" {
			return []byte(mockLast), nil
		}
		return nil, fmt.Errorf("unknown command")
	})

	err := col.Collect(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	events := col.Snapshot().([]metrics.LoginEvent)
	if len(events) < 2 {
		t.Fatalf("expected at least 2 login events, got %d", len(events))
	}
}
