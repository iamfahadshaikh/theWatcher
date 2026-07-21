package panel

type DockerPanel struct{}

func NewDockerPanel() *DockerPanel {
	return &DockerPanel{}
}

func (p *DockerPanel) Title() string {
	return "Docker Containers"
}

func (p *DockerPanel) Render(width, height int, active bool) string {
	content := "CONTAINER ID   NAME   CPU %   MEM USAGE / LIMIT   NET I/O"
	return RenderContainer(p.Title(), content, width, height, active)
}
