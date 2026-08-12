package metrics

import "time"

type InterfaceStat struct {
	Name      string
	RxBytes   uint64
	TxBytes   uint64
	RxPackets uint64
	TxPackets uint64
	RxRateKB  float64
	TxRateKB  float64
	RxErrors  uint64
	TxErrors  uint64
	RxDrop    uint64
	TxDrop    uint64
	Carrier   int
	MTU       int
	Speed     int
	Duplex    string
}

type ProtoCounters struct {
	TCPPackets  uint64
	UDPPackets  uint64
	ICMPPackets uint64
	IPv4Packets uint64
	IPv6Packets uint64
	Retransmits uint64
}

type NetworkMetric struct {
	Timestamp  time.Time
	RxBytes    uint64
	TxBytes    uint64
	RxPackets  uint64
	TxPackets  uint64
	RxErrors   uint64
	TxErrors   uint64
	RxDrop     uint64
	TxDrop     uint64
	RxRateKB   float64
	TxRateKB   float64
	Interfaces []InterfaceStat
	Proto      ProtoCounters
}
