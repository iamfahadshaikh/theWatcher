package collector

import (
	"testing"
	"theWatcher/internal/metrics"
)

func TestDockerCollectorMock(t *testing.T) {
	mockOutput := "a1b2c3d4e5f6\tweb-app\t12.5%\t256MiB / 2GiB\t12.8%\t1.2MB / 500KB\t10MB / 5MB\t4\n"

	col := NewDockerCollectorWithRunner(func() ([]byte, error) {
		return []byte(mockOutput), nil
	})

	err := col.Collect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stats := col.Snapshot().([]metrics.DockerStat)
	if len(stats) != 1 {
		t.Fatalf("expected 1 docker stat, got %d", len(stats))
	}

	if stats[0].Name != "web-app" || stats[0].CPUPct != 12.5 || stats[0].PIDs != 4 {
		t.Errorf("unexpected stat for web-app: %+v", stats[0])
	}
}
