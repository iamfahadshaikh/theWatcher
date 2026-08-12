package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"theWatcher/internal/metrics"
)

type Snapshot struct {
	Timestamp   time.Time            `json:"timestamp"`
	CPU         metrics.CPUMetric    `json:"cpu"`
	Memory      metrics.MemoryMetric `json:"memory"`
	Disk        metrics.DiskMetric   `json:"disk"`
	Network     metrics.NetworkMetric `json:"network"`
	Connections metrics.ConnSummary  `json:"connections"`
	Docker      []metrics.DockerStat `json:"docker"`
}

func ExportSnapshot(
	cpu metrics.CPUMetric,
	mem metrics.MemoryMetric,
	disk metrics.DiskMetric,
	net metrics.NetworkMetric,
	conn metrics.ConnSummary,
	docker []metrics.DockerStat,
) (string, error) {
	snap := Snapshot{
		Timestamp:   time.Now(),
		CPU:         cpu,
		Memory:      mem,
		Disk:        disk,
		Network:     net,
		Connections: conn,
		Docker:      docker,
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	exportDir := filepath.Join(home, ".local", "share", "thewatcher", "exports")
	_ = os.MkdirAll(exportDir, 0755)

	filename := fmt.Sprintf("snapshot_%s.json", time.Now().Format("20060102_150405"))
	path := filepath.Join(exportDir, filename)

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return "", err
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return "", err
	}

	return path, nil
}
