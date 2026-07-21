package collector

import (
	"os"
	"path/filepath"
	"testing"
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
	m1, err := col.Collect()
	if err != nil {
		t.Fatalf("unexpected error on first collect: %v", err)
	}
	if m1.UsagePct != 0.0 {
		t.Errorf("expected initial usage 0.0, got %f", m1.UsagePct)
	}

	// 1000 ticks passed: 200 non-idle, 800 idle -> 20% usage
	stat2 := "cpu  1100 0 600 9300 0 0 0 0\n"
	err = os.WriteFile(statFile, []byte(stat2), 0644)
	if err != nil {
		t.Fatalf("failed writing stat2: %v", err)
	}

	m2, err := col.Collect()
	if err != nil {
		t.Fatalf("unexpected error on second collect: %v", err)
	}
	if m2.UsagePct < 19.9 || m2.UsagePct > 20.1 {
		t.Errorf("expected ~20%% usage, got %f", m2.UsagePct)
	}
}
