package app

import (
	"testing"
	"time"
)

func TestSchedulerGracefulShutdown(t *testing.T) {
	c := NewCollectors()
	store := NewMetricsStore()
	sched := NewScheduler(c, store)

	// Set very fast intervals for testing
	sched.SetInterval("CPU", 10*time.Millisecond)
	sched.SetInterval("Memory", 10*time.Millisecond)

	sched.Start()
	time.Sleep(50 * time.Millisecond)
	sched.Stop()
}
