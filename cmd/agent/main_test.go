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

func TestLoadParamsUsesDefaultsWithoutDotEnv(t *testing.T) {
	withoutDotEnv(t)
	unsetAgentEnv(t)

	params, err := loadParams()
	if err != nil {
		t.Fatalf("load params: %v", err)
	}

	if params.Port != 9800 {
		t.Fatalf("expected default port 9800, got %d", params.Port)
	}
	if params.Debug {
		t.Fatal("expected default debug false")
	}
	if params.LogFormat != "console" {
		t.Fatalf("expected default log format %q, got %q", "console", params.LogFormat)
	}
}

func TestLoadParamsUsesEnv(t *testing.T) {
	withoutDotEnv(t)
	unsetAgentEnv(t)
	t.Setenv("PORT", "1234")
	t.Setenv("DEBUG", "true")
	t.Setenv("LOG_FORMAT", "json")

	params, err := loadParams()
	if err != nil {
		t.Fatalf("load params: %v", err)
	}

	if params.Port != 1234 {
		t.Fatalf("expected port 1234, got %d", params.Port)
	}
	if !params.Debug {
		t.Fatal("expected debug true")
	}
	if params.LogFormat != "json" {
		t.Fatalf("expected log format %q, got %q", "json", params.LogFormat)
	}
}

func TestLoadParamsRejectsInvalidEnv(t *testing.T) {
	tests := []struct {
		name    string
		envKey  string
		envVal  string
		wantErr string
	}{
		{
			name:    "port below range",
			envKey:  "PORT",
			envVal:  "0",
			wantErr: "validate params",
		},
		{
			name:    "port above range",
			envKey:  "PORT",
			envVal:  "65536",
			wantErr: "validate params",
		},
		{
			name:    "port not numeric",
			envKey:  "PORT",
			envVal:  "not-a-number",
			wantErr: "parse env",
		},
		{
			name:    "invalid log format",
			envKey:  "LOG_FORMAT",
			envVal:  "text",
			wantErr: "validate params",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withoutDotEnv(t)
			unsetAgentEnv(t)
			t.Setenv(tt.envKey, tt.envVal)

			_, err := loadParams()
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func TestCommandPrintsVersionAndSkipsRuntime(t *testing.T) {
	withoutDotEnv(t)
	unsetAgentEnv(t)
	t.Setenv("PORT", "not-a-number")

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
	runErr := errors.New("runtime should not start")
	cmd := newCommand(&out, &bytes.Buffer{}, func(*Params) error {
		return runErr
	})

	if err := cmd.Run(context.Background(), []string{"fleetglance-agent", "--version"}); err != nil {
		t.Fatalf("run command: %v", err)
	}

	want := "Fleetglance Agent\nversion=v1.2.3\ncommit=abc1234\nbuilt_at=2026-05-11T10:00:00Z\n"
	if out.String() != want {
		t.Fatalf("expected version output %q, got %q", want, out.String())
	}
}

func TestCommandLoadsParamsBeforeRuntime(t *testing.T) {
	withoutDotEnv(t)
	unsetAgentEnv(t)
	t.Setenv("PORT", "1234")
	t.Setenv("DEBUG", "true")
	t.Setenv("LOG_FORMAT", "json")

	var got *Params
	cmd := newCommand(&bytes.Buffer{}, &bytes.Buffer{}, func(params *Params) error {
		got = params
		return nil
	})

	if err := cmd.Run(context.Background(), []string{"fleetglance-agent"}); err != nil {
		t.Fatalf("run command: %v", err)
	}

	if got == nil {
		t.Fatal("expected runtime to receive params")
	}
	if got.Port != 1234 {
		t.Fatalf("expected port 1234, got %d", got.Port)
	}
	if !got.Debug {
		t.Fatal("expected debug true")
	}
	if got.LogFormat != "json" {
		t.Fatalf("expected log format %q, got %q", "json", got.LogFormat)
	}
}

func TestCommandReturnsRuntimeError(t *testing.T) {
	withoutDotEnv(t)
	unsetAgentEnv(t)

	runErr := errors.New("runtime failed")
	cmd := newCommand(&bytes.Buffer{}, &bytes.Buffer{}, func(*Params) error {
		return runErr
	})

	err := cmd.Run(context.Background(), []string{"fleetglance-agent"})
	if !errors.Is(err, runErr) {
		t.Fatalf("expected runtime error %v, got %v", runErr, err)
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

func unsetAgentEnv(t *testing.T) {
	t.Helper()

	for _, key := range []string{"PORT", "DEBUG", "LOG_FORMAT"} {
		value, ok := os.LookupEnv(key)
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("unset %s: %v", key, err)
		}

		t.Cleanup(func() {
			if ok {
				if err := os.Setenv(key, value); err != nil {
					t.Fatalf("restore %s: %v", key, err)
				}
				return
			}

			if err := os.Unsetenv(key); err != nil {
				t.Fatalf("restore unset %s: %v", key, err)
			}
		})
	}
}
