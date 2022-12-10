package chrono_test

import (
	"testing"

	"github.com/go-chrono/chrono"
)

func TestPeriodFormat(t *testing.T) {
	for _, tt := range []struct {
		name     string
		input    chrono.Period
		expected string
	}{
		{
			name:     "positive values",
			input:    chrono.Period{Years: 1, Months: 2, Weeks: 3, Days: 4},
			expected: "P1Y2M3W4D",
		},
		{
			name:     "negative values",
			input:    chrono.Period{Years: -1, Months: -2, Weeks: -3, Days: -4},
			expected: "P1Y2M3W4D",
		},
		{
			name:     "zero value",
			input:    chrono.Period{},
			expected: "P0D",
		},
		{
			name:     "single component",
			input:    chrono.Period{Weeks: 3},
			expected: "P3W",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if out := tt.input.Format(); out != tt.expected {
				t.Fatalf("formatted period = %s, want %s", out, tt.expected)
			}
		})
	}
}

func TestPeriodParse(t *testing.T) {
	for _, tt := range []struct {
		name     string
		input    string
		expected chrono.Period
	}{
		{
			name:     "valid YMD",
			input:    "P14Y8M2D",
			expected: chrono.Period{Years: 14, Months: 8, Days: 2},
		},
		{
			name:     "valid YM",
			input:    "P14Y8M",
			expected: chrono.Period{Years: 14, Months: 8},
		},
		{
			name:     "valid YD",
			input:    "P14Y2D",
			expected: chrono.Period{Years: 14, Days: 2},
		},
		{
			name:     "valid Y",
			input:    "P14Y",
			expected: chrono.Period{Years: 14},
		},
		{
			name:     "valid MD",
			input:    "P8M2D",
			expected: chrono.Period{Months: 8, Days: 2},
		},
		{
			name:     "valid M",
			input:    "P8M",
			expected: chrono.Period{Months: 8},
		},
		{
			name:     "valid W",
			input:    "P6W",
			expected: chrono.Period{Weeks: 6},
		},
		{
			name:     "valid D",
			input:    "P2D",
			expected: chrono.Period{Days: 2},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var p chrono.Period
			if err := p.Parse(tt.input); err != nil {
				t.Fatalf("failed to parse period: %v", err)
			} else if !p.Equal(tt.expected) {
				t.Fatalf("parsed period = %v, want %v", p, tt.expected)
			}
		})
	}

	t.Run("mixing YMD and W", func(t *testing.T) {
		for _, tt := range []string{
			"P1Y2W",
			"P6W2D",
		} {
			t.Run(tt, func(t *testing.T) {
				var p chrono.Period
				if err := p.Parse(tt); err == nil {
					t.Fatalf("expecting error but got nil: %v", err)
				}
			})
		}
	})
}

func TestParseDuration(t *testing.T) {
	for _, tt := range []struct {
		name     string
		input    string
		period   chrono.Period
		duration chrono.Duration
	}{
		{
			name:     "period only",
			input:    "P2.5Y",
			period:   chrono.Period{Years: 2.5},
			duration: chrono.Duration{},
		},
		{
			name:     "duration only",
			input:    "PT6.5H",
			period:   chrono.Period{},
			duration: chrono.DurationOf((6 * chrono.Hour) + (30 * chrono.Minute)),
		},
		{
			name:     "both period and duration",
			input:    "P2.5YT6.5H",
			period:   chrono.Period{Years: 2.5},
			duration: chrono.DurationOf((6 * chrono.Hour) + (30 * chrono.Minute)),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if p, d, err := chrono.ParseDuration(tt.input); err != nil {
				t.Fatalf("failed to parse period & duration: %v", err)
			} else if !p.Equal(tt.period) {
				t.Fatalf("parsed period = %v, want %v", p, tt.period)
			} else if !d.Equal(tt.duration) {
				t.Fatalf("parsed duration = %v, want %v", d, tt.duration)
			}
		})
	}

	t.Run("invalid strings", func(t *testing.T) {
		for _, tt := range []string{
			"P",
			"PT",
		} {
			if _, _, err := chrono.ParseDuration(tt); err == nil {
				t.Fatal("expecting error but got nil")
			}
		}
	})
}
