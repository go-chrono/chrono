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
				t.Errorf("d.Nanoseconds() = %f, want %f", out, tt.nsec)
			}
		})
	}
}

func TestDuration_units(t *testing.T) {
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
				t.Errorf("%v() = %f, want %f",
					runtime.FuncForPC(reflect.ValueOf(tt.f).Pointer()).Name(),
					out, tt.expected)
			}
		})
	}
}

func TestDuration_Compare(t *testing.T) {
	for _, tt := range []struct {
		name     string
		d        chrono.Duration
		d2       chrono.Duration
		expected int
	}{
		{"seconds less", chrono.DurationOf(1 * chrono.Hour), chrono.DurationOf(2 * chrono.Hour), -1},
		{"seconds more", chrono.DurationOf(2 * chrono.Hour), chrono.DurationOf(1 * chrono.Hour), 1},
		{"nanos less", chrono.DurationOf(1 * chrono.Nanosecond), chrono.DurationOf(2 * chrono.Nanosecond), -1},
		{"nanos more", chrono.DurationOf(2 * chrono.Nanosecond), chrono.DurationOf(1 * chrono.Nanosecond), 1},
		{"equal", chrono.DurationOf(chrono.Minute), chrono.DurationOf(chrono.Minute), 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if v := tt.d.Compare(tt.d2); v != tt.expected {
				t.Errorf("d.Compare(d2) = %d, want %d", v, tt.expected)
			}
		})
	}
}

func TestDuration_Add(t *testing.T) {
	for _, tt := range []struct {
		name     string
		d1       chrono.Duration
		d2       chrono.Duration
		expected chrono.Duration
	}{
		{
			name:     "add seconds component",
			d1:       chrono.DurationOf(1 * chrono.Hour),
			d2:       chrono.DurationOf(1 * chrono.Hour),
			expected: chrono.DurationOf(2 * chrono.Hour),
		},
		{
			name:     "add nanoseconds component",
			d1:       chrono.DurationOf(750 * chrono.Millisecond),
			d2:       chrono.DurationOf(550 * chrono.Millisecond),
			expected: chrono.DurationOf((1 * chrono.Second) + (300 * chrono.Millisecond)),
		},
		{
			name:     "add both components",
			d1:       chrono.DurationOf((1 * chrono.Hour) + (750 * chrono.Millisecond)),
			d2:       chrono.DurationOf((1 * chrono.Hour) + (550 * chrono.Millisecond)),
			expected: chrono.DurationOf((2 * chrono.Hour) + (1 * chrono.Second) + (300 * chrono.Millisecond)),
		},
		{
			name:     "minus seconds component",
			d1:       chrono.DurationOf(2 * chrono.Hour),
			d2:       chrono.DurationOf(-1 * chrono.Hour),
			expected: chrono.DurationOf(1 * chrono.Hour),
		},
		{
			name:     "minus nanoseconds component",
			d1:       chrono.DurationOf(750 * chrono.Millisecond),
			d2:       chrono.DurationOf(-550 * chrono.Millisecond),
			expected: chrono.DurationOf(200 * chrono.Millisecond),
		},
		{
			name:     "minus both components",
			d1:       chrono.DurationOf((2 * chrono.Hour) + (750 * chrono.Millisecond)),
			d2:       chrono.DurationOf(-((1 * chrono.Hour) + (550 * chrono.Millisecond))),
			expected: chrono.DurationOf((1 * chrono.Hour) + (200 * chrono.Millisecond)),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("d1.Add(d2)", func(t *testing.T) {
				if ok := tt.d1.CanAdd(tt.d2); !ok {
					t.Error("d1.CanAdd(d2) = false, want true")
				}

				if added := tt.d1.Add(tt.d2); added.Compare(tt.expected) != 0 {
					t.Errorf("d1.Add(d2) = %v, want %v", added, tt.expected)
				}
			})

			t.Run("d2.Add(d1)", func(t *testing.T) {
				if ok := tt.d2.CanAdd(tt.d1); !ok {
					t.Error("d2.CanAdd(d1) = false, want true")
				}

				if added := tt.d2.Add(tt.d1); added.Compare(tt.expected) != 0 {
					t.Errorf("d2.Add(d1) = %v, want %v", added, tt.expected)
				}
			})
		})
	}

	for _, tt := range []struct {
		name string
		d1   chrono.Duration
		d2   chrono.Duration
	}{
		{
			name: "overflow on seconds component",
			d1:   chrono.MaxDuration(),
			d2:   chrono.DurationOf(1 * chrono.Second),
		},
		{
			name: "overflow on nanoseconds component",
			d1:   chrono.MaxDuration().Add(chrono.DurationOf(-500 * chrono.Millisecond)),
			d2:   chrono.DurationOf(501 * chrono.Millisecond),
		},
		{
			name: "underflow on seconds component",
			d1:   chrono.MinDuration(),
			d2:   chrono.DurationOf(-1 * chrono.Second),
		},
		{
			name: "underflow on nanoseconds component",
			d1:   chrono.MinDuration().Add(chrono.DurationOf(500 * chrono.Millisecond)),
			d2:   chrono.DurationOf(-501 * chrono.Millisecond),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("d1.Add(d2)", func(t *testing.T) {
				if ok := tt.d1.CanAdd(tt.d2); ok {
					t.Error("d1.CanAdd(d2) = true, want false")
				}

				func() {
					defer func() {
						if r := recover(); r == nil {
							t.Error("expecting panic that didn't occur")
						}
					}()

					tt.d1.Add(tt.d2)
				}()
			})

			t.Run("d2.Add(d1)", func(t *testing.T) {
				if ok := tt.d2.CanAdd(tt.d1); ok {
					t.Error("d2.CanAdd(d1) = true, want false")
				}

				func() {
					defer func() {
						if r := recover(); r == nil {
							t.Error("expecting panic that didn't occur")
						}
					}()

					tt.d2.Add(tt.d1)
				}()
			})
		})
	}
}

