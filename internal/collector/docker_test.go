package collector

import (
	"fmt"
	"testing"
)

func TestDockerCollectorMock(t *testing.T) {
	mockOutput := "a1b2c3d4e5f6\tweb-app\t12.5%\t256MiB / 2GiB\t1.2MB / 500KB\n" +
		"f6e5d4c3b2a1\tdb-postgres\t2.1%\t512MiB / 4GiB\t5.4MB / 10.2MB\n"

	col := NewDockerCollectorWithRunner(func() ([]byte, error) {
		return []byte(mockOutput), nil
	})

	stats, err := col.Collect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(stats) != 2 {
		t.Fatalf("expected 2 docker stats, got %d", len(stats))
	}

	if stats[0].Name != "web-app" || stats[0].CPUPct != 12.5 {
		t.Errorf("unexpected stat for web-app: %+v", stats[0])
	}
}

func TestDockerCollectorUnavailable(t *testing.T) {
	col := NewDockerCollectorWithRunner(func() ([]byte, error) {
		return nil, fmt.Errorf("docker command not found")
	})

	_, err := col.Collect()
	if err == nil {
		t.Errorf("expected error when docker command fails")
	}
}
