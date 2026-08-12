package collector

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"theWatcher/internal/metrics"
)

type DiskCollector struct {
	mu            sync.Mutex
	diskstatsPath string
	prevTime      time.Time
	prevReadSec   uint64
	prevWriteSec  uint64
	prevReadIO    uint64
	prevWriteIO   uint64
	prevIOTicks   uint64
	hasSample     bool
	snapshot      metrics.DiskMetric
}

func NewDiskCollector() *DiskCollector {
	return &DiskCollector{diskstatsPath: "/proc/diskstats"}
}

func NewDiskCollectorWithPath(path string) *DiskCollector {
	return &DiskCollector{diskstatsPath: path}
}

func (c *DiskCollector) Name() string {
	return "Disk"
}

func (c *DiskCollector) ResetSession() {
	c.mu.Lock()
	defer c.mu.Unlock()
}

func (c *DiskCollector) Interval() time.Duration {
	return 2 * time.Second
}

func (c *DiskCollector) Collect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Open(c.diskstatsPath)
	if err != nil {
		return fmt.Errorf("failed to open diskstats file: %w", err)
	}
	defer file.Close()

	var totalReadBytes, totalWriteBytes uint64
	var totalReadIOPS, totalWriteIOPS uint64
	var totalIOTicks uint64
	var queueDepth int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 14 {
			continue
		}

		devName := fields[2]
		if strings.HasPrefix(devName, "loop") || strings.HasPrefix(devName, "ram") {
			continue
		}

		rIO, _ := strconv.ParseUint(fields[3], 10, 64)
		rSectors, _ := strconv.ParseUint(fields[5], 10, 64)
		wIO, _ := strconv.ParseUint(fields[7], 10, 64)
		wSectors, _ := strconv.ParseUint(fields[9], 10, 64)
		inProgress, _ := strconv.Atoi(fields[11])
		ioTicks, _ := strconv.ParseUint(fields[12], 10, 64)

		totalReadIOPS += rIO
		totalReadBytes += rSectors * 512
		totalWriteIOPS += wIO
		totalWriteBytes += wSectors * 512
		totalIOTicks += ioTicks
		queueDepth += inProgress
	}

	now := time.Now()
	readRateKB := 0.0
	writeRateKB := 0.0
	readIOPS := uint64(0)
	writeIOPS := uint64(0)
	latencyMs := 0.0
	utilPct := 0.0

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

			if totalIOTicks >= c.prevIOTicks {
				utilPct = float64(totalIOTicks-c.prevIOTicks) / (elapsedSec * 1000.0) * 100.0
				if utilPct > 100.0 {
					utilPct = 100.0
				}
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

	filesystems := readFilesystems()
	diskTemp := readDiskTemp()

	c.snapshot = metrics.DiskMetric{
		Timestamp:   now,
		ReadBytes:   totalReadBytes,
		WriteBytes:  totalWriteBytes,
		ReadIOPS:    readIOPS,
		WriteIOPS:   writeIOPS,
		ReadRateKB:  readRateKB,
		WriteRateKB: writeRateKB,
		LatencyMs:   latencyMs,
		QueueDepth:  queueDepth,
		DiskUtilPct: utilPct,
		DiskTempC:   diskTemp,
		Filesystems: filesystems,
	}
	return nil
}

func (c *DiskCollector) Snapshot() any {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.snapshot
}

func isRealFS(fsType string, device string) bool {
	return strings.HasPrefix(device, "/dev/") || fsType == "ext4" || fsType == "xfs" || fsType == "btrfs" || fsType == "zfs" || fsType == "vfat" || fsType == "ntfs" || fsType == "drvfs" || strings.Contains(fsType, "fuse")
}

func readFilesystems() []metrics.FilesystemStat {
	file, err := os.Open("/proc/mounts")
	if err != nil {
		return nil
	}
	defer file.Close()

	var filesystems []metrics.FilesystemStat
	seen := make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			continue
		}
		device := fields[0]
		mountPoint := fields[1]
		fsType := fields[2]
		opts := fields[3]

		if !isRealFS(fsType, device) || seen[mountPoint] {
			continue
		}
		seen[mountPoint] = true

		var stat syscall.Statfs_t
		if err := syscall.Statfs(mountPoint, &stat); err == nil {
			total := stat.Blocks * uint64(stat.Bsize)
			free := stat.Bfree * uint64(stat.Bsize)
			avail := stat.Bavail * uint64(stat.Bsize)
			if total == 0 {
				continue
			}
			used := total - free
			usedPct := (float64(used) / float64(total)) * 100.0

			var inodesTotal, inodesUsed uint64
			var inodesPct float64
			if stat.Files > 0 {
				inodesTotal = stat.Files
				inodesUsed = stat.Files - stat.Ffree
				inodesPct = (float64(inodesUsed) / float64(inodesTotal)) * 100.0
			}

			readOnly := false
			for _, opt := range strings.Split(opts, ",") {
				if opt == "ro" {
					readOnly = true
					break
				}
			}

			filesystems = append(filesystems, metrics.FilesystemStat{
				Device:       device,
				MountPoint:   mountPoint,
				FSType:       fsType,
				TotalBytes:   total,
				UsedBytes:    used,
				AvailBytes:   avail,
				UsedPct:      usedPct,
				ReadOnly:     readOnly,
				MountOptions: opts,
				InodesTotal:  inodesTotal,
				InodesUsed:   inodesUsed,
				InodesPct:    inodesPct,
			})
		}
	}
	return filesystems
}

func readDiskTemp() float64 {
	for i := 0; i < 10; i++ {
		namePath := fmt.Sprintf("/sys/class/hwmon/hwmon%d/name", i)
		if nameData, err := os.ReadFile(namePath); err == nil {
			name := strings.TrimSpace(string(nameData))
			if strings.Contains(name, "nvme") || strings.Contains(name, "drivetemp") {
				tempPath := fmt.Sprintf("/sys/class/hwmon/hwmon%d/temp1_input", i)
				if tempData, err := os.ReadFile(tempPath); err == nil {
					var t float64
					if _, err := fmt.Sscanf(string(tempData), "%f", &t); err == nil {
						return t / 1000.0
					}
				}
			}
		}
	}
	return 0.0
}
