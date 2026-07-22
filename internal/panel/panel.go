package panel

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	"theWatcher/internal/theme"
)

type Panel interface {
	Title() string
	Render(width, height int, active bool, expanded bool) string
}

func RenderContainer(title, content string, width, height int, active bool) string {
	var borderColor lipgloss.Color
	var titleStyle lipgloss.Style
	var borderType lipgloss.Border

	if active {
		borderColor = theme.ColorGold
		borderType = theme.MaximizedBorder
		titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.ColorBase).
			Background(theme.ColorGold)
		title = fmt.Sprintf(" %s %s ", title, theme.IconFocus)
	} else {
		borderColor = theme.ColorSubtle
		borderType = theme.DefaultBorder
		titleStyle = theme.PanelTitleStyle
		title = fmt.Sprintf(" %s ", title)
	}

	maxContentLines := height - 3
	if maxContentLines < 0 {
		maxContentLines = 0
	}

	contentLines := strings.Split(content, "\n")
	if len(contentLines) > maxContentLines {
		contentLines = contentLines[:maxContentLines]
	}

	lineLimit := width - 4
	if lineLimit < 1 {
		lineLimit = 1
	}
	for i, line := range contentLines {
		contentLines[i] = AnsiSafeTruncate(line, lineLimit)
	}

	content = strings.Join(contentLines, "\n")

	style := lipgloss.NewStyle().
		Border(borderType).
		BorderForeground(borderColor).
		Width(width - 2).
		Height(height - 2)

	renderedTitle := titleStyle.Render(title)
	body := fmt.Sprintf("%s\n%s", renderedTitle, content)
	return style.Render(body)
}

func AnsiSafeTruncate(s string, limit int) string {
	if limit <= 0 {
		return ""
	}
	var result strings.Builder
	visibleLen := 0
	inEscape := false

	i := 0
	for i < len(s) {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError {
			result.WriteByte(s[i])
			i++
			continue
		}

		if r == '\x1b' {
			inEscape = true
			result.WriteString(s[i : i+size])
			i += size
			continue
		}

		if inEscape {
			result.WriteString(s[i : i+size])
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				inEscape = false
			}
			i += size
			continue
		}

		if visibleLen >= limit {
			i += size
			continue
		}

		result.WriteString(s[i : i+size])
		visibleLen++
		i += size
	}
	return result.String()
}
