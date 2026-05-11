package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	panePaddingX = 1
	panePaddingY = 1
)

type Pane struct {
	Ship    Ship
	Index   int
	Width   int
	Height  int
	Compact bool
}

func (p Pane) Render() string {
	accent := ShipAccentByIndex(p.Index)
	paddingY := paneVerticalPadding(p.Height, p.Compact)
	contentWidth := max(p.Width-2-(panePaddingX*2), 1)
	contentHeight := max(p.Height-2-(paddingY*2), 1)

	lines := renderLines(p.Ship, contentWidth, accent, p.Compact)
	if p.Compact && len(lines)+1 <= contentHeight {
		lines = addHeaderGap(lines, contentWidth)
	}
	if len(lines) > contentHeight {
		lines = lines[:contentHeight]
	}

	return PanelStyle.
		Width(max(p.Width-2, 1)).
		Height(max(p.Height-2, 1)).
		MaxHeight(max(p.Height, 1)).
		Padding(paddingY, panePaddingX).
		BorderForeground(lipgloss.Color(accent)).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

func addHeaderGap(lines []string, width int) []string {
	if len(lines) == 0 {
		return lines
	}

	spaced := make([]string, 0, len(lines)+1)
	spaced = append(spaced, lines[0], RowStyle.Width(width).MaxWidth(width).Render(""))
	spaced = append(spaced, lines[1:]...)

	return spaced
}

func paneVerticalPadding(outerHeight int, compact bool) int {
	if compact && outerHeight < compactPaneHeight {
		return 0
	}

	return panePaddingY
}

func renderLines(ship Ship, width int, accent string, compact bool) []string {
	if compact {
		return []string{
			renderHeader(ship, width, accent),
			renderStatusRow(ship, width),
			renderMetricRow(Icons.CPU, "CPU", metricValue(ship, "cpu"), ColorCPU, width),
			renderMetricRow(Icons.RAM, "RAM", metricValue(ship, "ram"), ColorRAM, width),
			renderMetricRow(Icons.Disk, "DISK", metricValue(ship, "disk"), ColorDisk, width),
			renderUptimeRow(ship, width),
			renderContainersRow(ship, width),
		}
	}

	return []string{
		renderHeader(ship, width, accent),
		renderDivider(width),
		renderStatusRow(ship, width),
		renderDivider(width),
		renderMetricRow(Icons.CPU, "CPU", metricValue(ship, "cpu"), ColorCPU, width),
		renderDivider(width),
		renderMetricRow(Icons.RAM, "RAM", metricValue(ship, "ram"), ColorRAM, width),
		renderDivider(width),
		renderMetricRow(Icons.Disk, "DISK", metricValue(ship, "disk"), ColorDisk, width),
		renderDivider(width),
		renderUptimeRow(ship, width),
		renderDivider(width),
		renderContainersRow(ship, width),
	}
}

func renderHeader(ship Ship, width int, accent string) string {
	icon := IconCell(Icons.Ship, lipgloss.NewStyle().
		Foreground(lipgloss.Color(accent)))
	name := HeaderStyle.
		Foreground(lipgloss.Color(accent)).
		Bold(true).
		Render(ship.Name)
	status := StatusStyle(ship.Status).Render(StatusLabel(ship.Status))

	left := icon + " " + name
	return HeaderStyle.Width(width).Render(joinLeftRight(left, status, width))
}

func renderStatusRow(ship Ship, width int) string {
	return Row{
		Kind:  RowKindValue,
		Icon:  Icons.Status,
		Label: "STATUS",
		Value: StatusStyle(ship.Status).Render(StatusValue(ship.Status)),
		Width: width,
	}.Render()
}

func renderUptimeRow(ship Ship, width int) string {
	value := DimValueStyle.Render("--")
	if ship.Telemetry != nil {
		value = ValueStyle.Render(FormatUptime(ship.Telemetry.UptimeSeconds))
	}

	return Row{
		Kind:  RowKindValue,
		Icon:  Icons.Uptime,
		Label: "UPTIME",
		Value: value,
		Width: width,
	}.Render()
}

func renderContainersRow(ship Ship, width int) string {
	value := DimValueStyle.Render("--/--")
	if ship.Telemetry != nil && ship.Telemetry.Containers != nil {
		value = ValueStyle.Render(FormatContainers(ship.Telemetry.Containers.Running, ship.Telemetry.Containers.Total))
	}

	return Row{
		Kind:  RowKindValue,
		Icon:  Icons.Containers,
		Label: "CONT",
		Value: value,
		Width: width,
	}.Render()
}

func renderMetricRow(icon string, label string, value *float64, color string, width int) string {
	return Row{
		Kind:   RowKindMetric,
		Icon:   icon,
		Label:  label,
		Metric: value,
		Color:  color,
		Width:  width,
	}.Render()
}

func renderDivider(width int) string {
	return DividerStyle.Render(strings.Repeat("─", max(width, 0)))
}

func metricValue(ship Ship, metric string) *float64 {
	if ship.Telemetry == nil {
		return nil
	}

	switch metric {
	case "cpu":
		if ship.Telemetry.CPU == nil {
			return nil
		}
		return &ship.Telemetry.CPU.UsagePercent
	case "ram":
		if ship.Telemetry.Memory == nil {
			return nil
		}
		return &ship.Telemetry.Memory.UsagePercent
	case "disk":
		if ship.Telemetry.Storage == nil {
			return nil
		}
		return &ship.Telemetry.Storage.UsagePercent
	default:
		return nil
	}
}
