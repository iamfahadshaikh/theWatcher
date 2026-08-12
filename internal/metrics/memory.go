package metrics

import "time"

type MemoryMetric struct {
	Timestamp      time.Time
	TotalMB        float64
	UsedMB         float64
	FreeMB         float64
	CachedMB       float64
	BufferMB       float64
	UsedPct        float64
	SessionMin     float64
	SessionMax     float64
	SessionAvg     float64
	SessionSamples int
	SwapTotalMB    float64
	SwapUsedMB     float64
	SwapFreeMB     float64
	SwapCachedMB   float64
	DirtyMB        float64
	WritebackMB    float64
	AvailableMB    float64
	MemPressure    string
	OOMKills       int
	HugePagesTotal int
	HugePagesFree  int
	THPStatus      string
}
