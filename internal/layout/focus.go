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

func (f *FocusManager) Navigate(dir string, rows RowConfig) {
	currID := f.ActiveID
	nextID := currID

	switch dir {
	case "left":
		switch currID {
		case PanelMemory:
			nextID = PanelCPU
		case PanelDisk:
			nextID = PanelMemory
		case PanelNetTX:
			nextID = PanelNetRX
		case PanelProcBirth:
			nextID = PanelNetTX
		case PanelConnections:
			nextID = PanelPorts
		}
	case "right":
		switch currID {
		case PanelCPU:
			nextID = PanelMemory
		case PanelMemory:
			nextID = PanelDisk
		case PanelNetRX:
			nextID = PanelNetTX
		case PanelNetTX:
			nextID = PanelProcBirth
		case PanelPorts:
			nextID = PanelConnections
		}
	case "up":
		nextID = moveVertical(currID, -1, rows)
	case "down":
		nextID = moveVertical(currID, 1, rows)
	}

	if isPanelVisible(nextID, rows) {
		f.ActiveID = nextID
	}
}

func moveVertical(currID PanelID, step int, rows RowConfig) PanelID {
	rowMap := map[PanelID]int{
		PanelCPU: 1, PanelMemory: 1, PanelDisk: 1,
		PanelNetRX: 2, PanelNetTX: 2, PanelProcBirth: 2,
		PanelPorts: 3, PanelConnections: 3,
		PanelProcesses: 4,
		PanelDocker:    5,
		PanelLogins:    6,
		PanelEvents:    7,
	}

	colMap := map[PanelID]int{
		PanelCPU: 0, PanelMemory: 1, PanelDisk: 2,
		PanelNetRX: 0, PanelNetTX: 1, PanelProcBirth: 2,
		PanelPorts: 0, PanelConnections: 1,
		PanelProcesses: 0,
		PanelDocker:    0,
		PanelLogins:    0,
		PanelEvents:    0,
	}

	currRow := rowMap[currID]
	currCol := colMap[currID]

	targetRow := currRow
	for {
		targetRow += step
		if targetRow < 1 {
			targetRow = 7
		}
		if targetRow > 7 {
			targetRow = 1
		}
		if targetRow == currRow {
			break
		}

		if isRowVisible(targetRow, rows) {
			return getPanelInRowAndCol(targetRow, currCol)
		}
	}
	return currID
}

func isRowVisible(row int, rows RowConfig) bool {
	switch row {
	case 1:
		return rows.ShowRow1
	case 2:
		return rows.ShowRow2
	case 3:
		return rows.ShowRow3
	case 4:
		return rows.ShowRow4
	case 5:
		return rows.ShowRow5
	case 6:
		return rows.ShowRow6
	case 7:
		return rows.ShowRow7
	}
	return false
}

func getPanelInRowAndCol(row int, col int) PanelID {
	switch row {
	case 1:
		if col == 0 {
			return PanelCPU
		}
		if col == 1 {
			return PanelMemory
		}
		return PanelDisk
	case 2:
		if col == 0 {
			return PanelNetRX
		}
		if col == 1 {
			return PanelNetTX
		}
		return PanelProcBirth
	case 3:
		if col == 0 {
			return PanelPorts
		}
		return PanelConnections
	case 4:
		return PanelProcesses
	case 5:
		return PanelDocker
	case 6:
		return PanelLogins
	case 7:
		return PanelEvents
	}
	return PanelCPU
}
