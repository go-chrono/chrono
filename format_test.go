package chrono_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/go-chrono/chrono"
)

// TODO test notes 2-5.

const (
	formatYear  = 807
	formatMonth = chrono.February
	formatDay   = 9
	formatHour  = 1
	formatMin   = 5
	formatSec   = 2
	formatNsec  = 0
)

func setupCenturyParsing() {
	chrono.SetupCenturyParsing(800)
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
	}
)

func TestLocalDate_Parse_supported_specifiers(t *testing.T) {
	setupCenturyParsing()

	for _, tt := range dateSpecifiers {
		t.Run(fmt.Sprintf("%s (%q)", tt.specifier, tt.textToParse), func(t *testing.T) {
			var date chrono.LocalDate
			if err := date.Parse(tt.specifier, tt.textToParse); err != nil {
				t.Errorf("failed to parse date: %v", err)
			}

			tt.checkParse(t, date)
		})
	}

	for _, tt := range timeSpecifiers {
		t.Run(tt.specifier, func(t *testing.T) {
			func() {
				defer func() {
					if r := recover(); r == nil {
						t.Error("expecting panic that didn't occur")
					}
				}()

				var date chrono.LocalDate
				date.Format(tt.specifier)
			}()
		})
	}
}

func TestLocalTime_Parse_supported_specifiers(t *testing.T) {
	setupCenturyParsing()

	for _, tt := range dateSpecifiers {
		t.Run(fmt.Sprintf("%s (%q)", tt.specifier, tt.textToParse), func(t *testing.T) {
			func() {
				defer func() {
					if r := recover(); r == nil {
						t.Error("expecting panic that didn't occur")
					}
				}()

				var time chrono.LocalTime
				time.Format(tt.specifier)
			}()
		})
	}

	for _, tt := range timeSpecifiers {
		t.Run(tt.specifier, func(t *testing.T) {
			var time chrono.LocalTime
			if err := time.Parse(tt.specifier, tt.text); err != nil {
				t.Errorf("failed to parse time: %v", err)
			}

			tt.checkParse(t, time)
		})
	}
}

func TestLocalDateTime_Parse_supported_specifiers(t *testing.T) {
	setupCenturyParsing()

	for _, tt := range dateSpecifiers {
		t.Run(fmt.Sprintf("%s (%q)", tt.specifier, tt.textToParse), func(t *testing.T) {
			var dt chrono.LocalDateTime
			if err := dt.Parse(tt.specifier, tt.textToParse); err != nil {
				t.Errorf("failed to parse date: %v", err)
			}

			date, _ := dt.Split()
			tt.checkParse(t, date)
		})
	}

	for _, tt := range timeSpecifiers {
		t.Run(tt.specifier, func(t *testing.T) {
			var dt chrono.LocalDateTime
			if err := dt.Parse(tt.specifier, tt.text); err != nil {
				t.Errorf("failed to parse time: %v", err)
			}

			_, time := dt.Split()
			tt.checkParse(t, time)
		})
	}
}

func TestLocalDate_Format_supported_specifiers(t *testing.T) {
	for _, tt := range dateSpecifiers {
		t.Run(fmt.Sprintf("%s (%q)", tt.specifier, tt.textToParse), func(t *testing.T) {
			if formatted := chrono.LocalDateOf(formatYear, formatMonth, formatDay).Format(tt.specifier); formatted != tt.expectedFormatted {
				t.Errorf("date.Format(%s) = %s, want %q", tt.specifier, formatted, tt.expectedFormatted)
			}
		})
	}

	for _, tt := range timeSpecifiers {
		t.Run(tt.specifier, func(t *testing.T) {
			func() {
				defer func() {
					if r := recover(); r == nil {
						t.Error("expecting panic that didn't occur")
					}
				}()

				chrono.LocalDateOf(formatYear, formatMonth, formatDay).Format(tt.specifier)
			}()
		})
	}
}

func TestLocalTime_Format_supported_specifiers(t *testing.T) {
	for _, tt := range dateSpecifiers {
		t.Run(fmt.Sprintf("%s (%q)", tt.specifier, tt.textToParse), func(t *testing.T) {
			func() {
				defer func() {
					if r := recover(); r == nil {
						t.Error("expecting panic that didn't occur")
					}
				}()

				chrono.LocalTimeOf(formatHour, formatMin, formatSec, formatNsec).Format(tt.specifier)
			}()
		})
	}

	for _, tt := range timeSpecifiers {
		t.Run(tt.specifier, func(t *testing.T) {
			if formatted := chrono.LocalTimeOf(formatHour, formatMin, formatSec, formatNsec).Format(tt.specifier); formatted != tt.text {
				t.Errorf("time.Format(%s) = %s, want %q", tt.specifier, formatted, tt.text)
			}
		})
	}
}

