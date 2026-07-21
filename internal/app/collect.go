package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"theWatcher/internal/collector"
	"theWatcher/internal/model"
)

type MetricsMsg struct {
	CPU         model.CPUMetric
	Memory      model.MemoryMetric
	Network     model.NetworkMetric
	Disk        model.DiskMetric
	Processes   []model.ProcessInfo
	BirthRate   float64
	Ports       []model.PortInfo
	Conn        model.ConnSummary
	Docker      []model.DockerStat
	DockerErr   error
	Logins      []model.LoginEvent
	CollectedAt time.Time
}

type Collectors struct {
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

		if cpuM, err := c.CPU.Collect(); err == nil {
			msg.CPU = cpuM
		}
		if memM, err := c.Memory.Collect(); err == nil {
			msg.Memory = memM
		}
		if netM, err := c.Network.Collect(); err == nil {
			msg.Network = netM
		}
		if diskM, err := c.Disk.Collect(); err == nil {
			msg.Disk = diskM
		}
		if procs, birth, err := c.Process.Collect(); err == nil {
			msg.Processes = procs
			msg.BirthRate = birth
		}
		if ports, err := c.Ports.Collect(); err == nil {
			msg.Ports = ports
		}
		if conn, err := c.Conn.Collect(); err == nil {
			msg.Conn = conn
		}
		dStats, dErr := c.Docker.Collect()
		msg.Docker = dStats
		msg.DockerErr = dErr

		if logins, err := c.Login.Collect(); err == nil {
			msg.Logins = logins
		}

		return msg
	}
}
