package theme

import "github.com/charmbracelet/lipgloss"

var (
	HeaderTitleStyle       lipgloss.Style
	HeaderItemStyle        lipgloss.Style
	HeaderValueStyle       lipgloss.Style
	FooterStyle            lipgloss.Style
	PanelTitleStyle        lipgloss.Style
	PanelTitleFocusedStyle lipgloss.Style
)

func init() {
	ReinitializeStyles()
}

func ReinitializeStyles() {
	HeaderTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorLove).
		Background(ColorSurface).
		Padding(0, 1)

	HeaderItemStyle = lipgloss.NewStyle().
		Foreground(ColorFoam).
		Background(ColorSurface).
		Padding(0, 1)

	HeaderValueStyle = lipgloss.NewStyle().
		Foreground(ColorGold).
		Background(ColorSurface).
		Padding(0, 1)

	FooterStyle = lipgloss.NewStyle().
		Foreground(ColorSubtle).
		Background(ColorBase).
		Padding(0, 1)

	PanelTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorText)

	PanelTitleFocusedStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorRose)
}
