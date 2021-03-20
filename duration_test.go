package chrono_test

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/go-chrono/chrono"
)

func TestDurationOf(t *testing.T) {
	for _, tt := range []struct {
		name string
		of   chrono.Extent
		nsec float64
	}{
		{
			name: "of nanoseconds",
			of:   9000000000000 * chrono.Nanosecond,
			nsec: 9000000000000,
		},
		{
			name: "of microsecond",
			of:   9000000000 * chrono.Microsecond,
			nsec: 9000000000000,
		},
		{
			name: "of milliseconds",
			of:   9000000 * chrono.Millisecond,
			nsec: 9000000000000,
		},
		{
			name: "of seconds",
			of:   9000 * chrono.Second,
			nsec: 9000000000000,
		},
		{
			name: "of minutes",
			of:   150 * chrono.Minute,
			nsec: 9000000000000,
		},
		{
			name: "of hours",
			of:   2 * chrono.Hour,
			nsec: 7200000000000,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			d := chrono.DurationOf(tt.of)
			if out := d.Nanoseconds(); out != tt.nsec {
				t.Fatalf("d.Nanoseconds() = %f, want %f", out, tt.nsec)
			}
		})
	}
}

func TestDurationUnits(t *testing.T) {
	for _, tt := range []struct {
		name     string
		of       chrono.Extent
		f        func(chrono.Duration) float64
		expected float64
	}{
		{
			name:     "nanoseconds",
			of:       9000000000000 * chrono.Nanosecond,
			f:        chrono.Duration.Nanoseconds,
			expected: 9000000000000,
		},
		{
			name:     "microseconds",
			of:       9000000000000 * chrono.Nanosecond,
			f:        chrono.Duration.Microseconds,
			expected: 9000000000,
		},
		{
			name:     "milliseconds",
			of:       9000000000000 * chrono.Nanosecond,
			f:        chrono.Duration.Milliseconds,
			expected: 9000000,
		},
		{
			name:     "seconds",
			of:       9000000000000 * chrono.Nanosecond,
			f:        chrono.Duration.Seconds,
			expected: 9000,
		},
		{
			name:     "minutes",
			of:       9000000000000 * chrono.Nanosecond,
			f:        chrono.Duration.Minutes,
			expected: 150,
		},
		{
			name:     "hours",
			of:       9000000000000 * chrono.Nanosecond,
			f:        chrono.Duration.Hours,
			expected: 2.5,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			d := chrono.DurationOf(tt.of)
			if out := tt.f(d); out != tt.expected {
				t.Fatalf("%v() = %f, want %f",
					runtime.FuncForPC(reflect.ValueOf(tt.f).Pointer()).Name(),
					out, tt.expected)
			}
		})
	}
}

func TestDurationFormat(t *testing.T) {
	for _, tt := range []struct {
		name      string
		of        chrono.Extent
		exclusive []chrono.Designator
		expected  string
	}{
		{
			name:      "exclusive hms",
			of:        1*chrono.Hour + 15*chrono.Minute + 30*chrono.Second + 500*chrono.Millisecond,
			exclusive: []chrono.Designator{chrono.Hours, chrono.Minutes, chrono.Seconds},
			expected:  "PT1H15M30.5S",
		},
		{
			name:      "exclusive hm",
			of:        1*chrono.Hour + 15*chrono.Minute + 30*chrono.Second + 600*chrono.Millisecond,
			exclusive: []chrono.Designator{chrono.Hours, chrono.Minutes},
			expected:  "PT1H15.51M",
		},
		{
			name:      "exclusive hs",
			of:        12*chrono.Hour + 1*chrono.Minute + 30*chrono.Second + 500*chrono.Millisecond,
			exclusive: []chrono.Designator{chrono.Hours, chrono.Seconds},
			expected:  "PT12H90.5S",
		},
		{
			name:      "exclusive h",
			of:        1*chrono.Hour + 30*chrono.Minute + 36*chrono.Second + 36*chrono.Millisecond,
			exclusive: []chrono.Designator{chrono.Hours},
			expected:  "PT1.51001H",
		},
		{
			name:      "exclusive ms",
			of:        1*chrono.Hour + 15*chrono.Minute + 30*chrono.Second + 500*chrono.Millisecond,
			exclusive: []chrono.Designator{chrono.Minutes, chrono.Seconds},
			expected:  "PT75M30.5S",
		},
		{
			name:      "exclusive m",
			of:        1*chrono.Hour + 15*chrono.Minute + 30*chrono.Second + 600*chrono.Millisecond,
			exclusive: []chrono.Designator{chrono.Minutes},
			expected:  "PT75.51M",
		},
		{
			name:      "exclusive s",
			of:        1*chrono.Hour + 15*chrono.Minute + 30*chrono.Second + 500*chrono.Millisecond,
			exclusive: []chrono.Designator{chrono.Seconds},
			expected:  "PT4530.5S",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			d := chrono.DurationOf(tt.of)
			if out := d.Format(tt.exclusive...); out != tt.expected {
				t.Fatalf("formatted duration = %s, want %s", out, tt.expected)
			}
		})
	}
}
