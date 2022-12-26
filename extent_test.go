package chrono_test

import (
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
