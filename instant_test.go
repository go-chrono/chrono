package chrono_test

import (
	"testing"

	"github.com/go-chrono/chrono"
	chronotest "github.com/go-chrono/chrono/test"
)

func TestInstant_String(t *testing.T) {
	i := chronotest.InstantOf(11e14)
	if str := i.String(); str != "1100000000000000" {
		t.Fatalf("i.String() = %s, want 1100000000000000", str)
	}
}

func TestInstant_Compare(t *testing.T) {
	for _, tt := range []struct {
		name     string
		i        chrono.Instant
		i2       chrono.Instant
		expected int
	}{
		{"earlier", chronotest.InstantOf(11e14), chronotest.InstantOf(12e14), -1},
		{"later", chronotest.InstantOf(12e14), chronotest.InstantOf(11e14), 1},
		{"equal", chronotest.InstantOf(11e14), chronotest.InstantOf(11e14), 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if v := tt.i.Compare(tt.i2); v != tt.expected {
				t.Errorf("i.Compare(i2) = %d, want %d", v, tt.expected)
			}
		})
	}
}

func TestInstant_Until(t *testing.T) {
	i := chronotest.InstantOf(11e14)
	d := i.Until(chronotest.InstantOf(12e14))
	if nsec := d.Nanoseconds(); nsec != 1e14 {
		t.Fatalf("d.Nanoseconds() = %f, want 1e14", nsec)
	}
}
