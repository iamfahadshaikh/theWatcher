package collector

import (
	"bytes"
	"os/exec"
	"strings"

	"theWatcher/internal/model"
)

type CommandRunner func(name string, args ...string) ([]byte, error)

type LoginCollector struct {
	runner CommandRunner
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

func (c *LoginCollector) Collect() ([]model.LoginEvent, error) {
	var events []model.LoginEvent

	// 1. Check active logins with 'who'
	if out, err := c.runner("who"); err == nil && len(out) > 0 {
		lines := strings.Split(string(out), "\n")
		for _, l := range lines {
			fields := strings.Fields(l)
			if len(fields) >= 3 {
				u := fields[0]
				tty := fields[1]
				tStr := strings.Join(fields[2:], " ")
				events = append(events, model.LoginEvent{
					Time: tStr,
					User: u,
					Type: "SESSION",
					From: tty,
				})
			}
		}
	}

	// 2. Check recent logins with 'last -n 5'
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
				events = append(events, model.LoginEvent{
					Time: tStr,
					User: u,
					Type: "HISTORIC",
					From: fmtFrom(tty, from),
				})
			}
		}
	}

	return events, nil
}

func fmtFrom(tty, from string) string {
	if from != "" && from != ":" {
		return from
	}
	return tty
}
