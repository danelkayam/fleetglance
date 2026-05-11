package components

import (
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func TestRenderIsSingleLine(t *testing.T) {
	topBar := TopBar{
		Width: 80,
		Summary: Summary{
			OnlineShips:       1,
			TotalShips:        2,
			RunningContainers: 3,
			TotalContainers:   4,
		},
		Now: time.Date(2026, 5, 11, 9, 0, 0, 0, time.UTC),
	}.Render()

	if height := lipgloss.Height(topBar); height != 1 {
		t.Fatalf("top bar height = %d, want 1", height)
	}
	if strings.Contains(topBar, "CONSOLE") {
		t.Fatalf("top bar should not include subtitle; got %q", topBar)
	}
}
