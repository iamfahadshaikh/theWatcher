package collector

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"theWatcher/internal/model"
)

type PortsCollector struct {
	tcpPath string
	udpPath string
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

func (c *PortsCollector) Collect() ([]model.PortInfo, error) {
	var ports []model.PortInfo

	tcpPorts := parseProcNetSockets(c.tcpPath, "TCP", "0A") // 0A is LISTEN state in hex
	udpPorts := parseProcNetSockets(c.udpPath, "UDP", "07") // 07 is CLOSE/UNCONN for UDP

	ports = append(ports, tcpPorts...)
	ports = append(ports, udpPorts...)

	return ports, nil
}

func parseProcNetSockets(path, proto, stateFilter string) []model.PortInfo {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	var results []model.PortInfo
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

		results = append(results, model.PortInfo{
			Protocol: proto,
			Port:     strconv.FormatInt(portNum, 10),
			Process:  "-",
		})
	}

	return results
}
