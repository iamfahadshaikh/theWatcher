package panel

type ProcessPanel struct {
	isBirth bool
}

func NewProcessPanel(isBirth bool) *ProcessPanel {
	return &ProcessPanel{isBirth: isBirth}
}

func (p *ProcessPanel) Title() string {
	if p.isBirth {
		return "Process Birth"
	}
	return "Top Processes"
}

func (p *ProcessPanel) Render(width, height int, active bool) string {
	if p.isBirth {
		return RenderContainer(p.Title(), "  Live Graph Placeholder\n  New PIDs: 0/s", width, height, active)
	}
	content := "PID    USER    CPU%   MEM%   COMMAND\n---    ----    ----   ----   -------"
	return RenderContainer(p.Title(), content, width, height, active)
}
