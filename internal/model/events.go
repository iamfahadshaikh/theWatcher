package model

import "time"

type AlertLevel string

const (
	LevelInfo    AlertLevel = "INFO"
	LevelWarning AlertLevel = "WARN"
	LevelAlert   AlertLevel = "ALERT"
)

type SystemEvent struct {
	Timestamp time.Time
	Level     AlertLevel
	Component string
	Message   string
}
