package ui

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	maxPaneColumns = 4
)

const (
	preferredPaneWidth = 40
	minPaneWidth       = 17
	detailedPaneHeight = 17
	compactPaneHeight  = 12
)

const (
	screenPaddingX = 4
	screenPaddingY = 2
	gridGapX       = 2
	gridGapY       = 1
	headerGridGap  = 2
	panePaddingX   = 2
	panePaddingY   = 1
)

const (
	iconColWidth       = 2
	metricLabelWidth   = 8
	metricValueWidth   = 8
	metricPercentWidth = 6
	metricBarGap       = 2
	minProgressBarLen  = 3
)

const (
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
	gridHeight := max(contentHeight-headerHeight-headerGridGap, 1)
	layout := chooseGridLayout(contentWidth, gridHeight, len(m.shipNames))
	grid := m.renderGrid(contentWidth, layout)

	contentParts := []string{header}
	if grid != "" && gridHeight > 0 {
		for range headerGridGap {
			contentParts = append(contentParts, "")
		}
		contentParts = append(contentParts, grid)
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

	shipsIcon := iconCell(icons.ships, lipgloss.NewStyle().
		Background(lipgloss.Color(colorBackground)).
		Foreground(lipgloss.Color(colorOnline)))
	containersIcon := iconCell(icons.containers, lipgloss.NewStyle().
		Background(lipgloss.Color(colorBackground)).
		Foreground(lipgloss.Color(colorContainers)))

	parts := []string{
		fmt.Sprintf("%s %s SHIPS", shipsIcon, formatContainers(summary.onlineShips, summary.totalShips)),
		fmt.Sprintf("%s %s CONTAINERS", containersIcon, formatContainers(summary.runningContainers, summary.totalContainers)),
		m.now.Format("03:04 PM"),
		m.now.Format("2006-01-02"),
	}
	summaryLine := summaryStyle.Render(strings.Join(parts, "  |  "))

	title := titleStyle.Render("FLEETGLANCE")

	maxSummaryWidth := max(width-lipgloss.Width(title)-2, 0)
	if maxSummaryWidth == 0 {
		summaryLine = ""
	} else if lipgloss.Width(summaryLine) > maxSummaryWidth {
		summaryLine = summaryStyle.MaxWidth(maxSummaryWidth).Render(summaryLine)
	}

	spaces := max(width-lipgloss.Width(title)-lipgloss.Width(summaryLine), 1)
	top := title + strings.Repeat(" ", spaces) + summaryLine

	return topBarStyle.Width(width).Render(top)
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
		for offset, name := range m.shipNames[i:end] {
			panes = append(panes, renderPane(m.ships[name], i+offset, layout.paneWidth, layout.paneHeight, layout.compact))
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

func renderPane(ship shipState, index int, width int, outerHeight int, compact bool) string {
	accent := shipAccentByIndex(index)
	paddingY := paneVerticalPadding(outerHeight, compact)
	contentWidth := max(width-2-(panePaddingX*2), 1)
	contentHeight := max(outerHeight-2-(paddingY*2), 1)

	lines := renderPaneLines(ship, contentWidth, accent, compact)
	if compact && len(lines)+1 <= contentHeight {
		lines = addPaneHeaderGap(lines, contentWidth)
	}
	if len(lines) > contentHeight {
		lines = lines[:contentHeight]
	}

	return panelStyle.
		Width(max(width-2, 1)).
		Height(max(outerHeight-2, 1)).
		MaxHeight(max(outerHeight, 1)).
		Padding(paddingY, panePaddingX).
		BorderForeground(lipgloss.Color(accent)).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

func addPaneHeaderGap(lines []string, width int) []string {
	if len(lines) == 0 {
		return lines
	}

	spaced := make([]string, 0, len(lines)+1)
	spaced = append(spaced, lines[0], fillLine("", width))
	spaced = append(spaced, lines[1:]...)

	return spaced
}

func paneVerticalPadding(outerHeight int, compact bool) int {
	if compact && outerHeight < compactPaneHeight {
		return 0
	}

	return panePaddingY
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
	icon := iconCell(icons.ship, lipgloss.NewStyle().
		Background(lipgloss.Color(colorPanelBackground)).
		Foreground(lipgloss.Color(accent)))
	name := headerStyle.
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
		value = valueStyle.Render(formatUptime(ship.telemetry.UptimeSeconds))
	}

	return renderValueRow(icons.uptime, "UPTIME", value, width)
}

func renderContainersRow(ship shipState, width int) string {
	value := dimValueStyle.Render("--/--")
	if ship.telemetry != nil && ship.telemetry.Containers != nil {
		value = valueStyle.Render(formatContainers(ship.telemetry.Containers.Running, ship.telemetry.Containers.Total))
	}

	return renderValueRow(icons.containers, "CONT", value, width)
}

func renderValueRow(icon string, label string, value string, width int) string {
	left := rowLabel(icon, label, valueLabelWidth(label, width))
	right := valueCell(value, valueWidth(width, lipgloss.Width(left), lipgloss.Width(value)))
	return fillLine(joinLeftRight(left, right, width), width)
}

func renderMetricRow(icon string, label string, value *float64, color string, width int) string {
	if value == nil {
		left := rowLabel(icon, label, metricLabelWidth)
		right := valueCell(dimValueStyle.Render("--"), valueWidth(width, lipgloss.Width(left), 2))
		return fillLine(joinLeftRight(left, right, width), width)
	}

	formatted := formatPercent(*value)
	labelWidth := metricLabelColumnWidth(label, width)
	leftWidth := rowLabelWidth(labelWidth)
	left := rowLabel(icon, label, labelWidth)
	left = rowStyle.Inline(true).Width(leftWidth).MaxWidth(leftWidth).Render(left)
	percent := valueCell(valueStyle.Render(formatted), metricPercentWidth)
	barWidth := width - leftWidth - metricBarGap - 1 - metricPercentWidth
	if barWidth < minProgressBarLen {
		return fillLine(joinLeftRight(left, percent, width), width)
	}

	bar := renderProgressBar(*value, barWidth, color)
	right := strings.Repeat(" ", metricBarGap) + bar + " " + percent
	return fillLine(left+right, width)
}

func rowLabel(icon string, label string, width int) string {
	iconText := iconCell(icon, neutralIconStyle)
	labelText := labelStyle.Inline(true).Width(width).MaxWidth(width).Render(label)
	return iconText + " " + labelText
}

func iconCell(icon string, style lipgloss.Style) string {
	return style.
		Inline(true).
		Width(iconColWidth).
		MaxWidth(iconColWidth).
		Align(lipgloss.Center).
		Render(icon)
}

func joinPanes(panes []string) string {
	if len(panes) == 0 {
		return ""
	}

	gap := contentStyle.Render(strings.Repeat(" ", gridGapX))
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
		Background(lipgloss.Color(color)).
		Foreground(lipgloss.Color(color)).
		Render(strings.Repeat("█", filled))
	empty := lipgloss.NewStyle().
		Background(lipgloss.Color(colorChartBackground)).
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
	if !ship.seen {
		return mutedValueStyle
	}

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
	return rowStyle.Width(width).MaxWidth(width).Render(line)
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

func valueLabelWidth(label string, rowWidth int) int {
	if rowWidth >= iconColWidth+1+metricLabelWidth+1+metricValueWidth {
		return metricLabelWidth
	}

	return max(min(metricLabelWidth, max(lipgloss.Width(label), 1)), 1)
}

func metricLabelColumnWidth(label string, rowWidth int) int {
	fullLeftWidth := rowLabelWidth(metricLabelWidth)
	if rowWidth-fullLeftWidth-metricBarGap-1-metricPercentWidth >= minProgressBarLen {
		return metricLabelWidth
	}

	return max(min(metricLabelWidth, lipgloss.Width(label)), 1)
}

func rowLabelWidth(labelWidth int) int {
	return iconColWidth + 1 + labelWidth
}

func valueWidth(rowWidth int, leftWidth int, valueTextWidth int) int {
	available := rowWidth - leftWidth - 1
	if available <= 0 {
		return 1
	}

	if available >= metricValueWidth {
		return metricValueWidth
	}

	return max(min(available, valueTextWidth), 1)
}

func valueCell(value string, width int) string {
	return rowStyle.
		Inline(true).
		Width(max(width, 1)).
		MaxWidth(max(width, 1)).
		Align(lipgloss.Right).
		Render(value)
}
