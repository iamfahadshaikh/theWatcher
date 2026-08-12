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

type ifacePrev struct {
	RxBytes uint64
	TxBytes uint64
}

type NetworkCollector struct {
	mu         sync.Mutex
	devPath    string
	prevTime   time.Time
	prevRx     uint64
	prevTx     uint64
	prevIfaces map[string]ifacePrev
	hasSample  bool
	snapshot   metrics.NetworkMetric
}

func NewNetworkCollector() *NetworkCollector {
	return &NetworkCollector{
		devPath:    "/proc/net/dev",
		prevIfaces: make(map[string]ifacePrev),
	}
}

func NewNetworkCollectorWithPath(path string) *NetworkCollector {
	return &NetworkCollector{
		devPath:    path,
		prevIfaces: make(map[string]ifacePrev),
	}
}

func (c *NetworkCollector) Name() string {
	return "Network"
}

func (c *NetworkCollector) ResetSession() {
	c.mu.Lock()
	defer c.mu.Unlock()
}

func (c *NetworkCollector) Interval() time.Duration {
	return 1 * time.Second
}

func (c *NetworkCollector) Collect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Open(c.devPath)
	if err != nil {
		return fmt.Errorf("failed to open net dev file: %w", err)
	}
	defer file.Close()

	var totalRxBytes, totalTxBytes uint64
	var totalRxPkts, totalTxPkts uint64
	var totalRxErrs, totalTxErrs uint64
	var totalRxDrop, totalTxDrop uint64

	now := time.Now()
	elapsedSec := now.Sub(c.prevTime).Seconds()
	var interfaces []metrics.InterfaceStat

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
			continue
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

		// Compute rates per interface
		ifaceRxRate := 0.0
		ifaceTxRate := 0.0
		if c.hasSample && elapsedSec > 0 {
			prev, ok := c.prevIfaces[iface]
			if ok {
				if rxB >= prev.RxBytes {
					ifaceRxRate = float64(rxB-prev.RxBytes) / 1024.0 / elapsedSec
				}
				if txB >= prev.TxBytes {
					ifaceTxRate = float64(txB-prev.TxBytes) / 1024.0 / elapsedSec
				}
			}
		}
		c.prevIfaces[iface] = ifacePrev{RxBytes: rxB, TxBytes: txB}

		interfaces = append(interfaces, metrics.InterfaceStat{
			Name:      iface,
			RxBytes:   rxB,
			TxBytes:   txB,
			RxPackets: rxP,
			TxPackets: txP,
			RxRateKB:  ifaceRxRate,
			TxRateKB:  ifaceTxRate,
			RxErrors:  rxE,
			TxErrors:  txE,
			RxDrop:    rxD,
			TxDrop:    txD,
			Carrier:   readIfaceSysInt(iface, "carrier"),
			MTU:       readIfaceSysInt(iface, "mtu"),
			Speed:     readIfaceSysInt(iface, "speed"),
			Duplex:    readIfaceSysStr(iface, "duplex"),
		})
	}

	rxRateKB := 0.0
	txRateKB := 0.0
	if c.hasSample && elapsedSec > 0 {
		if totalRxBytes >= c.prevRx {
			rxRateKB = float64(totalRxBytes-c.prevRx) / 1024.0 / elapsedSec
		}
		if totalTxBytes >= c.prevTx {
			txRateKB = float64(totalTxBytes-c.prevTx) / 1024.0 / elapsedSec
		}
	}

	c.prevTime = now
	c.prevRx = totalRxBytes
	c.prevTx = totalTxBytes
	c.hasSample = true

	proto := readProtoCounters()

	c.snapshot = metrics.NetworkMetric{
		Timestamp:  now,
		RxBytes:    totalRxBytes,
		TxBytes:    totalTxBytes,
		RxPackets:  totalRxPkts,
		TxPackets:  totalTxPkts,
		RxErrors:   totalRxErrs,
		TxErrors:   totalTxErrs,
		RxDrop:     totalRxDrop,
		TxDrop:     totalTxDrop,
		RxRateKB:   rxRateKB,
		TxRateKB:   txRateKB,
		Interfaces: interfaces,
		Proto:      proto,
	}
	return nil
}

func (c *NetworkCollector) Snapshot() any {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.snapshot
}

func readIfaceSysInt(iface string, file string) int {
	path := fmt.Sprintf("/sys/class/net/%s/%s", iface, file)
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	var val int
	_, _ = fmt.Sscanf(string(data), "%d", &val)
	return val
}

func readIfaceSysStr(iface string, file string) string {
	path := fmt.Sprintf("/sys/class/net/%s/%s", iface, file)
	data, err := os.ReadFile(path)
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(data))
}

func readProtoCounters() metrics.ProtoCounters {
	var c metrics.ProtoCounters
	if data, err := os.ReadFile("/proc/net/snmp"); err == nil {
		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Fields(line)
			if len(fields) < 2 {
				continue
			}
			if fields[0] == "Ip:" && !strings.Contains(fields[1], "Forwarding") {
				if len(fields) >= 11 {
					inRecv, _ := strconv.ParseUint(fields[3], 10, 64)
					outReq, _ := strconv.ParseUint(fields[10], 10, 64)
					c.IPv4Packets = inRecv + outReq
				}
			}
			if fields[0] == "Icmp:" && !strings.Contains(fields[1], "InMsgs") {
				if len(fields) >= 15 {
					inMsgs, _ := strconv.ParseUint(fields[1], 10, 64)
					outMsgs, _ := strconv.ParseUint(fields[14], 10, 64)
					c.ICMPPackets = inMsgs + outMsgs
				}
			}
			if fields[0] == "Tcp:" && !strings.Contains(fields[1], "RtoAlgorithm") {
				if len(fields) >= 13 {
					inSegs, _ := strconv.ParseUint(fields[10], 10, 64)
					outSegs, _ := strconv.ParseUint(fields[11], 10, 64)
					retrans, _ := strconv.ParseUint(fields[12], 10, 64)
					c.TCPPackets = inSegs + outSegs
					c.Retransmits = retrans
				}
			}
			if fields[0] == "Udp:" && !strings.Contains(fields[1], "InDatagrams") {
				if len(fields) >= 5 {
					inDat, _ := strconv.ParseUint(fields[1], 10, 64)
					outDat, _ := strconv.ParseUint(fields[4], 10, 64)
					c.UDPPackets = inDat + outDat
				}
			}
		}
	}
	if data, err := os.ReadFile("/proc/net/snmp6"); err == nil {
		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		var ip6In, ip6Out uint64
		for scanner.Scan() {
			fields := strings.Fields(scanner.Text())
			if len(fields) == 2 {
				if fields[0] == "Ip6InReceives" {
					ip6In, _ = strconv.ParseUint(fields[1], 10, 64)
				}
				if fields[0] == "Ip6OutRequests" {
					ip6Out, _ = strconv.ParseUint(fields[1], 10, 64)
				}
			}
		}
		c.IPv6Packets = ip6In + ip6Out
	}
	return c
}
