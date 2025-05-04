package chrono_test

import (
	"fmt"
	"testing"

	"github.com/go-chrono/chrono"
)

func TestOffsetDateTimeOf(t *testing.T) {
	for _, tt := range []struct {
		datetime     chrono.OffsetDateTime
		expectedDate chrono.LocalDate
		expectedTime chrono.OffsetTime
	}{
		{
			datetime:     chrono.OffsetDateTime{},
			expectedDate: chrono.LocalDateOf(1970, chrono.January, 1),
			expectedTime: chrono.OffsetTimeOf(0, 0, 0, 0, 0, 0),
		},
		{
			datetime:     chrono.OffsetDateTimeOf(2020, chrono.March, 18, 12, 30, 0, 100000000, 0, 0),
			expectedDate: chrono.LocalDateOf(2020, chrono.March, 18),
			expectedTime: chrono.OffsetTimeOf(12, 30, 0, 100000000, 0, 0),
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

func TestOfLocalDateOffsetTime(t *testing.T) {
	datetime := chrono.OfLocalDateOffsetTime(
		chrono.LocalDateOf(2020, chrono.March, 18),
		chrono.OffsetTimeOf(12, 30, 0, 100000000, -2, 30),
	)

	date, time := datetime.Split()
	if expected := chrono.LocalDateOf(2020, chrono.March, 18); date != expected {
		t.Errorf("datetime.Split() date = %s, want %s", date, expected)
	}

	if expected := chrono.OffsetTimeOf(12, 30, 0, 100000000, -2, 30); time != expected {
		t.Errorf("datetime.Split() time = %s, want %s", time, expected)
	}
}

func TestOfLocalDateTimeOffset(t *testing.T) {
	datetime := chrono.OfLocalDateTimeOffset(
		chrono.LocalDateOf(2020, chrono.March, 18),
		chrono.LocalTimeOf(12, 30, 0, 100000000),
		chrono.Extent(chrono.OffsetOf(-2, 30)),
	)

	date, time := datetime.Split()
	if expected := chrono.LocalDateOf(2020, chrono.March, 18); date != expected {
		t.Errorf("datetime.Split() date = %s, want %s", date, expected)
	}

	if expected := chrono.OffsetTimeOf(12, 30, 0, 100000000, -2, 30); time != expected {
		t.Errorf("datetime.Split() time = %s, want %s", time, expected)
	}
}

func TestOffsetDateTime_String(t *testing.T) {
	for _, tt := range []struct {
		name     string
		datetime chrono.OffsetDateTime
		expected string
	}{
		{"simple", chrono.OffsetDateTimeOf(2020, chrono.March, 18, 9, 0, 0, 0, 2, 30), "2020-03-18 09:00:00+02:30"},
		{"micros", chrono.OffsetDateTimeOf(2020, chrono.March, 18, 9, 0, 0, 1e3, 2, 30), "2020-03-18 09:00:00.000001+02:30"},
		{"millis", chrono.OffsetDateTimeOf(2020, chrono.March, 18, 9, 0, 0, 1e6, 2, 30), "2020-03-18 09:00:00.001+02:30"},
		{"nanos", chrono.OffsetDateTimeOf(2020, chrono.March, 18, 9, 0, 0, 12345678, 2, 30), "2020-03-18 09:00:00.012345678+02:30"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if output := tt.datetime.String(); output != tt.expected {
				t.Errorf("OffsetDateTime.String() = %s, want %s", output, tt.expected)
			}
		})
	}
}

func TestOffsetDateTime_Compare(t *testing.T) {
	for _, tt := range []struct {
		name     string
		d        chrono.OffsetDateTime
		d2       chrono.OffsetDateTime
		expected int
	}{
		{"earlier", chrono.OffsetDateTimeOf(2020, chrono.March, 18, 11, 0, 0, 0, -2, 30), chrono.OffsetDateTimeOf(2020, chrono.March, 18, 12, 0, 0, 0, -2, 30), -1},
		{"later", chrono.OffsetDateTimeOf(2020, chrono.March, 18, 13, 30, 0, 0, -2, 30), chrono.OffsetDateTimeOf(2020, chrono.March, 18, 13, 29, 55, 0, -2, 30), 1},
		{"equal", chrono.OffsetDateTimeOf(2020, chrono.March, 18, 15, 0, 0, 1000, -2, 30), chrono.OffsetDateTimeOf(2020, chrono.March, 18, 15, 0, 0, 1000, -2, 30), 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if v := tt.d.Compare(tt.d2); v != tt.expected {
				t.Errorf("t.Compare(t2) = %d, want %d", v, tt.expected)
			}
		})
	}
}

func TestOffsetDateTime_Add(t *testing.T) {
	t.Run("valid add", func(t *testing.T) {
		datetime := chrono.OffsetDateTimeOf(2020, chrono.March, 18, 12, 30, 0, 0, -2, 0)
		duration := chrono.DurationOf(48*chrono.Hour + 1000*chrono.Nanosecond)

		if !datetime.CanAdd(duration) {
			t.Errorf("datetime = %s, datetime.CanAdd(%s) = false, want true", datetime, duration)
		}

		added := datetime.Add(duration)
		expected := chrono.OffsetDateTimeOf(2020, chrono.March, 20, 12, 30, 0, 1000, -2, 30)

		if added.Compare(expected) != 0 {
			t.Errorf("datetime = %s, datetime.Add(%s) = %s, want %s", datetime, duration, added, expected)
		}
	})

	t.Run("invalid add to low", func(t *testing.T) {
		datetime := chrono.MinLocalDateTime().UTC()
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
		datetime := chrono.MaxLocalDateTime().UTC()
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

func TestOffsetDateTime_AddDate(t *testing.T) {
	for _, tt := range []struct {
		name      string
		datetime  chrono.OffsetDateTime
		addYears  int
		addMonths int
		addDays   int
		expected  chrono.OffsetDateTime
	}{
		{"nothing", chrono.OffsetDateTimeOf(2020, chrono.March, 18, 0, 0, 0, 0, -2, 30), 0, 0, 0, chrono.OffsetDateTimeOf(2020, chrono.March, 18, 0, 0, 0, 0, -2, 30)},
		{"add years", chrono.OffsetDateTimeOf(2020, chrono.March, 18, 0, 0, 0, 0, -2, 30), 105, 0, 0, chrono.OffsetDateTimeOf(2125, chrono.March, 18, 0, 0, 0, 0, -2, 30)},
		{"sub years", chrono.OffsetDateTimeOf(2020, chrono.March, 18, 0, 0, 0, 0, -2, 30), -280, 0, 0, chrono.OffsetDateTimeOf(1740, chrono.March, 18, 0, 0, 0, 0, -2, 30)},
		{"add months", chrono.OffsetDateTimeOf(2020, chrono.March, 18, 0, 0, 0, 0, -2, 30), 0, 6, 0, chrono.OffsetDateTimeOf(2020, chrono.September, 18, 0, 0, 0, 0, -2, 30)},
		{"sub months", chrono.OffsetDateTimeOf(2020, chrono.March, 18, 0, 0, 0, 0, -2, 30), 0, -2, 0, chrono.OffsetDateTimeOf(2020, chrono.January, 18, 0, 0, 0, 0, -2, 30)},
		{"add days", chrono.OffsetDateTimeOf(2020, chrono.March, 18, 0, 0, 0, 0, -2, 30), 0, 0, 8, chrono.OffsetDateTimeOf(2020, chrono.March, 26, 0, 0, 0, 0, -2, 30)},
		{"sub days", chrono.OffsetDateTimeOf(2020, chrono.March, 18, 0, 0, 0, 0, -2, 30), 0, 0, -15, chrono.OffsetDateTimeOf(2020, chrono.March, 3, 0, 0, 0, 0, -2, 30)},
		{"time package example", chrono.OffsetDateTimeOf(2011, chrono.January, 1, 0, 0, 0, 0, -2, 30), -1, 2, 3, chrono.OffsetDateTimeOf(2010, chrono.March, 4, 0, 0, 0, 0, -2, 30)},
		{"normalized time package example", chrono.OffsetDateTimeOf(2011, chrono.October, 31, 0, 0, 0, 0, -2, 30), 0, 1, 0, chrono.OffsetDateTimeOf(2011, chrono.December, 1, 0, 0, 0, 0, -2, 30)},
		{"wrap around day", chrono.OffsetDateTimeOf(2020, chrono.March, 18, 0, 0, 0, 0, -2, 30), 0, 0, 20, chrono.OffsetDateTimeOf(2020, chrono.April, 7, 0, 0, 0, 0, -2, 30)},
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
		datetime chrono.OffsetDateTime
		addDays  int
	}{
		{"underflow", chrono.MinLocalDateTime().UTC(), -1},
		{"overflow", chrono.MaxLocalDateTime().UTC(), 1},
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

func TestOffsetDateTime_Sub(t *testing.T) {
	for _, tt := range []struct {
		dt1  chrono.OffsetDateTime
		dt2  chrono.OffsetDateTime
		diff chrono.Duration
	}{
		{
			chrono.OffsetDateTimeOf(2020, chrono.January, 5, 12, 0, 0, 0, 4, 30),
			chrono.OffsetDateTimeOf(2020, chrono.January, 3, 6, 0, 0, 0, 2, 0),
			chrono.DurationOf(56*chrono.Hour + 30*chrono.Minute),
		},
		{
			chrono.OffsetDateTimeOf(2020, chrono.January, 5, 12, 0, 0, 22, 0, 0),
			chrono.OffsetDateTimeOf(2020, chrono.January, 5, 12, 0, 0, 40, 0, 0),
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
