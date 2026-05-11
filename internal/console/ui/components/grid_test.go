package components

import "testing"

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
