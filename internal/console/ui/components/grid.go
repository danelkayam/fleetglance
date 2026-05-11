package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	maxPaneColumns = 4
	gridGapX       = 2
	gridGapY       = 1
)

const (
	preferredPaneWidth = 40
	minPaneWidth       = 17
	detailedPaneHeight = 17
	compactPaneHeight  = 12
)

type Grid struct {
	Width  int
	Height int
	Ships  []Ship
}

type gridLayout struct {
	columns    int
	paneWidth  int
	paneHeight int
	compact    bool
}

func (g Grid) Render() string {
	if len(g.Ships) == 0 {
		return ""
	}

	layout := chooseGridLayout(g.Width, g.Height, len(g.Ships))
	columns := max(layout.columns, 1)
	rows := make([]string, 0, (len(g.Ships)+columns-1)/columns)

	for i := 0; i < len(g.Ships); i += columns {
		end := min(i+columns, len(g.Ships))
		panes := make([]string, 0, end-i)
		for offset, ship := range g.Ships[i:end] {
			panes = append(panes, Pane{
				Ship:    ship,
				Index:   i + offset,
				Width:   layout.paneWidth,
				Height:  layout.paneHeight,
				Compact: layout.compact,
			}.Render())
		}

		rows = append(rows, joinPanes(panes))
	}

	return strings.Join(rows, strings.Repeat("\n", gridGapY+1))
}

func chooseGridLayout(width int, height int, count int) gridLayout {
	if count <= 0 {
		return gridLayout{columns: 1, paneWidth: max(width, minPaneWidth), paneHeight: compactPaneHeight, compact: true}
	}

	for _, compact := range []bool{false, true} {
		paneHeight := detailedPaneHeight
		if compact {
			paneHeight = compactPaneHeight
		}

		for columns := min(maxPaneColumns, count); columns >= 1; columns-- {
			paneWidth := paneWidthForColumns(width, columns)
			if paneWidth < minPaneWidth {
				continue
			}

			rows := rowsFor(count, columns)
			neededHeight := rows*paneHeight + max(rows-1, 0)*gridGapY
			if neededHeight <= height {
				return gridLayout{
					columns:    columns,
					paneWidth:  paneWidth,
					paneHeight: paneHeight,
					compact:    compact,
				}
			}
		}
	}

	columns := bestColumnsForWidth(width, count)
	rows := rowsFor(count, columns)
	gaps := max(rows-1, 0) * gridGapY
	paneHeight := max((height-gaps)/rows, 1)

	return gridLayout{
		columns:    columns,
		paneWidth:  paneWidthForColumns(width, columns),
		paneHeight: paneHeight,
		compact:    true,
	}
}

func joinPanes(panes []string) string {
	if len(panes) == 0 {
		return ""
	}

	gap := ContentStyle.Render(strings.Repeat(" ", gridGapX))
	row := panes[0]
	for _, pane := range panes[1:] {
		row = lipgloss.JoinHorizontal(lipgloss.Top, row, gap, pane)
	}

	return row
}

func bestColumnsForWidth(width int, count int) int {
	for columns := min(maxPaneColumns, count); columns > 1; columns-- {
		if paneWidthForColumns(width, columns) >= minPaneWidth {
			return columns
		}
	}

	return 1
}

func paneWidthForColumns(width int, columns int) int {
	if columns <= 1 {
		return max(min(width, preferredPaneWidth), 1)
	}

	return max(min((width-gridGapX*(columns-1))/columns, preferredPaneWidth), 1)
}

func rowsFor(count int, columns int) int {
	return (count + columns - 1) / columns
}
