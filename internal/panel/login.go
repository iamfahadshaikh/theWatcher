package panel

type LoginPanel struct{}

func NewLoginPanel() *LoginPanel {
	return &LoginPanel{}
}

func (p *LoginPanel) Title() string {
	return "Login Timeline"
}

func (p *LoginPanel) Render(width, height int, active bool) string {
	content := "TIME               USER     TYPE   FROM"
	return RenderContainer(p.Title(), content, width, height, active)
}
