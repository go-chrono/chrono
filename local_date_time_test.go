package chrono_test

import (
	"fmt"
	"testing"

	"github.com/go-chrono/chrono"
)

func TestLocalDateTimeOf(t *testing.T) {
	for _, tt := range []struct {
		datetime     chrono.LocalDateTime
		expectedDate chrono.LocalDate
		expectedTime chrono.LocalTime
	}{
		{
			datetime:     chrono.LocalDateTime{},
			expectedDate: chrono.LocalDateOf(1970, chrono.January, 1),
			expectedTime: chrono.LocalTimeOf(0, 0, 0, 0),
		},
		{
			datetime:     chrono.LocalDateTimeOf(2020, chrono.March, 18, 12, 30, 0, 100000000),
			expectedDate: chrono.LocalDateOf(2020, chrono.March, 18),
			expectedTime: chrono.LocalTimeOf(12, 30, 0, 100000000),
		},
	} {
		t.Run(tt.datetime.String(), func(t *testing.T) {
			date, time := tt.datetime.Split()
			if date != tt.expectedDate {
				t.Errorf("datetime.Split() date = %s, want %s", date, tt.expectedDate)
			}

			if time.Compare(tt.expectedTime) != 0 {
				t.Errorf("datetime.Split() time = %s, want %s", time, tt.expectedTime)
			}
		})
	}
}

func TestOfLocalDateTime(t *testing.T) {
	datetime := chrono.OfLocalDateTime(
		chrono.LocalDateOf(2020, chrono.March, 18),
		chrono.LocalTimeOf(12, 30, 0, 100000000),
	)

	date, time := datetime.Split()
	if expected := chrono.LocalDateOf(2020, chrono.March, 18); date != expected {
		t.Errorf("datetime.Split() date = %s, want %s", date, expected)
	}

	if expected := chrono.LocalTimeOf(12, 30, 0, 100000000); time != expected {
		t.Errorf("datetime.Split() time = %s, want %s", time, expected)
	}
}

func TestUnix(t *testing.T) {
	for _, tt := range []struct {
		name             string
		secs             int64
		nsecs            int64
		expected         chrono.LocalDateTime
		expectedUnix     int64
		expectedUnixNano int64
	}{
		{"zero", 0, 0, chrono.LocalDateTime{}, 0, 0},
		{"+secs", 1136171045, 0, chrono.LocalDateTimeOf(2006, chrono.January, 2, 3, 4, 5, 0), 1136171045, 1136171045000000000},
		{"-secs", 820551845, 0, chrono.LocalDateTimeOf(1996, chrono.January, 2, 3, 4, 5, 0), 820551845, 820551845000000000},
		{"+nsecs", 1136171045, 6e10 + 1e6, chrono.LocalDateTimeOf(2006, chrono.January, 2, 3, 5, 5, 1000000), 1136171105, 1136171105001000000},
		{"-nsecs", 1136171045, -6e10 - 1e6, chrono.LocalDateTimeOf(2006, chrono.January, 2, 3, 3, 4, 999000000), 1136170984, 1136170984999000000},
	} {
		t.Run(tt.name, func(t *testing.T) {
			dt := chrono.Unix(tt.secs, tt.nsecs)
			if dt.Compare(tt.expected) != 0 {
				t.Errorf("Unix(%d, %d) = %s, want %s", tt.secs, tt.nsecs, dt, tt.expected)
			}

			if unix := dt.Unix(); unix != tt.expectedUnix {
				t.Errorf("dt.Unix() = %d, want %d", unix, tt.expectedUnix)
			}

			if unix := dt.UnixNano(); unix != tt.expectedUnixNano {
				t.Errorf("dt.UnixNano() = %d, want %d", unix, tt.expectedUnixNano)
			}
		})
	}
}

func TestUnixMilli(t *testing.T) {
	for _, tt := range []struct {
		name     string
		msecs    int64
		expected chrono.LocalDateTime
	}{
		{"zero", 0, chrono.LocalDateTime{}},
		{"msecs", 1136171045000, chrono.LocalDateTimeOf(2006, chrono.January, 2, 3, 4, 5, 0)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			dt := chrono.UnixMilli(tt.msecs)
			if dt.Compare(tt.expected) != 0 {
				t.Errorf("UnixMilli(%d) = %s, want %s", tt.msecs, dt, tt.expected)
			}

			if unix := dt.UnixMilli(); unix != tt.msecs {
				t.Errorf("dt.UnixMilli() = %d, want %d", unix, tt.msecs)
			}
		})
	}
}

func TestUnixMicro(t *testing.T) {
	for _, tt := range []struct {
		name     string
		usecs    int64
		expected chrono.LocalDateTime
	}{
		{"zero", 0, chrono.LocalDateTime{}},
		{"usecs", 1136171045000000, chrono.LocalDateTimeOf(2006, chrono.January, 2, 3, 4, 5, 0)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			dt := chrono.UnixMicro(tt.usecs)
			if dt.Compare(tt.expected) != 0 {
				t.Errorf("UnixMicro(%d) = %s, want %s", tt.usecs, dt, tt.expected)
			}

			if unix := dt.UnixMicro(); unix != tt.usecs {
				t.Errorf("dt.UnixMicro() = %d, want %d", unix, tt.usecs)
			}
		})
	}
}

