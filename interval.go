package chrono

import (
	"fmt"
	"strconv"
	"strings"
)

// Interval represents the intervening time between two time points.
type Interval struct {
	s *OffsetDateTime
	e *OffsetDateTime
	d *periodDuration
	r int
}

type periodDuration struct {
	Period
	Duration
}

// IntervalOfStartEnd creates an [Interval] from the provided start and end time points.
func IntervalOfStartEnd(start, end OffsetDateTime, repetitions int) Interval {
	return Interval{s: &start, e: &end, r: repetitions}
}

// IntervalOfStartDuration creates an [Interval] from the provided start time point and duration.
func IntervalOfStartDuration(start OffsetDateTime, period Period, duration Duration, repetitions int) Interval {
	return Interval{s: &start, d: &periodDuration{Period: period, Duration: duration}, r: repetitions}
}

// IntervalOfDurationEnd creates an [Interval] from the provided duration and end time point.
func IntervalOfDurationEnd(period Period, duration Duration, end OffsetDateTime, repetitions int) Interval {
	return Interval{e: &end, d: &periodDuration{Period: period, Duration: duration}, r: repetitions}
}

// ParseInterval parses an ISO 8601 time interval, or a repeating time interval.
// Time intervals can be expressed in the following forms:
//   - <start>/<end>
//   - <start>/<duration>
//   - <duration>/<end>
//   - <duration>
//
// where <start> and <end> is any string that can be parsed by [Parse],
// and <duration> is any string that can be parsed by [ParseDuration].
//
// Repeating time intervals are expressed as such:
//   - Rn/<interval>
//   - R/<interval>
//
// where <interval> is one the forms in the first list, R is the character itself,
// and n is a signed integer value of the number of repetitions. Additionally, the following values have special meanings:
//   - 0 is no repetitions, equivalent to not including the repeat expression at all (i.e. using the forms in the first list as-is);
//   - <=-1 is unbounded number of repetitions, equivalent to not specifying a value at all (i.e. 'R-1' is the same as 'R').
//
// Additionally, '--' can be used as the separator, instead of the default '/' character.
func ParseInterval(s string) (Interval, error) {
	start, end, pd, r, err := parseInterval(s)
	return Interval{
		s: start,
		e: end,
		d: pd,
		r: r,
	}, err
}

// String returns the formatted Interval that can be parsed by i.Parse().
func (i Interval) String() string {
	return i.string("/")
}

func (i Interval) string(sep string) string {
	var out string
	r := i.Repetitions()
	switch r {
	case 0:
		// Omit R.
	case -1:
		out = "R" + sep
	default:
		out = "R" + strconv.Itoa(r) + sep
	}

	switch {
	case i.s != nil && i.e != nil:
		return out + i.s.Format(ISO8601) + sep + i.e.Format(ISO8601)
	case i.s != nil && i.d != nil:
		return out + i.s.Format(ISO8601) + sep + FormatDuration(i.d.Period, i.d.Duration)
	case i.d != nil && i.e != nil:
		return out + FormatDuration(i.d.Period, i.d.Duration) + sep + i.e.Format(ISO8601)
	case i.d != nil:
		return out + FormatDuration(i.d.Period, i.d.Duration)
	default:
		return out
	}
}

// Start returns the start time point if present, or a calculated time point if possible
// by subtracting i.Duration() from i.End().
// If neither are possible (i.e. only a duration is present),
// [ErrUnsupportedRepresentation] is returned instead.
func (i Interval) Start() (OffsetDateTime, error) {
	switch {
	case i.s != nil:
		return *i.s, nil
	case i.e != nil:
		d, err := i.d.mul(-1)
		if err != nil {
			return OffsetDateTime{}, err
		}
		return i.e.AddDate(-int(i.d.Years), -int(i.d.Months), -int(i.d.Days)).Add(d), nil
	default:
		return OffsetDateTime{}, ErrUnsupportedRepresentation
	}
}