func TestDuration_Format(t *testing.T) {
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
			name:      "default zero value",
			of:        0,
			exclusive: []chrono.Designator{},
			expected:  "PT0S",
		},
		{
			name:      "exclusive HMS",
			of:        1*chrono.Hour + 15*chrono.Minute + 30*chrono.Second + 500*chrono.Millisecond,
			exclusive: []chrono.Designator{chrono.Hours, chrono.Minutes, chrono.Seconds},
			expected:  "PT1H15M30.5S",
		},
		{
			name:      "exclusive HMS zero value",
			of:        0,
			exclusive: []chrono.Designator{chrono.Hours, chrono.Minutes, chrono.Seconds},
			expected:  "PT0H0M0S",
		},
		{
			name:      "exclusive HM",
			of:        1*chrono.Hour + 15*chrono.Minute + 30*chrono.Second + 600*chrono.Millisecond,
			exclusive: []chrono.Designator{chrono.Hours, chrono.Minutes},
			expected:  "PT1H15.51M",
		},
		{
			name:      "exclusive HM zero value",
			of:        0,
			exclusive: []chrono.Designator{chrono.Hours, chrono.Minutes},
			expected:  "PT0H0M",
		},
		{
			name:      "exclusive HS",
			of:        12*chrono.Hour + 1*chrono.Minute + 30*chrono.Second + 500*chrono.Millisecond,
			exclusive: []chrono.Designator{chrono.Hours, chrono.Seconds},
			expected:  "PT12H90.5S",
		},
		{
			name:      "exclusive HS zero value",
			of:        0,
			exclusive: []chrono.Designator{chrono.Hours, chrono.Seconds},
			expected:  "PT0H0S",
		},
		{
			name:      "exclusive H",
			of:        1*chrono.Hour + 30*chrono.Minute + 36*chrono.Second + 36*chrono.Millisecond,
			exclusive: []chrono.Designator{chrono.Hours},
			expected:  "PT1.51001H",
		},
		{
			name:      "exclusive H zero value",
			of:        0,
			exclusive: []chrono.Designator{chrono.Hours},
			expected:  "PT0H",
		},
		{
			name:      "exclusive MS",
			of:        1*chrono.Hour + 15*chrono.Minute + 30*chrono.Second + 500*chrono.Millisecond,
			exclusive: []chrono.Designator{chrono.Minutes, chrono.Seconds},
			expected:  "PT75M30.5S",
		},
		{
			name:      "exclusive MS zero value",
			of:        0,
			exclusive: []chrono.Designator{chrono.Minutes, chrono.Seconds},
			expected:  "PT0M0S",
		},
		{
			name:      "exclusive M",
			of:        1*chrono.Hour + 15*chrono.Minute + 30*chrono.Second + 600*chrono.Millisecond,
			exclusive: []chrono.Designator{chrono.Minutes},
			expected:  "PT75.51M",
		},
		{
			name:      "exclusive M zero value",
			of:        0,
			exclusive: []chrono.Designator{chrono.Minutes},
			expected:  "PT0M",
		},
		{
			name:      "exclusive S",
			of:        1*chrono.Hour + 15*chrono.Minute + 30*chrono.Second + 500*chrono.Millisecond,
			exclusive: []chrono.Designator{chrono.Seconds},
			expected:  "PT4530.5S",
		},
		{
			name:      "exclusive S zero value",
			of:        0,
			exclusive: []chrono.Designator{chrono.Seconds},
			expected:  "PT0S",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			d := chrono.DurationOf(tt.of)
			if out := d.Format(tt.exclusive...); out != tt.expected {
				t.Errorf("formatted duration = %s, want %s", out, tt.expected)
			}
		})
	}
}

func TestDuration_Parse(t *testing.T) {
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
			if err := d.Parse(tt.input); err != nil {
				t.Errorf("failed to parse duation: %v", err)
			} else if d.Compare(tt.expected) != 0 {
				t.Errorf("parsed duration = %v, want %v", d, tt.expected)
			}
		})
	}

	t.Run("overflows", func(t *testing.T) {
		var d chrono.Duration
		if err := d.Parse("PT2562047788015216H"); err == nil {
			t.Error("expecting error but got nil")
		}
	})

	t.Run("underflows", func(t *testing.T) {
		var d chrono.Duration
		if err := d.Parse("PT-2562047788015215H"); err == nil {
			t.Error("expecting error but got nil")
		}
	})
}
