package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadFleet(t *testing.T) {
	path := filepath.Join(t.TempDir(), "fleetglance.yaml")
	data := []byte(`version: 1
ships:
  donnager:
    url: http://donnager:9800
  nostromo:
    url: https://nostromo:9800
`)

	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	fleet, err := LoadFleet(path)
	if err != nil {
		t.Fatalf("load fleet: %v", err)
	}

	if fleet.Version != 1 {
		t.Fatalf("expected version 1, got %d", fleet.Version)
	}

	if got := fleet.Ships["donnager"].URL; got != "http://donnager:9800" {
		t.Fatalf("expected donnager url, got %q", got)
	}

	if got := fleet.Ships["nostromo"].URL; got != "https://nostromo:9800" {
		t.Fatalf("expected nostromo url, got %q", got)
	}
}

func TestLoadFleetInvalidYAML(t *testing.T) {
	path := filepath.Join(t.TempDir(), "fleetglance.yaml")
	data := []byte(`version: 1
ships:
  donnager:
    url: [not closed
`)

	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	_, err := LoadFleet(path)
	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), "parse fleet config") {
		t.Fatalf("expected parse error, got %q", err.Error())
	}
}
