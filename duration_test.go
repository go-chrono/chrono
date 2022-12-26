package chrono_test

import (
	"reflect"
	"runtime"
	"strings"
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
			of:   90000000 * chrono.Nanosecond,
			nsec: 90000000,
		},
		{
			name: "of positive microsecond",
			of:   90000 * chrono.Microsecond,
			nsec: 90000000,
		},
		{
			name: "of positive milliseconds",
			of:   90 * chrono.Millisecond,
			nsec: 90000000,
		},
		{
			name: "of positive seconds",
			of:   9 * chrono.Second,
			nsec: 9000000000,
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
			of:   -90000000 * chrono.Nanosecond,
			nsec: -90000000,
		},
		{
			name: "of negative microsecond",
			of:   -90000 * chrono.Microsecond,
			nsec: -90000000,
		},
		{
			name: "of negative milliseconds",
			of:   -90 * chrono.Millisecond,
			nsec: -90000000,
		},
		{
			name: "of negative seconds",
			of:   -9 * chrono.Second,
			nsec: -9000000000,
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
			of:       9500 * chrono.Nanosecond,
			f:        chrono.Duration.Microseconds,
			expected: 9.5,
		},
		{
			name:     "positive milliseconds",
			of:       9500 * chrono.Microsecond,
			f:        chrono.Duration.Milliseconds,
			expected: 9.5,
		},
		{
			name:     "positive seconds",
			of:       9500 * chrono.Millisecond,
			f:        chrono.Duration.Seconds,
			expected: 9.5,
		},
		{
			name:     "positive minutes",
			of:       90 * chrono.Second,
			f:        chrono.Duration.Minutes,
			expected: 1.5,
		},
		{
			name:     "positive hours",
			of:       90 * chrono.Minute,
			f:        chrono.Duration.Hours,
			expected: 1.5,
		},
		{
			name:     "negative nanoseconds",
			of:       -9000000000000 * chrono.Nanosecond,
			f:        chrono.Duration.Nanoseconds,
			expected: -9000000000000,
		},
		{
			name:     "negative microseconds",
			of:       -9500 * chrono.Nanosecond,
			f:        chrono.Duration.Microseconds,
			expected: -9.5,
		},
		{
			name:     "negative milliseconds",
			of:       -9500 * chrono.Microsecond,
			f:        chrono.Duration.Milliseconds,
			expected: -9.5,
		},
		{
			name:     "negative seconds",
			of:       -9500 * chrono.Millisecond,
			f:        chrono.Duration.Seconds,
			expected: -9.5,
		},
		{
			name:     "negative minutes",
			of:       -90 * chrono.Second,
			f:        chrono.Duration.Minutes,
			expected: -1.5,
		},
		{
			name:     "negative hours",
			of:       -90 * chrono.Minute,
			f:        chrono.Duration.Hours,
			expected: -1.5,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			d := chrono.DurationOf(tt.of)
			if out := tt.f(d); out != tt.expected {
				t.Errorf("%v() = %f, want %f", runtime.FuncForPC(reflect.ValueOf(tt.f).Pointer()).Name(), out, tt.expected)
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
			name:     "add",
			d1:       chrono.DurationOf((1 * chrono.Hour) + (750 * chrono.Millisecond)),
			d2:       chrono.DurationOf((1 * chrono.Hour) + (550 * chrono.Millisecond)),
			expected: chrono.DurationOf((2 * chrono.Hour) + (1 * chrono.Second) + (300 * chrono.Millisecond)),
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
			name: "overflow",
			d1:   chrono.MaxDuration(),
			d2:   chrono.DurationOf(1 * chrono.Nanosecond),
		},
		{
			name: "underflow",
			d1:   chrono.MinDuration(),
			d2:   chrono.DurationOf(-1 * chrono.Nanosecond),
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
			of:        500 * chrono.Millisecond,
			exclusive: []chrono.Designator{},
			expected:  "PT0.5S",
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
			t.Run("positive", func(t *testing.T) {
				d := chrono.DurationOf(tt.of)
				if out := d.Format(tt.exclusive...); out != tt.expected {
					t.Errorf("formatted duration = %s, want %s", out, tt.expected)
				}
			})

			t.Run("negative", func(t *testing.T) {
				expected := tt.expected
				if tt.of != 0 {
					expected = "-" + expected
				}

				d := chrono.DurationOf(tt.of * -1)
				if out := d.Format(tt.exclusive...); out != expected {
					t.Errorf("formatted duration = %s, want %s", out, expected)
				}
			})
		})
	}
}

