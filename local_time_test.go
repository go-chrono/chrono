package chrono_test

import (
	"fmt"
	"testing"

	"github.com/go-chrono/chrono"
)

func TestLocalTime(t *testing.T) {
	time := chrono.LocalTimeOf(12, 30, 59, 12345678)

	hour, min, sec := time.Clock()
	if hour != 12 {
		t.Errorf("time.Hour() = %d, want 12", hour)
	}

	if min != 30 {
		t.Errorf("time.Minute() = %d, want 30", min)
	}

	if sec != 59 {
		t.Errorf("time.Second() = %d, want 59", sec)
	}

	if nsec := time.Nanosecond(); nsec != 12345678 {
		t.Errorf("time.Nanosecond() = %d, want 12345678", nsec)
	}
}

func TestOfDayOfYear(t *testing.T) {
	for _, tt := range []struct {
		year     int
		day      int
		expected chrono.LocalDate
	}{
		{2020, 60, chrono.LocalDateOf(2020, chrono.February, 29)},
		{2021, 60, chrono.LocalDateOf(2021, chrono.March, 1)},
		{2020, 120, chrono.LocalDateOf(2020, chrono.April, 29)},
		{2021, 120, chrono.LocalDateOf(2021, chrono.April, 30)},
	} {
		t.Run(fmt.Sprintf("%04d-%03d", tt.year, tt.day), func(t *testing.T) {
			if date := chrono.OfDayOfYear(tt.year, tt.day); date != tt.expected {
				t.Errorf("OfDayOfYear(%d, %d) = %s, want %s", tt.year, tt.day, date, tt.expected)
			}
		})
	}
}

func TestOfFirstWeekday(t *testing.T) {
	for _, tt := range []struct {
		year     int
		month    chrono.Month
		weekday  chrono.Weekday
		expected chrono.LocalDate
	}{
		{2020, chrono.January, chrono.Wednesday, chrono.LocalDateOf(2020, chrono.January, 1)},
		{2020, chrono.January, chrono.Monday, chrono.LocalDateOf(2020, chrono.January, 6)},
		{2020, chrono.March, chrono.Sunday, chrono.LocalDateOf(2020, chrono.March, 1)},
	} {
		t.Run(fmt.Sprintf("%04d-%s %s", tt.year, tt.month, tt.weekday), func(t *testing.T) {
			if date := chrono.OfFirstWeekday(tt.year, tt.month, tt.weekday); date != tt.expected {
				t.Errorf("OfFirstWeekday(%d, %s, %s) = %s, want %s", tt.year, tt.month, tt.weekday, date, tt.expected)
			} else if weekday := date.Weekday(); weekday != tt.weekday {
				t.Errorf("weekday = %s, want %s", weekday, tt.weekday)
			}
		})
	}
}

func TestLocalTime_BusinessHour(t *testing.T) {
	time := chrono.LocalTimeOf(25, 0, 0, 0)

	if hour := time.BusinessHour(); hour != 25 {
		t.Errorf("time.Hour() = %d, want 25", hour)
	}

	if hour, _, _ := time.Clock(); hour != 1 {
		t.Errorf("time.Hour() = %d, want 1", hour)
	}
}

func TestLocalTime_Sub(t *testing.T) {
	for _, tt := range []struct {
		t1   chrono.LocalTime
		t2   chrono.LocalTime
		diff chrono.Extent
	}{
		{chrono.LocalTimeOf(12, 0, 0, 0), chrono.LocalTimeOf(6, 0, 0, 0), 6 * chrono.Hour},
		{chrono.LocalTimeOf(12, 0, 0, 22), chrono.LocalTimeOf(12, 0, 0, 40), -18 * chrono.Nanosecond},
	} {
		t.Run(fmt.Sprintf("%s - %s", tt.t1, tt.t2), func(t *testing.T) {
			if d := tt.t1.Sub(tt.t2); d.Compare(chrono.DurationOf(tt.diff)) != 0 {
				t.Errorf("t1.Sub(t2) = %v, want %v", d, tt.diff)
			}
		})
	}
}

func TestLocalTime_Add(t *testing.T) {
	for _, tt := range []struct {
		t        chrono.LocalTime
		e        chrono.Extent
		expected chrono.LocalTime
	}{
		{chrono.LocalTimeOf(12, 0, 0, 0), 29 * chrono.Minute, chrono.LocalTimeOf(12, 29, 0, 0)},
		{chrono.LocalTimeOf(14, 45, 0, 0), -22 * chrono.Minute, chrono.LocalTimeOf(14, 23, 0, 0)},
		{chrono.LocalTimeOf(5, 0, 0, 0), -7 * chrono.Hour, chrono.LocalTimeOf(22, 0, 0, 0)},
		{chrono.LocalTimeOf(5, 0, 0, 0), -31 * chrono.Hour, chrono.LocalTimeOf(22, 0, 0, 0)},
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
		t    chrono.LocalTime
		e    chrono.Extent
	}{
		{"invalid duration", chrono.LocalTimeOf(0, 0, 0, 0), 200 * chrono.Hour},
		{"invalid time", chrono.LocalTimeOf(90, 0, 0, 0), 20 * chrono.Hour},
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

func TestLocalTime_Compare(t *testing.T) {
	for _, tt := range []struct {
		name     string
		t        chrono.LocalTime
		t2       chrono.LocalTime
		expected int
	}{
		{"earlier", chrono.LocalTimeOf(11, 0, 0, 0), chrono.LocalTimeOf(12, 0, 0, 0), -1},
		{"later", chrono.LocalTimeOf(13, 30, 0, 0), chrono.LocalTimeOf(13, 29, 55, 0), 1},
		{"equal", chrono.LocalTimeOf(15, 0, 0, 1000), chrono.LocalTimeOf(15, 0, 0, 1000), 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if v := tt.t.Compare(tt.t2); v != tt.expected {
				t.Errorf("t.Compare(t2) = %d, want %d", v, tt.expected)
			}
		})
	}
}
