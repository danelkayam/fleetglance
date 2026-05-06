package main

import (
	"strings"
	"testing"
)

func TestLoadParamsRequiresConfigPath(t *testing.T) {
	t.Setenv("FLEETGLANCE_CONFIG_PATH", "")

	_, err := loadParams(nil)
	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), "fleet config path is required") {
		t.Fatalf("expected config path error, got %q", err.Error())
	}
}

func TestLoadParamsUsesEnvConfigPath(t *testing.T) {
	t.Setenv("FLEETGLANCE_CONFIG_PATH", "env.yaml")

	params, err := loadParams(nil)
	if err != nil {
		t.Fatalf("load params: %v", err)
	}

	if params.ConfigPath != "env.yaml" {
		t.Fatalf("expected env config path, got %q", params.ConfigPath)
	}
}

func TestLoadParamsFlagOverridesEnvConfigPath(t *testing.T) {
	t.Setenv("FLEETGLANCE_CONFIG_PATH", "env.yaml")

	params, err := loadParams([]string{"-f", "flag.yaml"})
	if err != nil {
		t.Fatalf("load params: %v", err)
	}

	if params.ConfigPath != "flag.yaml" {
		t.Fatalf("expected flag config path, got %q", params.ConfigPath)
	}
}
