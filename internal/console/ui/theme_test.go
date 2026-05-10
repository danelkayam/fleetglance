package ui

import (
	"os"
	"strings"
	"testing"
)

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

func TestProgressBarEmptyAreaUsesSpaces(t *testing.T) {
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
