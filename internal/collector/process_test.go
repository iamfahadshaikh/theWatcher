package collector

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProcessCollector(t *testing.T) {
	tmpDir := t.TempDir()

	// Mock pid 100
	pid100 := filepath.Join(tmpDir, "100")
	_ = os.MkdirAll(pid100, 0755)
	_ = os.WriteFile(filepath.Join(pid100, "cmdline"), []byte("nginx -g daemon off;\x00"), 0644)
	_ = os.WriteFile(filepath.Join(pid100, "status"), []byte("Name:\tnginx\nVmRSS:\t10240 kB\n"), 0644)

	col := NewProcessCollectorWithPath(tmpDir)
	procs, birthRate, err := col.Collect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(procs) != 1 {
		t.Fatalf("expected 1 process, got %d", len(procs))
	}

	if procs[0].PID != 100 || procs[0].Mem != 10.0 || procs[0].Command != "nginx -g daemon off;" {
		t.Errorf("unexpected process stats: %+v", procs[0])
	}

	_ = birthRate
}
