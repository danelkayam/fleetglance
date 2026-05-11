package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	IconColWidth = 2
)

func IconCell(icon string, style lipgloss.Style) string {
	return style.
		Inline(true).
		Width(IconColWidth).
		MaxWidth(IconColWidth).
		Align(lipgloss.Center).
		Render(icon)
}

func PadRight(s string, width int) string {
	currentWidth := lipgloss.Width(s)
	if currentWidth >= width {
		return s
	}

	return s + strings.Repeat(" ", width-currentWidth)
}

func TruncatePlain(value string, width int) string {
	if width <= 0 {
		return ""
	}

	var builder strings.Builder
	currentWidth := 0
	for _, r := range value {
		runeWidth := lipgloss.Width(string(r))
		if currentWidth+runeWidth > width {
			break
		}

		builder.WriteRune(r)
		currentWidth += runeWidth
	}

	return builder.String()
}
