package chrono_test

import (
	"testing"

	"github.com/go-chrono/chrono"
)

const (
	formatYear  = 807
	formatMonth = chrono.February
	formatDay   = 9
	formatHour  = 1
	formatMin   = 5
	formatSec   = 0
	formatNsec  = 0
)

var (
	dateSpecifiers = []struct {
		specifier string
		expected  string
	}{
		{"%Y", "0807"},
		{"%-Y", "807"},
		{"%EY", "0807"},
		{"%-EY", "807"},
		{"%y", "07"},
		{"%-y", "7"},
		{"%Ey", "07"},
		{"%-Ey", "7"},
		{"%j", "040"},
		{"%-j", "40"},
		{"%m", "02"},
		{"%-m", "2"},
		{"%B", "February"},
		{"%b", "Feb"},
		{"%d", "09"},
		{"%-d", "9"},
		{"%u", "5"},
		{"%-u", "5"},
		{"%A", "Friday"},
		{"%a", "Fri"},
		{"%G", "0807"},
		{"%-G", "807"},
		{"%V", "06"},
		{"%-V", "6"},
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

func TestLocalDate_Format_supported_specifiers(t *testing.T) {
	for _, tt := range dateSpecifiers {
		t.Run(tt.specifier, func(t *testing.T) {
			if formatted := chrono.LocalDateOf(formatYear, formatMonth, formatDay).Format(tt.specifier); formatted != tt.expected {
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

				chrono.LocalDateOf(formatYear, formatMonth, formatDay).Format(tt.specifier)
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

				chrono.LocalTimeOf(formatHour, formatMin, formatSec, formatNsec).Format(tt.specifier)
			}()
		})
	}

	for _, tt := range timeSpecifiers {
		t.Run(tt.specifier, func(t *testing.T) {
			if formatted := chrono.LocalTimeOf(formatHour, formatMin, formatSec, formatNsec).Format(tt.specifier); formatted != tt.expected {
				t.Errorf("time.Format(%s) = %s, want %s", tt.specifier, formatted, tt.expected)
			}
		})
	}
}

func TestLocalDateTime_Format_supported_specifiers(t *testing.T) {
	for _, tt := range dateSpecifiers {
		t.Run(tt.specifier, func(t *testing.T) {
			if formatted := chrono.LocalDateTimeOf(formatYear, formatMonth, formatDay, formatHour, formatMin, formatSec, formatNsec).Format(tt.specifier); formatted != tt.expected {
				t.Errorf("datetime.Format(%s) = %s, want %s", tt.specifier, formatted, tt.expected)
			}
		})
	}

	for _, tt := range timeSpecifiers {
		t.Run(tt.specifier, func(t *testing.T) {
			if formatted := chrono.LocalDateTimeOf(formatYear, formatMonth, formatDay, formatHour, formatMin, formatSec, formatNsec).Format(tt.specifier); formatted != tt.expected {
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
