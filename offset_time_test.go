package chrono_test

import (
	"fmt"
	"testing"

	"github.com/go-chrono/chrono"
)

func TestOffsetTime(t *testing.T) {
	time := chrono.OffsetTimeOf(12, 30, 59, 12345678, 2, 30)

	hour, min, sec := time.Clock()
	if hour != 12 {
		t.Errorf("time.Clock() hour = %d, want 12", hour)
	}

	if min != 30 {
		t.Errorf("time.Clock() min = %d, want 30", min)
	}

	if sec != 59 {
		t.Errorf("time.Clock() sec = %d, want 59", sec)
	}

	if nsec := time.Nanosecond(); nsec != 12345678 {
		t.Errorf("time.Nanosecond() = %d, want 12345678", nsec)
	}

	expectedOffset := chrono.OffsetOf(2, 30)
	if offset := time.Offset(); offset != expectedOffset {
		t.Errorf("time.Offset() = %s, want %s", offset, expectedOffset)
	}
}

func TestOfTimeOffset(t *testing.T) {
	expectedLocalTime := chrono.LocalTimeOf(12, 30, 59, 12345678)
	expectedOffset := chrono.OffsetOf(3, 30)

	offsetTime := chrono.OfTimeOffset(expectedLocalTime, expectedOffset)
	localTime := offsetTime.Local()
	offset := offsetTime.Offset()

	if localTime.Compare(expectedLocalTime) != 0 {
		t.Errorf("LocalTime = %s, want %s", localTime, expectedLocalTime)
	}
	if offset != expectedOffset {
		t.Errorf("Offset = %s, want %s", offset, expectedOffset)
	}
}

func TestOffsetTime_String(t *testing.T) {
	for _, tt := range []struct {
		name     string
		time     chrono.OffsetTime
		expected string
	}{
		{"simple", chrono.OffsetTimeOf(9, 0, 0, 0, 2, 30), "09:00:00+02:30"},
		{"micros", chrono.OffsetTimeOf(9, 0, 0, 1e3, 2, 30), "09:00:00.000001+02:30"},
		{"millis", chrono.OffsetTimeOf(9, 0, 0, 1e6, 2, 30), "09:00:00.001+02:30"},
		{"nanos", chrono.OffsetTimeOf(9, 0, 0, 12345678, 2, 30), "09:00:00.012345678+02:30"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if output := tt.time.String(); output != tt.expected {
				t.Errorf("OffsetTime.String() = %s, want %s", output, tt.expected)
			}
		})
	}
}

func TestOffsetTime_BusinessHour(t *testing.T) {
	time := chrono.OffsetTimeOf(25, 0, 0, 0, 2, 30)

	if hour := time.BusinessHour(); hour != 25 {
		t.Errorf("time.Hour() = %d, want 25", hour)
	}

	if hour, _, _ := time.Clock(); hour != 1 {
		t.Errorf("time.Hour() = %d, want 1", hour)
	}
}

func TestOffsetTime_Sub(t *testing.T) {
	for _, tt := range []struct {
		t1   chrono.OffsetTime
		t2   chrono.OffsetTime
		diff chrono.Extent
	}{
		{chrono.OffsetTimeOf(12, 0, 0, 0, 0, 0), chrono.OffsetTimeOf(6, 0, 0, 0, 0, 0), 6 * chrono.Hour},
		{chrono.OffsetTimeOf(12, 0, 0, 22, 0, 0), chrono.OffsetTimeOf(12, 0, 0, 40, 0, 0), -18 * chrono.Nanosecond},
		{chrono.OffsetTimeOf(12, 0, 0, 0, 2, 30), chrono.OffsetTimeOf(6, 0, 0, 0, 1, 0), 4*chrono.Hour + 30*chrono.Minute},
		{chrono.OffsetTimeOf(12, 0, 0, 0, -2, 30), chrono.OffsetTimeOf(6, 0, 0, 0, -1, 0), 7*chrono.Hour + 30*chrono.Minute},
	} {
		t.Run(fmt.Sprintf("%s - %s", tt.t1, tt.t2), func(t *testing.T) {
			if d := tt.t1.Sub(tt.t2); d != tt.diff {
				t.Errorf("t1.Sub(t2) = %v, want %v", d, tt.diff)
			}
		})
	}
}

