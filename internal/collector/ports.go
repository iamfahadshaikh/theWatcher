package collector

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"theWatcher/internal/metrics"
)

type PortsCollector struct {
	tcpPath  string
	udpPath  string
	snapshot []metrics.PortInfo
}

func NewPortsCollector() *PortsCollector {
	return &PortsCollector{
		tcpPath: "/proc/net/tcp",
		udpPath: "/proc/net/udp",
	}
}

func NewPortsCollectorWithPaths(tcpPath, udpPath string) *PortsCollector {
	return &PortsCollector{
		tcpPath: tcpPath,
		udpPath: udpPath,
	}
}

func (c *PortsCollector) Name() string {
	return "Ports"
}

func (c *PortsCollector) Collect() error {
	var ports []metrics.PortInfo

	tcpPorts := parseProcNetSockets(c.tcpPath, "TCP", "0A")
	udpPorts := parseProcNetSockets(c.udpPath, "UDP", "07")

	ports = append(ports, tcpPorts...)
	ports = append(ports, udpPorts...)

	c.snapshot = ports
	return nil
}

func (c *PortsCollector) Snapshot() any {
	return c.snapshot
}

func parseProcNetSockets(path, proto, stateFilter string) []metrics.PortInfo {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	var results []metrics.PortInfo
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

		localAddr := fields[1]
		state := fields[3]

		if stateFilter != "" && state != stateFilter && proto == "TCP" {
			continue
		}

		addrParts := strings.Split(localAddr, ":")
		if len(addrParts) < 2 {
			continue
		}

		portHex := addrParts[1]
		portNum, err := strconv.ParseInt(portHex, 16, 64)
		if err != nil {
			continue
		}

		results = append(results, metrics.PortInfo{
			Protocol: proto,
			Port:     strconv.FormatInt(portNum, 10),
			Process:  "-",
		})
	}

	return results
}
