package ui

import (
	"errors"
	"fleetglance/internal/console/config"
	"fleetglance/internal/console/engine"
	"fleetglance/internal/protocol"
	"reflect"
	"regexp"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func TestNewModelIncludesConfiguredShipsSorted(t *testing.T) {
	model := NewModel(&config.Fleet{
		Version: 1,
		Ships: map[string]config.Ship{
			"rocinante": {URL: "http://rocinante:9800"},
			"donnager":  {URL: "http://donnager:9800"},
			"serenity":  {URL: "http://serenity:9800"},
		},
	})

	want := []string{"donnager", "rocinante", "serenity"}
	if !reflect.DeepEqual(model.shipNames, want) {
		t.Fatalf("ship names = %v, want %v", model.shipNames, want)
	}

	summary := model.summary()
	if summary.onlineShips != 0 || summary.totalShips != 3 {
		t.Fatalf("summary ships = %d/%d, want 0/3", summary.onlineShips, summary.totalShips)
	}
}

func TestNewModelDoesNotTruncateShips(t *testing.T) {
	model := NewModel(testFleetWithShips(9))

	if len(model.shipNames) != 9 {
		t.Fatalf("ship count = %d, want %d", len(model.shipNames), 9)
	}
}

func TestPendingShipStatus(t *testing.T) {
	ship := shipState{}

	if got := ship.statusLabel(); got != "PENDING" {
		t.Fatalf("status label = %q, want %q", got, "PENDING")
	}
	if got := ship.statusValue(); got != "--" {
		t.Fatalf("status value = %q, want %q", got, "--")
	}
}

func TestTelemetryEventMarksShipOnline(t *testing.T) {
	model := newTestModel()
	model = updateModel(t, model, TelemetryEventMsg(engine.Event{
		ShipName: "donnager",
		Telemetry: &protocol.Telemetry{
			UptimeSeconds: 5 * 24 * 60 * 60,
			CPU:           &protocol.CPU{UsagePercent: 2.1},
			Memory:        &protocol.Memory{UsagePercent: 23.4},
			Storage:       &protocol.Storage{UsagePercent: 52.8},
			Containers:    &protocol.Containers{Running: 4, Total: 5},
		},
	}))

	ship := model.ships["donnager"]
	if !ship.online() {
		t.Fatal("ship should be online")
	}
	if got := ship.statusValue(); got != "OK" {
		t.Fatalf("status value = %q, want %q", got, "OK")
	}

	summary := model.summary()
	if summary.onlineShips != 1 || summary.totalShips != 2 {
		t.Fatalf("summary ships = %d/%d, want 1/2", summary.onlineShips, summary.totalShips)
	}
	if summary.runningContainers != 4 || summary.totalContainers != 5 {
		t.Fatalf("summary containers = %d/%d, want 4/5", summary.runningContainers, summary.totalContainers)
	}
}

func TestErrorEventMarksShipFailed(t *testing.T) {
	model := newTestModel()
	model = updateModel(t, model, TelemetryEventMsg(engine.Event{
		ShipName:  "donnager",
		Telemetry: &protocol.Telemetry{Containers: &protocol.Containers{Running: 4, Total: 5}},
	}))
	model = updateModel(t, model, TelemetryEventMsg(engine.Event{
		ShipName: "donnager",
		Error:    errors.New("unreachable"),
	}))

	ship := model.ships["donnager"]
	if ship.online() {
		t.Fatal("ship should be failed")
	}
	if got := ship.statusLabel(); got != "FAILED" {
		t.Fatalf("status label = %q, want %q", got, "FAILED")
	}

	summary := model.summary()
	if summary.onlineShips != 0 || summary.runningContainers != 0 || summary.totalContainers != 0 {
		t.Fatalf("summary = %+v, want no online ships or containers", summary)
	}
}

func TestNilTelemetrySectionsRenderWithoutPanic(t *testing.T) {
	model := newTestModel()
	model = updateModel(t, model, TelemetryEventMsg(engine.Event{
		ShipName:  "donnager",
		Telemetry: &protocol.Telemetry{},
	}))

	view := model.View()
	if !strings.Contains(view, "--") {
		t.Fatalf("view should include missing metric placeholders; got %q", view)
	}
	if !strings.Contains(view, "--/--") {
		t.Fatalf("view should include missing containers placeholder; got %q", view)
	}
}

func TestPendingShipsRenderPendingStatus(t *testing.T) {
	model := newTestModel()

	view := model.View()
	if !strings.Contains(view, "PENDING") {
		t.Fatalf("view should include pending status; got %q", view)
	}
	if !strings.Contains(view, "--") {
		t.Fatalf("view should include pending status value; got %q", view)
	}
}

func TestTopBarIsSingleLine(t *testing.T) {
	model := newTestModel()

	topBar := model.renderTopBar(80)
	if height := lipgloss.Height(topBar); height != 1 {
		t.Fatalf("top bar height = %d, want 1", height)
	}
	if strings.Contains(topBar, "CONSOLE") {
		t.Fatalf("top bar should not include subtitle; got %q", topBar)
	}
}

func TestConstrainedCompactPaneShowsAllSummaryRows(t *testing.T) {
	pane := renderPane(shipState{
		name: "ship-a",
		telemetry: &protocol.Telemetry{
			UptimeSeconds: 90,
			CPU:           &protocol.CPU{UsagePercent: 50},
			Memory:        &protocol.Memory{UsagePercent: 50},
			Storage:       &protocol.Storage{UsagePercent: 50},
			Containers:    &protocol.Containers{Running: 1, Total: 2},
		},
		seen: true,
	}, 0, 27, 9, true)

	for _, want := range []string{"STATUS", "CPU", "RAM", "DISK", "UPTIME", "CONT"} {
		if !strings.Contains(pane, want) {
			t.Fatalf("compact pane should include %q; got %q", want, pane)
		}
	}
}

func TestCompactPaneHeaderGap(t *testing.T) {
	lines := addPaneHeaderGap([]string{"header", "STATUS"}, 10)

	if len(lines) != 3 {
		t.Fatalf("line count = %d, want 3", len(lines))
	}
	if width := lipgloss.Width(lines[1]); width != 10 {
		t.Fatalf("header gap width = %d, want 10", width)
	}
	if strings.TrimSpace(stripANSI(lines[1])) != "" {
		t.Fatalf("header gap should be blank; got %q", lines[1])
	}
}

func TestNarrowMetricRowKeepsProgressBar(t *testing.T) {
	value := 50.0
	row := renderMetricRow(icons.cpu, "CPU", &value, colorCPU, 21)

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
		renderMetricRow(icons.cpu, "CPU", &values[0], colorCPU, width),
		renderMetricRow(icons.ram, "RAM", &values[1], colorRAM, width),
		renderMetricRow(icons.disk, "DISK", &values[2], colorDisk, width),
		renderMetricRow(icons.disk, "DISK", &values[3], colorDisk, width),
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

	wantBarWidth := width - rowLabelWidth(metricLabelWidth) - metricBarGap - 1 - metricPercentWidth
	if got := strings.Count(plain, "█"); got != wantBarWidth {
		t.Fatalf("full bar width = %d, want %d in row %q", got, wantBarWidth, rows[3])
	}
}

func TestViewDoesNotExceedModelWidth(t *testing.T) {
	model := NewModel(testFleetWithShips(7))
	model.width = 80
	model.height = 24

	view := model.View()
	for line := range strings.SplitSeq(view, "\n") {
		if width := lipgloss.Width(line); width > model.width {
			t.Fatalf("line width = %d, want <= %d: %q", width, model.width, line)
		}
	}

	if height := lipgloss.Height(view); height > model.height {
		t.Fatalf("view height = %d, want <= %d", height, model.height)
	}
}

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;:]*m`)

func stripANSI(value string) string {
	return ansiPattern.ReplaceAllString(value, "")
}

func TestCompactLayoutKeepsSevenShipsInFourColumnsAt80x24(t *testing.T) {
	layout := chooseGridLayout(76, 16, 7)

	if layout.columns != 4 {
		t.Fatalf("columns = %d, want 4", layout.columns)
	}
	if layout.paneHeight != 7 {
		t.Fatalf("pane height = %d, want %d", layout.paneHeight, 7)
	}
	if !layout.compact {
		t.Fatal("layout should be compact")
	}
}

func TestEightShipViewFitsExternalTerminal(t *testing.T) {
	model := NewModel(testFleetWithShips(8))
	model.width = 132
	model.height = 36

	view := model.View()
	if height := lipgloss.Height(view); height > model.height {
		t.Fatalf("view height = %d, want <= %d", height, model.height)
	}

	for line := range strings.SplitSeq(view, "\n") {
		if width := lipgloss.Width(line); width > model.width {
			t.Fatalf("line width = %d, want <= %d: %q", width, model.width, line)
		}
	}
}

func newTestModel() Model {
	return NewModel(&config.Fleet{
		Version: 1,
		Ships: map[string]config.Ship{
			"donnager":  {URL: "http://donnager:9800"},
			"rocinante": {URL: "http://rocinante:9800"},
		},
	})
}

func updateModel(t *testing.T, model Model, msg tea.Msg) Model {
	t.Helper()

	updated, _ := model.Update(msg)
	next, ok := updated.(Model)
	if !ok {
		t.Fatalf("updated model has type %T, want ui.Model", updated)
	}

	return next
}

func testFleetWithShips(count int) *config.Fleet {
	ships := make(map[string]config.Ship, count)
	for i := range count {
		name := "ship" + string(rune('a'+i))
		ships[name] = config.Ship{URL: "http://" + name + ":9800"}
	}

	return &config.Fleet{
		Version: 1,
		Ships:   ships,
	}
}
