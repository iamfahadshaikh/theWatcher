package metrics

import "time"

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
