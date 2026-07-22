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

func (f *FocusManager) NextFocus(rows RowConfig) {
	for i := 0; i < int(TotalPanels); i++ {
		f.ActiveID = (f.ActiveID + 1) % TotalPanels
		if isPanelVisible(f.ActiveID, rows) {
			return
		}
	}
}

func (f *FocusManager) PrevFocus(rows RowConfig) {
	for i := 0; i < int(TotalPanels); i++ {
		f.ActiveID = (f.ActiveID - 1 + TotalPanels) % TotalPanels
		if isPanelVisible(f.ActiveID, rows) {
			return
		}
	}
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

func isPanelVisible(id PanelID, rows RowConfig) bool {
	switch id {
	case PanelCPU, PanelMemory, PanelDisk:
		return rows.ShowRow1
	case PanelNetRX, PanelNetTX, PanelProcBirth:
		return rows.ShowRow2
	case PanelPorts, PanelConnections:
		return rows.ShowRow3
	case PanelProcesses:
		return rows.ShowRow4
	case PanelDocker:
		return rows.ShowRow5
	case PanelLogins:
		return rows.ShowRow6
	case PanelEvents:
		return rows.ShowRow7
	}
	return false
}
