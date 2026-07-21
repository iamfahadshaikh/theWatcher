package metrics

// ProcessInfo holds rich process metrics.
type ProcessInfo struct {
	PID      int
	PPID     int
	User     string
	TTY      string
	State    string
	Nice     int
	Priority int
	Threads  int
	CPU      float64
	Mem      float64
	Command  string
	Runtime  string
}
