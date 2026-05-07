package ui

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	minPaneWidth          = 17
	preferredPaneWidth    = 40
	maxPaneColumns        = 4
	gridGap               = 2
	gridRowGap            = 1
	detailedPaneHeight    = 15
	compactPaneHeight     = 9
	defaultTerminalWidth  = 80
	defaultTerminalHeight = 24
)

type gridLayout struct {
	columns    int
	paneWidth  int
	paneHeight int
	compact    bool
}

func (m Model) View() string {
	width, height := m.screenSize()
	horizontalPadding := horizontalPadding(width)
	verticalPadding := verticalPadding(height)
	contentWidth := max(width-horizontalPadding*2, 1)
	contentHeight := max(height-verticalPadding*2, 1)

	header := m.renderTopBar(contentWidth)
	headerHeight := lipgloss.Height(header)
	gridHeight := max(contentHeight-headerHeight-1, 1)
	layout := chooseGridLayout(contentWidth, gridHeight, len(m.shipNames))
	grid := m.renderGrid(contentWidth, layout)

	contentParts := []string{header}
	if grid != "" && gridHeight > 0 {
		contentParts = append(contentParts, "", grid)
	}
	content := contentStyle.Render(lipgloss.JoinVertical(lipgloss.Left, contentParts...))
	content = contentStyle.
		Width(contentWidth).
		Height(contentHeight).
		MaxHeight(contentHeight).
		Render(content)
	content = contentStyle.
		Padding(verticalPadding, horizontalPadding).
		Render(content)

	return backgroundStyle.Width(width).Height(height).Render(lipgloss.Place(
		width,
		height,
		lipgloss.Left,
		lipgloss.Top,
		content,
		lipgloss.WithWhitespaceBackground(lipgloss.Color(colorBackground)),
	))
}

func (m Model) renderTopBar(width int) string {
	summary := m.summary()

	shipsIcon := lipgloss.NewStyle().
		Background(lipgloss.Color(colorBackground)).
		Foreground(lipgloss.Color(colorOnline)).
		Render(icons.ships)
	containersIcon := lipgloss.NewStyle().
		Background(lipgloss.Color(colorBackground)).
		Foreground(lipgloss.Color(namedShipAccents["romulus"])).
		Render(icons.containers)

	parts := []string{
		fmt.Sprintf("%s %s SHIPS", shipsIcon, formatContainers(summary.onlineShips, summary.totalShips)),
		fmt.Sprintf("%s %s CONTAINERS", containersIcon, formatContainers(summary.runningContainers, summary.totalContainers)),
		m.now.Format("03:04 PM"),
		m.now.Format("2006-01-02"),
	}
	summaryLine := summaryStyle.Render(strings.Join(parts, "  |  "))

	title := titleStyle.Render("FLEETGLANCE")
	subtitle := subtitleStyle.Render("CONSOLE")

	maxSummaryWidth := max(width-lipgloss.Width(title)-2, 0)
	if maxSummaryWidth == 0 {
		summaryLine = ""
	} else if lipgloss.Width(summaryLine) > maxSummaryWidth {
		summaryLine = summaryStyle.MaxWidth(maxSummaryWidth).Render(summaryLine)
	}

	spaces := max(width-lipgloss.Width(title)-lipgloss.Width(summaryLine), 1)
	top := title + strings.Repeat(" ", spaces) + summaryLine

	return topBarStyle.Width(width).Render(lipgloss.JoinVertical(lipgloss.Left, top, subtitle))
}

func (m Model) renderGrid(width int, layout gridLayout) string {
	if len(m.shipNames) == 0 {
		return ""
	}

	columns := max(layout.columns, 1)
	rows := make([]string, 0, (len(m.shipNames)+columns-1)/columns)

	for i := 0; i < len(m.shipNames); i += columns {
		end := min(i+columns, len(m.shipNames))
		panes := make([]string, 0, end-i)
		for _, name := range m.shipNames[i:end] {
			panes = append(panes, renderPane(m.ships[name], layout.paneWidth, layout.paneHeight, layout.compact))
		}

		rows = append(rows, joinPanes(panes))
	}

	return strings.Join(rows, strings.Repeat("\n", gridRowGap+1))
}