// End returns the end time point if present, or a calculated time point if possible
// by adding i.Duration() to i.Start().
// If neither are possible, (i.e. only a duration is present),
// then [ErrUnsupportedRepresentation] is returned instead.
func (i Interval) End() (OffsetDateTime, error) {
	switch {
	case i.e != nil:
		return *i.e, nil
	case i.s != nil:
		return i.s.AddDate(int(i.d.Years), int(i.d.Months), int(i.d.Days)).Add(i.d.Duration), nil
	default:
		return OffsetDateTime{}, ErrUnsupportedRepresentation
	}
}

// Duration returns the [Period] and [Duration] if present, or a calculated [Duration]
// if possible by substracting i.Start() from i.End(). Note that the latter case,
// the [Period] returned will always be the zero value.
func (i Interval) Duration() (Period, Duration, error) {
	switch {
	case i.d != nil:
		return i.d.Period, i.d.Duration, nil
	case i.s != nil && i.e != nil:
		return Period{}, i.e.Sub(*i.s), nil
	default:
		return Period{}, Duration{}, ErrUnsupportedRepresentation
	}
}

// Repetitions returns the number of repetitions of a repeating interval.
// Any negative number, meaning an unbounded number of repitions, is normalized to -1.
func (i Interval) Repetitions() int {
	if i.r <= -1 {
		return -1
	}
	return i.r
}

func cutAB(s, sepA, sepB string) (before, after string, found int) {
	if i := strings.Index(s, sepA); i >= 0 {
		return s[:i], s[i+len(sepA):], 1
	} else if i := strings.Index(s, sepB); i >= 0 {
		return s[:i], s[i+len(sepB):], -1
	}
	return s, "", 0
}

func parseInterval(s string) (start, end *OffsetDateTime, pd *periodDuration, repeat int, err error) {
	if len(s) == 0 {
		return nil, nil, nil, 0, fmt.Errorf("empty string")
	}

	var sep int

	if s[0] == 'R' {
		s = s[1:]

		var r string
		if r, s, sep = cutAB(s, "/", "--"); sep == 0 {
			return nil, nil, nil, 0, fmt.Errorf("parsing interval: missing separator")
		}

		if len(r) == 0 {
			repeat = -1
		} else if repeat, err = strconv.Atoi(r); err != nil {
			return nil, nil, nil, 0, fmt.Errorf("parsing interval: invalid repeat")
		}
	}

	s1, s2, found := cutAB(s, "/", "--")
	if found != 0 && sep != 0 && found != sep {
		return nil, nil, nil, 0, fmt.Errorf("inconsistent separators")
	}

	if s1[0] >= '0' && s1[0] <= '9' { // <start>/<end> or <start>/<duation>
		if s2 == "" { // <start> is invalid
			return nil, nil, nil, 0, fmt.Errorf("invalid interval")
		}

		if start, err = parseOffsetDateTime(s1); err != nil {
			return nil, nil, nil, 0, err
		}
	} else { // <duration>/<end> or <duration>
		p, d, err := ParseDuration(s1)
		if err != nil {
			return nil, nil, nil, 0, err
		}
		pd = &periodDuration{p, d}
	}

	if s2 != "" && s2[0] >= '0' && s2[0] <= '9' { // <start>/<end> or <duration>/<end>
		if end, err = parseOffsetDateTime(s2); err != nil {
			return nil, nil, nil, 0, err
		}
	} else if s2 != "" { // <start>/<duation>
		p, d, err := ParseDuration(s2)
		if err != nil {
			return nil, nil, nil, 0, err
		}
		pd = &periodDuration{p, d}
	}

	return start, end, pd, repeat, nil
}

func parseOffsetDateTime(value string) (*OffsetDateTime, error) {
	var date, time, offset int64
	if err := parseDateAndTime(ISO8601, value, &date, &time, &offset); err != nil {
		return nil, err
	}

	return &OffsetDateTime{
		v: makeDateTime(date, time),
		o: offset,
	}, nil
}
