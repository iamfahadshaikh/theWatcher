package collector

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"theWatcher/internal/model"
)

type NetworkCollector struct {
	devPath   string
	prevTime  time.Time
	prevRx    uint64
	prevTx    uint64
	hasSample bool
}

func NewNetworkCollector() *NetworkCollector {
	return &NetworkCollector{devPath: "/proc/net/dev"}
}

func NewNetworkCollectorWithPath(path string) *NetworkCollector {
	return &NetworkCollector{devPath: path}
}

func (c *NetworkCollector) Collect() (model.NetworkMetric, error) {
	file, err := os.Open(c.devPath)
	if err != nil {
		return model.NetworkMetric{Timestamp: time.Now()}, fmt.Errorf("failed to open net dev file: %w", err)
	}
	defer file.Close()

	var totalRxBytes, totalTxBytes uint64
	var totalRxPkts, totalTxPkts uint64
	var totalRxErrs, totalTxErrs uint64
	var totalRxDrop, totalTxDrop uint64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, ":") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}

		iface := strings.TrimSpace(parts[0])
		if iface == "lo" {
			continue // Skip loopback
		}

		fields := strings.Fields(parts[1])
		if len(fields) < 16 {
			continue
		}

		rxB, _ := strconv.ParseUint(fields[0], 10, 64)
		rxP, _ := strconv.ParseUint(fields[1], 10, 64)
		rxE, _ := strconv.ParseUint(fields[2], 10, 64)
		rxD, _ := strconv.ParseUint(fields[3], 10, 64)

		txB, _ := strconv.ParseUint(fields[8], 10, 64)
		txP, _ := strconv.ParseUint(fields[9], 10, 64)
		txE, _ := strconv.ParseUint(fields[10], 10, 64)
		txD, _ := strconv.ParseUint(fields[11], 10, 64)

		totalRxBytes += rxB
		totalRxPkts += rxP
		totalRxErrs += rxE
		totalRxDrop += rxD

		totalTxBytes += txB
		totalTxPkts += txP
		totalTxErrs += txE
		totalTxDrop += txD
	}

	now := time.Now()
	rxRateKB := 0.0
	txRateKB := 0.0

	if c.hasSample {
		elapsedSec := now.Sub(c.prevTime).Seconds()
		if elapsedSec > 0 {
			if totalRxBytes >= c.prevRx {
				rxRateKB = float64(totalRxBytes-c.prevRx) / 1024.0 / elapsedSec
			}
			if totalTxBytes >= c.prevTx {
				txRateKB = float64(totalTxBytes-c.prevTx) / 1024.0 / elapsedSec
			}
		}
	}

	c.prevTime = now
	c.prevRx = totalRxBytes
	c.prevTx = totalTxBytes
	c.hasSample = true

	return model.NetworkMetric{
		Timestamp: now,
		RxBytes:   totalRxBytes,
		TxBytes:   totalTxBytes,
		RxPackets: totalRxPkts,
		TxPackets: totalTxPkts,
		RxErrors:  totalRxErrs,
		TxErrors:  totalTxErrs,
		RxDrop:    totalRxDrop,
		TxDrop:    totalTxDrop,
		RxRateKB:  rxRateKB,
		TxRateKB:  txRateKB,
	}, nil
}
