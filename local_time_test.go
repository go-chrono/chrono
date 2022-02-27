package chrono_test

import (
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
