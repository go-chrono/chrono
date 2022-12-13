package chrono_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/go-chrono/chrono"
)

const (
	year  = 807
	month = chrono.February
	day   = 9
	hour  = 1
	min   = 5
	sec   = 0
	nsec  = 0
)

func setupCenturyParsing() {
	chrono.SetupCenturyParsing(800)
}

func checkYear(t *testing.T, date chrono.LocalDate) {
	if y, _, _ := date.Date(); y != year {
		t.Errorf("date.Date() year = %d, want %d", y, year)
	}
}

func checkYearDay(t *testing.T, date chrono.LocalDate) {
	if d := date.YearDay(); d != 40 {
		t.Errorf("date.YearDay() = %d, want %d", d, 40)
	}
}

func checkMonth(t *testing.T, date chrono.LocalDate) {
	if _, m, _ := date.Date(); m != month {
		t.Errorf("date.Date() month = %d, want %d", m, month)
	}
}

func checkDay(t *testing.T, date chrono.LocalDate) {
	if _, _, d := date.Date(); d != day {
		t.Errorf("date.Date() day = %d, want %d", d, day)
	}
}

func checkWeekday(t *testing.T, date chrono.LocalDate) {
	if d := date.Weekday(); d != chrono.Friday {
		t.Errorf("date.Weekday() = %s, want %s", d, chrono.Friday)
	}
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

var (
	dateSpecifiers = []struct {
		specifier string
		text      string
		check     func(*testing.T, chrono.LocalDate)
	}{
		{"%Y", "0807", checkYear},
		{"%-Y", "807", checkYear},
		{"%EY", "0807", checkYear},
		{"%-EY", "807", checkYear},
		{"%y", "07", checkYear},
		{"%-y", "7", checkYear},
		{"%Ey", "07", checkYear},
		{"%-Ey", "7", checkYear},
		{"%j", "040", checkYearDay},
		{"%-j", "40", checkYearDay},
		{"%m", "02", checkMonth},
		{"%-m", "2", checkMonth},
		{"%B", "February", checkMonth},
		{"%b", "Feb", checkMonth},
		{"%d", "09", checkDay},
		{"%-d", "9", checkDay},
		{"%u", "5", checkWeekday},
		{"%-u", "5", checkWeekday},
		{"%A", "Friday", checkWeekday},
		{"%a", "Fri", checkWeekday},
		{"%G", "0807", checkISOYear},
		{"%-G", "807", checkISOYear},
		{"%V", "06", checkISOWeek},
		{"%-V", "6", checkISOWeek},
	}

	timeSpecifiers = []struct {
		specifier string
		expected  string
	}{
		{"%P", "am"},
		{"%p", "AM"},
		{"%I", "01"},
		{"%-I", "1"},
		{"%H", "01"},
		{"%-H", "1"},
		{"%M", "05"},
		{"%-M", "5"},
		{"%S", "00"},
		{"%-S", "0"},
	}
)

func TestLocalDate_Parse_supported_specifiers(t *testing.T) {
	setupCenturyParsing()

	for _, tt := range dateSpecifiers {
		t.Run(tt.specifier, func(t *testing.T) {
			var d chrono.LocalDate
			if err := d.Parse(tt.specifier, tt.text); err != nil {
				t.Errorf("failed to parse date: %v", err)
			}

			tt.check(t, d)
		})
	}
}

func TestLocalDate_Format_supported_specifiers(t *testing.T) {
	for _, tt := range dateSpecifiers {
		t.Run(tt.specifier, func(t *testing.T) {
			if formatted := chrono.LocalDateOf(year, month, day).Format(tt.specifier); formatted != tt.text {
				t.Errorf("date.Format(%s) = %s, want %s", tt.specifier, formatted, tt.text)
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

				chrono.LocalDateOf(year, month, day).Format(tt.specifier)
			}()
		})
	}
}

func TestLocalTime_Format_supported_specifiers(t *testing.T) {
	for _, tt := range dateSpecifiers {
		t.Run(tt.specifier, func(t *testing.T) {
			func() {
				defer func() {
					if r := recover(); r == nil {
						t.Error("expecting panic that didn't occur")
					}
				}()

				chrono.LocalTimeOf(hour, min, sec, nsec).Format(tt.specifier)
			}()
		})
	}

	for _, tt := range timeSpecifiers {
		t.Run(tt.specifier, func(t *testing.T) {
			if formatted := chrono.LocalTimeOf(hour, min, sec, nsec).Format(tt.specifier); formatted != tt.expected {
				t.Errorf("time.Format(%s) = %s, want %s", tt.specifier, formatted, tt.expected)
			}
		})
	}
}

func TestLocalDateTime_Format_supported_specifiers(t *testing.T) {
	for _, tt := range dateSpecifiers {
		t.Run(tt.specifier, func(t *testing.T) {
			if formatted := chrono.LocalDateTimeOf(year, month, day, hour, min, sec, nsec).Format(tt.specifier); formatted != tt.text {
				t.Errorf("datetime.Format(%s) = %s, want %s", tt.specifier, formatted, tt.text)
			}
		})
	}

	for _, tt := range timeSpecifiers {
		t.Run(tt.specifier, func(t *testing.T) {
			if formatted := chrono.LocalDateTimeOf(year, month, day, hour, min, sec, nsec).Format(tt.specifier); formatted != tt.expected {
				t.Errorf("datetime.Format(%s) = %s, want %s", tt.specifier, formatted, tt.expected)
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
				t.Errorf("datetime.Format(%s) = %s, want %s", tt.layout, actual, tt.expected)
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
						t.Errorf("expecting '%s' error but got '%s'", tt.expected, err.Error())
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
				t.Errorf("expecting '%s' error but got '%s'", expected, err.Error())
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
				t.Errorf("datetime.Format(%s) = %s, want %s", tt.layout, formatted, tt.expected)
			}
		})
	}
}