func TestLocalDateTime_Format_supported_specifiers(t *testing.T) {
	for _, tt := range dateSpecifiers {
		t.Run(fmt.Sprintf("%s (%q)", tt.specifier, tt.textToParse), func(t *testing.T) {
			if formatted := chrono.LocalDateTimeOf(formatYear, formatMonth, formatDay, formatHour, formatMin, formatSec, formatNsec).Format(tt.specifier); formatted != tt.expectedFormatted {
				t.Errorf("datetime.Format(%s) = %s, want %q", tt.specifier, formatted, tt.expectedFormatted)
			}
		})
	}

	for _, tt := range timeSpecifiers {
		t.Run(tt.specifier, func(t *testing.T) {
			if formatted := chrono.LocalDateTimeOf(formatYear, formatMonth, formatDay, formatHour, formatMin, formatSec, formatNsec).Format(tt.specifier); formatted != tt.text {
				t.Errorf("datetime.Format(%s) = %s, want %q", tt.specifier, formatted, tt.text)
			}
		})
	}
}

func Test_format_literals(t *testing.T) {
	for _, tt := range []struct {
		name     string
		value    interface{ Format(string) string }
		layout   string
		expected string
	}{
		{
			name:     "date",
			value:    chrono.LocalDateOf(2020, chrono.March, 18),
			layout:   "str1 %Y str2 100%%foo",
			expected: "str1 2020 str2 100%foo",
		},
		{
			name:     "time",
			value:    chrono.LocalTimeOf(12, 30, 15, 0),
			layout:   "str1 %H str2 100%%foo",
			expected: "str1 12 str2 100%foo",
		},
		{
			name:     "datetime",
			value:    chrono.LocalDateTimeOf(2020, chrono.March, 18, 12, 30, 15, 0),
			layout:   "str1 %Y str2 100%%foo",
			expected: "str1 2020 str2 100%foo",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if actual := tt.value.Format(tt.layout); actual != tt.expected {
				t.Errorf("datetime.Format(%s) = %s, want %q", tt.layout, actual, tt.expected)
			}
		})
	}
}

func Test_parse_cannot_parse_error(t *testing.T) {
	for _, tt := range []struct {
		name     string
		layout   string
		value    string
		expected string
	}{
		{"none", "foo bar", "foo", "parsing time \"foo\" as \"foo bar\": cannot parse \"foo\" as \"foo bar\""},
		{"partial", "bar", "foo", "parsing time \"foo\" as \"bar\": cannot parse \"foo\" as \"bar\""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var date chrono.LocalDate
			var time chrono.LocalTime
			var datetime chrono.LocalDateTime

			for _, v := range []interface {
				Parse(layout, value string) error
			}{
				&date,
				&time,
				&datetime,
			} {
				t.Run(reflect.TypeOf(v).Elem().Name(), func(t *testing.T) {
					if err := v.Parse(tt.layout, tt.value); err == nil {
						t.Errorf("expecting error but got nil")
					} else if !strings.Contains(err.Error(), tt.expected) {
						t.Errorf("expecting %q error but got %q", tt.expected, err.Error())
					}
				})
			}
		})
	}
}

func Test_parse_extra_text_error(t *testing.T) {
	var date chrono.LocalDate
	var time chrono.LocalTime
	var datetime chrono.LocalDateTime

	for _, v := range []interface {
		Parse(layout, value string) error
	}{
		&date,
		&time,
		&datetime,
	} {
		t.Run(reflect.TypeOf(v).Elem().Name(), func(t *testing.T) {
			expected := "parsing time \"foo bar\": extra text: \" bar\""
			if err := v.Parse("foo", "foo bar"); err == nil {
				t.Errorf("expecting error but got nil")
			} else if !strings.Contains(err.Error(), expected) {
				t.Errorf("expecting %q error but got %q", expected, err.Error())
			}
		})
	}
}

