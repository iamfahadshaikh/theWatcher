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

type DiskCollector struct {
	diskstatsPath string
	prevTime      time.Time
	prevReadSec   uint64
	prevWriteSec  uint64
	prevReadIO    uint64
	prevWriteIO   uint64
	prevIOTicks   uint64
	hasSample     bool
}

func NewDiskCollector() *DiskCollector {
	return &DiskCollector{diskstatsPath: "/proc/diskstats"}
}

func NewDiskCollectorWithPath(path string) *DiskCollector {
	return &DiskCollector{diskstatsPath: path}
}

func (c *DiskCollector) Collect() (model.DiskMetric, error) {
	file, err := os.Open(c.diskstatsPath)
	if err != nil {
		return model.DiskMetric{Timestamp: time.Now()}, fmt.Errorf("failed to open diskstats file: %w", err)
	}
	defer file.Close()

	var totalReadBytes, totalWriteBytes uint64
	var totalReadIOPS, totalWriteIOPS uint64
	var totalIOTicks uint64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 14 {
			continue
		}

		devName := fields[2]
		// Filter physical disks (sd*, vd*, nvme*n1, etc.)
		if strings.HasPrefix(devName, "loop") || strings.HasPrefix(devName, "ram") {
			continue
		}

		rIO, _ := strconv.ParseUint(fields[3], 10, 64)
		rSectors, _ := strconv.ParseUint(fields[5], 10, 64)
		wIO, _ := strconv.ParseUint(fields[7], 10, 64)
		wSectors, _ := strconv.ParseUint(fields[9], 10, 64)
		ioTicks, _ := strconv.ParseUint(fields[12], 10, 64)

		totalReadIOPS += rIO
		totalReadBytes += rSectors * 512
		totalWriteIOPS += wIO
		totalWriteBytes += wSectors * 512
		totalIOTicks += ioTicks
	}

	now := time.Now()
	readRateKB := 0.0
	writeRateKB := 0.0
	readIOPS := uint64(0)
	writeIOPS := uint64(0)
	latencyMs := 0.0

	if c.hasSample {
		elapsedSec := now.Sub(c.prevTime).Seconds()
		if elapsedSec > 0 {
			if totalReadBytes >= c.prevReadSec {
				readRateKB = float64(totalReadBytes-c.prevReadSec) / 1024.0 / elapsedSec
			}
			if totalWriteBytes >= c.prevWriteSec {
				writeRateKB = float64(totalWriteBytes-c.prevWriteSec) / 1024.0 / elapsedSec
			}
			if totalReadIOPS >= c.prevReadIO {
				readIOPS = uint64(float64(totalReadIOPS-c.prevReadIO) / elapsedSec)
			}
			if totalWriteIOPS >= c.prevWriteIO {
				writeIOPS = uint64(float64(totalWriteIOPS-c.prevWriteIO) / elapsedSec)
			}

			totalIOEvents := (totalReadIOPS - c.prevReadIO) + (totalWriteIOPS - c.prevWriteIO)
			if totalIOEvents > 0 && totalIOTicks >= c.prevIOTicks {
				latencyMs = float64(totalIOTicks-c.prevIOTicks) / float64(totalIOEvents)
			}
		}
	}

	c.prevTime = now
	c.prevReadSec = totalReadBytes
	c.prevWriteSec = totalWriteBytes
	c.prevReadIO = totalReadIOPS
	c.prevWriteIO = totalWriteIOPS
	c.prevIOTicks = totalIOTicks
	c.hasSample = true

	return model.DiskMetric{
		Timestamp:   now,
		ReadBytes:   totalReadBytes,
		WriteBytes:  totalWriteBytes,
		ReadIOPS:    readIOPS,
		WriteIOPS:   writeIOPS,
		ReadRateKB:  readRateKB,
		WriteRateKB: writeRateKB,
		LatencyMs:   latencyMs,
	}, nil
}
