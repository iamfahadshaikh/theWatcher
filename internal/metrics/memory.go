package metrics

import "time"

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
