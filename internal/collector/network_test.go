package collector

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"theWatcher/internal/metrics"
	"time"
)

func TestNetworkCollector(t *testing.T) {
	tmpDir := t.TempDir()
	netFile := filepath.Join(tmpDir, "dev")

	content1 := `Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
  eth0: 1024000    1000    0    0    0     0          0         0  2048000    2000    0    0    0     0       0          0
`
	err := os.WriteFile(netFile, []byte(content1), 0644)
	if err != nil {
		t.Fatalf("failed writing net dev 1: %v", err)
	}

	col := NewNetworkCollectorWithPath(netFile)
	err = col.Collect(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m1 := col.Snapshot().(metrics.NetworkMetric)
	if m1.RxBytes != 1024000 || m1.TxBytes != 2048000 {
		t.Errorf("expected 1024000 / 2048000 bytes, got %d / %d", m1.RxBytes, m1.TxBytes)
	}

	time.Sleep(10 * time.Millisecond)

	content2 := `Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
  eth0: 2048000    2000    0    0    0     0          0         0  4096000    4000    0    0    0     0       0          0
`
	_ = os.WriteFile(netFile, []byte(content2), 0644)
	err = col.Collect(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m2 := col.Snapshot().(metrics.NetworkMetric)
	if m2.RxRateKB <= 0 || m2.TxRateKB <= 0 {
		t.Errorf("expected positive transfer rates, got rx=%f tx=%f", m2.RxRateKB, m2.TxRateKB)
	}
}
