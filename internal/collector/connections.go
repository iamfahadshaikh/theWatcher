package collector

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"theWatcher/internal/metrics"
)

type ConnectionsCollector struct {
	mu        sync.Mutex
	tcpPaths  []string
	prevOpens uint64
	prevTime  time.Time
	hasSample bool
	snapshot  metrics.ConnSummary
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

func (c *ConnectionsCollector) ResetSession() {
	c.mu.Lock()
	defer c.mu.Unlock()
}

func (c *ConnectionsCollector) Interval() time.Duration {
	return 1 * time.Second
}

func (c *ConnectionsCollector) Collect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var summary metrics.ConnSummary
	ipMap := make(map[string]int)
	portMap := make(map[int]int)

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
			case "04", "05":
				summary.FinWait++
			}

			// Remote address is in fields[2], Local address is fields[1]
			remoteHex := fields[2]
			if ip, port := decodeHexAddr(remoteHex); ip != "" && ip != "0.0.0.0" {
				ipMap[ip]++
				portMap[port]++
			}
		}
		file.Close()
	}

	// Read UDP count
	summary.UDP = countProcLines([]string{"/proc/net/udp", "/proc/net/udp6"})
	summary.UNIX = countProcLines([]string{"/proc/net/unix"})
	summary.RAW = countProcLines([]string{"/proc/net/raw", "/proc/net/raw6"})

	// Connection rate calculation from SNMP
	now := time.Now()
	rate := 0.0
	curOpens := readTCPTotalOpens()
	if c.hasSample {
		elapsed := now.Sub(c.prevTime).Seconds()
		if elapsed > 0 && curOpens >= c.prevOpens {
			rate = float64(curOpens-c.prevOpens) / elapsed
		}
	}
	c.prevOpens = curOpens
	c.prevTime = now
	c.hasSample = true
	summary.ConnRate = rate

	// Sort top remote IPs
	var topIPs []metrics.RemoteIPStat
	for ip, count := range ipMap {
		topIPs = append(topIPs, metrics.RemoteIPStat{IP: ip, Count: count})
	}
	sort.Slice(topIPs, func(i, j int) bool {
		return topIPs[i].Count > topIPs[j].Count
	})
	if len(topIPs) > 5 {
		topIPs = topIPs[:5]
	}
	summary.TopIPs = topIPs

	// Sort top dest ports
	var topPorts []metrics.DestPortStat
	for port, count := range portMap {
		topPorts = append(topPorts, metrics.DestPortStat{Port: port, Count: count})
	}
	sort.Slice(topPorts, func(i, j int) bool {
		return topPorts[i].Count > topPorts[j].Count
	})
	if len(topPorts) > 5 {
		topPorts = topPorts[:5]
	}
	summary.TopPorts = topPorts

	c.snapshot = summary
	return nil
}

func (c *ConnectionsCollector) Snapshot() any {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.snapshot
}

func countProcLines(paths []string) int {
	count := 0
	for _, path := range paths {
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
			count++
		}
		file.Close()
	}
	return count
}

func readTCPTotalOpens() uint64 {
	file, err := os.Open("/proc/net/snmp")
	if err != nil {
		return 0
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Tcp:") && !strings.Contains(line, "RtoAlgorithm") {
			fields := strings.Fields(line)
			if len(fields) >= 7 {
				active, _ := strconv.ParseUint(fields[5], 10, 64)
				passive, _ := strconv.ParseUint(fields[6], 10, 64)
				return active + passive
			}
		}
	}
	return 0
}

func decodeHexAddr(s string) (string, int) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return "", 0
	}
	ipHex := parts[0]
	portHex := parts[1]
	port, _ := strconv.ParseUint(portHex, 16, 16)

	if len(ipHex) == 8 {
		b0, _ := strconv.ParseUint(ipHex[6:8], 16, 8)
		b1, _ := strconv.ParseUint(ipHex[4:6], 16, 8)
		b2, _ := strconv.ParseUint(ipHex[2:4], 16, 8)
		b3, _ := strconv.ParseUint(ipHex[0:2], 16, 8)
		return fmt.Sprintf("%d.%d.%d.%d", b0, b1, b2, b3), int(port)
	}
	return "", int(port)
}
