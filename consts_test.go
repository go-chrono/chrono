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
			day:      chrono.Weekday(1),
			expected: "Monday",
		},
		{
			day:      chrono.Weekday(7),
			expected: "Sunday",
		},
		{
			day:      chrono.Weekday(0),
			expected: "%!Weekday(0)",
		},
		{
			day:      chrono.Weekday(8),
			expected: "%!Weekday(8)",
		},
	} {
		t.Run(tt.expected, func(t *testing.T) {
			if out := tt.day.String(); out != tt.expected {
				t.Errorf("stringified day = %s, want %s", out, tt.expected)
			}
		})
	}
}

func TestMonth_String(t *testing.T) {
	for _, tt := range []struct {
		month    chrono.Month
		expected string
	}{
		{
			month:    chrono.Month(0),
			expected: "%!Month(0)",
		},
		{
			month:    chrono.Month(1),
			expected: "January",
		},
		{
			month:    chrono.Month(12),
			expected: "December",
		},
		{
			month:    chrono.Month(13),
			expected: "%!Month(13)",
		},
	} {
		t.Run(tt.expected, func(t *testing.T) {
			if out := tt.month.String(); out != tt.expected {
				t.Errorf("stringified month = %s, want %s", out, tt.expected)
			}
		})
	}
}
