package config

import "time"

const (
	Version     = "v0.1.0"
	AppName     = "theWatcher"
	RefreshRate = 500 * time.Millisecond
	HistorySize = 300
	MaxEvents   = 200
)
