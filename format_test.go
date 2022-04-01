package chrono_test

import (
	"testing"

	"github.com/go-chrono/chrono"
)

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
