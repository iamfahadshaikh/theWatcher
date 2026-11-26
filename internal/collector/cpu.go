package collector

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"theWatcher/internal/metrics"
)

type CPUCollector struct {
	mu             sync.Mutex
	statPath       string
	prevIdle       uint64
	prevTotal      uint64
	prevCoreIdles  []uint64
	prevCoreTotals []uint64
	sumPct         float64
	samples        uint64
	peakPct        float64
	sessionMin     float64
	sessionMax     float64
	prevEnergyUj   uint64
	prevEnergyTime time.Time
	snapshot       metrics.CPUMetric
}

func NewCPUCollector() *CPUCollector {
	return &CPUCollector{statPath: "/proc/stat"}
}

func NewCPUCollectorWithPath(path string) *CPUCollector {
	return &CPUCollector{statPath: path}
}

func (c *CPUCollector) Name() string {
	return "CPU"
}

func (c *CPUCollector) ResetSession() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.samples = 0
	c.sumPct = 0
	c.peakPct = 0
	c.sessionMin = 0
	c.sessionMax = 0
	c.snapshot.SessionMin = 0
	c.snapshot.SessionMax = 0
	c.snapshot.SessionAvg = 0
	c.snapshot.SessionSamples = 0
	c.snapshot.LastSpike = metrics.CPUSpike{}
	c.snapshot.HasSpikeActive = false
}

func (c *CPUCollector) Interval() time.Duration {
	return 250 * time.Millisecond
}

func (c *CPUCollector) Collect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Open(c.statPath)
	if err != nil {
		return fmt.Errorf("failed to open stat file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return fmt.Errorf("empty stat file")
	}

	line := scanner.Text()
	fields := strings.Fields(line)
	if len(fields) < 5 || fields[0] != "cpu" {
		return fmt.Errorf("invalid cpu line: %s", line)
	}

	var values []uint64
	for _, f := range fields[1:] {
		val, err := strconv.ParseUint(f, 10, 64)
		if err != nil {
			break
		}
		values = append(values, val)
	}

	if len(values) < 4 {
		return fmt.Errorf("insufficient cpu values")
	}

	user, nice, system, idle := values[0], values[1], values[2], values[3]
	var iowait, irq, softirq, steal uint64
	if len(values) > 4 {
		iowait = values[4]
	}
	if len(values) > 5 {
		irq = values[5]
	}
	if len(values) > 6 {
		softirq = values[6]
	}
	if len(values) > 7 {
		steal = values[7]
	}

	idleTotal := idle + iowait
	nonIdleTotal := user + nice + system + irq + softirq + steal
	total := idleTotal + nonIdleTotal

	if c.prevTotal == 0 {
		c.prevIdle = idleTotal
		c.prevTotal = total
		c.snapshot = metrics.CPUMetric{Timestamp: time.Now(), UsagePct: 0.0}
		return nil
	}

	totalDelta := float64(total - c.prevTotal)
	idleDelta := float64(idleTotal - c.prevIdle)

	c.prevIdle = idleTotal
	c.prevTotal = total

	usagePct := 0.0
	if totalDelta > 0 {
		usagePct = ((totalDelta - idleDelta) / totalDelta) * 100.0
	}
	if usagePct < 0 {
		usagePct = 0
	}
	if usagePct > 100 {
		usagePct = 100
	}

	c.samples++
	c.sumPct += usagePct
	avgPct := c.sumPct / float64(c.samples)

	if usagePct > c.peakPct {
		c.peakPct = usagePct
	}

	if c.sessionMin == 0 || usagePct < c.sessionMin {
		c.sessionMin = usagePct
	}
	if usagePct > c.sessionMax {
		c.sessionMax = usagePct
	}

	// CPUTemperature, power, per-core frequencies
	cpuTemp := readCPUTemp()
	raplPower := c.readIntelRaplPower()
	fAvg, fMin, fMax, perCore := readCPUFrequencies()

	// CPU Spike detection (>80%)
	spikeActive := c.snapshot.HasSpikeActive
	spikeStart := c.snapshot.SpikeStartTime
	spikePeak := c.snapshot.SpikePeakValue
	lastSpike := c.snapshot.LastSpike

	if usagePct > 80.0 {
		if !spikeActive {
			spikeActive = true
			spikeStart = time.Now()
			spikePeak = usagePct
		} else {
			if usagePct > spikePeak {
				spikePeak = usagePct
			}
		}
	} else {
		if spikeActive {
			spikeActive = false
			lastSpike = metrics.CPUSpike{
				StartTime: spikeStart,
				EndTime:   time.Now(),
				Duration:  time.Since(spikeStart),
				PeakValue: spikePeak,
			}
		}
	}

	// Parse per-core cpu0, cpu1... lines from /proc/stat
	var perCoreUsage []float64
	var currCoreIdles []uint64
	var currCoreTotals []uint64
	coreIdx := 0

	for scanner.Scan() {
		cLine := scanner.Text()
		cFields := strings.Fields(cLine)
		if len(cFields) < 5 {
			continue
		}
		prefix := fmt.Sprintf("cpu%d", coreIdx)
		if cFields[0] != prefix {
			if strings.HasPrefix(cFields[0], "cpu") && len(cFields[0]) > 3 {
				// Non-sequential cpu line or stopped
			} else {
				break
			}
		}

		var cVals []uint64
		for _, f := range cFields[1:] {
			v, err := strconv.ParseUint(f, 10, 64)
			if err != nil {
				break
			}
			cVals = append(cVals, v)
		}
		if len(cVals) >= 4 {
			cIdle := cVals[3]
			if len(cVals) > 4 {
				cIdle += cVals[4] // + iowait
			}
			cNonIdle := cVals[0] + cVals[1] + cVals[2]
			if len(cVals) > 5 {
				cNonIdle += cVals[5] // + irq
			}
			if len(cVals) > 6 {
				cNonIdle += cVals[6] // + softirq
			}
			if len(cVals) > 7 {
				cNonIdle += cVals[7] // + steal
			}
			cTot := cIdle + cNonIdle
			currCoreIdles = append(currCoreIdles, cIdle)
			currCoreTotals = append(currCoreTotals, cTot)

			cUsage := 0.0
			if len(c.prevCoreTotals) > coreIdx && len(c.prevCoreIdles) > coreIdx {
				cTotDelta := float64(cTot - c.prevCoreTotals[coreIdx])
				cIdleDelta := float64(cIdle - c.prevCoreIdles[coreIdx])
				if cTotDelta > 0 {
					cUsage = ((cTotDelta - cIdleDelta) / cTotDelta) * 100.0
				}
				if cUsage < 0 {
					cUsage = 0
				}
				if cUsage > 100 {
					cUsage = 100
				}
			}
			perCoreUsage = append(perCoreUsage, cUsage)
			coreIdx++
		}
	}
	c.prevCoreIdles = currCoreIdles
	c.prevCoreTotals = currCoreTotals

	c.snapshot = metrics.CPUMetric{
		Timestamp:       time.Now(),
		UsagePct:        usagePct,
		AveragePct:      avgPct,
		PeakPct:         c.peakPct,
		SessionMin:      c.sessionMin,
		SessionMax:      c.sessionMax,
		SessionAvg:      avgPct,
		SessionSamples:  int(c.samples),
		LastSpike:       lastSpike,
		HasSpikeActive:  spikeActive,
		SpikeStartTime:  spikeStart,
		SpikePeakValue:  spikePeak,
		TempC:           cpuTemp,
		PackagePowerW:   raplPower,
		FreqAvgMHz:      fAvg,
		FreqMinMHz:      fMin,
		FreqMaxMHz:      fMax,
		PerCoreFreqMHz:  perCore,
		PerCoreUsagePct: perCoreUsage,
	}

	return nil
}

