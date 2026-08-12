package collector

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestProcessCollector(t *testing.T) {
	tmpDir := t.TempDir()

	pid100 := filepath.Join(tmpDir, "100")
	_ = os.MkdirAll(pid100, 0755)
	_ = os.WriteFile(filepath.Join(pid100, "cmdline"), []byte("nginx -g daemon off;\x00"), 0644)
	_ = os.WriteFile(filepath.Join(pid100, "status"), []byte("Name:\tnginx\nVmRSS:\t10240 kB\n"), 0644)
	_ = os.WriteFile(filepath.Join(pid100, "stat"), []byte("100 (nginx) S 1 100 100 0 -1 4194304 100 0 0 0 10 20 0 0 20 0 1 0 1000\n"), 0644)

	col := NewProcessCollectorWithPath(tmpDir)
	err := col.Collect(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snap := col.Snapshot().(ProcessSnapshot)
	if len(snap.Processes) != 1 {
		t.Fatalf("expected 1 process, got %d", len(snap.Processes))
	}

	p := snap.Processes[0]
	if p.PID != 100 || p.Mem != 10.0 || p.PPID != 1 {
		t.Errorf("unexpected process stats: %+v", p)
	}
}
