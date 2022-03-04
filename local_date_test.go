package chrono_test

import (
	"fmt"
	"testing"

	"github.com/go-chrono/chrono"
)

func TestLocalDate(t *testing.T) {
	for _, tt := range []struct {
		year       int
		month      chrono.Month
		day        int
		weekday    chrono.Weekday
		isLeapYear bool
		yearDay    int
		isoYear    int
		isoWeek    int
	}{
		{-4713, chrono.November, 24, chrono.Monday, false, 328, -4713, 48}, // 4714 BCE
		{+5874898, chrono.June, 3, chrono.Tuesday, false, 154, +5874898, 23},
		{+1858, chrono.November, 17, chrono.Wednesday, false, 321, +1858, 46},
		{+1968, chrono.May, 24, chrono.Friday, true, 145, +1968, 21},
		{+1950, chrono.January, 1, chrono.Sunday, false, 1, +1949, 52},
		{+1958, chrono.January, 1, chrono.Wednesday, false, 1, +1958, 1},
		{+1582, chrono.October, 15, chrono.Friday, false, 288, +1582, 41},
		{+1, chrono.January, 1, chrono.Monday, false, 1, +1, 1},
		{+1970, chrono.January, 1, chrono.Thursday, false, 1, +1970, 1},
		{+200, chrono.March, 1, chrono.Saturday, false, 60, +200, 9},
		{+2020, chrono.December, 31, chrono.Thursday, true, 366, +2020, 53},
		{+2021, chrono.January, 1, chrono.Friday, false, 1, +2020, 53},
		{+2000, chrono.February, 29, chrono.Tuesday, true, 60, +2000, 9},
		{+2000, chrono.March, 1, chrono.Wednesday, true, 61, +2000, 9},
	} {
		t.Run(fmt.Sprintf("%+05d-%02d-%02d", tt.year, tt.month, tt.day), func(t *testing.T) {
			date := chrono.LocalDateOf(tt.year, tt.month, tt.day)

			year, month, day := date.Date()
			if year != tt.year {
				t.Errorf("date.Date() year = %d, want %d", year, tt.year)
			}

			if month != tt.month {
				t.Errorf("date.Date() month = %s, want %s", month, tt.month)
			}

			if day != tt.day {
				t.Errorf("date.Date() day = %d, want %d", day, tt.day)
			}

			if weekday := date.Weekday(); weekday != tt.weekday {
				t.Errorf("date.Weekday() = %s, want %s", weekday, tt.weekday)
			}

			if isLeapYear := date.IsLeapYear(); isLeapYear != tt.isLeapYear {
				t.Errorf("date.YearDay() = %t, want %t", isLeapYear, tt.isLeapYear)
			}

			if yearDay := date.YearDay(); yearDay != tt.yearDay {
				t.Errorf("date.YearDay() = %d, want %d", yearDay, tt.yearDay)
			}

			isoYear, isoWeek := date.ISOWeek()
			if isoYear != tt.isoYear {
				t.Errorf("date.ISOWeek() year = %d, want %d", isoYear, tt.isoYear)
			}

			if isoWeek != tt.isoWeek {
				t.Errorf("date.ISOWeek() week = %d, want %d", isoWeek, tt.isoWeek)
			}
		})
	}
}

func TestLocalDateOf(t *testing.T) {
	for _, tt := range []struct {
		name  string
		year  int
		month chrono.Month
		day   int
	}{
		{"year underflows", -4714, chrono.January, 1},
		{"year & month underflows", -4713, chrono.October, 1},
		{"year & month & day underflows", -4713, chrono.November, 23},
		{"year overflows", 5874899, chrono.January, 1},
		{"year & month overflows", 5874898, chrono.July, 1},
		{"year & month & day overflows", 5874898, chrono.June, 4},
		{"month underflows", 2001, 0, 1},
		{"month overflows", 2001, 13, 1},
		{"day underflows", 2001, chrono.February, 0},
		{"day overflows", 2001, chrono.February, 29},
		{"day overflows leap year", 2004, chrono.February, 30},
	} {
		t.Run(tt.name, func(t *testing.T) {
			func() {
				defer func() {
					if r := recover(); r == nil {
						t.Error("expecting panic that didn't occur")
					}
				}()

				chrono.LocalDateOf(tt.year, tt.month, tt.day)
			}()
		})
	}
}

func TestLocalDate_Date(t *testing.T) {
	for _, tt := range []struct {
		name  string
		date  chrono.LocalDate
		year  int
		month chrono.Month
		day   int
	}{
		{"default value", chrono.LocalDate(0), 1970, chrono.January, 1},
		{"minimum value", chrono.MinLocalDate(), -4713, chrono.November, 24},
		{"maximum value", chrono.MaxLocalDate(), 5874898, chrono.June, 3},
	} {
		t.Run(tt.name, func(t *testing.T) {
			year, month, day := tt.date.Date()
			if year != tt.year {
				t.Errorf("date.Date() year = %d, want %d", year, tt.year)
			}

			if month != tt.month {
				t.Errorf("date.Date() month = %s, want %s", month, tt.month)
			}

			if day != tt.day {
				t.Errorf("date.Date() day = %d, want %d", day, tt.day)
			}
		})
	}
}

func TestLocalDate_Add(t *testing.T) {
	for _, tt := range []struct {
		name      string
		date      chrono.LocalDate
		addYears  int
		addMonths int
		addDays   int
		expected  chrono.LocalDate
	}{
		{"nothing", chrono.LocalDateOf(2020, chrono.March, 18), 0, 0, 0, chrono.LocalDateOf(2020, chrono.March, 18)},
		{"add years", chrono.LocalDateOf(2020, chrono.March, 18), 105, 0, 0, chrono.LocalDateOf(2125, chrono.March, 18)},
		{"sub years", chrono.LocalDateOf(2020, chrono.March, 18), -280, 0, 0, chrono.LocalDateOf(1740, chrono.March, 18)},
		{"add months", chrono.LocalDateOf(2020, chrono.March, 18), 0, 6, 0, chrono.LocalDateOf(2020, chrono.September, 18)},
		{"sub months", chrono.LocalDateOf(2020, chrono.March, 18), 0, -2, 0, chrono.LocalDateOf(2020, chrono.January, 18)},
		{"add days", chrono.LocalDateOf(2020, chrono.March, 18), 0, 0, 8, chrono.LocalDateOf(2020, chrono.March, 26)},
		{"sub days", chrono.LocalDateOf(2020, chrono.March, 18), 0, 0, -15, chrono.LocalDateOf(2020, chrono.March, 3)},
		{"time package example", chrono.LocalDateOf(2011, chrono.January, 1), -1, 2, 3, chrono.LocalDateOf(2010, chrono.March, 4)},
		{"normalized time package example", chrono.LocalDateOf(2011, chrono.October, 31), 0, 1, 0, chrono.LocalDateOf(2011, chrono.November, 31)},
		{"overflow day", chrono.LocalDateOf(2020, chrono.March, 18), 0, 0, 20, chrono.LocalDateOf(2020, chrono.April, 7)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if date := tt.date.Add(tt.addYears, tt.addMonths, tt.addDays); date != tt.expected {
				t.Errorf("date = %s, date.Add(%d, %d, %d) = %s, want %s", tt.date, tt.addYears, tt.addMonths, tt.addDays, date, tt.expected)
			}
		})
	}
}
