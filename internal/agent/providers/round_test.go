package providers

import "testing"

func TestRoundPercent(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		want  float64
	}{
		{
			name:  "rounds down",
			value: 12.34,
			want:  12.3,
		},
		{
			name:  "rounds up",
			value: 12.35,
			want:  12.4,
		},
		{
			name:  "keeps integer",
			value: 50,
			want:  50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := roundPercent(tt.value); got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
