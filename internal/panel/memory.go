package panel

type MemoryPanel struct{}

func NewMemoryPanel() *MemoryPanel {
	return &MemoryPanel{}
}

func (p *MemoryPanel) Title() string {
	return "Memory"
}

func (p *MemoryPanel) Render(width, height int, active bool) string {
	content := "  Live Graph Placeholder\n  Used: 0MB Free: 0MB"
	return RenderContainer(p.Title(), content, width, height, active)
}
