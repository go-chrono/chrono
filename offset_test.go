package chrono_test

import (
	"testing"

	"github.com/go-chrono/chrono"
)

func TestOffsetOf(t *testing.T) {
	for _, tt := range []struct {
		name     string
		hours    int
		mins     int
		expected chrono.Extent
	}{
		{"UTC", 0, 0, 0},
		{"positive hours", 3, 0, 3 * chrono.Hour},
		{"negative hours", -3, 0, -3 * chrono.Hour},
		{"positive minutes", 0, 30, 30 * chrono.Minute},
		{"negative minutes", 0, -30, -30 * chrono.Minute},
		{"positive hours and minutes", 3, 30, 3*chrono.Hour + 30*chrono.Minute},
		{"negative hours and minutes", -3, -30, -3*chrono.Hour - 30*chrono.Minute},
		{"positive hours and negative minutes", 3, -30, 3*chrono.Hour + 30*chrono.Minute},
		{"negative hours and positive minutes", -3, 30, -3*chrono.Hour - 30*chrono.Minute},
	} {
		t.Run(tt.name, func(t *testing.T) {
			offset := chrono.OffsetOf(tt.hours, tt.mins)
			if offset != chrono.Offset(tt.expected) {
				t.Errorf("OffsetOf(%d, %d) = %s, want %s", tt.hours, tt.mins, offset, tt.expected)
			}
		})
	}
}

func TestOffset_String(t *testing.T) {
	for _, tt := range []struct {
		name     string
		value    chrono.Extent
		expected string
	}{
		{"UTC", 0, "Z"},
		{"UTC truncated", 30 * chrono.Second, "Z"},
		{"positive hours", 3 * chrono.Hour, "+03:00"},
		{"negative hours", -3 * chrono.Hour, "-03:00"},
		{"positive hours and minutes", 3*chrono.Hour + 30*chrono.Minute, "+03:30"},
		{"negative hours and minutes", -(3*chrono.Hour + 30*chrono.Minute), "-03:30"},
		{"positive hours and minutes truncated", 3*chrono.Hour + 30*chrono.Minute + 59*chrono.Second, "+03:30"},
		{"negative hours and minutes truncated", -(3*chrono.Hour + 30*chrono.Minute + 59*chrono.Second), "-03:30"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			offset := chrono.Offset(tt.value)
			if out := offset.String(); out != tt.expected {
				t.Errorf("stringified offset = %s, want %s", out, tt.expected)
			}
		})
	}
}
