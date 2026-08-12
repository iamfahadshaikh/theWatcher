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

func FetchMetricsCmd(s *MetricsStore) tea.Cmd {
	return func() tea.Msg {
		dockerSnap := s.GetDocker()
		procSnap := s.GetProcess()
		return MetricsMsg{
			System:      s.GetSystem(),
			CPU:         s.GetCPU(),
			Memory:      s.GetMemory(),
			Network:     s.GetNetwork(),
			Disk:        s.GetDisk(),
			Processes:   procSnap.Processes,
			BirthRate:   procSnap.BirthRate,
			Ports:       s.GetPorts(),
			Conn:        s.GetConn(),
			Docker:      dockerSnap.Stats,
			DockerErr:   dockerSnap.Err,
			Logins:      s.GetLogin(),
			CollectedAt: s.GetCollectedAt(),
		}
	}
}

func (c *Collectors) ResetSession() {
	c.CPU.ResetSession()
	c.Memory.ResetSession()
	c.Disk.ResetSession()
	c.Network.ResetSession()
	c.Process.ResetSession()
	c.Conn.ResetSession()
	c.Docker.ResetSession()
}
