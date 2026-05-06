package main

import (
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
