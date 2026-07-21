package model

import "time"

// CPUMetric holds a single CPU sample.
type CPUMetric struct {
	Timestamp  time.Time
	UsagePct   float64
	AveragePct float64
	PeakPct    float64
}

// MemoryMetric holds memory breakdown metrics.
type MemoryMetric struct {
	Timestamp time.Time
	TotalMB   float64
	UsedMB    float64
	FreeMB    float64
	CachedMB  float64
	BufferMB  float64
	UsedPct   float64
}

// NetworkMetric holds network interface stats.
type NetworkMetric struct {
	Timestamp time.Time
	RxBytes   uint64
	TxBytes   uint64
	RxPackets uint64
	TxPackets uint64
	RxErrors  uint64
	TxErrors  uint64
	RxDrop    uint64
	TxDrop    uint64
	RxRateKB  float64
	TxRateKB  float64
}

// DiskMetric holds disk read/write throughput and IOPS stats.
type DiskMetric struct {
	Timestamp   time.Time
	ReadBytes   uint64
	WriteBytes  uint64
	ReadIOPS    uint64
	WriteIOPS   uint64
	ReadRateKB  float64
	WriteRateKB float64
	LatencyMs   float64
}

// ProcessInfo holds process metrics.
type ProcessInfo struct {
	PID     int
	User    string
	CPU     float64
	Mem     float64
	Command string
	Runtime string
}

// PortInfo holds listening port information.
type PortInfo struct {
	Protocol string
	Port     string
	Process  string
}

// ConnSummary holds socket state counts.
type ConnSummary struct {
	Established int
	Listen      int
	TimeWait    int
	CloseWait   int
	SynSent     int
	SynRecv     int
}

// DockerStat holds container stats.
type DockerStat struct {
	ID       string
	Name     string
	CPUPct   float64
	MemUsage string
	NetIO    string
}

// LoginEvent holds system login logs.
type LoginEvent struct {
	Time string
	User string
	Type string
	From string
}
