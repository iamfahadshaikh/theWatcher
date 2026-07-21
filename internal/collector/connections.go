package collector

import (
	"bufio"
	"os"
	"strings"

	"theWatcher/internal/metrics"
)

type ConnectionsCollector struct {
	tcpPaths []string
	snapshot metrics.ConnSummary
}

func NewConnectionsCollector() *ConnectionsCollector {
	return &ConnectionsCollector{
		tcpPaths: []string{"/proc/net/tcp", "/proc/net/tcp6"},
	}
}

func NewConnectionsCollectorWithPaths(paths []string) *ConnectionsCollector {
	return &ConnectionsCollector{
		tcpPaths: paths,
	}
}

func (c *ConnectionsCollector) Name() string {
	return "Connections"
}

func (c *ConnectionsCollector) Collect() error {
	var summary metrics.ConnSummary

	for _, path := range c.tcpPaths {
		file, err := os.Open(path)
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(file)
		headerSkipped := false

		for scanner.Scan() {
			if !headerSkipped {
				headerSkipped = true
				continue
			}

			fields := strings.Fields(scanner.Text())
			if len(fields) < 4 {
				continue
			}

			stateHex := fields[3]
			switch stateHex {
			case "01":
				summary.Established++
			case "0A":
				summary.Listen++
			case "06":
				summary.TimeWait++
			case "08":
				summary.CloseWait++
			case "02":
				summary.SynSent++
			case "03":
				summary.SynRecv++
			}
		}
		file.Close()
	}

	c.snapshot = summary
	return nil
}

func (c *ConnectionsCollector) Snapshot() any {
	return c.snapshot
}
