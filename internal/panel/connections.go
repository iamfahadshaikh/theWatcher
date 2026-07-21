package panel

type ConnectionsPanel struct{}

func NewConnectionsPanel() *ConnectionsPanel {
	return &ConnectionsPanel{}
}

func (p *ConnectionsPanel) Title() string {
	return "Active Connections"
}

func (p *ConnectionsPanel) Render(width, height int, active bool) string {
	content := "ESTAB: 0 | LISTEN: 0 | TIME_WAIT: 0 | CLOSE_WAIT: 0"
	return RenderContainer(p.Title(), content, width, height, active)
}
