package chrono_test

import (
	"strings"
	"testing"

	"github.com/go-chrono/chrono"
)

func TestExtent_Truncate(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		e := 1*chrono.Hour + 20*chrono.Minute + 5*chrono.Second + 42*chrono.Millisecond + 307*chrono.Microsecond

		if e2 := e.Truncate(200 * chrono.Microsecond); e2 != 4805042200000 {
			t.Errorf("e.Truncate() = %d, want 4805042200000", e2)
		}
	})

	t.Run("negative", func(t *testing.T) {
		e := -1*chrono.Hour - 20*chrono.Minute - 5*chrono.Second - 42*chrono.Millisecond - 307*chrono.Microsecond

		if e2 := e.Truncate(200 * chrono.Microsecond); e2 != -4805042200000 {
			t.Errorf("e.Truncate() = %d, want -4805042200000", e2)
		}
	})
}

func TestExtent_Units(t *testing.T) {
	e := (12 * chrono.Hour) + (34 * chrono.Minute) + (56 * chrono.Second) + (7 * chrono.Nanosecond)

	hours, mins, secs, nsec := e.Units()
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

func TestExtent_Format(t *testing.T) {
	for _, tt := range formatDurationCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("positive", func(t *testing.T) {
				if out := tt.of.Format(tt.exclusive...); out != tt.expected {
					t.Errorf("formatted extent = %s, want %s", out, tt.expected)
				}
			})

			t.Run("negative", func(t *testing.T) {
				expected := tt.expected
				if tt.of != 0 {
					expected = "-" + expected
				}

				if out := (tt.of * -1).Format(tt.exclusive...); out != expected {
					t.Errorf("formatted extent = %s, want %s", out, expected)
				}
			})
		})
	}
}

func TestExtent_Parse(t *testing.T) {
	for _, tt := range parseDurationCases {
		t.Run(tt.name, func(t *testing.T) {
			for _, sign := range []string{"", "+", "-"} {
				t.Run(sign, func(t *testing.T) {
					input := sign + tt.input
					expected := tt.expected
					if sign == "-" {
						expected *= -1
					}

					run := func() {
						var e chrono.Extent
						if err := e.Parse(input); err != nil {
							t.Errorf("failed to parse duation: %v", err)
						} else if e != expected {
							t.Errorf("parsed extent = %v, want %v", e, expected)
						}
					}

					t.Run("dots", func(_ *testing.T) {
						run()
					})

					t.Run("commas", func(_ *testing.T) {
						tt.input = strings.ReplaceAll(tt.input, ".", ",")
						run()
					})
				})
			}
		})
	}

	t.Run("overflows", func(t *testing.T) {
		var e chrono.Extent
		if err := e.Parse("PT9223372037S"); err == nil {
			t.Error("expecting error but got nil")
		}
	})

	t.Run("underflows", func(t *testing.T) {
		var d chrono.Duration
		if err := d.Parse("PT-9223372035S"); err == nil {
			t.Error("expecting error but got nil")
		}
	})
}
