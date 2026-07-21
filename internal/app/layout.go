package app

// LayoutDimensions stores calculated UI component bounds.
type LayoutDimensions struct {
	Width        int
	Height       int
	HeaderHeight int
	FooterHeight int
	BodyHeight   int
}

// ComputeLayout calculates heights and widths for dashboard regions.
func ComputeLayout(width, height int) LayoutDimensions {
	headerH := 1
	footerH := 1
	bodyH := height - headerH - footerH
	if bodyH < 1 {
		bodyH = 1
	}
	return LayoutDimensions{
		Width:        width,
		Height:       height,
		HeaderHeight: headerH,
		FooterHeight: footerH,
		BodyHeight:   bodyH,
	}
}
