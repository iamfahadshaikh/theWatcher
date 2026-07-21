package collector

// Collector defines the lifecycle interface for all system data collectors.
type Collector interface {
	Name() string
	Collect() error
	Snapshot() any
}
