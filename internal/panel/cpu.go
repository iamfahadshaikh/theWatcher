package panel

type CPUPanel struct{}

func NewCPUPanel() *CPUPanel {
	return &CPUPanel{}
}

func (p *CPUPanel) Title() string {
	return "CPU"
}

func (p *CPUPanel) Render(width, height int, active bool) string {
	content := "  Live Graph Placeholder\n  Avg: 0.0%  Peak: 0.0%"
	return RenderContainer(p.Title(), content, width, height, active)
}
