package collector

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"theWatcher/internal/metrics"
)

func TestCPUCollector(t *testing.T) {
	tmpDir := t.TempDir()
	statFile := filepath.Join(tmpDir, "stat")

	stat1 := "cpu  1000 0 500 8500 0 0 0 0\n"
	err := os.WriteFile(statFile, []byte(stat1), 0644)
	if err != nil {
		t.Fatalf("failed writing stat1: %v", err)
	}

	col := NewCPUCollectorWithPath(statFile)
	err = col.Collect(context.Background())
	if err != nil {
		t.Fatalf("unexpected error on first collect: %v", err)
	}
	m1 := col.Snapshot().(metrics.CPUMetric)
	if m1.UsagePct != 0.0 {
		t.Errorf("expected initial usage 0.0, got %f", m1.UsagePct)
	}

	stat2 := "cpu  1100 0 600 9300 0 0 0 0\n"
	err = os.WriteFile(statFile, []byte(stat2), 0644)
	if err != nil {
		t.Fatalf("failed writing stat2: %v", err)
	}

	err = col.Collect(context.Background())
	if err != nil {
		t.Fatalf("unexpected error on second collect: %v", err)
	}
	m2 := col.Snapshot().(metrics.CPUMetric)
	if m2.UsagePct < 19.9 || m2.UsagePct > 20.1 {
		t.Errorf("expected ~20%% usage, got %f", m2.UsagePct)
	}
}
