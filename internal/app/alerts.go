package app

import (
	"fmt"
	"time"

	"theWatcher/internal/config"
	"theWatcher/internal/metrics"
)

type Alert struct {
	Type      string
	Message   string
	Timestamp time.Time
}

type AlertEngine struct {
	ActiveAlerts []Alert
}

func NewAlertEngine() *AlertEngine {
	return &AlertEngine{}
}

func (ae *AlertEngine) Update(
	cpu metrics.CPUMetric,
	mem metrics.MemoryMetric,
	disk metrics.DiskMetric,
	net metrics.NetworkMetric,
	containers []metrics.DockerStat,
	cpuTemp float64,
	cpuPower float64,
) {
	ae.ActiveAlerts = nil
	limits := config.GlobalConfig.AlertThresholds

	// CPU Alert
	if cpu.UsagePct >= limits.CPU {
		ae.ActiveAlerts = append(ae.ActiveAlerts, Alert{
			Type:      "CPU",
			Message:   fmt.Sprintf("CPU usage at %.1f%%", cpu.UsagePct),
			Timestamp: time.Now(),
		})
	}

	// Memory Alert
	if mem.UsedPct >= limits.Mem {
		ae.ActiveAlerts = append(ae.ActiveAlerts, Alert{
			Type:      "MEM",
			Message:   fmt.Sprintf("Mem usage at %.1f%%", mem.UsedPct),
			Timestamp: time.Now(),
		})
	}

	// Disk Alert (if any filesystem exceeds threshold)
	for _, fs := range disk.Filesystems {
		if fs.UsedPct >= limits.Disk {
			ae.ActiveAlerts = append(ae.ActiveAlerts, Alert{
				Type:      "DISK",
				Message:   fmt.Sprintf("FS %s full (%.1f%%)", fs.MountPoint, fs.UsedPct),
				Timestamp: time.Now(),
			})
		}
	}

	// Network Alert (e.g. throughput > threshold, let's say limit network is in MB/s, if RX or TX exceeds limits.Net)
	// We'll default to alert if packet error rates are high, or just let users configure it.
	// Since limits has cpu, mem, disk, temp, power:
	// Temperature Alert
	if cpuTemp >= limits.Temp {
		ae.ActiveAlerts = append(ae.ActiveAlerts, Alert{
			Type:      "TEMP",
			Message:   fmt.Sprintf("CPU Temp at %.1f°C", cpuTemp),
			Timestamp: time.Now(),
		})
	}

	// Power Alert
	if cpuPower >= limits.Power {
		ae.ActiveAlerts = append(ae.ActiveAlerts, Alert{
			Type:      "PWR",
			Message:   fmt.Sprintf("CPU Power at %.1fW", cpuPower),
			Timestamp: time.Now(),
		})
	}

	// Docker health Alerts
	for _, c := range containers {
		if c.Health == "unhealthy" {
			ae.ActiveAlerts = append(ae.ActiveAlerts, Alert{
				Type:      "DOCKER",
				Message:   fmt.Sprintf("Container %s unhealthy", c.Name),
				Timestamp: time.Now(),
			})
		}
	}
}

func (ae *AlertEngine) RenderFooterAlerts() string {
	if len(ae.ActiveAlerts) == 0 {
		return ""
	}
	// Return the most recent active alert message formatted as warning
	alert := ae.ActiveAlerts[len(ae.ActiveAlerts)-1]
	return fmt.Sprintf("  ALERT: [%s] %s  ", alert.Type, alert.Message)
}
