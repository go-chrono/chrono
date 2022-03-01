package chrono

import "fmt"

// LocalTime is a time without a time zone or date component.
// It represents a time within the 24-hour clock system with nanosecond precision, according to ISO 8601.
//
// Additional flexibility is provided whereby times after 23:59:59.999999999 are also considered valid.
// This feature supports various usecases where times such as 25:00 (instead of 01:00) represent
// business hours that extend beyond midnight. LocalTime supports a maximum hour of 99.
type LocalTime struct {
	v Extent
}

// LocalTimeOf returns a LocalTime that represents the specified hour, minute, second and nanosecond offset within the specified second.
// A valid time is between 00:00:00 and 99:59:59.999999999. If an invalid time is specified, this function panics.
func LocalTimeOf(hour, min, sec, nsec int) LocalTime {
	if hour < 0 || hour > 99 || min < 0 || min > 59 || sec < 0 || sec > 59 || nsec < 0 || nsec > 999999999 {
		panic("invalid time")
	}
	return LocalTime{v: Extent(hour)*Hour + Extent(min)*Minute + Extent(sec)*Second + Extent(nsec)}
}

// BusinessHour returns the hour specified by t.
// If the hour is greater than 23, that hour is returned without normalization.
func (t LocalTime) BusinessHour() int {
	return int(t.v / Hour)
}

// Hour returns the hour specified by t.
// If hour is greater than 23, the returned value is normalized so as to fit within
// the 24-hour clock as specified by ISO 8601, e.g. 25 is returned as 01.
func (t LocalTime) Hour() int {
	return int((t.v % (24 * Hour)) / Hour)
}

// Minute returns the minute specified by t.
func (t LocalTime) Minute() int {
	return int(t.v % Hour / Minute)
}

// Second returns the second specified by t.
func (t LocalTime) Second() int {
	return int(t.v % Minute / Second)
}

// Nanosecond returns the nanosecond offset within the second specified by t, in the range [0, 999999999].
func (t LocalTime) Nanosecond() int {
	return int(t.v % Second)
}

// Sub returns the duration t-u.
func (t LocalTime) Sub(u LocalTime) Duration {
	return DurationOf(t.v - u.v)
}

// Add returns the time t+v. If the result exceeds the maximum representable time, this function panics.
// If the result would be earlier than 00:00:00, the returned time is moved to the previous day,
// e.g. 01:00:00 - 2 hours = 23:00:00.
func (t LocalTime) Add(v Extent) LocalTime {
	if v > maxLocalTime {
		panic("invalid duration v")
	}

	out := t.v + v
	if out > maxLocalTime {
		panic("invalid time t+v")
	}

	if out < 0 {
		return LocalTime{v: 24*Hour + (out % (24 * Hour))}
	}
	return LocalTime{v: out}
}

// Compare compares t with u. If t is before u, it returns -1;
// if t is after u, it returns 1; if they're the same, it returns 0.
func (t LocalTime) Compare(u LocalTime) int {
	switch {
	case t.v < u.v:
		return -1
	case t.v > u.v:
		return 1
	default:
		return 0
	}
}

func (t LocalTime) String() string {
	out := fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())
	if nsec := t.Nanosecond(); nsec != 0 {
		out += fmt.Sprintf(".%09d", nsec)
	}
	return out
}

const (
	maxLocalTime Extent = 359999999999999
)
