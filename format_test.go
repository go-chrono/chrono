package chrono_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/go-chrono/chrono"
)

func TestLocalDate_Parse_supported_specifiers(t *testing.T) {
	setupCenturyParsing()
	defer tearDownCenturyParsing()

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
	defer tearDownCenturyParsing()

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
	defer tearDownCenturyParsing()

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

func TestOffsetTime_Parse_suported_specifiers(t *testing.T) {
	setupCenturyParsing()
	defer tearDownCenturyParsing()

	for _, tt := range dateSpecifiers {
		t.Run(fmt.Sprintf("%s (%q)", tt.specifier, tt.textToParse), func(t *testing.T) {
			func() {
				defer func() {
					if r := recover(); r == nil {
						t.Error("expecting panic that didn't occur")
					}
				}()

				var time chrono.OffsetTime
				time.Format(tt.specifier)
			}()
		})
	}

	for _, tt := range timeSpecifiers {
		t.Run(tt.specifier, func(t *testing.T) {
			var time chrono.OffsetTime
			if err := time.Parse(tt.specifier, tt.text); err != nil {
				t.Errorf("failed to parse time: %v", err)
			}

			tt.checkParse(t, time)
		})
	}

	for _, tt := range offsetTimeSpecifiersUTC {
		t.Run(tt.specifier, func(t *testing.T) {
			var time chrono.OffsetTime
			if err := time.Parse(tt.specifier, tt.text); err != nil {
				t.Errorf("failed to parse time: %v", err)
			}

			tt.checkParse(t, time)
		})
	}
}

func TestOffsetDateTime_Parse_supported_specifiers(t *testing.T) {
	setupCenturyParsing()
	defer tearDownCenturyParsing()

	for _, tt := range dateSpecifiers {
		t.Run(fmt.Sprintf("%s (%q)", tt.specifier, tt.textToParse), func(t *testing.T) {
			var dt chrono.OffsetDateTime
			if err := dt.Parse(tt.specifier, tt.textToParse); err != nil {
				t.Errorf("failed to parse date: %v", err)
			}

			date, _ := dt.Split()
			tt.checkParse(t, date)
		})
	}

	for _, tt := range timeSpecifiers {
		t.Run(tt.specifier, func(t *testing.T) {
			var dt chrono.OffsetDateTime
			if err := dt.Parse(tt.specifier, tt.text); err != nil {
				t.Errorf("failed to parse time: %v", err)
			}

			_, time := dt.Split()
			tt.checkParse(t, time)
		})
	}

	for _, tt := range offsetTimeSpecifiersUTC {
		t.Run(tt.specifier, func(t *testing.T) {
			var dt chrono.OffsetDateTime
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

				chrono.LocalTimeOf(formatHour, formatMin, formatSec, formatNanos).Format(tt.specifier)
			}()
		})
	}

	for _, tt := range timeSpecifiers {
		t.Run(tt.specifier, func(t *testing.T) {
			if formatted := chrono.LocalTimeOf(formatHour, formatMin, formatSec, formatNanos).Format(tt.specifier); formatted != tt.text {
				t.Errorf("time.Format(%s) = %s, want %q", tt.specifier, formatted, tt.text)
			}
		})
	}
}

func TestLocalDateTime_Format_supported_specifiers(t *testing.T) {
	for _, tt := range dateSpecifiers {
		t.Run(fmt.Sprintf("%s (%q)", tt.specifier, tt.textToParse), func(t *testing.T) {
			if formatted := chrono.LocalDateTimeOf(formatYear, formatMonth, formatDay, formatHour, formatMin, formatSec, formatNanos).
				Format(tt.specifier); formatted != tt.expectedFormatted {
				t.Errorf("datetime.Format(%s) = %s, want %q", tt.specifier, formatted, tt.expectedFormatted)
			}
		})
	}

	for _, tt := range timeSpecifiers {
		t.Run(tt.specifier, func(t *testing.T) {
			if formatted := chrono.LocalDateTimeOf(formatYear, formatMonth, formatDay, formatHour, formatMin, formatSec, formatNanos).
				Format(tt.specifier); formatted != tt.text {
				t.Errorf("datetime.Format(%s) = %s, want %q", tt.specifier, formatted, tt.text)
			}
		})
	}
}

