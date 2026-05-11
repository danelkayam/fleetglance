package main

import (
	"bytes"
	"context"
	"errors"
	"fleetglance/internal/version"
	"os"
	"strings"
	"testing"
)

func TestLoadParamsRequiresConfigPath(t *testing.T) {
	withoutDotEnv(t)
	t.Setenv("FLEETGLANCE_CONFIG_PATH", "")

	cmd := newCommand(&bytes.Buffer{}, &bytes.Buffer{}, func(*Params) error {
		return errors.New("runtime should not start")
	})
	err := cmd.Run(context.Background(), []string{"fleetglance-console"})
	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), "fleet config path is required") {
		t.Fatalf("expected config path error, got %q", err.Error())
	}
}

func TestLoadParamsUsesEnvConfigPath(t *testing.T) {
	withoutDotEnv(t)
	t.Setenv("FLEETGLANCE_CONFIG_PATH", "env.yaml")

	var params *Params
	cmd := newCommand(&bytes.Buffer{}, &bytes.Buffer{}, func(got *Params) error {
		params = got
		return nil
	})

	if err := cmd.Run(context.Background(), []string{"fleetglance-console"}); err != nil {
		t.Fatalf("run command: %v", err)
	}
	if params == nil {
		t.Fatal("expected runtime to receive params")
	}
	if params.ConfigPath != "env.yaml" {
		t.Fatalf("expected env config path, got %q", params.ConfigPath)
	}
}

func TestLoadParamsFlagOverridesEnvConfigPath(t *testing.T) {
	withoutDotEnv(t)
	t.Setenv("FLEETGLANCE_CONFIG_PATH", "env.yaml")

	var params *Params
	cmd := newCommand(&bytes.Buffer{}, &bytes.Buffer{}, func(got *Params) error {
		params = got
		return nil
	})

	if err := cmd.Run(context.Background(), []string{"fleetglance-console", "-f", "flag.yaml"}); err != nil {
		t.Fatalf("run command: %v", err)
	}
	if params == nil {
		t.Fatal("expected runtime to receive params")
	}
	if params.ConfigPath != "flag.yaml" {
		t.Fatalf("expected flag config path, got %q", params.ConfigPath)
	}
}

func TestCommandPrintsVersionWithoutConfigPath(t *testing.T) {
	withoutDotEnv(t)
	t.Setenv("FLEETGLANCE_CONFIG_PATH", "")

	oldVersion := version.Version
	oldCommit := version.Commit
	oldBuiltAt := version.BuiltAt
	t.Cleanup(func() {
		version.Version = oldVersion
		version.Commit = oldCommit
		version.BuiltAt = oldBuiltAt
	})

	version.Version = "v1.2.3"
	version.Commit = "abc1234"
	version.BuiltAt = "2026-05-11T10:00:00Z"

	var out bytes.Buffer
	cmd := newCommand(&out, &bytes.Buffer{}, func(*Params) error {
		return errors.New("runtime should not start")
	})

	if err := cmd.Run(context.Background(), []string{"fleetglance-console", "--version"}); err != nil {
		t.Fatalf("run command: %v", err)
	}

	want := "Fleetglance Console\nversion=v1.2.3\ncommit=abc1234\nbuilt_at=2026-05-11T10:00:00Z\n"
	if out.String() != want {
		t.Fatalf("expected version output %q, got %q", want, out.String())
	}
}

func TestCommandRejectsInvalidFlag(t *testing.T) {
	withoutDotEnv(t)

	cmd := newCommand(&bytes.Buffer{}, &bytes.Buffer{}, func(*Params) error {
		return errors.New("runtime should not start")
	})
	err := cmd.Run(context.Background(), []string{"fleetglance-console", "--unknown"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestStartConsoleReturnsConfigLoadError(t *testing.T) {
	withoutDotEnv(t)

	err := startConsole(&Params{
		ConfigPath: "missing.yaml",
		LogFormat:  "json",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "load fleet config") {
		t.Fatalf("expected load fleet config error, got %q", err.Error())
	}
}

func withoutDotEnv(t *testing.T) {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}

	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatalf("change working directory: %v", err)
	}

	t.Cleanup(func() {
		if err := os.Chdir(wd); err != nil {
			t.Fatalf("restore working directory: %v", err)
		}
	})
}
