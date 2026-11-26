package collector

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"theWatcher/internal/metrics"
)

type PortsCollector struct {
	mu       sync.Mutex
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

func (c *PortsCollector) ResetSession() {
	c.mu.Lock()
	defer c.mu.Unlock()
}

func (c *PortsCollector) Interval() time.Duration {
	return 5 * time.Second
}

func (c *PortsCollector) Collect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var ports []metrics.PortInfo

	tcpPorts := parseProcNetSockets(c.tcpPath, "TCP", "0A")
	udpPorts := parseProcNetSockets(c.udpPath, "UDP", "07")

	ports = append(ports, tcpPorts...)
	ports = append(ports, udpPorts...)

	c.snapshot = ports
	return nil
}

func (c *PortsCollector) Snapshot() any {
	c.mu.Lock()
	defer c.mu.Unlock()
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

	inodeMap := buildSocketInodeProcessMap()

	for scanner.Scan() {
		if !headerSkipped {
			headerSkipped = true
			continue
		}

		fields := strings.Fields(scanner.Text())
		if len(fields) < 10 {
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

		inode := fields[9]
		procName := "-"
		if name, ok := inodeMap[inode]; ok {
			procName = name
		} else if os.Geteuid() != 0 {
			procName = "(sudo required)"
		}

		results = append(results, metrics.PortInfo{
			Protocol: proto,
			Port:     strconv.FormatInt(portNum, 10),
			Process:  procName,
		})
	}

	return results
}

func buildSocketInodeProcessMap() map[string]string {
	res := make(map[string]string)
	procDirs, err := os.ReadDir("/proc")
	if err != nil {
		return res
	}

	for _, d := range procDirs {
		if !d.IsDir() {
			continue
		}
		pid := d.Name()
		if _, err := strconv.Atoi(pid); err != nil {
			continue
		}

		commBytes, _ := os.ReadFile("/proc/" + pid + "/comm")
		comm := strings.TrimSpace(string(commBytes))
		if comm == "" {
			cmdlineBytes, _ := os.ReadFile("/proc/" + pid + "/cmdline")
			parts := strings.Split(string(cmdlineBytes), "\x00")
			if len(parts) > 0 && parts[0] != "" {
				comm = parts[0]
				if idx := strings.LastIndex(comm, "/"); idx >= 0 {
					comm = comm[idx+1:]
				}
			}
		}

		fdDir := "/proc/" + pid + "/fd"
		fds, err := os.ReadDir(fdDir)
		if err != nil {
			continue
		}

		for _, fd := range fds {
			link, err := os.Readlink(fdDir + "/" + fd.Name())
			if err != nil {
				continue
			}
			if strings.HasPrefix(link, "socket:[") && strings.HasSuffix(link, "]") {
				inode := link[8 : len(link)-1]
				res[inode] = fmt.Sprintf("%s (%s)", comm, pid)
			}
		}
	}
	return res
}
