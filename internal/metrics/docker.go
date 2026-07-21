package metrics

// DockerStat holds container stats matching docker stats format.
type DockerStat struct {
	ID         string
	Name       string
	CPUPct     float64
	MemUsage   string
	MemLimit   string
	MemPct     string
	NetIO      string
	NetRx      string
	NetTx      string
	BlockIO    string
	BlockRead  string
	BlockWrite string
	PIDs       int
}
