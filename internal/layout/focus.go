package layout

type PanelID int

const (
	PanelCPU PanelID = iota
	PanelMemory
	PanelDisk
	PanelNetRX
	PanelNetTX
	PanelProcBirth
	PanelPorts
	PanelConnections
	PanelProcesses
	PanelDocker
	PanelLogins
	PanelEvents
	TotalPanels
)

type FocusManager struct {
	ActiveID  PanelID
	Maximized bool
}

func NewFocusManager() *FocusManager {
	return &FocusManager{
		ActiveID:  PanelCPU,
		Maximized: false,
	}
}

func (f *FocusManager) NextFocus() {
	f.ActiveID = (f.ActiveID + 1) % TotalPanels
}

func (f *FocusManager) PrevFocus() {
	f.ActiveID = (f.ActiveID - 1 + TotalPanels) % TotalPanels
}

func (f *FocusManager) ToggleMaximize() {
	f.Maximized = !f.Maximized
}

func (f *FocusManager) ResetMaximize() {
	f.Maximized = false
}

func (f *FocusManager) IsActive(id PanelID) bool {
	return f.ActiveID == id
}
