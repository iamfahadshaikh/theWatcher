package metrics

type PortInfo struct {
	Protocol string
	Port     string
	Process  string
}

type RemoteIPStat struct {
	IP    string
	Count int
}

type DestPortStat struct {
	Port  int
	Count int
}

type ConnSummary struct {
	Established int
	Listen      int
	TimeWait    int
	CloseWait   int
	SynSent     int
	SynRecv     int
	FinWait     int
	UDP         int
	UNIX        int
	RAW         int
	ConnRate    float64
	TopIPs      []RemoteIPStat
	TopPorts    []DestPortStat
}

type LoginEvent struct {
	Time string
	User string
	Type string
	From string
}
