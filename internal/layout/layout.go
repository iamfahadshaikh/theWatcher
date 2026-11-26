package layout

// LayoutDimensions stores calculated UI component bounds and row configurations.
type LayoutDimensions struct {
	Width        int
	Height       int
	HeaderHeight int
	FooterHeight int
	BodyHeight   int
	Rows         RowConfig
}

type RowConfig struct {
	ShowRow1   bool
	ShowRow2   bool
	ShowRow3   bool
	ShowRow4   bool
	ShowRow5   bool
	ShowRow6   bool
	ShowRow7   bool
	Row1Height int
	Row2Height int
	Row3Height int
	Row4Height int
	Row5Height int
	Row6Height int
	Row7Height int
}

// ComputeLayout calculates heights and widths for dashboard regions.
func ComputeLayout(width, height int) LayoutDimensions {
	headerH := 2
	footerH := 2
	bodyH := height - headerH - footerH
	if bodyH < 1 {
		bodyH = 1
	}

	rows := RowConfig{}

	if bodyH >= 45 {
		rows.ShowRow1 = true
		rows.ShowRow2 = true
		rows.ShowRow3 = true
		rows.ShowRow4 = true
		rows.ShowRow5 = true
		rows.ShowRow6 = true
		rows.ShowRow7 = true
		rows.Row1Height = 8
		rows.Row2Height = 8
		rows.Row3Height = 6
		rows.Row6Height = 5 // Login panel
		rows.Row7Height = 5 // Events & Alerts panel (compacted to prevent screen waste)
		rows.Row5Height = 5 // Docker / Disk Mounts
		rows.Row4Height = bodyH - (rows.Row1Height + rows.Row2Height + rows.Row3Height + rows.Row5Height + rows.Row6Height + rows.Row7Height)
		if rows.Row4Height < 8 {
			rows.Row4Height = 8
		}
	} else if bodyH >= 35 {
		rows.ShowRow1 = true
		rows.ShowRow2 = true
		rows.ShowRow3 = true
		rows.ShowRow4 = true
		rows.ShowRow5 = true
		rows.ShowRow6 = false
		rows.ShowRow7 = false
		rows.Row1Height = 7
		rows.Row2Height = 7
		rows.Row3Height = 5
		rows.Row5Height = 5
		rows.Row4Height = bodyH - (rows.Row1Height + rows.Row2Height + rows.Row3Height + rows.Row5Height)
	} else if bodyH >= 22 {
		rows.ShowRow1 = true
		rows.ShowRow2 = true
		rows.ShowRow3 = true
		rows.ShowRow4 = true
		rows.ShowRow5 = false
		rows.ShowRow6 = false
		rows.ShowRow7 = false
		rows.Row1Height = 5
		rows.Row2Height = 5
		rows.Row3Height = 4
		rows.Row4Height = bodyH - (5 + 5 + 4)
	} else {
		rows.ShowRow1 = true
		rows.ShowRow2 = false
		rows.ShowRow3 = false
		rows.ShowRow4 = true
		rows.ShowRow5 = false
		rows.ShowRow6 = false
		rows.ShowRow7 = false
		rows.Row1Height = bodyH / 2
		if rows.Row1Height < 4 {
			rows.Row1Height = 4
		}
		if rows.Row1Height > bodyH-4 {
			rows.Row1Height = bodyH - 4
		}
		if rows.Row1Height < 1 {
			rows.Row1Height = 1
		}
		rows.Row4Height = bodyH - rows.Row1Height
		if rows.Row4Height < 1 {
			rows.Row4Height = 1
		}
	}

	return LayoutDimensions{
		Width:        width,
		Height:       height,
		HeaderHeight: headerH,
		FooterHeight: footerH,
		BodyHeight:   bodyH,
		Rows:         rows,
	}
}
