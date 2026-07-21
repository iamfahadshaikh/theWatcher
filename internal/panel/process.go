package panel

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"theWatcher/internal/graph"
	"theWatcher/internal/history"
	"theWatcher/internal/model"
)

type ProcessPanel struct {
	isBirth   bool
	history   *history.RingBuffer[float64]
	procs     []model.ProcessInfo
	birthRate float64
}

func NewProcessPanel(isBirth bool) *ProcessPanel {
	return &ProcessPanel{
		isBirth: isBirth,
		history: history.NewRingBuffer[float64](300),
	}
}

func (p *ProcessPanel) UpdateProcesses(procs []model.ProcessInfo, birthRate float64) {
	p.procs = procs
	p.birthRate = birthRate
	if p.isBirth {
		p.history.Push(birthRate)
	}
}

func (p *ProcessPanel) Title() string {
	if p.isBirth {
		return "Process Birth"
	}
	return "Top Processes"
}

func (p *ProcessPanel) Render(width, height int, active bool) string {
	if p.isBirth {
		slice := p.history.Slice()
		graphH := height - 4
		if graphH < 1 {
			graphH = 1
		}
		graphW := width - 4
		if graphW < 1 {
			graphW = 1
		}

		lineGraph := graph.RenderLineGraph(slice, graphW, graphH, 0, 50, lipgloss.NewStyle().Foreground(lipgloss.Color("#f6c177")))
		statsStr := fmt.Sprintf("New PIDs/sec: %.1f", p.birthRate)
		return RenderContainer(p.Title(), fmt.Sprintf("%s\n%s", lineGraph, statsStr), width, height, active)
	}

	header := fmt.Sprintf("%-6s %-8s %-6s %-8s %s", "PID", "USER", "CPU%", "MEM(MB)", "COMMAND")
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
		if len(cmd) > width-35 && width > 35 {
			cmd = cmd[:width-38] + "..."
		}
		row := fmt.Sprintf("%-6d %-8s %-6.1f %-8.1f %s", proc.PID, proc.User, proc.CPU, proc.Mem, cmd)
		lines = append(lines, row)
	}

	content := strings.Join(lines, "\n")
	return RenderContainer(p.Title(), content, width, height, active)
}
