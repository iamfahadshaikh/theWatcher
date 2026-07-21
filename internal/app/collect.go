package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"theWatcher/internal/collector"
	"theWatcher/internal/metrics"
)

type MetricsMsg struct {
	System      metrics.SystemInfo
	CPU         metrics.CPUMetric
	Memory      metrics.MemoryMetric
	Network     metrics.NetworkMetric
	Disk        metrics.DiskMetric
	Processes   []metrics.ProcessInfo
	BirthRate   float64
	Ports       []metrics.PortInfo
	Conn        metrics.ConnSummary
	Docker      []metrics.DockerStat
	DockerErr   error
	Logins      []metrics.LoginEvent
	CollectedAt time.Time
}

type Collectors struct {
	System  *collector.SystemCollector
	CPU     *collector.CPUCollector
	Memory  *collector.MemoryCollector
	Network *collector.NetworkCollector
	Disk    *collector.DiskCollector
	Process *collector.ProcessCollector
	Ports   *collector.PortsCollector
	Conn    *collector.ConnectionsCollector
	Docker  *collector.DockerCollector
	Login   *collector.LoginCollector
}

func NewCollectors() *Collectors {
	return &Collectors{
		System:  collector.NewSystemCollector(),
		CPU:     collector.NewCPUCollector(),
		Memory:  collector.NewMemoryCollector(),
		Network: collector.NewNetworkCollector(),
		Disk:    collector.NewDiskCollector(),
		Process: collector.NewProcessCollector(),
		Ports:   collector.NewPortsCollector(),
		Conn:    collector.NewConnectionsCollector(),
		Docker:  collector.NewDockerCollector(),
		Login:   collector.NewLoginCollector(),
	}
}

func FetchMetricsCmd(c *Collectors) tea.Cmd {
	return func() tea.Msg {
		msg := MetricsMsg{CollectedAt: time.Now()}

		_ = c.System.Collect()
		msg.System = c.System.Snapshot().(metrics.SystemInfo)

		_ = c.CPU.Collect()
		msg.CPU = c.CPU.Snapshot().(metrics.CPUMetric)

		_ = c.Memory.Collect()
		msg.Memory = c.Memory.Snapshot().(metrics.MemoryMetric)

		_ = c.Network.Collect()
		msg.Network = c.Network.Snapshot().(metrics.NetworkMetric)

		_ = c.Disk.Collect()
		msg.Disk = c.Disk.Snapshot().(metrics.DiskMetric)

		_ = c.Process.Collect()
		procSnap := c.Process.Snapshot().(collector.ProcessSnapshot)
		msg.Processes = procSnap.Processes
		msg.BirthRate = procSnap.BirthRate

		_ = c.Ports.Collect()
		msg.Ports = c.Ports.Snapshot().([]metrics.PortInfo)

		_ = c.Conn.Collect()
		msg.Conn = c.Conn.Snapshot().(metrics.ConnSummary)

		dErr := c.Docker.Collect()
		if dErr == nil {
			msg.Docker = c.Docker.Snapshot().([]metrics.DockerStat)
		} else {
			msg.DockerErr = dErr
		}

		_ = c.Login.Collect()
		msg.Logins = c.Login.Snapshot().([]metrics.LoginEvent)

		return msg
	}
}
