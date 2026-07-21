package collector

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"theWatcher/internal/model"
)

type DockerCmdRunner func() ([]byte, error)

type DockerCollector struct {
	cmdRunner DockerCmdRunner
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
	cmd := exec.Command("docker", "stats", "--no-stream", "--format", "{{.ID}}\t{{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func (c *DockerCollector) Collect() ([]model.DockerStat, error) {
	output, err := c.cmdRunner()
	if err != nil || len(output) == 0 {
		return nil, fmt.Errorf("docker stats unavailable")
	}

	var stats []model.DockerStat
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Split(line, "\t")
		if len(fields) < 5 {
			continue
		}

		cpuStr := strings.TrimSuffix(fields[2], "%")
		cpuPct, _ := strconv.ParseFloat(cpuStr, 64)

		stats = append(stats, model.DockerStat{
			ID:       fields[0],
			Name:     fields[1],
			CPUPct:   cpuPct,
			MemUsage: fields[3],
			NetIO:    fields[4],
		})
	}

	return stats, nil
}
