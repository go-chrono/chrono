package chrono_test

import (
	"testing"

	"github.com/go-chrono/chrono"
)

const (
	formatYear   = 807
	formatMonth  = chrono.February
	formatDay    = 9
	formatHour   = 1
	formatMin    = 5
	formatSec    = 2
	formatMillis = 123
	formatMicros = 123457
	formatNanos  = 123456789
)

func setupCenturyParsing() {
	chrono.SetupCenturyParsing(800)
}

func tearDownCenturyParsing() {
	chrono.TearDownCenturyParsing()
}

func checkYear(t *testing.T, date chrono.LocalDate) {
	if y, _, _ := date.Date(); y != formatYear {
		t.Errorf("date.Date() year = %d, want %d", y, formatYear)
	}
}

func checkYearDay(t *testing.T, date chrono.LocalDate) {
	if d := date.YearDay(); d != 40 {
		t.Errorf("date.YearDay() = %d, want %d", d, 40)
	}
}

func checkMonth(t *testing.T, date chrono.LocalDate) {
	if _, m, _ := date.Date(); m != formatMonth {
		t.Errorf("date.Date() month = %d, want %d", m, formatMonth)
	}
}

func checkDay(t *testing.T, date chrono.LocalDate) {
	if _, _, d := date.Date(); d != formatDay {
		t.Errorf("date.Date() day = %d, want %d", d, formatDay)
	}
}

func checkWeekday(t *testing.T, date chrono.LocalDate) {
	// A parsed weekday is only checked for correctness - it does not affect the resulting LocalDate.
	// See note (3).
}

func checkISOYear(t *testing.T, date chrono.LocalDate) {
	if y, _ := date.ISOWeek(); y != 807 {
		t.Errorf("date.ISOWeek() year = %d, want %d", y, 807)
	}
}

func checkISOWeek(t *testing.T, date chrono.LocalDate) {
	if _, w := date.ISOWeek(); w != 6 {
		t.Errorf("date.ISOWeek() week = %d, want %d", w, 6)
	}
}

func checkTimeOfDay(t *testing.T, time chrono.LocalTime) {
	// Time of day is checked implicitly by checking the hour.
}

func checkHour12HourClock(t *testing.T, time chrono.LocalTime) {
	if h, _, _ := time.Clock(); h != formatHour {
		t.Errorf("time.Clock() hour = %d, want %d", h, formatHour)
	}
}

func checkHour(t *testing.T, time chrono.LocalTime) {
	if h, _, _ := time.Clock(); h != formatHour {
		t.Errorf("time.Clock() hour = %d, want %d", h, formatHour)
	}
}

func checkMinute(t *testing.T, time chrono.LocalTime) {
	if _, m, _ := time.Clock(); m != formatMin {
		t.Errorf("time.Clock() min = %d, want %d", m, formatMin)
	}
}

func checkSecond(t *testing.T, time chrono.LocalTime) {
	if _, _, s := time.Clock(); s != formatSec {
		t.Errorf("time.Clock() sec = %d, want %d", s, formatSec)
	}
}

func checkMillis(t *testing.T, time chrono.LocalTime) {
	if nanos := time.Nanosecond(); nanos != formatMillis*1000000 {
		t.Errorf("time.Nanosecond() = %d, want %d", nanos, formatMillis*1000000)
	}
}

func checkMicros(t *testing.T, time chrono.LocalTime) {
	if nanos := time.Nanosecond(); nanos != formatMicros*1000 {
		t.Errorf("time.Nanosecond() = %d, want %d", nanos, formatMillis*1000)
	}
}

func checkNanos(t *testing.T, time chrono.LocalTime) {
	if nanos := time.Nanosecond(); nanos != formatNanos {
		t.Errorf("time.Nanosecond() = %d, want %d", nanos, formatNanos)
	}
}

