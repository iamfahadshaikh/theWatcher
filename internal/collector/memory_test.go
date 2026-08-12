package collector

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"theWatcher/internal/metrics"
)

func TestMemoryCollector(t *testing.T) {
	tmpDir := t.TempDir()
	memFile := filepath.Join(tmpDir, "meminfo")

	content := `MemTotal:       16384000 kB
MemFree:         4096000 kB
MemAvailable:    8192000 kB
Buffers:         1024000 kB
Cached:          2048000 kB
`
	err := os.WriteFile(memFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed writing meminfo: %v", err)
	}

	col := NewMemoryCollectorWithPath(memFile)
	err = col.Collect(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m := col.Snapshot().(metrics.MemoryMetric)
	if m.TotalMB != 16000.0 {
		t.Errorf("expected TotalMB 16000, got %f", m.TotalMB)
	}
	if m.UsedMB != 8000.0 {
		t.Errorf("expected UsedMB 8000, got %f", m.UsedMB)
	}
}
