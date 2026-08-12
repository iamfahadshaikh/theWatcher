package metrics

import "time"

type CPUSpike struct {
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
	PeakValue float64
}

type CPUMetric struct {
	Timestamp      time.Time
	UsagePct       float64
	AveragePct     float64
	PeakPct        float64
	SessionMin     float64
	SessionMax     float64
	SessionAvg     float64
	SessionSamples int
	LastSpike      CPUSpike
	HasSpikeActive bool
	SpikeStartTime time.Time
	SpikePeakValue float64
	TempC          float64
	PackagePowerW  float64
	FreqAvgMHz     float64
	FreqMinMHz     float64
	FreqMaxMHz     float64
	PerCoreFreqMHz []float64
}
