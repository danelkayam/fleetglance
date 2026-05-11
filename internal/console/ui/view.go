package ui

import "fleetglance/internal/console/ui/components"

func (m Model) View() string {
	return components.Render(m.componentModel())
}

func (m Model) componentModel() components.Model {
	summary := m.summary()
	ships := make([]components.Ship, 0, len(m.shipNames))
	for _, name := range m.shipNames {
		ship := m.ships[name]
		ships = append(ships, components.Ship{
			Name:      ship.name,
			Telemetry: ship.telemetry,
			Status:    componentStatus(ship),
		})
	}

	return components.Model{
		Width:  m.width,
		Height: m.height,
		Now:    m.now,
		Summary: components.Summary{
			OnlineShips:       summary.onlineShips,
			TotalShips:        summary.totalShips,
			RunningContainers: summary.runningContainers,
			TotalContainers:   summary.totalContainers,
		},
		Ships: ships,
	}
}

func componentStatus(ship shipState) components.Status {
	if !ship.seen {
		return components.StatusPending
	}

	if ship.online() {
		return components.StatusOnline
	}

	return components.StatusFailed
}
