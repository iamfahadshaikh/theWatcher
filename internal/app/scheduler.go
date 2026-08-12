package app

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"theWatcher/internal/collector"
	"theWatcher/internal/metrics"
)

type SchedulerMetrics struct {
	LastRun             time.Time
	LastSuccess         time.Time
	LastError           time.Time
	Duration            time.Duration
	TotalDuration       time.Duration
	RunCount            int64
	ErrorCount          int64
	SkipCount           int64
	ConsecutiveFailures int64
}

type Job struct {
	Collector collector.Collector
	Running   int32
}

type Scheduler struct {
	jobs              []*Job
	store             *MetricsStore
	ctx               context.Context
	cancel            context.CancelFunc
	wg                sync.WaitGroup
	metricsMu         sync.Mutex
	metrics           map[string]*SchedulerMetrics
	customIntervals   map[string]time.Duration
	customIntervalsMu sync.RWMutex
}

func NewScheduler(c *Collectors, s *MetricsStore) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	sched := &Scheduler{
		store:           s,
		ctx:             ctx,
		cancel:          cancel,
		metrics:         make(map[string]*SchedulerMetrics),
		customIntervals: make(map[string]time.Duration),
	}

	// Order collectors to respect dependencies:
	// CPU, Memory, Process must run before Docker (since Docker inspect might use process names/pids if needed)
	sched.jobs = []*Job{
		{Collector: c.System},
		{Collector: c.CPU},
		{Collector: c.Memory},
		{Collector: c.Process},
		{Collector: c.Network},
		{Collector: c.Disk},
		{Collector: c.Ports},
		{Collector: c.Conn},
		{Collector: c.Docker},
		{Collector: c.Login},
	}

	for _, j := range sched.jobs {
		sched.metrics[j.Collector.Name()] = &SchedulerMetrics{}
	}

	return sched
}

func (s *Scheduler) Start() {
	// 1. Concurrent Warmup
	var warmupWg sync.WaitGroup
	for _, j := range s.jobs {
		warmupWg.Add(1)
		go func(job *Job) {
			defer warmupWg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = s.executeJobSync(ctx, job)
		}(j)
	}
	warmupWg.Wait()

	// 2. Start Ticker Goroutines for each job
	for _, j := range s.jobs {
		s.wg.Add(1)
		go s.runJobLoop(j)
	}
}

func (s *Scheduler) Stop() {
	s.cancel()
	s.wg.Wait()
}

func (s *Scheduler) runJobLoop(job *Job) {
	defer s.wg.Done()

	for {
		interval := s.getInterval(job.Collector)
		select {
		case <-s.ctx.Done():
			return
		case <-time.After(interval):
			s.runJobAsync(job)
		}
	}
}

func (s *Scheduler) getInterval(c collector.Collector) time.Duration {
	s.customIntervalsMu.RLock()
	defer s.customIntervalsMu.RUnlock()
	if val, ok := s.customIntervals[c.Name()]; ok {
		return val
	}
	return c.Interval()
}

func (s *Scheduler) SetInterval(name string, d time.Duration) {
	s.customIntervalsMu.Lock()
	defer s.customIntervalsMu.Unlock()
	s.customIntervals[name] = d
}

func (s *Scheduler) runJobAsync(job *Job) {
	// TryLock to avoid overlapping executions
	if !atomic.CompareAndSwapInt32(&job.Running, 0, 1) {
		s.metricsMu.Lock()
		s.metrics[job.Collector.Name()].SkipCount++
		s.metricsMu.Unlock()
		return
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer atomic.StoreInt32(&job.Running, 0)

		// Panic recovery
		defer func() {
			if r := recover(); r != nil {
				s.metricsMu.Lock()
				m := s.metrics[job.Collector.Name()]
				m.LastError = time.Now()
				m.ErrorCount++
				m.ConsecutiveFailures++
				s.metricsMu.Unlock()
			}
		}()

		timeout := s.getInterval(job.Collector)
		if timeout < 2*time.Second {
			timeout = 2 * time.Second
		}

		ctx, cancel := context.WithTimeout(s.ctx, timeout)
		defer cancel()

		start := time.Now()
		err := job.Collector.Collect(ctx)
		duration := time.Since(start)

		s.metricsMu.Lock()
		m := s.metrics[job.Collector.Name()]
		m.LastRun = time.Now()
		m.Duration = duration
		m.TotalDuration += duration
		m.RunCount++

		if err != nil {
			m.LastError = time.Now()
			m.ErrorCount++
			m.ConsecutiveFailures++
		} else {
			m.LastSuccess = time.Now()
			m.ConsecutiveFailures = 0
			s.publishSnapshot(job.Collector)
		}
		s.metricsMu.Unlock()
	}()
}

func (s *Scheduler) executeJobSync(ctx context.Context, job *Job) error {
	defer func() {
		if r := recover(); r != nil {
			s.metricsMu.Lock()
			m := s.metrics[job.Collector.Name()]
			m.ErrorCount++
			s.metricsMu.Unlock()
		}
	}()

	err := job.Collector.Collect(ctx)
	if err == nil {
		s.metricsMu.Lock()
		s.publishSnapshot(job.Collector)
		s.metrics[job.Collector.Name()].LastSuccess = time.Now()
		s.metricsMu.Unlock()
	}
	return err
}

func (s *Scheduler) publishSnapshot(col collector.Collector) {
	snap := col.Snapshot()
	switch col.Name() {
	case "System":
		s.store.SetSystem(snap.(metrics.SystemInfo))
	case "CPU":
		s.store.SetCPU(snap.(metrics.CPUMetric))
	case "Memory":
		s.store.SetMemory(snap.(metrics.MemoryMetric))
	case "Network":
		s.store.SetNetwork(snap.(metrics.NetworkMetric))
	case "Disk":
		s.store.SetDisk(snap.(metrics.DiskMetric))
	case "Process":
		// ProcessSnapshot needs to be constructed
		procSnap := snap.(collector.ProcessSnapshot)
		s.store.SetProcess(ProcessSnapshot{
			Processes: procSnap.Processes,
			BirthRate: procSnap.BirthRate,
		})
	case "Ports":
		s.store.SetPorts(snap.([]metrics.PortInfo))
	case "Conn":
		s.store.SetConn(snap.(metrics.ConnSummary))
	case "Docker":
		// Docker snapshot might have returned error or stats.
		// Wait! Let's check docker.go snapshot implementation.
		// If DockerCollector returns []metrics.DockerStat directly or something else, we check.
		// Usually it's []metrics.DockerStat. We'll handle both.
		switch v := snap.(type) {
		case []metrics.DockerStat:
			s.store.SetDocker(DockerSnapshot{Stats: v, Err: nil})
		case error:
			s.store.SetDocker(DockerSnapshot{Stats: nil, Err: v})
		}
	case "Login":
		s.store.SetLogin(snap.([]metrics.LoginEvent))
	}
	s.store.SetCollectedAt(time.Now())
}

func (s *Scheduler) GetMetrics() map[string]SchedulerMetrics {
	s.metricsMu.Lock()
	defer s.metricsMu.Unlock()

	res := make(map[string]SchedulerMetrics)
	for k, v := range s.metrics {
		res[k] = *v
	}
	return res
}

func (s *Scheduler) ResetSession() {
	for _, j := range s.jobs {
		j.Collector.ResetSession()
	}
}
