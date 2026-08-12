package metrics

type DockerStat struct {
	ID           string
	Name         string
	CPUPct       float64
	MemUsage     string
	MemLimit     string
	MemPct       string
	NetIO        string
	NetRx        string
	NetTx        string
	BlockIO      string
	BlockRead    string
	BlockWrite   string
	PIDs         int
	Uptime       string
	RestartCount int
	Health       string
	Image        string
	NetworkMode  string
	IPAddress    string
	VolumeCount  int
	Ports        string
	ExitCode     int
	OOMKilled    bool
	LogsCount    int
}
