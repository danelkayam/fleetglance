package components

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	screenPaddingX = 2
	screenPaddingY = 1
	headerGridGap  = 2
)

const (
	defaultTerminalWidth  = 80
	defaultTerminalHeight = 24
)

func Render(model Model) string {
	width, height := screenSize(model.Width, model.Height)
	horizontalPadding := horizontalPadding(width)
	verticalPadding := verticalPadding(height)
	contentWidth := max(width-horizontalPadding*2, 1)
	contentHeight := max(height-verticalPadding*2, 1)

	header := TopBar{
		Width:   contentWidth,
		Summary: model.Summary,
		Now:     model.Now,
	}.Render()
	headerHeight := lipgloss.Height(header)
	gridHeight := max(contentHeight-headerHeight-headerGridGap, 1)
	gridView := Grid{
		Width:  contentWidth,
		Height: gridHeight,
		Ships:  model.Ships,
	}.Render()

	contentParts := []string{header}
	if gridView != "" && gridHeight > 0 {
		for range headerGridGap {
			contentParts = append(contentParts, "")
		}
		contentParts = append(contentParts, gridView)
	}
	content := ContentStyle.Render(lipgloss.JoinVertical(lipgloss.Left, contentParts...))
	content = ContentStyle.
		Width(contentWidth).
		Height(contentHeight).
		MaxHeight(contentHeight).
		Render(content)
	content = ContentStyle.
		Padding(verticalPadding, horizontalPadding).
		Render(content)

	return BackgroundStyle.Width(width).Height(height).Render(lipgloss.Place(
		width,
		height,
		lipgloss.Left,
		lipgloss.Top,
		content,
	))
}

func screenSize(width int, height int) (int, int) {
	if width <= 0 {
		width = defaultTerminalWidth
	}

	if height <= 0 {
		height = defaultTerminalHeight
	}

	return max(width, 1), max(height, 1)
}

func horizontalPadding(width int) int {
	switch {
	case width >= 120:
		return screenPaddingX
	case width >= 80:
		return 2
	default:
		return 1
	}
}

func verticalPadding(height int) int {
	if height >= 24 {
		return screenPaddingY
	}

	if height >= 18 {
		return 1
	}

	return 0
}