func TestLocalDateTime_Parse_predefined_layouts(t *testing.T) {
	date := chrono.LocalDateOf(2022, chrono.June, 18)
	time := chrono.LocalTimeOf(21, 05, 30, 0)

	for _, tt := range []struct {
		layout   string
		value    string
		expected chrono.LocalDateTime
	}{
		{chrono.ISO8601Date, "20220618", chrono.OfLocalDateAndTime(date, chrono.LocalTime{})},
		{chrono.ISO8601DateExtended, "2022-06-18", chrono.OfLocalDateAndTime(date, chrono.LocalTime{})},
		{chrono.ISO8601Time, "T210530", chrono.OfLocalDateAndTime(chrono.LocalDate(0), time)},
		{chrono.ISO8601TimeExtended, "T21:05:30", chrono.OfLocalDateAndTime(chrono.LocalDate(0), time)},
		{chrono.ISO8601DateTime, "20220618T210530", chrono.OfLocalDateAndTime(date, time)},
		{chrono.ISO8601DateTimeExtended, "2022-06-18T21:05:30", chrono.OfLocalDateAndTime(date, time)},
	} {
		t.Run(tt.layout, func(t *testing.T) {
			var datetime chrono.LocalDateTime
			if err := datetime.Parse(tt.layout, tt.value); err != nil {
				t.Errorf("datetime.Parse(%s, %s) = %v, want nil", tt.layout, tt.value, err)
			} else if datetime.Compare(tt.expected) != 0 {
				t.Errorf("expecting %v, but got %v", tt.expected, datetime)
			}
		})
	}
}

func TestLocalDate_Parse_default_values(t *testing.T) {
	for _, tt := range []struct {
		name     string
		layout   string
		value    string
		expected chrono.LocalDate
	}{
		{"nothing", "", "", chrono.LocalDateOf(1970, chrono.January, 1)},
		{"only year", "%Y", "2020", chrono.LocalDateOf(2020, chrono.January, 1)},
		{"only month", "%m", "04", chrono.LocalDateOf(1970, chrono.April, 1)},
		{"only day", "%d", "22", chrono.LocalDateOf(1970, chrono.January, 22)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var date chrono.LocalDate
			if err := date.Parse(tt.layout, tt.value); err != nil {
				t.Errorf("date.Parse(%s, %s) = %v, want nil", tt.layout, tt.value, err)
			} else if date != tt.expected {
				t.Errorf("expecting %v, but got %v", tt.expected, date)
			}
		})
	}
}

func TestLocalDateTime_Format_predefined_layouts(t *testing.T) {
	date := chrono.LocalDateOf(2022, chrono.June, 18)
	time := chrono.LocalTimeOf(21, 05, 30, 0)

	for _, tt := range []struct {
		layout   string
		expected string
	}{
		{chrono.ISO8601Date, "20220618"},
		{chrono.ISO8601DateExtended, "2022-06-18"},
		{chrono.ISO8601Time, "T210530"},
		{chrono.ISO8601TimeExtended, "T21:05:30"},
		{chrono.ISO8601DateTime, "20220618T210530"},
		{chrono.ISO8601DateTimeExtended, "2022-06-18T21:05:30"},
	} {
		t.Run(tt.layout, func(t *testing.T) {
			if formatted := chrono.OfLocalDateAndTime(date, time).Format(tt.layout); formatted != tt.expected {
				t.Errorf("datetime.Format(%s) = %s, want %q", tt.layout, formatted, tt.expected)
			}
		})
	}
}

func TestLocalDate_Format_eras(t *testing.T) {
	t.Run("CE", func(t *testing.T) {
		date := chrono.LocalDateOf(2022, chrono.June, 18)
		if formatted := date.Format("%EY %EC"); formatted != "2022 CE" {
			t.Errorf("got %q, want '2022 CE'", formatted)
		}
	})

	t.Run("BCE", func(t *testing.T) {
		date := chrono.LocalDateOf(-2021, chrono.June, 18)
		if formatted := date.Format("%EY %EC"); formatted != "2022 BCE" {
			t.Errorf("got %q, want '2022 BCE'", formatted)
		}
	})

	t.Run("zero", func(t *testing.T) {
		date := chrono.LocalDateOf(0, chrono.June, 18)
		if formatted := date.Format("%EY %EC"); formatted != "0001 BCE" {
			t.Errorf("got %q, want '1 BCE'", formatted)
		}
	})
}
