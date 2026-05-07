package ui

import "testing"

func TestMainSurfacesUseBaseBackground(t *testing.T) {
	if colorPanelBackground != colorBackground {
		t.Fatalf("panel background = %s, want %s", colorPanelBackground, colorBackground)
	}
	if colorHeaderBackground != colorBackground {
		t.Fatalf("header background = %s, want %s", colorHeaderBackground, colorBackground)
	}
}

func TestProgressBarEmptyAreaUsesBackgroundOnly(t *testing.T) {
	bar := renderProgressBar(25, 8, colorCPU)

	if containsRune(bar, '░') {
		t.Fatal("empty progress bar area should not render a shaded glyph")
	}
}

func containsRune(value string, want rune) bool {
	for _, got := range value {
		if got == want {
			return true
		}
	}

	return false
}
