package panel

type DiskPanel struct{}

func NewDiskPanel() *DiskPanel {
	return &DiskPanel{}
}

func (p *DiskPanel) Title() string {
	return "Disk"
}

func (p *DiskPanel) Render(width, height int, active bool) string {
	content := "  Live Graph Placeholder\n  Read: 0 KB/s  Write: 0 KB/s"
	return RenderContainer(p.Title(), content, width, height, active)
}
