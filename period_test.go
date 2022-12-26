package chrono_test

import (
	"strings"
	"testing"

	"github.com/go-chrono/chrono"
)

func TestPeriod_Format(t *testing.T) {
	for _, tt := range []struct {
		name     string
		input    chrono.Period
		expected string
	}{
		{
			name:     "YMWD",
			input:    chrono.Period{Years: 14, Months: 8, Weeks: 4, Days: 2},
			expected: "P14Y8M4W2D",
		},
		{
			name:     "YMW",
			input:    chrono.Period{Years: 14, Months: 8, Weeks: 4},
			expected: "P14Y8M4W",
		},
		{
			name:     "YM",
			input:    chrono.Period{Years: 14, Months: 8},
			expected: "P14Y8M",
		},
		{
			name:     "Y",
			input:    chrono.Period{Years: 14},
			expected: "P14Y",
		},
		{
			name:     "YWD",
			input:    chrono.Period{Years: 14, Weeks: 4, Days: 2},
			expected: "P14Y4W2D",
		},
		{
			name:     "YW",
			input:    chrono.Period{Years: 14, Weeks: 4},
			expected: "P14Y4W",
		},
		{
			name:     "YMD",
			input:    chrono.Period{Years: 14, Months: 8, Days: 2},
			expected: "P14Y8M2D",
		},
		{
			name:     "YD",
			input:    chrono.Period{Years: 14, Days: 2},
			expected: "P14Y2D",
		},
		{
			name:     "MWD",
			input:    chrono.Period{Months: 8, Weeks: 4, Days: 2},
			expected: "P8M4W2D",
		},
		{
			name:     "MW",
			input:    chrono.Period{Months: 8, Weeks: 4},
			expected: "P8M4W",
		},
		{
			name:     "M",
			input:    chrono.Period{Months: 8},
			expected: "P8M",
		},
		{
			name:     "MD",
			input:    chrono.Period{Months: 8, Days: 2},
			expected: "P8M2D",
		},
		{
			name:     "WD",
			input:    chrono.Period{Weeks: 4, Days: 2},
			expected: "P4W2D",
		},
		{
			name:     "D",
			input:    chrono.Period{Days: 2},
			expected: "P2D",
		},
		{
			name:     "empty",
			input:    chrono.Period{},
			expected: "P0D",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if out := tt.input.Format(); out != tt.expected {
				t.Errorf("formatted period = %s, want %s", out, tt.expected)
			}
		})
	}
}

func TestPeriod_Parse(t *testing.T) {
	for _, tt := range []struct {
		name     string
		input    string
		expected chrono.Period
	}{
		{
			name:     "valid YMWD",
			input:    "P14Y8M5W2D",
			expected: chrono.Period{Years: 14, Months: 8, Weeks: 5, Days: 2},
		},
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
			name:     "valid MWD",
			input:    "P8M5W2D",
			expected: chrono.Period{Months: 8, Weeks: 5, Days: 2},
		},
		{
			name:     "valid MW",
			input:    "P8M5W",
			expected: chrono.Period{Months: 8, Weeks: 5},
		},
		{
			name:     "valid MD",
			input:    "P8M2D",
			expected: chrono.Period{Months: 8, Days: 2},
		},
		{
			name:     "valid WD",
			input:    "P5W2D",
			expected: chrono.Period{Weeks: 5, Days: 2},
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
				t.Errorf("failed to parse period: %v", err)
			} else if !p.Equal(tt.expected) {
				t.Errorf("parsed period = %v, want %v", p, tt.expected)
			}
		})
	}
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
			run := func() {
				if p, d, err := chrono.ParseDuration(tt.input); err != nil {
					t.Errorf("failed to parse period & duration: %v", err)
				} else if !p.Equal(tt.period) {
					t.Errorf("parsed period = %v, want %v", p, tt.period)
				} else if d.Compare(tt.duration) != 0 {
					t.Errorf("parsed duration = %v, want %v", d, tt.duration)
				}
			}

			t.Run("dots", func(t *testing.T) {
				run()
			})

			t.Run("commas", func(t *testing.T) {
				tt.input = strings.ReplaceAll(tt.input, ".", ",")
				run()
			})
		})
	}

	t.Run("invalid strings", func(t *testing.T) {
		for _, tt := range []string{
			"P",
			"PT",
		} {
			if _, _, err := chrono.ParseDuration(tt); err == nil {
				t.Error("expecting error but got nil")
			}
		}
	})
}
