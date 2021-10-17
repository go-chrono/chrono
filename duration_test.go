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
			name: "of positive nanoseconds",
			of:   9000000000000 * chrono.Nanosecond,
			nsec: 9000000000000,
		},
		{
			name: "of positive microsecond",
			of:   9000000000 * chrono.Microsecond,
			nsec: 9000000000000,
		},
		{
			name: "of positive milliseconds",
			of:   9000000 * chrono.Millisecond,
			nsec: 9000000000000,
		},
		{
			name: "of positive seconds",
			of:   9000 * chrono.Second,
			nsec: 9000000000000,
		},
		{
			name: "of positive minutes",
			of:   150 * chrono.Minute,
			nsec: 9000000000000,
		},
		{
			name: "of positive hours",
			of:   2 * chrono.Hour,
			nsec: 7200000000000,
		},
		{
			name: "of negative nanoseconds",
			of:   -9000000000000 * chrono.Nanosecond,
			nsec: -9000000000000,
		},
		{
			name: "of negative microsecond",
			of:   -9000000000 * chrono.Microsecond,
			nsec: -9000000000000,
		},
		{
			name: "of negative milliseconds",
			of:   -9000000 * chrono.Millisecond,
			nsec: -9000000000000,
		},
		{
			name: "of negative seconds",
			of:   -9000 * chrono.Second,
			nsec: -9000000000000,
		},
		{
			name: "of negative minutes",
			of:   -150 * chrono.Minute,
			nsec: -9000000000000,
		},
		{
			name: "of negative hours",
			of:   -2 * chrono.Hour,
			nsec: -7200000000000,
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
			name:     "positive nanoseconds",
			of:       9000000000000 * chrono.Nanosecond,
			f:        chrono.Duration.Nanoseconds,
			expected: 9000000000000,
		},
		{
			name:     "positive microseconds",
			of:       9000000000000 * chrono.Nanosecond,
			f:        chrono.Duration.Microseconds,
			expected: 9000000000,
		},
		{
			name:     "positive milliseconds",
			of:       9000000000000 * chrono.Nanosecond,
			f:        chrono.Duration.Milliseconds,
			expected: 9000000,
		},
		{
			name:     "positive seconds",
			of:       9000000000000 * chrono.Nanosecond,
			f:        chrono.Duration.Seconds,
			expected: 9000,
		},
		{
			name:     "positive minutes",
			of:       9000000000000 * chrono.Nanosecond,
			f:        chrono.Duration.Minutes,
			expected: 150,
		},
		{
			name:     "positive hours",
			of:       9000000000000 * chrono.Nanosecond,
			f:        chrono.Duration.Hours,
			expected: 2.5,
		},
		{
			name:     "negative nanoseconds",
			of:       -9000000000000 * chrono.Nanosecond,
			f:        chrono.Duration.Nanoseconds,
			expected: -9000000000000,
		},
		{
			name:     "negative microseconds",
			of:       -9000000000000 * chrono.Nanosecond,
			f:        chrono.Duration.Microseconds,
			expected: -9000000000,
		},
		{
			name:     "negative milliseconds",
			of:       -9000000000000 * chrono.Nanosecond,
			f:        chrono.Duration.Milliseconds,
			expected: -9000000,
		},
		{
			name:     "negative seconds",
			of:       -9000000000000 * chrono.Nanosecond,
			f:        chrono.Duration.Seconds,
			expected: -9000,
		},
		{
			name:     "negative minutes",
			of:       -9000000000000 * chrono.Nanosecond,
			f:        chrono.Duration.Minutes,
			expected: -150,
		},
		{
			name:     "negative hours",
			of:       -9000000000000 * chrono.Nanosecond,
			f:        chrono.Duration.Hours,
			expected: -2.5,
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
			name:      "default HMS",
			of:        1*chrono.Hour + 15*chrono.Minute + 30*chrono.Second + 500*chrono.Millisecond,
			exclusive: []chrono.Designator{},
			expected:  "PT1H15M30.5S",
		},
		{
			name:      "default HM",
			of:        1*chrono.Hour + 15*chrono.Minute,
			exclusive: []chrono.Designator{},
			expected:  "PT1H15M",
		},
		{
			name:      "default HS",
			of:        12*chrono.Hour + 30*chrono.Second + 500*chrono.Millisecond,
			exclusive: []chrono.Designator{},
			expected:  "PT12H0M30.5S",
		},
		{
			name:      "default H",
			of:        1 * chrono.Hour,
			exclusive: []chrono.Designator{},
			expected:  "PT1H",
		},
		{
			name:      "default MS",
			of:        15*chrono.Minute + 30*chrono.Second + 500*chrono.Millisecond,
			exclusive: []chrono.Designator{},
			expected:  "PT15M30.5S",
		},
		{
			name:      "default M",
			of:        15 * chrono.Minute,
			exclusive: []chrono.Designator{},
			expected:  "PT15M",
		},
		{
			name:      "default S",
			of:        30*chrono.Second + 500*chrono.Millisecond,
			exclusive: []chrono.Designator{},
			expected:  "PT30.5S",
		},
		{
			name:      "exclusive HMS",
			of:        1*chrono.Hour + 15*chrono.Minute + 30*chrono.Second + 500*chrono.Millisecond,
			exclusive: []chrono.Designator{chrono.Hours, chrono.Minutes, chrono.Seconds},
			expected:  "PT1H15M30.5S",
		},
		{
			name:      "exclusive HM",
			of:        1*chrono.Hour + 15*chrono.Minute + 30*chrono.Second + 600*chrono.Millisecond,
			exclusive: []chrono.Designator{chrono.Hours, chrono.Minutes},
			expected:  "PT1H15.51M",
		},
		{
			name:      "exclusive HS",
			of:        12*chrono.Hour + 1*chrono.Minute + 30*chrono.Second + 500*chrono.Millisecond,
			exclusive: []chrono.Designator{chrono.Hours, chrono.Seconds},
			expected:  "PT12H90.5S",
		},
		{
			name:      "exclusive H",
			of:        1*chrono.Hour + 30*chrono.Minute + 36*chrono.Second + 36*chrono.Millisecond,
			exclusive: []chrono.Designator{chrono.Hours},
			expected:  "PT1.51001H",
		},
		{
			name:      "exclusive MS",
			of:        1*chrono.Hour + 15*chrono.Minute + 30*chrono.Second + 500*chrono.Millisecond,
			exclusive: []chrono.Designator{chrono.Minutes, chrono.Seconds},
			expected:  "PT75M30.5S",
		},
		{
			name:      "exclusive M",
			of:        1*chrono.Hour + 15*chrono.Minute + 30*chrono.Second + 600*chrono.Millisecond,
			exclusive: []chrono.Designator{chrono.Minutes},
			expected:  "PT75.51M",
		},
		{
			name:      "exclusive S",
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

func TestParseDuration(t *testing.T) {
	for _, tt := range []struct {
		name     string
		input    string
		expected chrono.Duration
	}{
		{
			name:     "valid HMS integers",
			input:    "PT5H3M1S",
			expected: chrono.DurationOf(5*chrono.Hour + 3*chrono.Minute + 1*chrono.Second),
		},
		{
			name:     "valid HMS floats",
			input:    "PT4.5H3.25M1.1S",
			expected: chrono.DurationOf(chrono.Extent(4.5*float64(chrono.Hour) + 3.25*float64(chrono.Minute) + 1.1*float64(chrono.Second))),
		},
		{
			name:     "valid HM integers",
			input:    "PT5H3M",
			expected: chrono.DurationOf(5*chrono.Hour + 3*chrono.Minute),
		},
		{
			name:     "valid HM floats",
			input:    "PT4.5H3.25M",
			expected: chrono.DurationOf(chrono.Extent(4.5*float64(chrono.Hour) + 3.25*float64(chrono.Minute))),
		},
		{
			name:     "valid HS integers",
			input:    "PT5H1S",
			expected: chrono.DurationOf(5*chrono.Hour + 1*chrono.Second),
		},
		{
			name:     "valid HS floats",
			input:    "PT4.5H1.1S",
			expected: chrono.DurationOf(chrono.Extent(4.5*float64(chrono.Hour) + 1.1*float64(chrono.Second))),
		},
		{
			name:     "valid H integer",
			input:    "PT5H",
			expected: chrono.DurationOf(5 * chrono.Hour),
		},
		{
			name:     "valid H float",
			input:    "PT4.5H",
			expected: chrono.DurationOf(chrono.Extent(4.5 * float64(chrono.Hour))),
		},
		{
			name:     "valid MS integers",
			input:    "PT3M1S",
			expected: chrono.DurationOf(3*chrono.Minute + 1*chrono.Second),
		},
		{
			name:     "valid MS floats",
			input:    "PT3.25M1.1S",
			expected: chrono.DurationOf(chrono.Extent(3.25*float64(chrono.Minute) + 1.1*float64(chrono.Second))),
		},
		{
			name:     "valid M integer",
			input:    "PT3M",
			expected: chrono.DurationOf(3 * chrono.Minute),
		},
		{
			name:     "valid M float",
			input:    "PT3.25M",
			expected: chrono.DurationOf(chrono.Extent(3.25 * float64(chrono.Minute))),
		},
		{
			name:     "valid S integer",
			input:    "PT1S",
			expected: chrono.DurationOf(1 * chrono.Second),
		},
		{
			name:     "valid S float",
			input:    "PT1.1S",
			expected: chrono.DurationOf(chrono.Extent(1.1 * float64(chrono.Second))),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var d chrono.Duration
			err := d.Parse(tt.input)
			if err != nil {
				t.Fatalf("failed to parse duation: %v", err)
			}

			if !d.Equal(tt.expected) {
				t.Fatalf("parsed duration = %v, want %v", d, tt.expected)
			}
		})
	}
}
