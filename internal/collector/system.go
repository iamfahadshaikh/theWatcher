package collector

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"theWatcher/internal/metrics"
)

type SystemCollector struct {
	mu       sync.Mutex
	hostname string
	snapshot metrics.SystemInfo
}

func NewSystemCollector() *SystemCollector {
	host, err := os.Hostname()
	if err != nil || host == "" {
		host = "localhost"
	}
	return &SystemCollector{hostname: host}
}

func (c *SystemCollector) Name() string {
	return "System"
}

func (c *SystemCollector) Interval() time.Duration {
	return 30 * time.Second
}

func (c *SystemCollector) ResetSession() {
	c.mu.Lock()
	defer c.mu.Unlock()
}

func (c *SystemCollector) Collect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var kern, upt, load string

	if kBytes, err := os.ReadFile("/proc/sys/kernel/osrelease"); err == nil {
		kern = "Linux " + string(kBytes)
	} else {
		kern = "Linux"
	}

	if lBytes, err := os.ReadFile("/proc/loadavg"); err == nil {
		fields := string(lBytes)
		if len(fields) > 15 {
			fields = fields[:15]
		}
		load = fields
	} else {
		load = "0.00 0.00 0.00"
	}

	if uBytes, err := os.ReadFile("/proc/uptime"); err == nil {
		var uptimeSec float64
		_, _ = fmt.Sscanf(string(uBytes), "%f", &uptimeSec)
		days := int(uptimeSec) / 86400
		hours := (int(uptimeSec) % 86400) / 3600
		mins := (int(uptimeSec) % 3600) / 60
		upt = fmt.Sprintf("%dd %dh %dm", days, hours, mins)
	} else {
		upt = "0m"
	}

	c.snapshot = metrics.SystemInfo{
		Hostname: c.hostname,
		Kernel:   kern,
		Uptime:   upt,
		LoadAvg:  load,
	}
	return nil
}

func (c *SystemCollector) Snapshot() any {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.snapshot
}