func TestLocalDateTime_String(t *testing.T) {
	for _, tt := range []struct {
		name     string
		time     chrono.LocalDateTime
		expected string
	}{
		{"simple", chrono.LocalDateTimeOf(2020, chrono.March, 18, 9, 0, 0, 0), "2020-03-18 09:00:00"},
		{"micros", chrono.LocalDateTimeOf(2020, chrono.March, 18, 9, 0, 0, 1e3), "2020-03-18 09:00:00.000001"},
		{"millis", chrono.LocalDateTimeOf(2020, chrono.March, 18, 9, 0, 0, 1e6), "2020-03-18 09:00:00.001"},
		{"nanos", chrono.LocalDateTimeOf(2020, chrono.March, 18, 9, 0, 0, 12345678), "2020-03-18 09:00:00.012345678"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if output := tt.time.String(); output != tt.expected {
				t.Errorf("LocalDateTime.String() = %s, want %s", output, tt.expected)
			}
		})
	}
}

func TestLocalDateTime_Compare(t *testing.T) {
	for _, tt := range []struct {
		name     string
		d        chrono.LocalDateTime
		d2       chrono.LocalDateTime
		expected int
	}{
		{"earlier", chrono.LocalDateTimeOf(2020, chrono.March, 18, 11, 0, 0, 0), chrono.LocalDateTimeOf(2020, chrono.March, 18, 12, 0, 0, 0), -1},
		{"later", chrono.LocalDateTimeOf(2020, chrono.March, 18, 13, 30, 0, 0), chrono.LocalDateTimeOf(2020, chrono.March, 18, 13, 29, 55, 0), 1},
		{"equal", chrono.LocalDateTimeOf(2020, chrono.March, 18, 15, 0, 0, 1000), chrono.LocalDateTimeOf(2020, chrono.March, 18, 15, 0, 0, 1000), 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if v := tt.d.Compare(tt.d2); v != tt.expected {
				t.Errorf("t.Compare(t2) = %d, want %d", v, tt.expected)
			}
		})
	}
}

func TestLocalDateTime_Add(t *testing.T) {
	t.Run("valid add", func(t *testing.T) {
		datetime := chrono.LocalDateTimeOf(2020, chrono.March, 18, 12, 30, 0, 0)
		duration := chrono.DurationOf(48*chrono.Hour + 1000*chrono.Nanosecond)

		if !datetime.CanAdd(duration) {
			t.Errorf("datetime = %s, datetime.CanAdd(%s) = false, want true", datetime, duration)
		}

		added := datetime.Add(duration)
		expected := chrono.LocalDateTimeOf(2020, chrono.March, 20, 12, 30, 0, 1000)

		if added.Compare(expected) != 0 {
			t.Errorf("datetime = %s, datetime.Add(%s) = %s, want %s", datetime, duration, added, expected)
		}
	})

	t.Run("invalid add to low", func(t *testing.T) {
		datetime := chrono.MinLocalDateTime()
		duration := chrono.DurationOf(-1 * chrono.Nanosecond)

		if datetime.CanAdd(duration) {
			t.Errorf("datetime = %s, datetime.CanAdd(%s) = true, want false", datetime, duration)
		}

		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Error("expecting panic that didn't occur")
				}
			}()

			datetime.Add(duration)
		}()
	})

	t.Run("invalid add to high", func(t *testing.T) {
		datetime := chrono.MaxLocalDateTime()
		duration := chrono.DurationOf(1 * chrono.Nanosecond)

		if datetime.CanAdd(duration) {
			t.Errorf("datetime = %s, datetime.CanAdd(%s) = true, want false", datetime, duration)
		}

		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Error("expecting panic that didn't occur")
				}
			}()

			datetime.Add(duration)
		}()
	})
}

