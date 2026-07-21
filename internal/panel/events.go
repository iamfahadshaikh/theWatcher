package panel

type EventsPanel struct{}

func NewEventsPanel() *EventsPanel {
	return &EventsPanel{}
}

func (p *EventsPanel) Title() string {
	return "Events"
}

func (p *EventsPanel) Render(width, height int, active bool) string {
	content := "[INFO] Dashboard initialized."
	return RenderContainer(p.Title(), content, width, height, active)
}
