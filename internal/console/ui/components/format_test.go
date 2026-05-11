package components

import "testing"

func TestFormatPercent(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		want  string
	}{
		{name: "one decimal place", value: 2.12, want: "2.1%"},
		{name: "zero", value: 0, want: "0.0%"},
		{name: "clamps negative", value: -1, want: "0.0%"},
		{name: "clamps over one hundred", value: 118.9, want: "100.0%"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatPercent(tt.value); got != tt.want {
				t.Fatalf("FormatPercent(%v) = %q, want %q", tt.value, got, tt.want)
			}
		})
	}
}

func TestFormatUptime(t *testing.T) {
	tests := []struct {
		name    string
		seconds int64
		want    string
	}{
		{name: "days and hours", seconds: 5*24*60*60 + 14*60*60 + 9*60, want: "5d 14h"},
		{name: "hours and minutes", seconds: 3*60*60 + 12*60, want: "3h 12m"},
		{name: "minutes", seconds: 42 * 60, want: "42m"},
		{name: "less than a minute", seconds: 30, want: "0m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatUptime(tt.seconds); got != tt.want {
				t.Fatalf("FormatUptime(%d) = %q, want %q", tt.seconds, got, tt.want)
			}
		})
	}
}

func TestFormatContainers(t *testing.T) {
	if got := FormatContainers(18, 24); got != "18/24" {
		t.Fatalf("FormatContainers(18, 24) = %q, want %q", got, "18/24")
	}
}
