package metrics

import "time"

type FilesystemStat struct {
	Device       string
	MountPoint   string
	FSType       string
	TotalBytes   uint64
	UsedBytes    uint64
	AvailBytes   uint64
	UsedPct      float64
	ReadOnly     bool
	MountOptions string
	InodesTotal  uint64
	InodesUsed   uint64
	InodesPct    float64
}

type DiskMetric struct {
	Timestamp   time.Time
	ReadBytes   uint64
	WriteBytes  uint64
	ReadIOPS    uint64
	WriteIOPS   uint64
	ReadRateKB  float64
	WriteRateKB float64
	LatencyMs   float64
	QueueDepth  int
	DiskUtilPct float64
	DiskTempC   float64
	Filesystems []FilesystemStat
}
