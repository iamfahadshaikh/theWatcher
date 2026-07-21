package panel

type PortsPanel struct{}

func NewPortsPanel() *PortsPanel {
	return &PortsPanel{}
}

func (p *PortsPanel) Title() string {
	return "Listening Ports"
}

func (p *PortsPanel) Render(width, height int, active bool) string {
	content := "PROTO   PORT   PROCESS\n-----   ----   -------"
	return RenderContainer(p.Title(), content, width, height, active)
}
