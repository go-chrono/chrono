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
				t.Errorf("d.Date() year = %d, want %d", year, tt.year)
				t.Fail()
			}

			if month != tt.month {
				t.Errorf("d.Date() month = %s, want %s", month, tt.month)
				t.Fail()
			}

			if day != tt.day {
				t.Errorf("d.Date() day = %d, want %d", day, tt.day)
				t.Fail()
			}

			if weekday := date.Weekday(); weekday != tt.weekday {
				t.Errorf("d.Weekday() = %s, want %s", weekday, tt.weekday)
				t.Fail()
			}

			if isLeapYear := date.IsLeapYear(); isLeapYear != tt.isLeapYear {
				t.Errorf("d.YearDay() = %t, want %t", isLeapYear, tt.isLeapYear)
				t.Fail()
			}

			if yearDay := date.YearDay(); yearDay != tt.yearDay {
				t.Errorf("d.YearDay() = %d, want %d", yearDay, tt.yearDay)
				t.Fail()
			}

			isoYear, isoWeek := date.ISOWeek()
			if isoYear != tt.isoYear {
				t.Errorf("d.ISOWeek() year = %d, want %d", isoYear, tt.isoYear)
				t.Fail()
			}

			if isoWeek != tt.isoWeek {
				t.Errorf("d.ISOWeek() week = %d, want %d", isoWeek, tt.isoWeek)
				t.Fail()
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
						t.Fatalf("expecting panic that didn't occur")
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
				t.Errorf("d.Date() year = %d, want %d", year, tt.year)
				t.Fail()
			}

			if month != tt.month {
				t.Errorf("d.Date() month = %s, want %s", month, tt.month)
				t.Fail()
			}

			if day != tt.day {
				t.Errorf("d.Date() day = %d, want %d", day, tt.day)
				t.Fail()
			}
		})
	}
}
