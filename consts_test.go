package chrono_test

import (
	"testing"

	"github.com/go-chrono/chrono"
)

func TestWeekday_String(t *testing.T) {
	for _, tt := range []struct {
		day      chrono.Weekday
		expected string
	}{
		{
			day:      chrono.Weekday(0),
			expected: "Monday",
		},
		{
			day:      chrono.Weekday(6),
			expected: "Sunday",
		},
		{
			day:      chrono.Weekday(7),
			expected: "%!Weekday(7)",
		},
	} {
		t.Run(tt.expected, func(t *testing.T) {
			if out := tt.day.String(); out != tt.expected {
				t.Fatalf("stringified day = %s, want %s", out, tt.expected)
			}
		})
	}
}

func TestMonth_String(t *testing.T) {
	for _, tt := range []struct {
		day      chrono.Month
		expected string
	}{
		{
			day:      chrono.Month(0),
			expected: "%!Month(0)",
		},
		{
			day:      chrono.Month(1),
			expected: "January",
		},
		{
			day:      chrono.Month(12),
			expected: "December",
		},
		{
			day:      chrono.Month(13),
			expected: "%!Month(13)",
		},
	} {
		t.Run(tt.expected, func(t *testing.T) {
			if out := tt.day.String(); out != tt.expected {
				t.Fatalf("stringified month = %s, want %s", out, tt.expected)
			}
		})
	}
}
