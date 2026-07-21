package collector

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDiskCollector(t *testing.T) {
	tmpDir := t.TempDir()
	diskFile := filepath.Join(tmpDir, "diskstats")

	content1 := `   8       0 sda 1000 0 2000 100 500 0 1000 50 0 150 150
`
	err := os.WriteFile(diskFile, []byte(content1), 0644)
	if err != nil {
		t.Fatalf("failed writing diskstats 1: %v", err)
	}

	col := NewDiskCollectorWithPath(diskFile)
	m1, err := col.Collect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m1.ReadBytes != 2000*512 || m1.WriteBytes != 1000*512 {
		t.Errorf("expected sector byte calculations, got read=%d write=%d", m1.ReadBytes, m1.WriteBytes)
	}

	time.Sleep(10 * time.Millisecond)

	content2 := `   8       0 sda 1500 0 4000 200 1000 0 2000 100 0 300 300
`
	_ = os.WriteFile(diskFile, []byte(content2), 0644)

	m2, err := col.Collect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m2.ReadRateKB <= 0 || m2.WriteRateKB <= 0 {
		t.Errorf("expected positive disk throughput rates, got read=%f write=%f", m2.ReadRateKB, m2.WriteRateKB)
	}
}
