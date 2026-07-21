package collector

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"theWatcher/internal/metrics"
)

type MemoryCollector struct {
	meminfoPath string
	snapshot    metrics.MemoryMetric
}

func NewMemoryCollector() *MemoryCollector {
	return &MemoryCollector{meminfoPath: "/proc/meminfo"}
}

func NewMemoryCollectorWithPath(path string) *MemoryCollector {
	return &MemoryCollector{meminfoPath: path}
}

func (c *MemoryCollector) Name() string {
	return "Memory"
}

func (c *MemoryCollector) Collect() error {
	file, err := os.Open(c.meminfoPath)
	if err != nil {
		return fmt.Errorf("failed to open meminfo file: %w", err)
	}
	defer file.Close()

	var totalKB, freeKB, availKB, buffersKB, cachedKB float64
	hasAvail := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		valStr := strings.TrimSpace(parts[1])
		valFields := strings.Fields(valStr)
		if len(valFields) == 0 {
			continue
		}

		valKB, err := strconv.ParseFloat(valFields[0], 64)
		if err != nil {
			continue
		}

		switch key {
		case "MemTotal":
			totalKB = valKB
		case "MemFree":
			freeKB = valKB
		case "MemAvailable":
			availKB = valKB
			hasAvail = true
		case "Buffers":
			buffersKB = valKB
		case "Cached":
			cachedKB = valKB
		}
	}

	totalMB := totalKB / 1024.0
	freeMB := freeKB / 1024.0
	buffersMB := buffersKB / 1024.0
	cachedMB := cachedKB / 1024.0

	var usedMB float64
	if hasAvail {
		usedMB = (totalKB - availKB) / 1024.0
	} else {
		usedMB = (totalKB - freeKB - buffersKB - cachedKB) / 1024.0
	}
	if usedMB < 0 {
		usedMB = 0
	}

	usedPct := 0.0
	if totalMB > 0 {
		usedPct = (usedMB / totalMB) * 100.0
	}

	c.snapshot = metrics.MemoryMetric{
		Timestamp: time.Now(),
		TotalMB:   totalMB,
		UsedMB:    usedMB,
		FreeMB:    freeMB,
		CachedMB:  cachedMB,
		BufferMB:  buffersMB,
		UsedPct:   usedPct,
	}
	return nil
}

func (c *MemoryCollector) Snapshot() any {
	return c.snapshot
}
