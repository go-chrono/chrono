package chrono_test

import (
	"fmt"
	"testing"

	"github.com/go-chrono/chrono"
)

func TestLocalDate(t *testing.T) {
	for _, tt := range []struct {
		year       int32
		month      chrono.Month
		day        uint8
		weekday    chrono.Weekday
		isLeapYear bool
		yearDay    uint32
		isoYear    int32
		isoWeek    uint8
	}{
		{-4713, chrono.November, 24, chrono.Monday, false, 328, -4713, 48}, // 4714 BCE
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
	} {
		t.Run(fmt.Sprintf("%+05d-%02d-%02d", tt.year, tt.month, tt.day), func(t *testing.T) {
			d := chrono.LocalDateOf(tt.year, tt.month, tt.day)

			year, month, day := d.Date()
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

			if weekday := d.Weekday(); weekday != tt.weekday {
				t.Errorf("d.Weekday() = %s, want %s", weekday, tt.weekday)
				t.Fail()
			}

			if isLeapYear := d.IsLeapYear(); isLeapYear != tt.isLeapYear {
				t.Errorf("d.YearDay() = %t, want %t", isLeapYear, tt.isLeapYear)
				t.Fail()
			}

			if yearDay := d.YearDay(); yearDay != tt.yearDay {
				t.Errorf("d.YearDay() = %d, want %d", yearDay, tt.yearDay)
				t.Fail()
			}

			isoYear, isoWeek := d.ISOWeek()
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

func TestDefaultLocalDate(t *testing.T) {
	var d chrono.LocalDate

	expectedYear, expectedMonth, expectedDay := int32(1970), chrono.January, uint8(1)

	year, month, day := d.Date()
	if year != expectedYear {
		t.Errorf("d.Date() year = %d, want %d", year, expectedYear)
		t.Fail()
	}

	if month != expectedMonth {
		t.Errorf("d.Date() month = %s, want %s", month, expectedMonth)
		t.Fail()
	}

	if day != expectedDay {
		t.Errorf("d.Date() day = %d, want %d", day, expectedDay)
		t.Fail()
	}
}