var (
	dateSpecifiers = []struct {
		specifier         string
		textToParse       string
		checkParse        func(*testing.T, chrono.LocalDate)
		expectedFormatted string
	}{
		{"%Y", "0807", checkYear, "0807"},
		{"%-Y", "807", checkYear, "807"},
		{"%EY", "0807", checkYear, "0807"},
		{"%-EY", "807", checkYear, "807"},
		{"%y", "07", checkYear, "07"},
		{"%-y", "7", checkYear, "7"},
		{"%Ey", "07", checkYear, "07"},
		{"%-Ey", "7", checkYear, "7"},
		{"%j", "040", checkYearDay, "040"},
		{"%-j", "40", checkYearDay, "40"},
		{"%m", "02", checkMonth, "02"},
		{"%-m", "2", checkMonth, "2"},
		{"%B", "February", checkMonth, "February"},
		{"%B", "february", checkMonth, "February"},
		{"%b", "Feb", checkMonth, "Feb"},
		{"%b", "feb", checkMonth, "Feb"},
		{"%d", "09", checkDay, "09"},
		{"%-d", "9", checkDay, "9"},
		{"%u", "5", checkWeekday, "5"},
		{"%-u", "5", checkWeekday, "5"},
		{"%A", "Friday", checkWeekday, "Friday"},
		{"%A", "friday", checkWeekday, "Friday"},
		{"%a", "Fri", checkWeekday, "Fri"},
		{"%a", "fri", checkWeekday, "Fri"},
		{"%G", "0807", checkISOYear, "0807"},
		{"%-G", "807", checkISOYear, "807"},
		{"%V", "06", checkISOWeek, "06"},
		{"%-V", "6", checkISOWeek, "6"},
	}

	timeSpecifiers = []struct {
		specifier  string
		text       string
		checkParse func(*testing.T, chrono.LocalTime)
	}{
		{"%P", "am", checkTimeOfDay},
		{"%p", "AM", checkTimeOfDay},
		{"%I", "01", checkHour12HourClock},
		{"%-I", "1", checkHour12HourClock},
		{"%H", "01", checkHour},
		{"%-H", "1", checkHour},
		{"%M", "05", checkMinute},
		{"%-M", "5", checkMinute},
		{"%S", "02", checkSecond},
		{"%-S", "2", checkSecond},
		{"%3f", "123", checkMillis},
		{"%6f", "123457", checkMicros},
		{"%9f", "123456789", checkNanos},
		{"%f", "123457", checkMicros},
	}

	predefinedLayouts = []struct {
		layout   string
		text     string
		datetime chrono.LocalDateTime
	}{
		{chrono.ISO8601DateSimple, "08070209", chrono.LocalDateTimeOf(formatYear, formatMonth, formatDay, 0, 0, 0, 0)},
		{chrono.ISO8601DateExtended, "0807-02-09", chrono.LocalDateTimeOf(formatYear, formatMonth, formatDay, 0, 0, 0, 0)},
		{chrono.ISO8601DateTruncated, "0807-02", chrono.LocalDateTimeOf(formatYear, formatMonth, 1, 0, 0, 0, 0)},
		{chrono.ISO8601TimeSimple, "T010502", chrono.LocalDateTimeOf(1970, chrono.January, 1, formatHour, formatMin, formatSec, 0)},
		{chrono.ISO8601TimeExtended, "T01:05:02", chrono.LocalDateTimeOf(1970, chrono.January, 1, formatHour, formatMin, formatSec, 0)},
		{chrono.ISO8601TimeMillisSimple, "T010502.123", chrono.LocalDateTimeOf(1970, chrono.January, 1, formatHour, formatMin, formatSec, 123000000)},
		{chrono.ISO8601TimeMillisExtended, "T01:05:02.123", chrono.LocalDateTimeOf(1970, chrono.January, 1, formatHour, formatMin, formatSec, 123000000)},
		{chrono.ISO8601TimeTruncatedMinsSimple, "T0105", chrono.LocalDateTimeOf(1970, chrono.January, 1, formatHour, formatMin, 0, 0)},
		{chrono.ISO8601TimeTruncatedMinsExtended, "T01:05", chrono.LocalDateTimeOf(1970, chrono.January, 1, formatHour, formatMin, 0, 0)},
		{chrono.ISO8601TimeTruncatedHours, "T01", chrono.LocalDateTimeOf(1970, chrono.January, 1, formatHour, 0, 0, 0)},
		{chrono.ISO8601DateTimeSimple, "08070209T010502", chrono.LocalDateTimeOf(formatYear, formatMonth, formatDay, formatHour, formatMin, formatSec, 0)},
		{chrono.ISO8601DateTimeExtended, "0807-02-09T01:05:02", chrono.LocalDateTimeOf(formatYear, formatMonth, formatDay, formatHour, formatMin, formatSec, 0)},
		{chrono.ISO8601WeekSimple, "0807W06", chrono.LocalDateTimeOf(formatYear, formatMonth, 5, 0, 0, 0, 0)},
		{chrono.ISO8601WeekExtended, "0807-W06", chrono.LocalDateTimeOf(formatYear, formatMonth, 5, 0, 0, 0, 0)},
		{chrono.ISO8601WeekDaySimple, "0807W065", chrono.LocalDateTimeOf(formatYear, formatMonth, formatDay, 0, 0, 0, 0)},
		{chrono.ISO8601WeekDayExtended, "0807-W06-5", chrono.LocalDateTimeOf(formatYear, formatMonth, formatDay, 0, 0, 0, 0)},
		{chrono.ISO8601OrdinalDateSimple, "0807040", chrono.LocalDateTimeOf(formatYear, formatMonth, formatDay, 0, 0, 0, 0)},
		{chrono.ISO8601OrdinalDateExtended, "0807-040", chrono.LocalDateTimeOf(formatYear, formatMonth, formatDay, 0, 0, 0, 0)},
		{chrono.ANSIC, "Fri Feb 09 01:05:02 0807", chrono.LocalDateTimeOf(formatYear, formatMonth, formatDay, formatHour, formatMin, formatSec, 0)},
		{chrono.Kitchen, "01:05AM", chrono.LocalDateTimeOf(1970, chrono.January, 1, formatHour, formatMin, 0, 0)},
	}
)