func TestOffsetTime_Add(t *testing.T) {
	for _, tt := range []struct {
		t        chrono.OffsetTime
		e        chrono.Extent
		expected chrono.OffsetTime
	}{
		{chrono.OffsetTimeOf(12, 0, 0, 0, 0, 0), 29 * chrono.Minute, chrono.OffsetTimeOf(12, 29, 0, 0, 0, 0)},
		{chrono.OffsetTimeOf(14, 45, 0, 0, 0, 0), -22 * chrono.Minute, chrono.OffsetTimeOf(14, 23, 0, 0, 0, 0)},
		{chrono.OffsetTimeOf(5, 0, 0, 0, 0, 0), -7 * chrono.Hour, chrono.OffsetTimeOf(22, 0, 0, 0, 0, 0)},
		{chrono.OffsetTimeOf(5, 0, 0, 0, 0, 0), -31 * chrono.Hour, chrono.OffsetTimeOf(22, 0, 0, 0, 0, 0)},
	} {
		t.Run(fmt.Sprintf("%s + %v", tt.t, tt.e), func(t *testing.T) {
			if ok := tt.t.CanAdd(tt.e); !ok {
				t.Error("t.CanAdd(e) = false, want true")
			}

			if added := tt.t.Add(tt.e); added.Compare(tt.expected) != 0 {
				t.Errorf("t.Add(e) = %s, want %s", added, tt.expected)
			}
		})
	}

	for _, tt := range []struct {
		name string
		t    chrono.OffsetTime
		e    chrono.Extent
	}{
		{"invalid duration", chrono.OffsetTimeOf(0, 0, 0, 0, 0, 0), 200 * chrono.Hour},
		{"invalid time", chrono.OffsetTimeOf(90, 0, 0, 0, 0, 0), 20 * chrono.Hour},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if ok := tt.t.CanAdd(tt.e); ok {
				t.Error("t.CanAdd(e) = true, want false")
			}

			func() {
				defer func() {
					if r := recover(); r == nil {
						t.Error("expecting panic that didn't occur")
					}
				}()

				tt.t.Add(tt.e)
			}()
		})
	}
}

func TestOffsetTime_Split(t *testing.T) {
	offsetTime := chrono.OffsetTimeOf(12, 30, 59, 12345678, 3, 30)
	localTime := offsetTime.Local()
	offset := offsetTime.Offset()

	expectedLocalTime := chrono.LocalTimeOf(12, 30, 59, 12345678)
	expectedOffset := chrono.OffsetOf(3, 30)
	if localTime.Compare(expectedLocalTime) != 0 {
		t.Errorf("LocalTime = %s, want %s", localTime, expectedLocalTime)
	}
	if offset != expectedOffset {
		t.Errorf("Offset = %s, want %s", offset, expectedOffset)
	}
}

func TestOffsetTime_UTC(t *testing.T) {
	offsetTime := chrono.OffsetTimeOf(12, 30, 59, 12345678, 3, 30)
	utc := offsetTime.UTC()
	localTime := utc.Local()
	offset := utc.Offset()

	expectedLocalTime := chrono.LocalTimeOf(9, 0, 59, 12345678)
	if localTime.Compare(expectedLocalTime) != 0 {
		t.Errorf("time.UTC().Split() time = %s, want %s", localTime, expectedLocalTime)
	}

	if offset != chrono.UTC {
		t.Errorf("time.UTC().Split() offset = %s, want %s", localTime, chrono.UTC)
	}
}

func TestOffsetTime_In(t *testing.T) {
	originalTime := chrono.OffsetTimeOf(12, 30, 0, 0, 3, 30)

	for _, tt := range []struct {
		name     string
		offset   chrono.Offset
		expected chrono.LocalTime
	}{
		{"UTC", chrono.UTC, chrono.LocalTimeOf(9, 0, 0, 0)},
		{"positive", chrono.OffsetOf(5, 30), chrono.LocalTimeOf(14, 30, 0, 0)},
		{"negative", chrono.OffsetOf(-5, 30), chrono.LocalTimeOf(3, 30, 0, 0)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			inOffset := originalTime.In(tt.offset)
			localTime := inOffset.Local()
			offset := inOffset.Offset()

			if localTime.Compare(tt.expected) != 0 {
				t.Errorf("time.In(%s).Split() time = %s, want %s", tt.offset, localTime, tt.expected)
			}

			if offset != tt.offset {
				t.Errorf("time.In(%s).Split() offset = %s, want %s", tt.offset, localTime, tt.offset)
			}
		})
	}
}
