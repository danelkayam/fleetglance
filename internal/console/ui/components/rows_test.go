package components

import (
	"regexp"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestNarrowMetricRowKeepsProgressBar(t *testing.T) {
	value := 50.0
	row := Row{
		Kind:   RowKindMetric,
		Icon:   Icons.CPU,
		Label:  "CPU",
		Metric: &value,
		Color:  ColorCPU,
		Width:  21,
	}.Render()

	if !strings.Contains(row, "█") {
		t.Fatalf("metric row should include progress bar; got %q", row)
	}
	if !strings.Contains(row, "50.0%") {
		t.Fatalf("metric row should include value; got %q", row)
	}
}

func TestMetricRowsAlignBarsAndReservePercentWidth(t *testing.T) {
	width := 25
	values := []float64{1.8, 11.2, 19.4, 100}
	rows := []string{
		Row{Kind: RowKindMetric, Icon: Icons.CPU, Label: "CPU", Metric: &values[0], Color: ColorCPU, Width: width}.Render(),
		Row{Kind: RowKindMetric, Icon: Icons.RAM, Label: "RAM", Metric: &values[1], Color: ColorRAM, Width: width}.Render(),
		Row{Kind: RowKindMetric, Icon: Icons.Disk, Label: "DISK", Metric: &values[2], Color: ColorDisk, Width: width}.Render(),
		Row{Kind: RowKindMetric, Icon: Icons.Disk, Label: "DISK", Metric: &values[3], Color: ColorDisk, Width: width}.Render(),
	}

	barStart := -1
	for _, row := range rows {
		plain := stripANSI(row)
		before, _, ok := strings.Cut(plain, "█")
		if !ok {
			t.Fatalf("metric row should include progress bar; got %q", row)
		}

		start := lipgloss.Width(before)
		if barStart == -1 {
			barStart = start
		}
		if start != barStart {
			t.Fatalf("bar start = %d, want %d in row %q", start, barStart, row)
		}
	}

	plain := stripANSI(rows[3])
	if !strings.Contains(plain, "100.0%") {
		t.Fatalf("metric row should reserve width for 100.0%%; got %q", rows[3])
	}

	wantBarWidth := width - rowLabelWidth(metricBarLabelWidth) - metricBarGap - 1 - metricPercentWidth
	if got := strings.Count(plain, "█"); got != wantBarWidth {
		t.Fatalf("full bar width = %d, want %d in row %q", got, wantBarWidth, rows[3])
	}
}

func TestProgressBarEmptyAreaUsesSpaces(t *testing.T) {
	bar := renderProgressBar(25, 8, ColorCPU)

	if containsRune(bar, '░') {
		t.Fatal("empty progress bar area should not render a shaded glyph")
	}
}

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;:]*m`)

func stripANSI(value string) string {
	return ansiPattern.ReplaceAllString(value, "")
}

func containsRune(value string, want rune) bool {
	for _, got := range value {
		if got == want {
			return true
		}
	}

	return false
}