func TestDuration_Parse(t *testing.T) {
	for _, tt := range []struct {
		name     string
		input    string
		expected chrono.Extent
	}{
		{
			name:     "valid HMS integers",
			input:    "PT5H3M1S",
			expected: 5*chrono.Hour + 3*chrono.Minute + 1*chrono.Second,
		},
		{
			name:     "valid HMS floats",
			input:    "PT4.5H3.25M1.1S",
			expected: chrono.Extent(4.5*float64(chrono.Hour) + 3.25*float64(chrono.Minute) + 1.1*float64(chrono.Second)),
		},
		{
			name:     "valid HM integers",
			input:    "PT5H3M",
			expected: 5*chrono.Hour + 3*chrono.Minute,
		},
		{
			name:     "valid HM floats",
			input:    "PT4.5H3.25M",
			expected: chrono.Extent(4.5*float64(chrono.Hour) + 3.25*float64(chrono.Minute)),
		},
		{
			name:     "valid HS integers",
			input:    "PT5H1S",
			expected: 5*chrono.Hour + 1*chrono.Second,
		},
		{
			name:     "valid HS floats",
			input:    "PT4.5H1.1S",
			expected: chrono.Extent(4.5*float64(chrono.Hour) + 1.1*float64(chrono.Second)),
		},
		{
			name:     "valid H integer",
			input:    "PT5H",
			expected: 5 * chrono.Hour,
		},
		{
			name:     "valid H float",
			input:    "PT4.5H",
			expected: chrono.Extent(4.5 * float64(chrono.Hour)),
		},
		{
			name:     "valid MS integers",
			input:    "PT3M1S",
			expected: 3*chrono.Minute + 1*chrono.Second,
		},
		{
			name:     "valid MS floats",
			input:    "PT3.25M1.1S",
			expected: chrono.Extent(3.25*float64(chrono.Minute) + 1.1*float64(chrono.Second)),
		},
		{
			name:     "valid M integer",
			input:    "PT3M",
			expected: 3 * chrono.Minute,
		},
		{
			name:     "valid M float",
			input:    "PT3.25M",
			expected: chrono.Extent(3.25 * float64(chrono.Minute)),
		},
		{
			name:     "valid S integer",
			input:    "PT1S",
			expected: 1 * chrono.Second,
		},
		{
			name:     "valid S float",
			input:    "PT1.1S",
			expected: chrono.Extent(1.1 * float64(chrono.Second)),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, sign := range []string{"", "+", "-"} {
				t.Run(sign, func(t *testing.T) {
					input := sign + tt.input
					expected := tt.expected
					if sign == "-" {
						expected *= -1
					}

					run := func() {
						var d chrono.Duration
						if err := d.Parse(input); err != nil {
							t.Errorf("failed to parse duation: %v", err)
						} else if d.Compare(chrono.DurationOf(expected)) != 0 {
							t.Errorf("parsed duration = %v, want %v", d, expected)
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

func TestDuration_Units(t *testing.T) {
	d := chrono.DurationOf(12*chrono.Hour + 34*chrono.Minute + 56*chrono.Second + 7*chrono.Nanosecond)

	hours, mins, secs, nsec := d.Units()
	if hours != 12 {
		t.Errorf("expecting 12 hours, got %d", hours)
	}

	if mins != 34 {
		t.Errorf("expecting 34 mins, got %d", mins)
	}

	if secs != 56 {
		t.Errorf("expecting 56 secs, got %d", secs)
	}

	if nsec != 7 {
		t.Errorf("expecting 7 nsecs, got %d", nsec)
	}
}
