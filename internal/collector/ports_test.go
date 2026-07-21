package collector

import (
	"os"
	"path/filepath"
	"testing"
	"theWatcher/internal/metrics"
)

func TestPortsAndConnectionsCollector(t *testing.T) {
	tmpDir := t.TempDir()
	tcpFile := filepath.Join(tmpDir, "tcp")
	udpFile := filepath.Join(tmpDir, "udp")

	tcpContent := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 00000000:0050 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 1000 1 0000000000000000 100 0 0 10 0
   1: 00000000:01BB 00000000:0000 01 00000000:00000000 00:00000000 00000000     0        0 1001 1 0000000000000000 100 0 0 10 0
`
	_ = os.WriteFile(tcpFile, []byte(tcpContent), 0644)
	_ = os.WriteFile(udpFile, []byte(""), 0644)

	portsCol := NewPortsCollectorWithPaths(tcpFile, udpFile)
	err := portsCol.Collect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ports := portsCol.Snapshot().([]metrics.PortInfo)
	if len(ports) != 1 || ports[0].Port != "80" {
		t.Errorf("expected 1 listening port on 80, got %+v", ports)
	}

	connCol := NewConnectionsCollectorWithPaths([]string{tcpFile})
	err = connCol.Collect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	conn := connCol.Snapshot().(metrics.ConnSummary)
	if conn.Listen != 1 || conn.Established != 1 {
		t.Errorf("expected 1 Listen and 1 Established, got %+v", conn)
	}
}
