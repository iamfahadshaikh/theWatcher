package config

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type AlertThresholdConfig struct {
	CPU   float64 `yaml:"cpu"`
	Mem   float64 `yaml:"memory"`
	Disk  float64 `yaml:"disk"`
	Temp  float64 `yaml:"temperature"`
	Power float64 `yaml:"power"`
}

// GraphConfig controls graph rendering behavior.
type GraphConfig struct {
	// Scale controls Y-axis scaling mode: "auto" (adaptive) or "fixed" (always 0–100).
	Scale string `yaml:"scale"`
}

type UserConfig struct {
	RefreshRate        string               `yaml:"refresh_rate"`
	Theme              string               `yaml:"theme"`
	HistorySize        int                  `yaml:"history_size"`
	StartupFullscreen  bool                 `yaml:"startup_fullscreen"`
	DockerEnabled      bool                 `yaml:"docker_enabled"`
	TemperatureEnabled bool                 `yaml:"temperature_enabled"`
	NetworkUnits       string               `yaml:"network_units"`
	GraphSmoothing     bool                 `yaml:"graph_smoothing"`
	Graph              GraphConfig          `yaml:"graph"`
	AlertThresholds    AlertThresholdConfig `yaml:"alert_thresholds"`
}

var GlobalConfig UserConfig

func LoadConfig() UserConfig {
	defaults := UserConfig{
		RefreshRate:        "500ms",
		Theme:              "default",
		HistorySize:        300,
		StartupFullscreen:  false,
		DockerEnabled:      true,
		TemperatureEnabled: true,
		NetworkUnits:       "KB",
		GraphSmoothing:     true,
		Graph: GraphConfig{
			Scale: "auto",
		},
		AlertThresholds: AlertThresholdConfig{
			CPU:   80.0,
			Mem:   85.0,
			Disk:  90.0,
			Temp:  80.0,
			Power: 150.0,
		},
	}

	home, err := os.UserHomeDir()
	if err != nil {
		GlobalConfig = defaults
		return defaults
	}

	configDir := filepath.Join(home, ".config", "thewatcher")
	configFile := filepath.Join(configDir, "config.yaml")

	// Ensure config directory exists
	_ = os.MkdirAll(configDir, 0755)

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Save default configuration file
		data, err := yaml.Marshal(defaults)
		if err == nil {
			_ = os.WriteFile(configFile, data, 0644)
		}
		GlobalConfig = defaults
		return defaults
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		GlobalConfig = defaults
		return defaults
	}

	var parsed UserConfig
	if err := yaml.Unmarshal(data, &parsed); err != nil {
		GlobalConfig = defaults
		return defaults
	}

	// Validate / fill missing defaults
	if parsed.RefreshRate == "" {
		parsed.RefreshRate = defaults.RefreshRate
	}
	if parsed.Theme == "" {
		parsed.Theme = defaults.Theme
	}
	if parsed.HistorySize <= 0 {
		parsed.HistorySize = defaults.HistorySize
	}
	if parsed.NetworkUnits == "" {
		parsed.NetworkUnits = defaults.NetworkUnits
	}
	if parsed.AlertThresholds.CPU <= 0 {
		parsed.AlertThresholds.CPU = defaults.AlertThresholds.CPU
	}
	if parsed.AlertThresholds.Mem <= 0 {
		parsed.AlertThresholds.Mem = defaults.AlertThresholds.Mem
	}
	if parsed.AlertThresholds.Disk <= 0 {
		parsed.AlertThresholds.Disk = defaults.AlertThresholds.Disk
	}
	if parsed.AlertThresholds.Temp <= 0 {
		parsed.AlertThresholds.Temp = defaults.AlertThresholds.Temp
	}
	if parsed.AlertThresholds.Power <= 0 {
		parsed.AlertThresholds.Power = defaults.AlertThresholds.Power
	}

	// Validate graph scale option
	if parsed.Graph.Scale != "auto" && parsed.Graph.Scale != "fixed" {
		parsed.Graph.Scale = defaults.Graph.Scale
	}

	GlobalConfig = parsed
	return parsed
}

func GetRefreshPeriod() time.Duration {
	d, err := time.ParseDuration(GlobalConfig.RefreshRate)
	if err != nil {
		return 500 * time.Millisecond
	}
	return d
}
