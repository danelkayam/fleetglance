package ui

import (
	"fleetglance/internal/console/config"
	"fleetglance/internal/console/engine"
	"fleetglance/internal/protocol"
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	width  int
	height int
	now    time.Time

	shipNames []string
	ships     map[string]shipState
}

type shipState struct {
	name      string
	telemetry *protocol.Telemetry
	err       error
	seen      bool
}

type fleetSummary struct {
	onlineShips       int
	totalShips        int
	runningContainers int
	totalContainers   int
}

func NewModel(fleet *config.Fleet) Model {
	names := []string{}
	if fleet != nil {
		names = make([]string, 0, len(fleet.Ships))
		for name := range fleet.Ships {
			names = append(names, name)
		}
	}
	sort.Strings(names)

	ships := make(map[string]shipState, len(names))
	for _, name := range names {
		ships[name] = shipState{name: name}
	}

	return Model{
		width:     80,
		height:    24,
		now:       time.Now(),
		shipNames: names,
		ships:     ships,
	}
}

func (m Model) Init() tea.Cmd {
	return tick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	case TelemetryEventMsg:
		m.applyEvent(engine.Event(msg))
	case tickMsg:
		m.now = msg.now
		return m, tick()
	}

	return m, nil
}

func (m Model) applyEvent(event engine.Event) {
	if _, ok := m.ships[event.ShipName]; !ok {
		m.shipNames = append(m.shipNames, event.ShipName)
		sort.Strings(m.shipNames)
	}

	m.ships[event.ShipName] = shipState{
		name:      event.ShipName,
		telemetry: event.Telemetry,
		err:       event.Error,
		seen:      true,
	}
}

func (m Model) summary() fleetSummary {
	summary := fleetSummary{totalShips: len(m.shipNames)}

	for _, name := range m.shipNames {
		ship := m.ships[name]
		if !ship.online() {
			continue
		}

		summary.onlineShips++
		if ship.telemetry != nil && ship.telemetry.Containers != nil {
			summary.runningContainers += ship.telemetry.Containers.Running
			summary.totalContainers += ship.telemetry.Containers.Total
		}
	}

	return summary
}

func (s shipState) online() bool {
	return s.seen && s.err == nil && s.telemetry != nil
}

func (s shipState) statusLabel() string {
	if !s.seen {
		return "PENDING"
	}

	if s.online() {
		return "ONLINE"
	}

	return "FAILED"
}

func (s shipState) statusValue() string {
	if !s.seen {
		return "--"
	}

	if s.online() {
		return "OK"
	}

	return "FAILED"
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(now time.Time) tea.Msg {
		return tickMsg{now: now}
	})
}
