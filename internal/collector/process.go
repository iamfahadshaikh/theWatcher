package collector

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
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

type ProcessCollector struct {
	procDir   string
	seenPIDs  map[int]bool
	prevTime  time.Time
	birthRate float64
	snapshot  ProcessSnapshot
}

func NewProcessCollector() *ProcessCollector {
	return &ProcessCollector{
		procDir:  "/proc",
		seenPIDs: make(map[int]bool),
		prevTime: time.Now(),
	}
}

func NewProcessCollectorWithPath(procDir string) *ProcessCollector {
	return &ProcessCollector{
		procDir:  procDir,
		seenPIDs: make(map[int]bool),
		prevTime: time.Now(),
	}
}

func (c *ProcessCollector) Name() string {
	return "Process"
}

func (c *ProcessCollector) Collect() error {
	entries, err := os.ReadDir(c.procDir)
	if err != nil {
		return fmt.Errorf("failed reading proc dir: %w", err)
	}

	var processes []metrics.ProcessInfo
	currentPIDs := make(map[int]bool)
	newPIDsThisInterval := 0
	var runCount, sleepCount, zombieCount, stopCount int

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

		pInfo := c.parseProcess(pid)
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

	now := time.Now()
	elapsed := now.Sub(c.prevTime).Seconds()
	if elapsed > 0 {
		c.birthRate = float64(newPIDsThisInterval) / elapsed
	}

	c.seenPIDs = currentPIDs
	c.prevTime = now

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
	return c.snapshot
}

func (c *ProcessCollector) parseProcess(pid int) metrics.ProcessInfo {
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

	if cmdBytes, err := os.ReadFile(filepath.Join(pidDir, "cmdline")); err == nil && len(cmdBytes) > 0 {
		cmdStr := strings.ReplaceAll(string(cmdBytes), "\x00", " ")
		cmdStr = strings.TrimSpace(cmdStr)
		if cmdStr != "" {
			info.Command = cmdStr
		}
	}

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
			}
		}
	}

	if statusFile, err := os.Open(filepath.Join(pidDir, "status")); err == nil {
		scanner := bufio.NewScanner(statusFile)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "VmRSS:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					rssKB, _ := strconv.ParseFloat(fields[1], 64)
					info.Mem = rssKB / 1024.0
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

	return info
}
