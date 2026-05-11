package components

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
	if got := ShipAccentByIndex(0); got != shipAccentColors[0] {
		t.Fatalf("first accent = %s, want %s", got, shipAccentColors[0])
	}

	index := len(shipAccentColors) + 2
	if got := ShipAccentByIndex(index); got != shipAccentColors[2] {
		t.Fatalf("wrapped accent = %s, want %s", got, shipAccentColors[2])
	}
}
