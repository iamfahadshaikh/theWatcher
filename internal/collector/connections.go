package collector

import (
	"bufio"
	"os"
	"strings"

	"theWatcher/internal/model"
)

type ConnectionsCollector struct {
	tcpPaths []string
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

func (c *ConnectionsCollector) Collect() (model.ConnSummary, error) {
	var summary model.ConnSummary

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

	return summary, nil
}
