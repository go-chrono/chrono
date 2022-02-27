package chrono_test

import (
	"fmt"
	"testing"

	"github.com/go-chrono/chrono"
)

func TestLocalTime(t *testing.T) {
	time := chrono.LocalTimeOf(12, 30, 59, 12345678)

	if hour := time.Hour(); hour != 12 {
		t.Errorf("time.Hour() = %d, want 12", hour)
	}

	if min := time.Minute(); min != 30 {
		t.Errorf("time.Minute() = %d, want 30", min)
	}

	if sec := time.Second(); sec != 59 {
		t.Errorf("time.Second() = %d, want 59", sec)
	}

	if nsec := time.Nanosecond(); nsec != 12345678 {
		t.Errorf("time.Nanosecond() = %d, want 12345678", nsec)
	}
}

func TestLocalTimeSub(t *testing.T) {
	for _, tt := range []struct {
		t1   chrono.LocalTime
		t2   chrono.LocalTime
		diff chrono.Extent
	}{
		{chrono.LocalTimeOf(12, 0, 0, 0), chrono.LocalTimeOf(6, 0, 0, 0), 6 * chrono.Hour},
		{chrono.LocalTimeOf(12, 0, 0, 22), chrono.LocalTimeOf(12, 0, 0, 40), -18 * chrono.Nanosecond},
	} {
		t.Run(fmt.Sprintf("%s - %s", tt.t1, tt.t2), func(t *testing.T) {
			if d := tt.t1.Sub(tt.t2); d != chrono.DurationOf(tt.diff) {
				t.Errorf("t1.Sub(t2) = %v, want %v", d, tt.diff)
			}
		})
	}
}
