package metrics

import "time"

// DiskMetric holds disk throughput and IOPS stats.
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