func TestLocalDateTime_AddDate(t *testing.T) {
	for _, tt := range []struct {
		name      string
		datetime  chrono.LocalDateTime
		addYears  int
		addMonths int
		addDays   int
		expected  chrono.LocalDateTime
	}{
		{"nothing", chrono.LocalDateTimeOf(2020, chrono.March, 18, 0, 0, 0, 0), 0, 0, 0, chrono.LocalDateTimeOf(2020, chrono.March, 18, 0, 0, 0, 0)},
		{"add years", chrono.LocalDateTimeOf(2020, chrono.March, 18, 0, 0, 0, 0), 105, 0, 0, chrono.LocalDateTimeOf(2125, chrono.March, 18, 0, 0, 0, 0)},
		{"sub years", chrono.LocalDateTimeOf(2020, chrono.March, 18, 0, 0, 0, 0), -280, 0, 0, chrono.LocalDateTimeOf(1740, chrono.March, 18, 0, 0, 0, 0)},
		{"add months", chrono.LocalDateTimeOf(2020, chrono.March, 18, 0, 0, 0, 0), 0, 6, 0, chrono.LocalDateTimeOf(2020, chrono.September, 18, 0, 0, 0, 0)},
		{"sub months", chrono.LocalDateTimeOf(2020, chrono.March, 18, 0, 0, 0, 0), 0, -2, 0, chrono.LocalDateTimeOf(2020, chrono.January, 18, 0, 0, 0, 0)},
		{"add days", chrono.LocalDateTimeOf(2020, chrono.March, 18, 0, 0, 0, 0), 0, 0, 8, chrono.LocalDateTimeOf(2020, chrono.March, 26, 0, 0, 0, 0)},
		{"sub days", chrono.LocalDateTimeOf(2020, chrono.March, 18, 0, 0, 0, 0), 0, 0, -15, chrono.LocalDateTimeOf(2020, chrono.March, 3, 0, 0, 0, 0)},
		{"time package example", chrono.LocalDateTimeOf(2011, chrono.January, 1, 0, 0, 0, 0), -1, 2, 3, chrono.LocalDateTimeOf(2010, chrono.March, 4, 0, 0, 0, 0)},
		{"normalized time package example", chrono.LocalDateTimeOf(2011, chrono.October, 31, 0, 0, 0, 0), 0, 1, 0, chrono.LocalDateTimeOf(2011, chrono.December, 1, 0, 0, 0, 0)},
		{"wrap around day", chrono.LocalDateTimeOf(2020, chrono.March, 18, 0, 0, 0, 0), 0, 0, 20, chrono.LocalDateTimeOf(2020, chrono.April, 7, 0, 0, 0, 0)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if ok := tt.datetime.CanAddDate(tt.addYears, tt.addMonths, tt.addDays); !ok {
				t.Errorf("date = %s, date.CanAddDate(%d, %d, %d) = false, want true", tt.datetime, tt.addYears, tt.addMonths, tt.addDays)
			}

			if date := tt.datetime.AddDate(tt.addYears, tt.addMonths, tt.addDays); date.Compare(tt.expected) != 0 {
				t.Errorf("date = %s, date.AddDate(%d, %d, %d) = %s, want %s", tt.datetime, tt.addYears, tt.addMonths, tt.addDays, date, tt.expected)
			}
		})
	}

	for _, tt := range []struct {
		name     string
		datetime chrono.LocalDateTime
		addDays  int
	}{
		{"underflow", chrono.MinLocalDateTime(), -1},
		{"overflow", chrono.MaxLocalDateTime(), 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if ok := tt.datetime.CanAddDate(0, 0, tt.addDays); ok {
				t.Errorf("date = %s, date.CanAddDate(0, 0, %d) = true, want false", tt.datetime, tt.addDays)
			}

			func() {
				defer func() {
					if r := recover(); r == nil {
						t.Error("expecting panic that didn't occur")
					}
				}()

				tt.datetime.AddDate(0, 0, tt.addDays)
			}()
		})
	}
}

func TestLocalDateTime_Sub(t *testing.T) {
	for _, tt := range []struct {
		dt1  chrono.LocalDateTime
		dt2  chrono.LocalDateTime
		diff chrono.Duration
	}{
		{
			chrono.LocalDateTimeOf(2020, chrono.January, 5, 12, 0, 0, 0),
			chrono.LocalDateTimeOf(2020, chrono.January, 3, 6, 0, 0, 0),
			chrono.DurationOf(54 * chrono.Hour),
		},
		{
			chrono.LocalDateTimeOf(2020, chrono.January, 5, 12, 0, 0, 22),
			chrono.LocalDateTimeOf(2020, chrono.January, 5, 12, 0, 0, 40),
			chrono.DurationOf(-18 * chrono.Nanosecond),
		},
	} {
		t.Run(fmt.Sprintf("%s - %s", tt.dt1, tt.dt2), func(t *testing.T) {
			if d := tt.dt1.Sub(tt.dt2); d.Compare(tt.diff) != 0 {
				t.Errorf("dt1.Sub(dt2) = %v, want %v", d, tt.diff)
			}
		})
	}
}

func TestLocalDateTime_In(t *testing.T) {
	dt := chrono.LocalDateTimeOf(2020, chrono.March, 18, 12, 30, 0, 0)
	output := dt.In(chrono.OffsetOf(2, 30))

	expected := chrono.OffsetDateTimeOf(2020, chrono.March, 18, 12, 30, 0, 0, 2, 30)
	if output.Compare(expected) != 0 {
		t.Errorf("dt.In = %s, want %s", output, expected)
	}
}

func TestLocalDateTime_UTC(t *testing.T) {
	dt := chrono.LocalDateTimeOf(2020, chrono.March, 18, 12, 30, 0, 0)
	output := dt.UTC()

	expected := chrono.OffsetDateTimeOf(2020, chrono.March, 18, 12, 30, 0, 0, 0, 0)
	if output.Compare(expected) != 0 {
		t.Errorf("dt.UTC() = %s, want %s", output, expected)
	}
}
