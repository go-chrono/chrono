package chrono_test

import (
	"testing"

	chronotest "github.com/go-chrono/chrono/test"
)

func TestInstantString(t *testing.T) {
	i := chronotest.InstantOf(11e14)
	if str := i.String(); str != "1100000000000000" {
		t.Fatalf("i.String() = %s, want 1100000000000000", str)
	}
}

func TestInstantUntil(t *testing.T) {
	i := chronotest.InstantOf(11e14)
	d := i.Until(chronotest.InstantOf(12e14))
	if nsec := d.Nanoseconds(); nsec != 1e14 {
		t.Fatalf("d.Nanoseconds() = %f, want 1e14", nsec)
	}
}
