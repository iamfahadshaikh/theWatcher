package collector

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"sync"
	"time"

	"theWatcher/internal/metrics"
)

type CommandRunner func(ctx context.Context, name string, args ...string) ([]byte, error)

type LoginCollector struct {
	mu       sync.Mutex
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

func defaultCommandRunner(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
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

func (c *LoginCollector) ResetSession() {
	c.mu.Lock()
	defer c.mu.Unlock()
}

func (c *LoginCollector) Interval() time.Duration {
	return 10 * time.Second
}

func (c *LoginCollector) Collect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var events []metrics.LoginEvent

	if out, err := c.runner(ctx, "who"); err == nil && len(out) > 0 {
		lines := strings.Split(string(out), "\n")
		for _, l := range lines {
			fields := strings.Fields(l)
			if len(fields) >= 3 {
				u := fields[0]
				tty := fields[1]
				tStr := strings.Join(fields[2:], " ")
				from := tty
				if len(fields) >= 5 {
					lastField := fields[len(fields)-1]
					if strings.HasPrefix(lastField, "(") && strings.HasSuffix(lastField, ")") {
						from = strings.Trim(lastField, "()")
					}
				}
				events = append(events, metrics.LoginEvent{
					User: u,
					Type: "SSH",
					Time: tStr,
					From: from,
				})
			}
		}
	}

	if out, err := c.runner(ctx, "last", "-n", "10"); err == nil && len(out) > 0 {
		lines := strings.Split(string(out), "\n")
		for _, l := range lines {
			if strings.TrimSpace(l) == "" || strings.HasPrefix(l, "wtmp") || strings.HasPrefix(l, "reboot") {
				continue
			}
			fields := strings.Fields(l)
			if len(fields) >= 4 {
				u := fields[0]
				from := fields[2]
				var tStr string
				if len(fields) >= 7 {
					tStr = strings.Join(fields[3:7], " ")
				} else {
					tStr = strings.Join(fields[3:], " ")
				}
				events = append(events, metrics.LoginEvent{
					User: u,
					Type: "LAST LOGIN",
					Time: tStr,
					From: from,
				})
			}
		}
	}

	c.snapshot = events
	return nil
}

func (c *LoginCollector) Snapshot() any {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.snapshot
}
