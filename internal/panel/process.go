package panel

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"theWatcher/internal/graph"
	"theWatcher/internal/history"
	"theWatcher/internal/metrics"
	"theWatcher/internal/theme"
)

type ProcessPanel struct {
	isBirth      bool
	history      *history.RingBuffer[history.Sample]
	procs        []metrics.ProcessInfo
	birthRate    float64
	renderer     graph.Renderer
	running      int
	sleeping     int
	zombie       int
	stopped      int
	scrollOffset int
	SortField    string
	SortAsc      bool
	SearchFilter string
	Timeframe    string
}

func NewProcessPanel(isBirth bool) *ProcessPanel {
	return &ProcessPanel{
		isBirth:   isBirth,
		history:   history.NewRingBuffer[history.Sample](300),
		renderer:  graph.NewBrailleRenderer(),
		SortField: "cpu",
		SortAsc:   false,
		Timeframe: "30s",
	}
}

func (p *ProcessPanel) ToggleSort(field string) {
	if p.SortField == field {
		p.SortAsc = !p.SortAsc
	} else {
		p.SortField = field
		p.SortAsc = false
		if field == "pid" || field == "name" {
			p.SortAsc = true
		}
	}
	p.Sort()
}

func (p *ProcessPanel) Sort() {
	sort.Slice(p.procs, func(i, j int) bool {
		var isLess bool
		pi := p.procs[i]
		pj := p.procs[j]
		switch p.SortField {
		case "mem":
			isLess = pi.Mem > pj.Mem
		case "pid":
			isLess = pi.PID < pj.PID
		case "runtime":
			isLess = pi.Runtime > pj.Runtime
		case "name":
			isLess = pi.Command < pj.Command
		default:
			isLess = pi.CPU > pj.CPU
		}
		if p.SortAsc {
			return !isLess
		}
		return isLess
	})
}

func (p *ProcessPanel) ScrollUp() {
	p.scrollOffset--
	if p.scrollOffset < 0 {
		p.scrollOffset = 0
	}
}

func (p *ProcessPanel) ScrollDown() {
	p.scrollOffset++
}

