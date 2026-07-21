package collector

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"theWatcher/internal/metrics"
)

type DockerCmdRunner func() ([]byte, error)

type DockerCollector struct {
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

func defaultDockerRunner() ([]byte, error) {
	cmd := exec.Command("docker", "stats", "--no-stream", "--format", "{{.ID}}\t{{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}\t{{.NetIO}}\t{{.BlockIO}}\t{{.PIDs}}")
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

func (c *DockerCollector) Collect() error {
	output, err := c.cmdRunner()
	if err != nil || len(output) == 0 {
		c.err = fmt.Errorf("docker stats unavailable")
		c.snapshot = nil
		return c.err
	}

	var stats []metrics.DockerStat
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Split(line, "\t")
		if len(fields) < 8 {
			continue
		}

		cpuStr := strings.TrimSuffix(fields[2], "%")
		cpuPct, _ := strconv.ParseFloat(cpuStr, 64)
		pids, _ := strconv.Atoi(fields[7])

		memUsageParts := strings.Split(fields[3], " / ")
		memUse := fields[3]
		memLim := ""
		if len(memUsageParts) == 2 {
			memUse = memUsageParts[0]
			memLim = memUsageParts[1]
		}

		netIOParts := strings.Split(fields[5], " / ")
		netRx := fields[5]
		netTx := ""
		if len(netIOParts) == 2 {
			netRx = netIOParts[0]
			netTx = netIOParts[1]
		}

		blockIOParts := strings.Split(fields[6], " / ")
		blockR := fields[6]
		blockW := ""
		if len(blockIOParts) == 2 {
			blockR = blockIOParts[0]
			blockW = blockIOParts[1]
		}

		stats = append(stats, metrics.DockerStat{
			ID:         fields[0],
			Name:       fields[1],
			CPUPct:     cpuPct,
			MemUsage:   memUse,
			MemLimit:   memLim,
			MemPct:     fields[4],
			NetIO:      fields[5],
			NetRx:      netRx,
			NetTx:      netTx,
			BlockIO:    fields[6],
			BlockRead:  blockR,
			BlockWrite: blockW,
			PIDs:       pids,
		})
	}

	c.snapshot = stats
	c.err = nil
	return nil
}

func (c *DockerCollector) Snapshot() any {
	return c.snapshot
}

func (c *DockerCollector) Err() error {
	return c.err
}
