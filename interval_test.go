package chrono_test

import (
	"strings"
	"testing"

	"github.com/go-chrono/chrono"
)

func TestInterval_Parse(t *testing.T) {
	for _, tr := range []struct {
		str      string
		expected string
		reps     int
	}{
		{"", "", 0},
		{"R1/", "R1/", 1},
		{"R10/", "R10/", 10},
		{"R/", "R/", -1},
		{"R-1/", "R/", -1},
		{"R-10/", "R/", -1},
	} {
		for _, tt := range []struct {
			str        string
			start      chrono.OffsetDateTime
			startOk    bool
			end        chrono.OffsetDateTime
			endOk      bool
			period     chrono.Period
			duration   chrono.Duration
			durationOk bool
		}{
			{
				str:     "2007-03-01T13:00:00Z/2008-05-11T15:30:00Z",
				start:   chrono.OffsetDateTimeOf(2007, chrono.March, 1, 13, 0, 0, 0, 0, 0),
				startOk: true,
				end:     chrono.OffsetDateTimeOf(2008, chrono.May, 11, 15, 30, 0, 0, 0, 0),
				endOk:   true,
			},
			{
				str:        "2007-03-01T13:00:00Z/P1Y2M10DT2H30M",
				start:      chrono.OffsetDateTimeOf(2007, chrono.March, 1, 13, 0, 0, 0, 0, 0),
				startOk:    true,
				period:     chrono.Period{Years: 1, Months: 2, Days: 10},
				duration:   chrono.DurationOf(2*chrono.Hour + 30*chrono.Minute),
				durationOk: true,
			},
			{
				str:        "P1Y2M10DT2H30M/2008-05-11T15:30:00Z",
				period:     chrono.Period{Years: 1, Months: 2, Days: 10},
				duration:   chrono.DurationOf(2*chrono.Hour + 30*chrono.Minute),
				durationOk: true,
				end:        chrono.OffsetDateTimeOf(2008, chrono.May, 11, 15, 30, 0, 0, 0, 0),
				endOk:      true,
			},
			{
				str:        "P1Y2M10DT2H30M",
				period:     chrono.Period{Years: 1, Months: 2, Days: 10},
				duration:   chrono.DurationOf(2*chrono.Hour + 30*chrono.Minute),
				durationOk: true,
			},
		} {
			str := tr.str + tt.str

			t.Run(str, func(t *testing.T) {
				in := str
				expected := tr.expected + tt.str

				checkParse := func(t *testing.T) chrono.Interval {
					i, err := chrono.ParseInterval(in)
					if err != nil {
						t.Errorf("failed to parse interval: %v", err)
					} else if dt, err := i.Start(); tt.startOk && (err != nil || dt.Compare(tt.start) != 0) {
						t.Errorf("i.Start() = %v, %v, want %v, true", dt, err, tt.start)
					} else if dt, err := i.End(); tt.endOk && (err != nil || dt.Compare(tt.end) != 0) {
						t.Errorf("i.End() = %v, %v, want %v, true", dt, err, tt.end)
					} else if p, d, err := i.Duration(); tt.durationOk && (err != nil || !p.Equal(tt.period) || d.Compare(tt.duration) != 0) {
						t.Errorf("i.Duration() = %v, %v, %v, want %v, %v, true", p, d, err, tt.period, tt.duration)
					} else if r := i.Repetitions(); r != tr.reps {
						t.Errorf("i.Repetions() = %v, want %v", r, tr.reps)
					}
					return i
				}

				t.Run("slash", func(t *testing.T) {
					i := checkParse(t)

					if formatted := i.String(); formatted != expected {
						t.Errorf("i.String() = %v, want %v", formatted, expected)
					}
				})

				t.Run("hyphens", func(t *testing.T) {
					in = strings.ReplaceAll(str, "/", "--")
					i := checkParse(t)

					if formatted := i.String(); formatted != expected {
						t.Errorf("i.String() = %v, want %v", formatted, expected)
					}
				})
			})
		}
	}
}