func chooseGridLayout(width int, height int, count int) gridLayout {
	count = min(count, MaxShips)
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
			neededHeight := rows*paneHeight + max(rows-1, 0)*gridRowGap
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
	gaps := max(rows-1, 0) * gridRowGap
	paneHeight := max((height-gaps)/rows, 1)

	return gridLayout{
		columns:    columns,
		paneWidth:  paneWidthForColumns(width, columns),
		paneHeight: paneHeight,
		compact:    true,
	}
}

func renderPane(ship shipState, width int, outerHeight int, compact bool) string {
	accent := shipAccent(ship.name)
	contentWidth := max(width-2, 1)
	contentHeight := max(outerHeight-2, 1)

	lines := renderPaneLines(ship, contentWidth, accent, compact)
	if len(lines) > contentHeight {
		lines = lines[:contentHeight]
	}

	return panelStyle.
		Width(contentWidth).
		Height(contentHeight).
		MaxHeight(contentHeight + 2).
		BorderForeground(lipgloss.Color(accent)).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

func renderPaneLines(ship shipState, width int, accent string, compact bool) []string {
	if compact {
		return []string{
			renderPaneHeader(ship, width, accent),
			renderStatusRow(ship, width),
			renderMetricRow(icons.cpu, "CPU", metricValue(ship, "cpu"), colorCPU, width),
			renderMetricRow(icons.ram, "RAM", metricValue(ship, "ram"), colorRAM, width),
			renderMetricRow(icons.disk, "DISK", metricValue(ship, "disk"), colorDisk, width),
			renderUptimeRow(ship, width),
			renderContainersRow(ship, width),
		}
	}

	return []string{
		renderPaneHeader(ship, width, accent),
		renderDivider(width),
		renderStatusRow(ship, width),
		renderDivider(width),
		renderMetricRow(icons.cpu, "CPU", metricValue(ship, "cpu"), colorCPU, width),
		renderDivider(width),
		renderMetricRow(icons.ram, "RAM", metricValue(ship, "ram"), colorRAM, width),
		renderDivider(width),
		renderMetricRow(icons.disk, "DISK", metricValue(ship, "disk"), colorDisk, width),
		renderDivider(width),
		renderUptimeRow(ship, width),
		renderDivider(width),
		renderContainersRow(ship, width),
	}
}

func renderPaneHeader(ship shipState, width int, accent string) string {
	icon := lipgloss.NewStyle().
		Background(lipgloss.Color(colorBackground)).
		Foreground(lipgloss.Color(accent)).
		Render(shipIcon(ship.name))
	name := lipgloss.NewStyle().
		Background(lipgloss.Color(colorBackground)).
		Foreground(lipgloss.Color(accent)).
		Bold(true).
		Render(ship.name)
	status := statusStyle(ship).Render(ship.statusLabel())

	left := icon + " " + name
	return headerStyle.Width(width).Render(joinLeftRight(left, status, width))
}

func renderStatusRow(ship shipState, width int) string {
	return renderValueRow(icons.status, "STATUS", statusStyle(ship).Render(ship.statusValue()), width)
}

func renderUptimeRow(ship shipState, width int) string {
	value := dimValueStyle.Render("--")
	if ship.telemetry != nil {
		value = mutedValueStyle.Render(formatUptime(ship.telemetry.UptimeSeconds))
	}

	return renderValueRow(icons.uptime, "UPTIME", value, width)
}

func renderContainersRow(ship shipState, width int) string {
	value := dimValueStyle.Render("--/--")
	if ship.telemetry != nil && ship.telemetry.Containers != nil {
		value = valueStyle.Render(formatContainers(ship.telemetry.Containers.Running, ship.telemetry.Containers.Total))
	}

	return renderValueRow(icons.containers, "CONTAINERS", value, width)
}

func renderValueRow(icon string, label string, value string, width int) string {
	left := rowLabel(icon, label, width)
	return fillLine(joinLeftRight(left, value, width), width)
}

func renderMetricRow(icon string, label string, value *float64, color string, width int) string {
	left := rowLabel(icon, label, width)
	if value == nil {
		return fillLine(joinLeftRight(left, dimValueStyle.Render("--"), width), width)
	}

	percent := valueStyle.Render(formatPercent(*value))
	barWidth := width - lipgloss.Width(left) - lipgloss.Width(percent) - 3
	if barWidth < 5 {
		return fillLine(joinLeftRight(left, percent, width), width)
	}

	right := renderProgressBar(*value, barWidth, color) + " " + percent
	return fillLine(joinLeftRight(left, right, width), width)
}

func rowLabel(icon string, label string, rowWidth int) string {
	iconText := neutralIconStyle.Render(icon)
	width := labelWidth(label, rowWidth)
	labelText := labelStyle.Width(width).MaxWidth(width).Render(labelTextForWidth(label, width))
	return iconText + " " + labelText
}

func joinPanes(panes []string) string {
	if len(panes) == 0 {
		return ""
	}

	gap := contentStyle.Render(strings.Repeat(" ", gridGap))
	row := panes[0]
	for _, pane := range panes[1:] {
		row = lipgloss.JoinHorizontal(lipgloss.Top, row, gap, pane)
	}

	return row
}

func renderDivider(width int) string {
	return dividerStyle.Render(strings.Repeat("─", max(width, 0)))
}

func renderProgressBar(value float64, width int, color string) string {
	if width <= 0 {
		return ""
	}

	if math.IsNaN(value) || math.IsInf(value, 0) {
		value = 0
	}

	value = min(max(value, 0), 100)
	filled := int(math.Round((value / 100) * float64(width)))
	if value > 0 && filled == 0 {
		filled = 1
	}
	if filled > width {
		filled = width
	}

	fill := lipgloss.NewStyle().
		Background(lipgloss.Color(colorBackground)).
		Foreground(lipgloss.Color(color)).
		Render(strings.Repeat("█", filled))
	empty := lipgloss.NewStyle().
		Background(lipgloss.Color(colorBackground)).
		Render(strings.Repeat(" ", width-filled))

	return fill + empty
}

func metricValue(ship shipState, metric string) *float64 {
	if ship.telemetry == nil {
		return nil
	}

	switch metric {
	case "cpu":
		if ship.telemetry.CPU == nil {
			return nil
		}
		return &ship.telemetry.CPU.UsagePercent
	case "ram":
		if ship.telemetry.Memory == nil {
			return nil
		}
		return &ship.telemetry.Memory.UsagePercent
	case "disk":
		if ship.telemetry.Storage == nil {
			return nil
		}
		return &ship.telemetry.Storage.UsagePercent
	default:
		return nil
	}
}

func statusStyle(ship shipState) lipgloss.Style {
	if ship.online() {
		return onlineStyle
	}

	return failedStyle
}

func joinLeftRight(left string, right string, width int) string {
	if width <= 0 {
		return ""
	}

	rightWidth := lipgloss.Width(right)
	if rightWidth >= width {
		return lipgloss.NewStyle().MaxWidth(width).Render(right)
	}

	leftWidth := max(width-rightWidth-1, 0)
	left = lipgloss.NewStyle().MaxWidth(leftWidth).Render(left)
	space := max(width-lipgloss.Width(left)-lipgloss.Width(right), 0)

	return left + strings.Repeat(" ", space) + right
}

func fillLine(line string, width int) string {
	return contentStyle.Width(width).MaxWidth(width).Render(line)
}

func (m Model) screenSize() (int, int) {
	width := m.width
	if width <= 0 {
		width = defaultTerminalWidth
	}

	height := m.height
	if height <= 0 {
		height = defaultTerminalHeight
	}

	return max(width, 1), max(height, 1)
}

func horizontalPadding(width int) int {
	switch {
	case width >= 120:
		return 4
	case width >= 70:
		return 2
	default:
		return 1
	}
}

func verticalPadding(height int) int {
	if height >= 18 {
		return 1
	}

	return 0
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

	return max(min((width-gridGap*(columns-1))/columns, preferredPaneWidth), 1)
}

func rowsFor(count int, columns int) int {
	return (count + columns - 1) / columns
}

func labelWidth(label string, rowWidth int) int {
	if rowWidth < 22 {
		return 4
	}

	if rowWidth < 30 {
		return 6
	}

	if len(label) > 6 {
		return 10
	}

	return 6
}

func labelTextForWidth(label string, width int) string {
	if width >= len(label) {
		return label
	}

	switch label {
	case "STATUS":
		return "STAT"
	case "UPTIME":
		return "UP"
	case "CONTAINERS":
		return "CONT"
	default:
		return label[:width]
	}
}
