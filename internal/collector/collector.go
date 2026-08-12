package collector

import (
	"context"
	"time"
)

// Collector defines the lifecycle interface for all system data collectors.
type Collector interface {
	Name() string
	Collect(ctx context.Context) error
	Snapshot() any
	ResetSession()
	Interval() time.Duration
}
