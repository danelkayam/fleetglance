package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type TopBar struct {
	Width   int
	Summary Summary
	Now     time.Time
}

func (t TopBar) Render() string {
	if t.Width <= 0 {
		return ""
	}

	shipsIcon := IconCell(Icons.Ships, lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorOnline)))
	containersIcon := IconCell(Icons.Containers, lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorContainers)))

	parts := []string{
		fmt.Sprintf("%s %s SHIPS", shipsIcon, FormatContainers(t.Summary.OnlineShips, t.Summary.TotalShips)),
		fmt.Sprintf("%s %s CONTAINERS", containersIcon, FormatContainers(t.Summary.RunningContainers, t.Summary.TotalContainers)),
		t.Now.Format("03:04 PM"),
		t.Now.Format("2006-01-02"),
	}
	summaryLine := SummaryStyle.Render(strings.Join(parts, "  |  "))
	title := TitleStyle.Render(TruncatePlain("FLEETGLANCE", t.Width))

	titleWidth := lipgloss.Width(title)
	if titleWidth >= t.Width {
		return TopBarStyle.Width(t.Width).MaxWidth(t.Width).Render(title)
	}

	maxSummaryWidth := t.Width - titleWidth - 1
	if maxSummaryWidth <= 0 {
		return TopBarStyle.Width(t.Width).MaxWidth(t.Width).Render(title)
	}

	if lipgloss.Width(summaryLine) > maxSummaryWidth {
		summaryLine = SummaryStyle.MaxWidth(maxSummaryWidth).Render(summaryLine)
	}

	spaces := max(t.Width-titleWidth-lipgloss.Width(summaryLine), 0)
	top := title + strings.Repeat(" ", spaces) + summaryLine

	return TopBarStyle.Width(t.Width).MaxWidth(t.Width).Render(top)
}
