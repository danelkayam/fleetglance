package ui

import (
	"os"
	"strings"
	"testing"
)

func TestBackgroundLayersAreDistinct(t *testing.T) {
	if colorPanelBackground == colorBackground {
		t.Fatalf("panel background should differ from app background: %s", colorPanelBackground)
	}
	if colorChartBackground == colorBackground {
		t.Fatalf("chart background should differ from app background: %s", colorChartBackground)
	}
	if colorChartBackground == colorPanelBackground {
		t.Fatalf("chart background should differ from panel background: %s", colorChartBackground)
	}
}

func TestThemeDoesNotHardcodeShipAccentNames(t *testing.T) {
	source, err := os.ReadFile("theme.go")
	if err != nil {
		t.Fatalf("read theme source: %v", err)
	}

	for _, name := range []string{"donnager", "rocinante", "romulus", "nostromo", "tycho", "betty", "serenity"} {
		if strings.Contains(string(source), name) {
			t.Fatalf("theme.go should not contain hardcoded ship name %q", name)
		}
	}
}

func TestShipAccentByIndexWrapsPalette(t *testing.T) {
	if got := shipAccentByIndex(0); got != shipAccentColors[0] {
		t.Fatalf("first accent = %s, want %s", got, shipAccentColors[0])
	}

	index := len(shipAccentColors) + 2
	if got := shipAccentByIndex(index); got != shipAccentColors[2] {
		t.Fatalf("wrapped accent = %s, want %s", got, shipAccentColors[2])
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
