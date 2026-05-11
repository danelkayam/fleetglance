package components

import (
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	metricLabelWidth    = 8
	metricBarLabelWidth = 4
	metricValueWidth    = 8
	metricPercentWidth  = 6
	metricBarGap        = 1
	minProgressBarLen   = 3
)

type RowKind int

const (
	RowKindValue RowKind = iota
	RowKindMetric
)

type Row struct {
	Kind   RowKind
	Icon   string
	Label  string
	Value  string
	Metric *float64
	Color  string
	Width  int
}

func (r Row) Render() string {
	if r.Kind == RowKindMetric {
		return renderMetric(r.Icon, r.Label, r.Metric, r.Color, r.Width)
	}

	return renderValue(r.Icon, r.Label, r.Value, r.Width)
}

func renderValue(icon string, label string, value string, width int) string {
	left := rowLabel(icon, label, valueLabelWidth(label, width))
	right := valueCell(value, valueWidth(width, lipgloss.Width(left), lipgloss.Width(value)))
	return fillLine(joinLeftRight(left, right, width), width)
}

func renderMetric(icon string, label string, value *float64, color string, width int) string {
	left := rowLabel(icon, label, metricBarLabelWidth)
	leftWidth := lipgloss.Width(left)

	if value == nil {
		right := valueCell(DimValueStyle.Render("--"), valueWidth(width, leftWidth, 2))
		return fillLine(joinLeftRight(left, right, width), width)
	}

	formatted := FormatPercent(*value)
	percent := valueCell(ValueStyle.Render(formatted), metricPercentWidth)
	barWidth := width - leftWidth - metricBarGap - 1 - metricPercentWidth
	if barWidth < minProgressBarLen {
		return fillLine(joinLeftRight(left, percent, width), width)
	}

	bar := renderProgressBar(*value, barWidth, color)
	right := strings.Repeat(" ", metricBarGap) + bar + " " + percent
	return fillLine(left+right, width)
}

func rowLabel(icon string, label string, width int) string {
	iconText := IconCell(icon, NeutralIconStyle)
	labelText := LabelStyle.Inline(true).Render(PadRight(label, width))
	return iconText + " " + labelText
}

func renderProgressBar(value float64, width int, color string) string {
	if width <= 0 {
		return ""
	}

	if math.IsNaN(value) || math.IsInf(value, 0) {
		value = 0
	}

	if value < 0 {
		value = 0
	}
	if value > 100 {
		value = 100
	}
	filled := int(math.Round((value / 100) * float64(width)))
	if value > 0 && filled == 0 {
		filled = 1
	}
	if filled > width {
		filled = width
	}

	fill := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Render(strings.Repeat("█", filled))
	empty := lipgloss.NewStyle().
		Render(strings.Repeat(" ", width-filled))

	return fill + empty
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
	return RowStyle.Width(width).MaxWidth(width).Render(line)
}

func valueLabelWidth(label string, rowWidth int) int {
	if rowWidth >= IconColWidth+1+metricLabelWidth+1+metricValueWidth {
		return metricLabelWidth
	}

	return max(min(metricLabelWidth, max(lipgloss.Width(label), 1)), 1)
}

func rowLabelWidth(labelWidth int) int {
	return IconColWidth + 1 + labelWidth
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
	return RowStyle.
		Inline(true).
		Width(max(width, 1)).
		MaxWidth(max(width, 1)).
		Align(lipgloss.Right).
		Render(value)
}
