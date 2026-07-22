package panel

import (
	"fmt"
	"strings"

	"theWatcher/internal/graph"
	"theWatcher/internal/history"
	"theWatcher/internal/metrics"
	"theWatcher/internal/theme"
)

type ProcessPanel struct {
	isBirth   bool
	history   *history.RingBuffer[float64]
	procs     []metrics.ProcessInfo
	birthRate float64
	renderer  graph.Renderer
	running   int
	sleeping  int
	zombie    int
	stopped   int
}

func NewProcessPanel(isBirth bool) *ProcessPanel {
	return &ProcessPanel{
		isBirth:  isBirth,
		history:  history.NewRingBuffer[float64](300),
		renderer: graph.NewBrailleRenderer(),
	}
}

func (p *ProcessPanel) UpdateProcesses(procs []metrics.ProcessInfo, birthRate float64) {
	p.procs = procs
	p.birthRate = birthRate
	if p.isBirth {
		p.history.Push(birthRate)
	}

	p.running = 0
	p.sleeping = 0
	p.zombie = 0
	p.stopped = 0
	for _, proc := range procs {
		switch proc.State {
		case "R":
			p.running++
		case "S", "D":
			p.sleeping++
		case "Z":
			p.zombie++
		case "T", "t":
			p.stopped++
		default:
			p.sleeping++
		}
	}
}

func (p *ProcessPanel) Title() string {
	if p.isBirth {
		return "Process Activity"
	}
	return "Top Processes"
}

func (p *ProcessPanel) Render(width, height int, active bool) string {
	if p.isBirth {
		graphH := height - 4
		if graphH < 1 {
			graphH = 1
		}
		graphW := width - 4
		if graphW < 1 {
			graphW = 1
		}

		slice := p.history.Slice()
		maxV := 5.0
		for _, v := range slice {
			if v > maxV {
				maxV = v
			}
		}

		brailleStr := graph.RenderHistory(p.history, p.renderer, graph.GraphData{
			Width:  graphW,
			Height: graphH,
			Min:    0,
			Max:    maxV,
			Color:  string(theme.ColorGold),
		})

		statsStr := fmt.Sprintf("PIDs/s: %.1f | Run: %d | Sleep: %d | Zombie: %d | Stop: %d", p.birthRate, p.running, p.sleeping, p.zombie, p.stopped)
		return RenderContainer(p.Title(), fmt.Sprintf("%s\n%s", brailleStr, statsStr), width, height, active)
	}

	header := fmt.Sprintf("%-6s %-6s %-8s %-4s %-2s %-4s %-4s %-3s %-6s %-8s %s",
		"PID", "PPID", "USER", "TTY", "ST", "NICE", "PRI", "THR", "CPU%", "MEM(MB)", "COMMAND")
	var lines []string
	lines = append(lines, header)

	maxLines := height - 4
	if maxLines < 1 {
		maxLines = 1
	}

	for i, proc := range p.procs {
		if i >= maxLines {
			break
		}
		cmd := proc.Command
		if len(cmd) > width-60 && width > 60 {
			cmd = cmd[:width-63] + "..."
		}
		row := fmt.Sprintf("%-6d %-6d %-8s %-4s %-2s %-4d %-4d %-3d %-6.1f %-8.1f %s",
			proc.PID, proc.PPID, proc.User, proc.TTY, proc.State, proc.Nice, proc.Priority, proc.Threads, proc.CPU, proc.Mem, cmd)
		lines = append(lines, row)
	}

	content := strings.Join(lines, "\n")
	return RenderContainer(p.Title(), content, width, height, active)
}
