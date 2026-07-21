package metrics

import "time"

// CPUMetric holds a single CPU sample.
type CPUMetric struct {
	Timestamp  time.Time
	UsagePct   float64
	AveragePct float64
	PeakPct    float64
}