func TestOffsetTime_Format_supported_specifiers(t *testing.T) {
	for _, tt := range dateSpecifiers {
		t.Run(fmt.Sprintf("%s (%q)", tt.specifier, tt.textToParse), func(t *testing.T) {
			func() {
				defer func() {
					if r := recover(); r == nil {
						t.Error("expecting panic that didn't occur")
					}
				}()

				chrono.OffsetTimeOf(formatHour, formatMin, formatSec, formatNanos,
					formatOffsetHours, formatOffsetMins).Format(tt.specifier)
			}()
		})
	}

	for _, tt := range timeSpecifiers {
		t.Run(tt.specifier, func(t *testing.T) {
			if formatted := chrono.OffsetTimeOf(formatHour, formatMin, formatSec, formatNanos,
				formatOffsetHours, formatOffsetMins).Format(tt.specifier); formatted != tt.text {
				t.Errorf("time.Format(%s) = %s, want %q", tt.specifier, formatted, tt.text)
			}
		})
	}

	for _, tt := range offsetTimeSpecifiersUTC {
		t.Run(tt.specifier, func(t *testing.T) {
			if formatted := chrono.OffsetTimeOf(formatHour, formatMin, formatSec, formatNanos,
				formatOffsetHours, formatOffsetMins).Format(tt.specifier); formatted != tt.expectedFormatted {
				t.Errorf("time.Format(%s) = %s, want %q", tt.specifier, formatted, tt.expectedFormatted)
			}
		})
	}
}

func TestOffsetDateTime_Format_supported_specifiers(t *testing.T) {
	for _, tt := range dateSpecifiers {
		t.Run(fmt.Sprintf("%s (%q)", tt.specifier, tt.textToParse), func(t *testing.T) {
			if formatted := chrono.OffsetDateTimeOf(formatYear, formatMonth, formatDay, formatHour, formatMin, formatSec, formatNanos,
				formatOffsetHours, formatOffsetMins).Format(tt.specifier); formatted != tt.expectedFormatted {
				t.Errorf("datetime.Format(%s) = %s, want %q", tt.specifier, formatted, tt.expectedFormatted)
			}
		})
	}

	for _, tt := range timeSpecifiers {
		t.Run(tt.specifier, func(t *testing.T) {
			if formatted := chrono.OffsetDateTimeOf(formatYear, formatMonth, formatDay, formatHour, formatMin, formatSec, formatNanos,
				formatOffsetHours, formatOffsetMins).Format(tt.specifier); formatted != tt.text {
				t.Errorf("datetime.Format(%s) = %s, want %q", tt.specifier, formatted, tt.text)
			}
		})
	}

	for _, tt := range offsetTimeSpecifiersUTC {
		t.Run(tt.specifier, func(t *testing.T) {
			if formatted := chrono.OffsetDateTimeOf(formatYear, formatMonth, formatDay, formatHour, formatMin, formatSec, formatNanos,
				formatOffsetHours, formatOffsetMins).Format(tt.specifier); formatted != tt.expectedFormatted {
				t.Errorf("datetime.Format(%s) = %s, want %q", tt.specifier, formatted, tt.expectedFormatted)
			}
		})
	}
}

func TestOffsetTime_Parse_offset_formats(t *testing.T) {
	for _, tt := range offsetTimeSpecifiers {
		t.Run(fmt.Sprintf("%s-%s", tt.specifier, tt.text), func(t *testing.T) {
			var time chrono.OffsetTime
			if err := time.Parse(tt.specifier, tt.text); !tt.expectErr && err != nil {
				t.Errorf("failed to parse time: %v", err)
			} else if err == nil && tt.expectErr {
				t.Errorf("expecting error")
			}

			if v := time.Offset(); v != chrono.Offset(tt.expected) {
				t.Errorf("time.Offset() = %d, want %d", v, tt.expected)
			}
		})
	}
}

func TestOffsetDateTime_Parse_offset_formats(t *testing.T) {
	for _, tt := range offsetTimeSpecifiers {
		t.Run(fmt.Sprintf("%s-%s", tt.specifier, tt.text), func(t *testing.T) {
			var dt chrono.OffsetDateTime
			if err := dt.Parse(tt.specifier, tt.text); !tt.expectErr && err != nil {
				t.Errorf("failed to parse time: %v", err)
			} else if err == nil && tt.expectErr {
				t.Errorf("expecting error")
			}

			if v := dt.Offset(); v != chrono.Offset(tt.expected) {
				t.Errorf("dt.Offset() = %d, want %d", v, tt.expected)
			}
		})
	}
}

func TestOffsetTime_Format_offset_formats(t *testing.T) {
	for _, tt := range offsetTimeSpecifiers {
		t.Run(fmt.Sprintf("%s-%s", tt.specifier, tt.text), func(t *testing.T) {
			time := chrono.OfTimeOffset(chrono.LocalTimeOf(formatHour, formatMin, formatSec, formatNanos), chrono.Offset(tt.expected))
			if formatted := time.Format(tt.specifier); formatted != tt.expectedFormatted {
				t.Errorf("time = %s, time.Format(%s) = %s, want %q", time, tt.specifier, formatted, tt.expectedFormatted)
			}
		})
	}
}

