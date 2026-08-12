package app

import (
	"sync/atomic"
	"time"

	"theWatcher/internal/metrics"
)

// ProcessSnapshot bundles process listing with process birth rate.
type ProcessSnapshot struct {
	Processes []metrics.ProcessInfo
	BirthRate float64
}

// DockerSnapshot bundles Docker stats listing with any collection error.
type DockerSnapshot struct {
	Stats []metrics.DockerStat
	Err   error
}

// MetricsStore provides thread-safe, lock-free access to all collected system metrics.
type MetricsStore struct {
	system      atomic.Value // metrics.SystemInfo
	cpu         atomic.Value // metrics.CPUMetric
	memory      atomic.Value // metrics.MemoryMetric
	network     atomic.Value // metrics.NetworkMetric
	disk        atomic.Value // metrics.DiskMetric
	process     atomic.Value // ProcessSnapshot
	ports       atomic.Value // []metrics.PortInfo
	conn        atomic.Value // metrics.ConnSummary
	docker      atomic.Value // DockerSnapshot
	login       atomic.Value // []metrics.LoginEvent
	collectedAt atomic.Value // time.Time
}

func NewMetricsStore() *MetricsStore {
	s := &MetricsStore{}
	// Initialize with zero values to prevent Load() from returning nil on start
	s.system.Store(metrics.SystemInfo{})
	s.cpu.Store(metrics.CPUMetric{})
	s.memory.Store(metrics.MemoryMetric{})
	s.network.Store(metrics.NetworkMetric{})
	s.disk.Store(metrics.DiskMetric{})
	s.process.Store(ProcessSnapshot{})
	s.ports.Store([]metrics.PortInfo{})
	s.conn.Store(metrics.ConnSummary{})
	s.docker.Store(DockerSnapshot{})
	s.login.Store([]metrics.LoginEvent{})
	s.collectedAt.Store(time.Now())
	return s
}

func (s *MetricsStore) SetSystem(val metrics.SystemInfo) {
	s.system.Store(val)
}

func (s *MetricsStore) GetSystem() metrics.SystemInfo {
	return s.system.Load().(metrics.SystemInfo)
}

func (s *MetricsStore) SetCPU(val metrics.CPUMetric) {
	s.cpu.Store(val)
}

func (s *MetricsStore) GetCPU() metrics.CPUMetric {
	return s.cpu.Load().(metrics.CPUMetric)
}

func (s *MetricsStore) SetMemory(val metrics.MemoryMetric) {
	s.memory.Store(val)
}

func (s *MetricsStore) GetMemory() metrics.MemoryMetric {
	return s.memory.Load().(metrics.MemoryMetric)
}

func (s *MetricsStore) SetNetwork(val metrics.NetworkMetric) {
	s.network.Store(val)
}

func (s *MetricsStore) GetNetwork() metrics.NetworkMetric {
	return s.network.Load().(metrics.NetworkMetric)
}

func (s *MetricsStore) SetDisk(val metrics.DiskMetric) {
	s.disk.Store(val)
}

func (s *MetricsStore) GetDisk() metrics.DiskMetric {
	return s.disk.Load().(metrics.DiskMetric)
}

func (s *MetricsStore) SetProcess(val ProcessSnapshot) {
	s.process.Store(val)
}

func (s *MetricsStore) GetProcess() ProcessSnapshot {
	return s.process.Load().(ProcessSnapshot)
}

func (s *MetricsStore) SetPorts(val []metrics.PortInfo) {
	s.ports.Store(val)
}

func (s *MetricsStore) GetPorts() []metrics.PortInfo {
	return s.ports.Load().([]metrics.PortInfo)
}

func (s *MetricsStore) SetConn(val metrics.ConnSummary) {
	s.conn.Store(val)
}

func (s *MetricsStore) GetConn() metrics.ConnSummary {
	return s.conn.Load().(metrics.ConnSummary)
}

func (s *MetricsStore) SetDocker(val DockerSnapshot) {
	s.docker.Store(val)
}

func (s *MetricsStore) GetDocker() DockerSnapshot {
	return s.docker.Load().(DockerSnapshot)
}

func (s *MetricsStore) SetLogin(val []metrics.LoginEvent) {
	s.login.Store(val)
}

func (s *MetricsStore) GetLogin() []metrics.LoginEvent {
	return s.login.Load().([]metrics.LoginEvent)
}

func (s *MetricsStore) SetCollectedAt(val time.Time) {
	s.collectedAt.Store(val)
}

func (s *MetricsStore) GetCollectedAt() time.Time {
	return s.collectedAt.Load().(time.Time)
}
