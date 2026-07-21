package metrics

// PortInfo holds listening socket port details.
type PortInfo struct {
	Protocol string
	Port     string
	Process  string
}

// ConnSummary holds socket state counts.
type ConnSummary struct {
	Established int
	Listen      int
	TimeWait    int
	CloseWait   int
	SynSent     int
	SynRecv     int
}

// LoginEvent holds login audit entries.
type LoginEvent struct {
	Time string
	User string
	Type string
	From string
}