func TestOffsetDateTime_Format_offset_formats(t *testing.T) {
	for _, tt := range offsetTimeSpecifiers {
		t.Run(fmt.Sprintf("%s-%s", tt.specifier, tt.text), func(t *testing.T) {
			dt := chrono.OfLocalDateTimeOffset(
				chrono.LocalDateOf(formatYear, formatMonth, formatDay),
				chrono.LocalTimeOf(formatHour, formatMin, formatSec, formatNanos),
				tt.expected,
			)

			if formatted := dt.Format(tt.specifier); formatted != tt.expectedFormatted {
				t.Errorf("datetime.Format(%s) = %s, want %q", tt.specifier, formatted, tt.expectedFormatted)
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
	for _, tt := range predefinedLayouts {
		t.Run(tt.layout, func(t *testing.T) {
			var datetime chrono.LocalDateTime
			if err := datetime.Parse(tt.layout, tt.text); err != nil {
				t.Errorf("datetime.Parse(%s, %s) = %v, want nil", tt.layout, tt.text, err)
			} else if datetime.Compare(tt.datetime.Local()) != 0 {
				t.Errorf("expecting %v, but got %v", tt.datetime, datetime)
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
	for _, tt := range predefinedLayouts {
		t.Run(tt.layout, func(t *testing.T) {
			if formatted := tt.datetime.Local().Format(tt.layout); formatted != tt.text {
				t.Errorf("datetime.Format(%s) = %s, want %q", tt.layout, formatted, tt.text)
			}
		})
	}
}

func TestLocalTime_Format_12HourClock(t *testing.T) {
	t.Run("am", func(t *testing.T) {
		time := chrono.LocalTimeOf(10, 0, 0, 0)
		if formatted := time.Format("%I %P"); formatted != "10 am" {
			t.Errorf("got %q, want '10 am'", formatted)
		}
	})

	t.Run("AM", func(t *testing.T) {
		time := chrono.LocalTimeOf(10, 0, 0, 0)
		if formatted := time.Format("%I %p"); formatted != "10 AM" {
			t.Errorf("got %q, want '10 AM'", formatted)
		}
	})

	t.Run("pm", func(t *testing.T) {
		time := chrono.LocalTimeOf(22, 0, 0, 0)
		if formatted := time.Format("%I %P"); formatted != "10 pm" {
			t.Errorf("got %q, want '10 pm'", formatted)
		}
	})

	t.Run("PM", func(t *testing.T) {
		time := chrono.LocalTimeOf(22, 0, 0, 0)
		if formatted := time.Format("%I %p"); formatted != "10 PM" {
			t.Errorf("got %q, want '10 PM'", formatted)
		}
	})

	t.Run("noon", func(t *testing.T) {
		time := chrono.LocalTimeOf(12, 0, 0, 0)
		if formatted := time.Format("%I %P"); formatted != "12 pm" {
			t.Errorf("got %q, want '12 pm'", formatted)
		}
	})

	t.Run("midnight", func(t *testing.T) {
		time := chrono.LocalTimeOf(0, 0, 0, 0)
		if formatted := time.Format("%I %P"); formatted != "12 am" {
			t.Errorf("got %q, want '12 am'", formatted)
		}
	})
}

func TestLocalTime_Parse_12HourClock(t *testing.T) {
	t.Run("am", func(t *testing.T) {
		var time chrono.LocalTime
		if err := time.Parse("%I %P", "10 am"); err != nil {
			t.Errorf("failed to parse time: %v", err)
		}

		if hour, _, _ := time.Clock(); hour != 10 {
			t.Errorf("got %d, want 10", hour)
		}
	})

	t.Run("AM", func(t *testing.T) {
		var time chrono.LocalTime
		if err := time.Parse("%I %p", "10 AM"); err != nil {
			t.Errorf("failed to parse time: %v", err)
		}

		if hour, _, _ := time.Clock(); hour != 10 {
			t.Errorf("got %d, want 10", hour)
		}
	})

	t.Run("pm", func(t *testing.T) {
		var time chrono.LocalTime
		if err := time.Parse("%I %P", "10 pm"); err != nil {
			t.Errorf("failed to parse time: %v", err)
		}

		if hour, _, _ := time.Clock(); hour != 22 {
			t.Errorf("got %d, want 22", hour)
		}
	})

	t.Run("PM", func(t *testing.T) {
		var time chrono.LocalTime
		if err := time.Parse("%I %p", "10 PM"); err != nil {
			t.Errorf("failed to parse time: %v", err)
		}

		if hour, _, _ := time.Clock(); hour != 22 {
			t.Errorf("got %d, want 22", hour)
		}
	})

	t.Run("noon", func(t *testing.T) {
		var time chrono.LocalTime
		if err := time.Parse("%I %P", "12 pm"); err != nil {
			t.Errorf("failed to parse time: %v", err)
		}

		if hour, _, _ := time.Clock(); hour != 12 {
			t.Errorf("got %d, want 12", hour)
		}
	})

	t.Run("midnight", func(t *testing.T) {
		var time chrono.LocalTime
		if err := time.Parse("%I %P", "12 am"); err != nil {
			t.Errorf("failed to parse time: %v", err)
		}

		if hour, _, _ := time.Clock(); hour != 0 {
			t.Errorf("got %d, want 0", hour)
		}
	})

	t.Run("invalid hour", func(t *testing.T) {
		var time chrono.LocalTime
		if err := time.Parse("%I %P", "14 am"); err == nil {
			t.Errorf("expecting error but got nil")
		}
	})
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

func TestLocalDate_Parse_eras(t *testing.T) {
	t.Run("CE", func(t *testing.T) {
		var date chrono.LocalDate
		if err := date.Parse("%EY %EC", "2022 CE"); err != nil {
			t.Errorf("failed to parse date: %v", err)
		}

		if year, _, _ := date.Date(); year != 2022 {
			t.Errorf("got %d, want 2022", year)
		}
	})

	t.Run("BCE", func(t *testing.T) {
		var date chrono.LocalDate
		if err := date.Parse("%EY %EC", "2022 BCE"); err != nil {
			t.Errorf("failed to parse date: %v", err)
		}

		if year, _, _ := date.Date(); year != -2021 {
			t.Errorf("got %d, want -2021", year)
		}
	})

	t.Run("zero", func(t *testing.T) {
		var date chrono.LocalDate
		if err := date.Parse("%EY %EC", "1 BCE"); err != nil {
			t.Errorf("failed to parse date: %v", err)
		}

		if year, _, _ := date.Date(); year != 0 {
			t.Errorf("got %d, want 0", year)
		}
	})
}

func TestLocalDate_Parse_century(t *testing.T) {
	t.Run("1900s", func(t *testing.T) {
		var date chrono.LocalDate
		if err := date.Parse("%y", "80"); err != nil {
			t.Errorf("failed to parse date: %v", err)
		}

		if year, _, _ := date.Date(); year != 1980 {
			t.Errorf("got %d, want 1980", year)
		}
	})

	t.Run("2000s", func(t *testing.T) {
		var date chrono.LocalDate
		if err := date.Parse("%y", "10"); err != nil {
			t.Errorf("failed to parse date: %v", err)
		}

		if year, _, _ := date.Date(); year != 2010 {
			t.Errorf("got %d, want 2010", year)
		}
	})
}

func TestLocalDate_Parse_invalidDayOfYear(t *testing.T) {
	var date chrono.LocalDate
	if err := date.Parse("%Y-%m-%d (day %j)", "2020-01-20 (day 21)"); err == nil {
		t.Errorf("expecting error but got nil")
	}
}

func TestLocalDate_Parse_invalidISOWeekYear(t *testing.T) {
	var date chrono.LocalDate
	if err := date.Parse("%Y-%m-%d (week %V)", "2020-01-20 (week 2)"); err == nil {
		t.Errorf("expecting error but got nil")
	}
}

func TestLocalDate_Parse_invalidDayOfWeek(t *testing.T) {
	var date chrono.LocalDate
	if err := date.Parse("%Y-%m-%d (weekday %A)", "2020-01-20 (weekday Thursday)"); err == nil {
		t.Errorf("expecting error but got nil")
	}
}

func TestLocalDate_Format_invalid_specifier(t *testing.T) {
	for _, specifier := range []string{
		"%C",
		"%Z",
	} {
		t.Run(specifier, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("expecting panic that didn't occur")
				}
			}()

			var date chrono.LocalDate
			_ = date.Format(specifier)
		})
	}
}

func TestLocalDate_Parse_invalid_specifier(t *testing.T) {
	for _, specifier := range []string{
		"%C",
		"%Z",
	} {
		t.Run(specifier, func(t *testing.T) {
			var date chrono.LocalDate
			if err := date.Parse(specifier, ""); err == nil {
				t.Errorf("expecting error but got nil")
			}
		})
	}
}
