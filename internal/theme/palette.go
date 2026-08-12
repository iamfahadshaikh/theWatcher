package theme

import "github.com/charmbracelet/lipgloss"

type ThemeColors struct {
	Love    lipgloss.Color
	Gold    lipgloss.Color
	Foam    lipgloss.Color
	Pine    lipgloss.Color
	Rose    lipgloss.Color
	Iris    lipgloss.Color
	Subtle  lipgloss.Color
	Text    lipgloss.Color
	Base    lipgloss.Color
	Surface lipgloss.Color
	Success lipgloss.Color
}

// Rose Pine (Default) colors
var DefaultTheme = ThemeColors{
	Love:    lipgloss.Color("#eb6f92"),
	Gold:    lipgloss.Color("#f6c177"),
	Foam:    lipgloss.Color("#9ccfd8"),
	Pine:    lipgloss.Color("#31748f"),
	Rose:    lipgloss.Color("#ebbcba"),
	Iris:    lipgloss.Color("#c4a7e7"),
	Subtle:  lipgloss.Color("#908caa"),
	Text:    lipgloss.Color("#e0def4"),
	Base:    lipgloss.Color("#191724"),
	Surface: lipgloss.Color("#1f1d2e"),
	Success: lipgloss.Color("#a6e3a1"),
}

// Dracula theme colors
var DraculaTheme = ThemeColors{
	Love:    lipgloss.Color("#ff5555"), // Red
	Gold:    lipgloss.Color("#f1fa8c"), // Yellow
	Foam:    lipgloss.Color("#8be9fd"), // Cyan
	Pine:    lipgloss.Color("#50fa7b"), // Green
	Rose:    lipgloss.Color("#ff79c6"), // Pink
	Iris:    lipgloss.Color("#bd93f9"), // Purple
	Subtle:  lipgloss.Color("#6272a4"), // Comment / Gray
	Text:    lipgloss.Color("#f8f8f2"), // Foreground
	Base:    lipgloss.Color("#282a36"), // Background
	Surface: lipgloss.Color("#44475a"), // Selection / Current line
	Success: lipgloss.Color("#50fa7b"), // Green
}

// Nord theme colors
var NordTheme = ThemeColors{
	Love:    lipgloss.Color("#bf616a"), // Red
	Gold:    lipgloss.Color("#ebcb8b"), // Yellow
	Foam:    lipgloss.Color("#88c0d0"), // Cyan
	Pine:    lipgloss.Color("#8fbcbb"), // Teal
	Rose:    lipgloss.Color("#d8dee9"), // Off-white
	Iris:    lipgloss.Color("#b48ead"), // Purple
	Subtle:  lipgloss.Color("#4c566a"), // Slate
	Text:    lipgloss.Color("#e5e9f0"), // Snow Storm
	Base:    lipgloss.Color("#2e3440"), // Polar Night
	Surface: lipgloss.Color("#3b4252"), // Polar Night Light
	Success: lipgloss.Color("#a3be8c"), // Green
}

// Gruvbox theme colors
var GruvboxTheme = ThemeColors{
	Love:    lipgloss.Color("#cc241d"), // Red
	Gold:    lipgloss.Color("#d79921"), // Yellow
	Foam:    lipgloss.Color("#83a598"), // Blue
	Pine:    lipgloss.Color("#8ec07c"), // Aqua
	Rose:    lipgloss.Color("#d5c4a1"), // Light gray
	Iris:    lipgloss.Color("#b16286"), // Purple
	Subtle:  lipgloss.Color("#928374"), // Gray
	Text:    lipgloss.Color("#ebdbb2"), // Foreground
	Base:    lipgloss.Color("#282828"), // Background
	Surface: lipgloss.Color("#3c3836"), // Dark surface
	Success: lipgloss.Color("#b8bb26"), // Green
}

var (
	ColorLove    = DefaultTheme.Love
	ColorGold    = DefaultTheme.Gold
	ColorFoam    = DefaultTheme.Foam
	ColorPine    = DefaultTheme.Pine
	ColorRose    = DefaultTheme.Rose
	ColorIris    = DefaultTheme.Iris
	ColorSubtle  = DefaultTheme.Subtle
	ColorText    = DefaultTheme.Text
	ColorBase    = DefaultTheme.Base
	ColorSurface = DefaultTheme.Surface
	ColorSuccess = DefaultTheme.Success
)

func SetTheme(name string) {
	var t ThemeColors
	switch name {
	case "dracula":
		t = DraculaTheme
	case "nord":
		t = NordTheme
	case "gruvbox":
		t = GruvboxTheme
	default:
		t = DefaultTheme
	}

	ColorLove = t.Love
	ColorGold = t.Gold
	ColorFoam = t.Foam
	ColorPine = t.Pine
	ColorRose = t.Rose
	ColorIris = t.Iris
	ColorSubtle = t.Subtle
	ColorText = t.Text
	ColorBase = t.Base
	ColorSurface = t.Surface
	ColorSuccess = t.Success

	ReinitializeStyles()
}
