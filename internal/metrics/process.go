package metrics

type ProcessInfo struct {
	PID           int
	PPID          int
	User          string
	TTY           string
	State         string
	Nice          int
	Priority      int
	Threads       int
	CPU           float64
	Mem           float64
	Command       string
	Runtime       string
	VirtMemBytes  uint64
	ResMemBytes   uint64
	Executable    string
	OpenFiles     int
	ReadBytes     uint64
	WriteBytes    uint64
	NetworkBytes  uint64
	ContainerName string
}
