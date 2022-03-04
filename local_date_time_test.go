package chrono_test

import (
	"testing"

	"github.com/go-chrono/chrono"
)

func TestLocalDateTime(t *testing.T) {
	for _, tt := range []struct {
		datetime     chrono.LocalDateTime
		expectedDate chrono.LocalDate
		expectedTime chrono.LocalTime
	}{
		{
			datetime:     chrono.LocalDateTime{},
			expectedDate: chrono.LocalDateOf(1970, chrono.January, 1),
			expectedTime: chrono.LocalTimeOf(0, 0, 0, 0),
		},
		{
			datetime:     chrono.LocalDateTimeOf(2020, chrono.March, 18, 12, 30, 0, 100000000),
			expectedDate: chrono.LocalDateOf(2020, chrono.March, 18),
			expectedTime: chrono.LocalTimeOf(12, 30, 0, 100000000),
		},
	} {
		t.Run(tt.datetime.String(), func(t *testing.T) {
			date, time := tt.datetime.Split()
			if date != tt.expectedDate {
				t.Errorf("datetime.Split() date = %s, want %s", date, tt.expectedDate)
			}

			if time.Compare(tt.expectedTime) != 0 {
				t.Errorf("datetime.Split() time = %s, want %s", time, tt.expectedTime)
			}
		})
	}
}

func TestOfLocalDateAndTime(t *testing.T) {
	datetime := chrono.OfLocalDateAndTime(
		chrono.LocalDateOf(2020, chrono.March, 18),
		chrono.LocalTimeOf(12, 30, 0, 100000000),
	)

	date, time := datetime.Split()
	if expected := chrono.LocalDateOf(2020, chrono.March, 18); date != expected {
		t.Errorf("datetime.Split() date = %s, want %s", date, expected)
	}

	if expected := chrono.LocalTimeOf(12, 30, 0, 100000000); time != expected {
		t.Errorf("datetime.Split() time = %s, want %s", time, expected)
	}
}
