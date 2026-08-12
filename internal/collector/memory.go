package collector

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"theWatcher/internal/metrics"
)

type MemoryCollector struct {
	mu          sync.Mutex
	meminfoPath string
	samples     uint64
	sumPct      float64
	sessionMin  float64
	sessionMax  float64
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

func (c *MemoryCollector) ResetSession() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.samples = 0
	c.sumPct = 0
	c.sessionMin = 0
	c.sessionMax = 0
	c.snapshot.SessionMin = 0
	c.snapshot.SessionMax = 0
	c.snapshot.SessionAvg = 0
	c.snapshot.SessionSamples = 0
}

func (c *MemoryCollector) Interval() time.Duration {
	return 500 * time.Millisecond
}

func (c *MemoryCollector) Collect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Open(c.meminfoPath)
	if err != nil {
		return fmt.Errorf("failed to open meminfo file: %w", err)
	}
	defer file.Close()

	var totalKB, freeKB, availKB, buffersKB, cachedKB float64
	var swapTotalKB, swapFreeKB, swapCachedKB float64
	var dirtyKB, writebackKB float64
	var hugeTotal, hugeFree int
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
		case "SwapTotal":
			swapTotalKB = valKB
		case "SwapFree":
			swapFreeKB = valKB
		case "SwapCached":
			swapCachedKB = valKB
		case "Dirty":
			dirtyKB = valKB
		case "Writeback":
			writebackKB = valKB
		case "HugePages_Total":
			hugeTotal = int(valKB)
		case "HugePages_Free":
			hugeFree = int(valKB)
		}
	}

	totalMB := totalKB / 1024.0
	freeMB := freeKB / 1024.0
	buffersMB := buffersKB / 1024.0
	cachedMB := cachedKB / 1024.0
	swapTotalMB := swapTotalKB / 1024.0
	swapFreeMB := swapFreeKB / 1024.0
	swapCachedMB := swapCachedKB / 1024.0
	dirtyMB := dirtyKB / 1024.0
	writebackMB := writebackKB / 1024.0
	availMB := availKB / 1024.0

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

	c.samples++
	c.sumPct += usedPct
	avgPct := c.sumPct / float64(c.samples)

	if c.sessionMin == 0 || usedPct < c.sessionMin {
		c.sessionMin = usedPct
	}
	if usedPct > c.sessionMax {
		c.sessionMax = usedPct
	}

	thpStatus := readTHPStatus()
	memPressure := readMemoryPressure()
	oomKills := readOOMKills()

	c.snapshot = metrics.MemoryMetric{
		Timestamp:      time.Now(),
		TotalMB:        totalMB,
		UsedMB:         usedMB,
		FreeMB:         freeMB,
		CachedMB:       cachedMB,
		BufferMB:       buffersMB,
		UsedPct:        usedPct,
		SessionMin:     c.sessionMin,
		SessionMax:     c.sessionMax,
		SessionAvg:     avgPct,
		SessionSamples: int(c.samples),
		SwapTotalMB:    swapTotalMB,
		SwapUsedMB:     swapTotalMB - swapFreeMB,
		SwapFreeMB:     swapFreeMB,
		SwapCachedMB:   swapCachedMB,
		DirtyMB:        dirtyMB,
		WritebackMB:    writebackMB,
		AvailableMB:    availMB,
		THPStatus:      thpStatus,
		MemPressure:    memPressure,
		OOMKills:       oomKills,
		HugePagesTotal: hugeTotal,
		HugePagesFree:  hugeFree,
	}
	return nil
}

func (c *MemoryCollector) Snapshot() any {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.snapshot
}

func readTHPStatus() string {
	path := "/sys/kernel/mm/transparent_hugepage/enabled"
	data, err := os.ReadFile(path)
	if err != nil {
		return "N/A"
	}
	s := string(data)
	start := strings.Index(s, "[")
	end := strings.Index(s, "]")
	if start != -1 && end != -1 && end > start {
		return s[start+1 : end]
	}
	return strings.TrimSpace(s)
}

func readMemoryPressure() string {
	path := "/proc/pressure/memory"
	data, err := os.ReadFile(path)
	if err != nil {
		return "N/A"
	}
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	if scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) >= 4 {
			return fmt.Sprintf("some %s %s %s", fields[1], fields[2], fields[3])
		}
	}
	return "N/A"
}

func readOOMKills() int {
	path := "/proc/vmstat"
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "oom_kill") {
			fields := strings.Fields(line)
			if len(fields) == 2 {
				if val, err := strconv.Atoi(fields[1]); err == nil {
					return val
				}
			}
		}
	}
	return 0
}
