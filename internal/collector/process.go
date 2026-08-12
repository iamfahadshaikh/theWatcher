package collector

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"theWatcher/internal/metrics"
)

type ProcessSnapshot struct {
	Processes []metrics.ProcessInfo
	BirthRate float64
	Running   int
	Sleeping  int
	Zombie    int
	Stopped   int
}

type procIO struct {
	ReadBytes  uint64
	WriteBytes uint64
}

type ProcessCollector struct {
	mu        sync.Mutex
	procDir   string
	seenPIDs  map[int]bool
	prevTicks map[int]uint64
	prevIO    map[int]procIO
	userCache map[string]string
	prevTime  time.Time
	birthRate float64
	snapshot  ProcessSnapshot
}

func NewProcessCollector() *ProcessCollector {
	return &ProcessCollector{
		procDir:   "/proc",
		seenPIDs:  make(map[int]bool),
		prevTicks: make(map[int]uint64),
		prevIO:    make(map[int]procIO),
		userCache: make(map[string]string),
		prevTime:  time.Now(),
	}
}

func NewProcessCollectorWithPath(procDir string) *ProcessCollector {
	return &ProcessCollector{
		procDir:   procDir,
		seenPIDs:  make(map[int]bool),
		prevTicks: make(map[int]uint64),
		prevIO:    make(map[int]procIO),
		userCache: make(map[string]string),
		prevTime:  time.Now(),
	}
}

func (c *ProcessCollector) Name() string {
	return "Process"
}

func (c *ProcessCollector) ResetSession() {
	c.mu.Lock()
	defer c.mu.Unlock()
}

func (c *ProcessCollector) Interval() time.Duration {
	return 500 * time.Millisecond
}

func (c *ProcessCollector) Collect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	entries, err := os.ReadDir(c.procDir)
	if err != nil {
		return fmt.Errorf("failed reading proc dir: %w", err)
	}

	var processes []metrics.ProcessInfo
	currentPIDs := make(map[int]bool)
	newPIDsThisInterval := 0
	var runCount, sleepCount, zombieCount, stopCount int

	now := time.Now()
	elapsed := now.Sub(c.prevTime).Seconds()

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		currentPIDs[pid] = true
		if len(c.seenPIDs) > 0 && !c.seenPIDs[pid] {
			newPIDsThisInterval++
		}

		pInfo := c.parseProcess(pid, elapsed)
		processes = append(processes, pInfo)

		switch pInfo.State {
		case "R":
			runCount++
		case "S", "D":
			sleepCount++
		case "Z":
			zombieCount++
		case "T", "t":
			stopCount++
		default:
			sleepCount++
		}
	}

	if elapsed > 0 {
		c.birthRate = float64(newPIDsThisInterval) / elapsed
	}

	// Clean up old PIDs from trackers
	for pid := range c.prevTicks {
		if !currentPIDs[pid] {
			delete(c.prevTicks, pid)
			delete(c.prevIO, pid)
		}
	}

	c.seenPIDs = currentPIDs
	c.prevTime = now

	// Default sort by CPU descending
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].CPU > processes[j].CPU
	})

	c.snapshot = ProcessSnapshot{
		Processes: processes,
		BirthRate: c.birthRate,
		Running:   runCount,
		Sleeping:  sleepCount,
		Zombie:    zombieCount,
		Stopped:   stopCount,
	}
	return nil
}

func (c *ProcessCollector) Snapshot() any {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.snapshot
}

