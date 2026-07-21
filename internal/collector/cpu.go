package collector

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"theWatcher/internal/model"
)

type CPUCollector struct {
	statPath  string
	prevIdle  uint64
	prevTotal uint64
	sumPct    float64
	samples   uint64
	peakPct   float64
}

func NewCPUCollector() *CPUCollector {
	return &CPUCollector{statPath: "/proc/stat"}
}

func NewCPUCollectorWithPath(path string) *CPUCollector {
	return &CPUCollector{statPath: path}
}

func (c *CPUCollector) Collect() (model.CPUMetric, error) {
	file, err := os.Open(c.statPath)
	if err != nil {
		return model.CPUMetric{Timestamp: time.Now()}, fmt.Errorf("failed to open stat file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return model.CPUMetric{Timestamp: time.Now()}, fmt.Errorf("empty stat file")
	}

	line := scanner.Text()
	fields := strings.Fields(line)
	if len(fields) < 5 || fields[0] != "cpu" {
		return model.CPUMetric{Timestamp: time.Now()}, fmt.Errorf("invalid cpu line: %s", line)
	}

	var values []uint64
	for _, f := range fields[1:] {
		val, err := strconv.ParseUint(f, 10, 64)
		if err != nil {
			break
		}
		values = append(values, val)
	}

	if len(values) < 4 {
		return model.CPUMetric{Timestamp: time.Now()}, fmt.Errorf("insufficient cpu values")
	}

	user := values[0]
	nice := values[1]
	system := values[2]
	idle := values[3]

	var iowait, irq, softirq, steal uint64
	if len(values) > 4 {
		iowait = values[4]
	}
	if len(values) > 5 {
		irq = values[5]
	}
	if len(values) > 6 {
		softirq = values[6]
	}
	if len(values) > 7 {
		steal = values[7]
	}

	idleTotal := idle + iowait
	nonIdleTotal := user + nice + system + irq + softirq + steal
	total := idleTotal + nonIdleTotal

	if c.prevTotal == 0 {
		c.prevIdle = idleTotal
		c.prevTotal = total
		return model.CPUMetric{Timestamp: time.Now(), UsagePct: 0.0}, nil
	}

	totalDelta := float64(total - c.prevTotal)
	idleDelta := float64(idleTotal - c.prevIdle)

	c.prevIdle = idleTotal
	c.prevTotal = total

	usagePct := 0.0
	if totalDelta > 0 {
		usagePct = ((totalDelta - idleDelta) / totalDelta) * 100.0
	}
	if usagePct < 0 {
		usagePct = 0
	}
	if usagePct > 100 {
		usagePct = 100
	}

	c.samples++
	c.sumPct += usagePct
	avgPct := c.sumPct / float64(c.samples)

	if usagePct > c.peakPct {
		c.peakPct = usagePct
	}

	return model.CPUMetric{
		Timestamp:  time.Now(),
		UsagePct:   usagePct,
		AveragePct: avgPct,
		PeakPct:    c.peakPct,
	}, nil
}
