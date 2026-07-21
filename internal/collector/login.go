package collector

import (
	"bytes"
	"os/exec"
	"strings"

	"theWatcher/internal/metrics"
)

type CommandRunner func(name string, args ...string) ([]byte, error)

type LoginCollector struct {
	runner   CommandRunner
	snapshot []metrics.LoginEvent
}

func NewLoginCollector() *LoginCollector {
	return &LoginCollector{
		runner: defaultCommandRunner,
	}
}

func NewLoginCollectorWithRunner(runner CommandRunner) *LoginCollector {
	return &LoginCollector{
		runner: runner,
	}
}

func defaultCommandRunner(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func (c *LoginCollector) Name() string {
	return "Login"
}

func (c *LoginCollector) Collect() error {
	var events []metrics.LoginEvent

	if out, err := c.runner("who"); err == nil && len(out) > 0 {
		lines := strings.Split(string(out), "\n")
		for _, l := range lines {
			fields := strings.Fields(l)
			if len(fields) >= 3 {
				u := fields[0]
				tty := fields[1]
				tStr := strings.Join(fields[2:], " ")
				events = append(events, metrics.LoginEvent{
					Time: tStr,
					User: u,
					Type: "SESSION",
					From: tty,
				})
			}
		}
	}

	if out, err := c.runner("last", "-n", "5"); err == nil && len(out) > 0 {
		lines := strings.Split(string(out), "\n")
		for _, l := range lines {
			if strings.HasPrefix(l, "wtmp begins") || strings.TrimSpace(l) == "" {
				continue
			}
			fields := strings.Fields(l)
			if len(fields) >= 4 {
				u := fields[0]
				tty := fields[1]
				from := fields[2]
				tStr := strings.Join(fields[3:], " ")
				events = append(events, metrics.LoginEvent{
					Time: tStr,
					User: u,
					Type: "HISTORIC",
					From: fmtFrom(tty, from),
				})
			}
		}
	}

	// Reverse to show newest events first
	for i, j := 0, len(events)-1; i < j; i, j = i+1, j-1 {
		events[i], events[j] = events[j], events[i]
	}

	c.snapshot = events
	return nil
}

func (c *LoginCollector) Snapshot() any {
	return c.snapshot
}

func fmtFrom(tty, from string) string {
	if from != "" && from != ":" {
		return from
	}
	return tty
}
