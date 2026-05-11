package components

import (
	"fleetglance/internal/protocol"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestConstrainedCompactPaneShowsAllSummaryRows(t *testing.T) {
	pane := Pane{
		Ship: Ship{
			Name: "ship-a",
			Telemetry: &protocol.Telemetry{
				UptimeSeconds: 90,
				CPU:           &protocol.CPU{UsagePercent: 50},
				Memory:        &protocol.Memory{UsagePercent: 50},
				Storage:       &protocol.Storage{UsagePercent: 50},
				Containers:    &protocol.Containers{Running: 1, Total: 2},
			},
			Status: StatusOnline,
		},
		Width:   27,
		Height:  9,
		Compact: true,
	}.Render()

	for _, want := range []string{"STATUS", "CPU", "RAM", "DISK", "UPTIME", "CONT"} {
		if !strings.Contains(pane, want) {
			t.Fatalf("compact pane should include %q; got %q", want, pane)
		}
	}
}

func TestCompactPaneHeaderGap(t *testing.T) {
	lines := addHeaderGap([]string{"header", "STATUS"}, 10)

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
