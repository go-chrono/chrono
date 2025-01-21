//go:build parse

package chrono_test

import (
	"testing"

	"github.com/go-chrono/chrono"
)

func TestParseToLayout(t *testing.T) {
	t.Run("OffsetDateTime", func(t *testing.T) {
		for _, tt := range []struct {
			name           string
			value          string
			conf           chrono.ParseConfig
			expectedC      chrono.OffsetDateTime
			expectedLayout string
			expectedErr    error
		}{
			{"%Y-%m-%d", "2006-04-09", chrono.ParseConfig{}, chrono.OffsetDateTimeOf(2006, 4, 9, 0, 0, 0, 0, 0, 0), "%Y-%m-%d", nil},
			{"%Y-%m", "2006-04", chrono.ParseConfig{}, chrono.OffsetDateTimeOf(2006, 4, 1, 0, 0, 0, 0, 0, 0), "%Y-%m", nil},
			{"%Y", "2006", chrono.ParseConfig{}, chrono.OffsetDateTimeOf(2006, 1, 1, 0, 0, 0, 0, 0, 0), "%Y", nil},
			{"%Y-%d", "2006-09", chrono.ParseConfig{DayFirst: true}, chrono.OffsetDateTimeOf(2006, 1, 9, 0, 0, 0, 0, 0, 0), "%Y-%d", nil},
			{"%m-%Y", "04-2006", chrono.ParseConfig{}, chrono.OffsetDateTimeOf(2006, 4, 1, 0, 0, 0, 0, 0, 0), "%m-%Y", nil},
			{"%m-%d", "04-09", chrono.ParseConfig{}, chrono.OffsetDateTimeOf(1970, 4, 9, 0, 0, 0, 0, 0, 0), "%m-%d", nil},
			{"%d-%m-%Y", "09-04-2006", chrono.ParseConfig{DayFirst: true}, chrono.OffsetDateTimeOf(2006, 4, 9, 0, 0, 0, 0, 0, 0), "%d-%m-%Y", nil},
			{"%d-%m", "09-04", chrono.ParseConfig{DayFirst: true}, chrono.OffsetDateTimeOf(1970, 4, 9, 0, 0, 0, 0, 0, 0), "%d-%m", nil},
			{"%Y-%d-%m", "2006-09-04", chrono.ParseConfig{DayFirst: true}, chrono.OffsetDateTimeOf(2006, 4, 9, 0, 0, 0, 0, 0, 0), "%Y-%d-%m", nil},
			{"%Y-%d", "2006-09", chrono.ParseConfig{DayFirst: true}, chrono.OffsetDateTimeOf(2006, 1, 9, 0, 0, 0, 0, 0, 0), "%Y-%d", nil},
			{"%Y", "2006", chrono.ParseConfig{}, chrono.OffsetDateTimeOf(2006, 1, 1, 0, 0, 0, 0, 0, 0), "%Y", nil},
			{"%Y-%m", "2006-04", chrono.ParseConfig{}, chrono.OffsetDateTimeOf(2006, 4, 1, 0, 0, 0, 0, 0, 0), "%Y-%m", nil},
			{"%d-%m", "09-04", chrono.ParseConfig{DayFirst: true}, chrono.OffsetDateTimeOf(1970, 4, 9, 0, 0, 0, 0, 0, 0), "%d-%m", nil},
			{"%m-%d-%Y", "04-09-2006", chrono.ParseConfig{}, chrono.OffsetDateTimeOf(2006, 4, 9, 0, 0, 0, 0, 0, 0), "%m-%d-%Y", nil},
			{"%m-%d", "04-09", chrono.ParseConfig{}, chrono.OffsetDateTimeOf(1970, 4, 9, 0, 0, 0, 0, 0, 0), "%m-%d", nil},
			{"%m-%Y", "04-2006", chrono.ParseConfig{}, chrono.OffsetDateTimeOf(2006, 4, 1, 0, 0, 0, 0, 0, 0), "%m-%Y", nil},
		} {
			t.Run(tt.name, func(t *testing.T) {
				var c chrono.OffsetDateTime

				if actual, err := chrono.ParseToLayout(tt.value, tt.conf, &c); err != nil {
					if tt.expectedErr == nil {
						t.Errorf("unexpected error: %v", err)
					} else if err.Error() != tt.expectedErr.Error() {
						t.Errorf("got unexpected error %v, want %v", err, tt.expectedErr)
					}
				} else {
					if tt.expectedErr != nil {
						t.Errorf("expecting error %v but got nil", tt.expectedErr)
					} else if actual != tt.expectedLayout {
						t.Errorf("got %q, want %q", actual, tt.expectedLayout)
					}
				}

				if c.Compare(tt.expectedC) != 0 {
					t.Errorf("got %q, want %q", c, tt.expectedC)
				}
			})
		}
	})
}
