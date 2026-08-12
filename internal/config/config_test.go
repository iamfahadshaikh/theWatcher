package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigDefaults(t *testing.T) {
	// Temporarily redirect home directory by setting env variable or mocking if needed,
	// but since LoadConfig uses os.UserHomeDir(), we can temporarily backup the config file
	// if it exists, run LoadConfig, and restore it.
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("skipping test; user home dir not available")
	}

	configDir := filepath.Join(home, ".config", "thewatcher")
	configFile := filepath.Join(configDir, "config.yaml")

	existed := false
	if _, err := os.Stat(configFile); err == nil {
		existed = true
		// Rename it temporarily
		err = os.Rename(configFile, configFile+".tmp")
		if err != nil {
			t.Fatalf("failed to rename config for backup: %v", err)
		}
	}

	defer func() {
		if existed {
			_ = os.Remove(configFile)
			_ = os.Rename(configFile+".tmp", configFile)
		} else {
			_ = os.Remove(configFile)
		}
	}()

	// Run LoadConfig with no config file present (should write defaults)
	cfg := LoadConfig()
	if cfg.Graph.Scale != "auto" {
		t.Errorf("expected default graph scale to be 'auto', got '%s'", cfg.Graph.Scale)
	}

	// Now write a custom config file without graph section
	customYAML := []byte(`
refresh_rate: 1s
theme: nord
`)
	err = os.WriteFile(configFile, customYAML, 0644)
	if err != nil {
		t.Fatalf("failed to write custom config: %v", err)
	}

	cfg2 := LoadConfig()
	if cfg2.RefreshRate != "1s" {
		t.Errorf("expected refresh rate to be '1s', got '%s'", cfg2.RefreshRate)
	}
	if cfg2.Graph.Scale != "auto" {
		t.Errorf("expected fallback graph scale to be 'auto' when missing from YAML, got '%s'", cfg2.Graph.Scale)
	}

	// Now write a custom config file with graph scale fixed
	customYAML2 := []byte(`
refresh_rate: 1s
graph:
  scale: fixed
`)
	err = os.WriteFile(configFile, customYAML2, 0644)
	if err != nil {
		t.Fatalf("failed to write custom config: %v", err)
	}

	cfg3 := LoadConfig()
	if cfg3.Graph.Scale != "fixed" {
		t.Errorf("expected graph scale to be 'fixed', got '%s'", cfg3.Graph.Scale)
	}
}