func (p *ProcessPanel) UpdateProcesses(procs []metrics.ProcessInfo, birthRate float64) {
	p.birthRate = birthRate
	if p.isBirth {
		p.history.Push(history.Sample{Time: time.Now(), Value: birthRate})
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

	// Filter
	p.procs = nil
	if p.SearchFilter != "" {
		q := strings.ToLower(p.SearchFilter)
		for _, proc := range procs {
			if strings.Contains(strings.ToLower(proc.Command), q) ||
				strings.Contains(strings.ToLower(proc.User), q) ||
				strings.Contains(strings.ToLower(proc.ContainerName), q) {
				p.procs = append(p.procs, proc)
			}
		}
	} else {
		// Slice copy
		p.procs = append(p.procs, procs...)
	}

	p.Sort()
}

func (p *ProcessPanel) Title() string {
	if p.isBirth {
		return "Process Activity"
	}
	sortDir := "desc"
	if p.SortAsc {
		sortDir = "asc"
	}
	title := fmt.Sprintf("Top Processes (Sorted by %s %s)", strings.ToUpper(p.SortField), sortDir)
	if p.SearchFilter != "" {
		title += fmt.Sprintf(" [Filter: %s]", p.SearchFilter)
	}
	return title
}

func (p *ProcessPanel) Render(width, height int, active bool, expanded bool) string {
	if p.isBirth {
		graphH := height - 5
		if expanded {
			graphH = height - 11
		}
		if graphH < 1 {
			graphH = 1
		}
		graphW := width - 12
		if graphW < 1 {
			graphW = 1
		}

		slice := p.history.Slice()
		maxV := 5.0
		for _, v := range slice {
			if v.Value > maxV {
				maxV = v.Value
			}
		}

		brailleStr := graph.RenderHistory(p.history, p.renderer, graph.GraphData{
			Width:  graphW,
			Height: graphH,
			Min:    0,
			Max:    maxV,
			Color:  string(theme.ColorGold),
		})

		timestamps := make([]time.Time, len(slice))
		for i, s := range slice {
			timestamps[i] = s.Time
		}

		brailleStr = graph.AddYAxis(brailleStr, graphH, 0, maxV, "")
		xAxisStr := "        " + graph.RenderXAxis(graphW, timestamps, p.Timeframe)

		var statsStr string
		if expanded {
			statsStr = fmt.Sprintf("Fork Rate    : %.1f PIDs/sec\nRunning      : %d\nSleeping     : %d\nZombie       : %d\nStopped      : %d\nTotal Procs  : %d",
				p.birthRate, p.running, p.sleeping, p.zombie, p.stopped, len(p.procs))
		} else {
			statsStr = fmt.Sprintf("PIDs/s: %.1f | Run: %d | Sleep: %d | Zombie: %d | Stop: %d", p.birthRate, p.running, p.sleeping, p.zombie, p.stopped)
		}
		return RenderContainer(p.Title(), fmt.Sprintf("%s\n%s\n%s", brailleStr, xAxisStr, statsStr), width, height, active)
	}

	var header string
	if expanded {
		header = fmt.Sprintf("%-6s %-6s %-8s %-4s %-2s %-4s %-4s %-3s %-6s %-8s %-8s %-8s %-4s %-8s %-8s %-12s %s",
			"PID", "PPID", "USER", "TTY", "ST", "NICE", "PRI", "THR", "CPU%", "MEM(MB)", "VIRT(M)", "RES(M)", "FDs", "R/s", "W/s", "CONTAINER", "COMMAND")
	} else {
		header = fmt.Sprintf("%-6s %-6s %-8s %-4s %-2s %-4s %-4s %-3s %-6s %-8s %s",
			"PID", "PPID", "USER", "TTY", "ST", "NICE", "PRI", "THR", "CPU%", "MEM(MB)", "COMMAND")
	}

	var lines []string
	lines = append(lines, header)

	maxLines := height - 4
	if maxLines < 1 {
		maxLines = 1
	}

	if p.scrollOffset > len(p.procs)-maxLines {
		p.scrollOffset = len(p.procs) - maxLines
	}
	if p.scrollOffset < 0 {
		p.scrollOffset = 0
	}

	for i := p.scrollOffset; i < len(p.procs); i++ {
		if len(lines)-1 >= maxLines {
			break
		}
		proc := p.procs[i]
		cmd := proc.Command
		if expanded {
			if len(cmd) > width-110 && width > 110 {
				cmd = cmd[:width-113] + "..."
			}
			rRate := formatRate(proc.ReadBytes)
			wRate := formatRate(proc.WriteBytes)
			virtM := float64(proc.VirtMemBytes) / (1024 * 1024)
			resM := float64(proc.ResMemBytes) / (1024 * 1024)
			row := fmt.Sprintf("%-6d %-6d %-8s %-4s %-2s %-4d %-4d %-3d %-6.1f %-8.1f %-8.1f %-8.1f %-4d %-8s %-8s %-12s %s",
				proc.PID, proc.PPID, truncate(proc.User, 8), proc.TTY, proc.State, proc.Nice, proc.Priority, proc.Threads, proc.CPU, proc.Mem, virtM, resM, proc.OpenFiles, rRate, wRate, truncate(proc.ContainerName, 12), cmd)
			lines = append(lines, row)
		} else {
			if len(cmd) > width-60 && width > 60 {
				cmd = cmd[:width-63] + "..."
			}
			row := fmt.Sprintf("%-6d %-6d %-8s %-4s %-2s %-4d %-4d %-3d %-6.1f %-8.1f %s",
				proc.PID, proc.PPID, truncate(proc.User, 8), proc.TTY, proc.State, proc.Nice, proc.Priority, proc.Threads, proc.CPU, proc.Mem, cmd)
			lines = append(lines, row)
		}
	}

	content := strings.Join(lines, "\n")
	return RenderContainer(p.Title(), content, width, height, active)
}

func formatRate(b uint64) string {
	if b == 0 {
		return "0 B/s"
	}
	return formatBytes(b) + "/s"
}