func (c *CPUCollector) Snapshot() any {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.snapshot
}

func readCPUTemp() float64 {
	for i := 0; i < 10; i++ {
		path := fmt.Sprintf("/sys/class/thermal/thermal_zone%d/temp", i)
		if data, err := os.ReadFile(path); err == nil {
			var t float64
			if _, err := fmt.Sscanf(string(data), "%f", &t); err == nil {
				if t > 1000 {
					return t / 1000.0
				}
				return t
			}
		}
	}
	return 0.0
}

func (c *CPUCollector) readIntelRaplPower() float64 {
	path := "/sys/class/powercap/intel-rapl/intel-rapl:0/energy_uj"
	data, err := os.ReadFile(path)
	if err != nil {
		return 0.0
	}
	var energy uint64
	if _, err := fmt.Sscanf(string(data), "%d", &energy); err != nil {
		return 0.0
	}
	now := time.Now()
	if c.prevEnergyUj == 0 {
		c.prevEnergyUj = energy
		c.prevEnergyTime = now
		return 0.0
	}
	timeDelta := now.Sub(c.prevEnergyTime).Seconds()
	energyDelta := float64(energy - c.prevEnergyUj)
	c.prevEnergyUj = energy
	c.prevEnergyTime = now
	if timeDelta > 0 {
		return (energyDelta / 1000000.0) / timeDelta
	}
	return 0.0
}

func readCPUFrequencies() (avg, min, max float64, perCore []float64) {
	var total float64
	for i := 0; i < 256; i++ {
		path := fmt.Sprintf("/sys/devices/system/cpu/cpu%d/cpufreq/scaling_cur_freq", i)
		if data, err := os.ReadFile(path); err == nil {
			var f float64
			if _, err := fmt.Sscanf(string(data), "%f", &f); err == nil {
				freqMHz := f / 1000.0
				perCore = append(perCore, freqMHz)
				total += freqMHz
				if min == 0 || freqMHz < min {
					min = freqMHz
				}
				if freqMHz > max {
					max = freqMHz
				}
			}
		} else {
			break
		}
	}
	if len(perCore) == 0 {
		if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
			scanner := bufio.NewScanner(strings.NewReader(string(data)))
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "cpu MHz") {
					parts := strings.Split(line, ":")
					if len(parts) == 2 {
						if f, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64); err == nil {
							perCore = append(perCore, f)
							total += f
							if min == 0 || f < min {
								min = f
							}
							if f > max {
								max = f
							}
						}
					}
				}
			}
		}
	}
	if len(perCore) > 0 {
		avg = total / float64(len(perCore))
	}
	return
}
