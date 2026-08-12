package collector

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"theWatcher/internal/metrics"
)

type DockerCmdRunner func(context.Context) ([]byte, error)

type inspectInfo struct {
	Status       string
	Health       string
	RestartCount int
	OOMKilled    bool
	ExitCode     int
	Uptime       string
	NetMode      string
	IPs          string
	VolumeCount  int
	Image        string
}

type DockerCollector struct {
	mu        sync.Mutex
	cmdRunner DockerCmdRunner
	snapshot  []metrics.DockerStat
	err       error
}

func NewDockerCollector() *DockerCollector {
	return &DockerCollector{
		cmdRunner: defaultDockerRunner,
	}
}

func NewDockerCollectorWithRunner(runner DockerCmdRunner) *DockerCollector {
	return &DockerCollector{
		cmdRunner: runner,
	}
}

func defaultDockerRunner(ctx context.Context) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "docker", "stats", "--no-stream", "--format", "{{.ID}}\t{{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}\t{{.NetIO}}\t{{.BlockIO}}\t{{.PIDs}}")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func (c *DockerCollector) Name() string {
	return "Docker"
}

func (c *DockerCollector) ResetSession() {
	c.mu.Lock()
	defer c.mu.Unlock()
}

func (c *DockerCollector) Interval() time.Duration {
	return 5 * time.Second
}

func (c *DockerCollector) Collect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	output, err := c.cmdRunner(ctx)
	if err != nil || len(output) == 0 {
		c.err = fmt.Errorf("docker stats unavailable")
		c.snapshot = nil
		return c.err
	}

	var stats []metrics.DockerStat
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) < 8 {
			continue
		}

		cpuStr := strings.TrimSuffix(fields[2], "%")
		cpuVal, _ := strconv.ParseFloat(cpuStr, 64)

		pidsVal, _ := strconv.Atoi(fields[7])

		stat := metrics.DockerStat{
			ID:         fields[0],
			Name:       fields[1],
			CPUPct:     cpuVal,
			MemUsage:   fields[3],
			MemLimit:   "N/A",
			MemPct:     fields[4],
			NetRx:      "N/A",
			NetTx:      "N/A",
			BlockRead:  "N/A",
			BlockWrite: "N/A",
			PIDs:       pidsVal,
		}

		if parts := strings.Split(fields[3], "/"); len(parts) == 2 {
			stat.MemUsage = strings.TrimSpace(parts[0])
			stat.MemLimit = strings.TrimSpace(parts[1])
		}

		if parts := strings.Split(fields[5], "/"); len(parts) == 2 {
			stat.NetRx = strings.TrimSpace(parts[0])
			stat.NetTx = strings.TrimSpace(parts[1])
		}

		if parts := strings.Split(fields[6], "/"); len(parts) == 2 {
			stat.BlockRead = strings.TrimSpace(parts[0])
			stat.BlockWrite = strings.TrimSpace(parts[1])
		}

		stats = append(stats, stat)
	}

	if inspects, err := getDockerInspectInfo(ctx); err == nil && len(inspects) > 0 {
		for i, stat := range stats {
			info, ok := inspects[stat.ID]
			if !ok {
				info, ok = inspects[stat.Name]
			}
			if !ok {
				continue
			}

			stats[i].RestartCount = info.RestartCount
			stats[i].Health = info.Health
			stats[i].Uptime = info.Uptime
			stats[i].Image = info.Image
			stats[i].NetworkMode = info.NetMode
			stats[i].IPAddress = strings.TrimSuffix(info.IPs, ",")
			stats[i].VolumeCount = info.VolumeCount
			stats[i].ExitCode = info.ExitCode
			stats[i].OOMKilled = info.OOMKilled
		}
	}

	c.snapshot = stats
	c.err = nil
	return nil
}

func (c *DockerCollector) Snapshot() any {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err != nil {
		return c.err
	}
	return c.snapshot
}

func getDockerInspectInfo(ctx context.Context) (map[string]inspectInfo, error) {
	cmdPs := exec.CommandContext(ctx, "docker", "ps", "-q")
	var outPs bytes.Buffer
	cmdPs.Stdout = &outPs
	if err := cmdPs.Run(); err != nil {
		return nil, err
	}
	ids := strings.Fields(outPs.String())
	if len(ids) == 0 {
		return nil, nil
	}

	args := []string{"inspect", "--format", "{{.Id}}\t{{.State.Status}}\t{{if .State.Health}}{{.State.Health.Status}}{{else}}none{{end}}\t{{.State.RestartCount}}\t{{.State.OOMKilled}}\t{{.State.ExitCode}}\t{{.Created}}\t{{.HostConfig.NetworkMode}}\t{{range .NetworkSettings.Networks}}{{.IPAddress}},{{end}}\t{{len .Mounts}}\t{{.Config.Image}}"}
	args = append(args, ids...)
	cmdInspect := exec.CommandContext(ctx, "docker", args...)
	var outInspect bytes.Buffer
	cmdInspect.Stdout = &outInspect
	if err := cmdInspect.Run(); err != nil {
		return nil, err
	}

	results := make(map[string]inspectInfo)
	scanner := bufio.NewScanner(&outInspect)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "\t")
		if len(fields) < 11 {
			continue
		}

		id := fields[0]
		if len(id) > 12 {
			id = id[:12]
		}

		restarts, _ := strconv.Atoi(fields[3])
		oom, _ := strconv.ParseBool(fields[4])
		exitVal, _ := strconv.Atoi(fields[5])
		vols, _ := strconv.Atoi(fields[9])

		createdTime := fields[6]
		uptimeStr := "N/A"
		if t, err := time.Parse(time.RFC3339Nano, createdTime); err == nil {
			dur := time.Since(t).Round(time.Second)
			uptimeStr = formatDurationCompact(dur)
		}

		results[id] = inspectInfo{
			Status:       fields[1],
			Health:       fields[2],
			RestartCount: restarts,
			OOMKilled:    oom,
			ExitCode:     exitVal,
			Uptime:       uptimeStr,
			NetMode:      fields[7],
			IPs:          fields[8],
			VolumeCount:  vols,
			Image:        fields[10],
		}
	}

	return results, nil
}

func formatDurationCompact(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	mins := int(d.Minutes()) % 60
	secs := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh", days, hours)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, mins)
	}
	if mins > 0 {
		return fmt.Sprintf("%dm %ds", mins, secs)
	}
	return fmt.Sprintf("%ds", secs)
}
