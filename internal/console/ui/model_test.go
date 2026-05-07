package ui

import (
	"errors"
	"fleetglance/internal/console/config"
	"fleetglance/internal/console/engine"
	"fleetglance/internal/protocol"
	"reflect"
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

func TestNewModelCapsShipsAtVersionLimit(t *testing.T) {
	model := NewModel(testFleetWithShips(9))

	if len(model.shipNames) != MaxShips {
		t.Fatalf("ship count = %d, want %d", len(model.shipNames), MaxShips)
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

func TestCompactLayoutKeepsSevenShipsInTwoRowsAt80x24(t *testing.T) {
	layout := chooseGridLayout(76, 19, 7)

	if layout.columns != 4 {
		t.Fatalf("columns = %d, want 4", layout.columns)
	}
	if layout.paneHeight != compactPaneHeight {
		t.Fatalf("pane height = %d, want %d", layout.paneHeight, compactPaneHeight)
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
