package panel

type NetworkPanel struct {
	isTX bool
}

func NewNetworkPanel(isTX bool) *NetworkPanel {
	return &NetworkPanel{isTX: isTX}
}

func (p *NetworkPanel) Title() string {
	if p.isTX {
		return "Network TX"
	}
	return "Network RX"
}

func (p *NetworkPanel) Render(width, height int, active bool) string {
	content := "  Live Graph Placeholder\n  Rate: 0 KB/s"
	return RenderContainer(p.Title(), content, width, height, active)
}
