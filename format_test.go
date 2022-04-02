package chrono_test

import (
	"testing"

	"github.com/go-chrono/chrono"
)

const (
	year  = 2020
	month = chrono.March
	day   = 13
	hour  = 12
	min   = 30
	sec   = 15
	nsec  = 0
)

var (
	dateSpecifiers = []struct {
		specifier string
		expected  string
	}{
		{"%Y", "2020"},
		{"%EY", "2020"},
		{"%y", "20"},
		{"%Ey", "20"},
		{"%j", "073"},
		{"%m", "03"},
		{"%B", "March"},
		{"%b", "Mar"},
		{"%d", "13"},
		{"%u", "5"},
		{"%A", "Friday"},
		{"%a", "Fri"},
		{"%G", "2020"},
		{"%V", "11"},
	}

	timeSpecifiers = []struct {
		specifier string
		expected  string
	}{
		{"%P", "pm"},
		{"%p", "PM"},
		{"%I", "12"},
		{"%H", "12"},
		{"%M", "30"},
		{"%S", "15"},
	}
)

func TestLocalDate_Format_supported_specifiers(t *testing.T) {
	for _, tt := range dateSpecifiers {
		t.Run(tt.specifier, func(t *testing.T) {
			if formatted := chrono.LocalDateOf(year, month, day).Format(tt.specifier); formatted != tt.expected {
				t.Errorf("date.Format(%s) = %s, want %s", tt.specifier, formatted, tt.expected)
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
			if formatted := chrono.LocalDateTimeOf(year, month, day, hour, min, sec, nsec).Format(tt.specifier); formatted != tt.expected {
				t.Errorf("datetime.Format(%s) = %s, want %s", tt.specifier, formatted, tt.expected)
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

func TestLocalDateTime_Format_literals(t *testing.T) {
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