func (c *ProcessCollector) parseProcess(pid int, elapsed float64) metrics.ProcessInfo {
	info := metrics.ProcessInfo{
		PID:      pid,
		PPID:     0,
		User:     "root",
		TTY:      "?",
		State:    "S",
		Nice:     0,
		Priority: 20,
		Threads:  1,
		CPU:      0.0,
		Mem:      0.0,
		Command:  fmt.Sprintf("[%d]", pid),
		Runtime:  "0:00",
	}

	pidDir := filepath.Join(c.procDir, strconv.Itoa(pid))

	// Owner Username Lookup with caching
	if fileInfo, err := os.Stat(pidDir); err == nil {
		if stat, ok := fileInfo.Sys().(*syscall.Stat_t); ok {
			uidStr := strconv.Itoa(int(stat.Uid))
			if username, cached := c.userCache[uidStr]; cached {
				info.User = username
			} else {
				if u, err := user.LookupId(uidStr); err == nil {
					c.userCache[uidStr] = u.Username
					info.User = u.Username
				} else {
					c.userCache[uidStr] = uidStr
					info.User = uidStr
				}
			}
		}
	}

	// Executable Path
	if link, err := os.Readlink(filepath.Join(pidDir, "exe")); err == nil {
		info.Executable = link
	}

	// Open FDs Count
	if fds, err := os.ReadDir(filepath.Join(pidDir, "fd")); err == nil {
		info.OpenFiles = len(fds)
	}

	// Command line
	if cmdBytes, err := os.ReadFile(filepath.Join(pidDir, "cmdline")); err == nil && len(cmdBytes) > 0 {
		cmdStr := strings.ReplaceAll(string(cmdBytes), "\x00", " ")
		cmdStr = strings.TrimSpace(cmdStr)
		if cmdStr != "" {
			info.Command = cmdStr
		}
	}

	// Stat details (CPU, state, PPID, priority, nice, threads)
	var totalTicks uint64
	if statBytes, err := os.ReadFile(filepath.Join(pidDir, "stat")); err == nil {
		sStr := string(statBytes)
		rIdx := strings.LastIndex(sStr, ")")
		if rIdx > 0 && rIdx+2 < len(sStr) {
			fields := strings.Fields(sStr[rIdx+2:])
			if len(fields) >= 18 {
				info.State = fields[0]
				info.PPID, _ = strconv.Atoi(fields[1])
				info.Priority, _ = strconv.Atoi(fields[15])
				info.Nice, _ = strconv.Atoi(fields[16])
				info.Threads, _ = strconv.Atoi(fields[17])

				utime, _ := strconv.ParseUint(fields[11], 10, 64)
				stime, _ := strconv.ParseUint(fields[12], 10, 64)
				totalTicks = utime + stime

				// TTY Name
				ttyNum, _ := strconv.Atoi(fields[4])
				if ttyNum > 0 {
					info.TTY = fmt.Sprintf("tty%d", ttyNum&0xff)
				}

				// Runtime start time
				startTimeTicks, _ := strconv.ParseUint(fields[19], 10, 64)
				if uptimeBytes, err := os.ReadFile("/proc/uptime"); err == nil {
					var uptime float64
					_, _ = fmt.Sscanf(string(uptimeBytes), "%f", &uptime)
					// Ticks since boot
					processAgeSec := uptime - (float64(startTimeTicks) / 100.0)
					if processAgeSec > 0 {
						h := int(processAgeSec) / 3600
						m := (int(processAgeSec) % 3600) / 60
						s := int(processAgeSec) % 60
						if h > 0 {
							info.Runtime = fmt.Sprintf("%d:%02d:%02d", h, m, s)
						} else {
							info.Runtime = fmt.Sprintf("%d:%02d", m, s)
						}
					}
				}
			}
		}
	}

	// CPU Rate calculation
	if elapsed > 0 {
		prev, ok := c.prevTicks[pid]
		if ok && totalTicks >= prev {
			ticksDelta := float64(totalTicks - prev)
			cpuSec := ticksDelta / 100.0 // assume 100 Hz ticks
			info.CPU = (cpuSec / elapsed) * 100.0 / float64(runtime.NumCPU())
		}
	}
	c.prevTicks[pid] = totalTicks

	// Read RSS / VmSize
	if statusFile, err := os.Open(filepath.Join(pidDir, "status")); err == nil {
		scanner := bufio.NewScanner(statusFile)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "VmRSS:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					rssKB, _ := strconv.ParseFloat(fields[1], 64)
					info.Mem = rssKB / 1024.0
					info.ResMemBytes = uint64(rssKB * 1024)
				}
			} else if strings.HasPrefix(line, "VmSize:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					sizeKB, _ := strconv.ParseFloat(fields[1], 64)
					info.VirtMemBytes = uint64(sizeKB * 1024)
				}
			} else if strings.HasPrefix(line, "Name:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 && info.Command == fmt.Sprintf("[%d]", pid) {
					info.Command = fields[1]
				}
			}
		}
		statusFile.Close()
	}

	// Read I/O Rate
	var curIO procIO
	if ioFile, err := os.Open(filepath.Join(pidDir, "io")); err == nil {
		scanner := bufio.NewScanner(ioFile)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "read_bytes:") {
				fields := strings.Fields(line)
				if len(fields) == 2 {
					curIO.ReadBytes, _ = strconv.ParseUint(fields[1], 10, 64)
				}
			} else if strings.HasPrefix(line, "write_bytes:") {
				fields := strings.Fields(line)
				if len(fields) == 2 {
					curIO.WriteBytes, _ = strconv.ParseUint(fields[1], 10, 64)
				}
			}
		}
		ioFile.Close()
	}

	if elapsed > 0 {
		prev, ok := c.prevIO[pid]
		if ok {
			if curIO.ReadBytes >= prev.ReadBytes {
				info.ReadBytes = uint64(float64(curIO.ReadBytes-prev.ReadBytes) / elapsed)
			}
			if curIO.WriteBytes >= prev.WriteBytes {
				info.WriteBytes = uint64(float64(curIO.WriteBytes-prev.WriteBytes) / elapsed)
			}
		}
	}
	c.prevIO[pid] = curIO

	// Parent Container ID/Name from Cgroup
	if cgBytes, err := os.ReadFile(filepath.Join(pidDir, "cgroup")); err == nil {
		cgStr := string(cgBytes)
		if idx := strings.Index(cgStr, "/docker/"); idx != -1 {
			id := cgStr[idx+8:]
			if len(id) > 12 {
				info.ContainerName = id[:12]
			}
		} else if idx := strings.Index(cgStr, "docker-"); idx != -1 {
			id := cgStr[idx+7:]
			if len(id) > 12 {
				info.ContainerName = id[:12]
			}
		}
	}

	return info
}
