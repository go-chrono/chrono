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

// LocalTimeOf returns a LocalTime that represents the specified hour, minute, second, and nanosecond offset within the specified second.
// A valid time is between 00:00:00 and 99:59:59.999999999. If an invalid time is specified, this function panics.
func LocalTimeOf(hour, min, sec, nsec int) LocalTime {
	out, err := makeLocalTime(hour, min, sec, nsec)
	if err != nil {
		panic(err.Error())
	}
	return LocalTime{v: Extent(out)}
}

func makeLocalTime(hour, min, sec, nsec int) (int64, error) {
	if hour < 0 || hour > 99 || min < 0 || min > 59 || sec < 0 || sec > 59 || nsec < 0 || nsec > 999999999 {
		return 0, fmt.Errorf("invalid time")
	}

	h, m, s, n := int64(hour), int64(min), int64(sec), int64(nsec)
	return h*int64(Hour) + m*int64(Minute) + s*int64(Second) + n, nil
}

// BusinessHour returns the hour specified by t.
// If the hour is greater than 23, that hour is returned without normalization.
func (t LocalTime) BusinessHour() int {
	return int(t.v / Hour)
}

// Clock returns the hour, minute and second represented by t.
// If hour is greater than 23, the returned value is normalized so as to fit within
// the 24-hour clock as specified by ISO 8601, e.g. 25 is returned as 01.
func (t LocalTime) Clock() (hour, min, sec int) {
	hour, min, sec, _ = fromLocalTime(int64(t.v))
	return
}

// Nanosecond returns the nanosecond offset within the second specified by t, in the range [0, 999999999].
func (t LocalTime) Nanosecond() int {
	return int(t.v % Second)
}

func fromLocalTime(v int64) (hour, min, sec, nsec int) {
	nsec = int(v) % int(Second)
	sec = int(v) / int(Second)

	hour = (sec / (60 * 60)) % 24
	sec -= hour * (60 * 60)

	min = sec / 60
	sec -= min * 60
	return
}

// Sub returns the duration t-u.
func (t LocalTime) Sub(u LocalTime) Duration {
	return DurationOf(t.v - u.v)
}

// Add returns the time t+v. If the result exceeds the maximum representable time, this function panics.
// If the result would be earlier than 00:00:00, the returned time is moved to the previous day,
// e.g. 01:00:00 - 2 hours = 23:00:00.
func (t LocalTime) Add(v Extent) LocalTime {
	out, err := t.add(v)
	if err != nil {
		panic(err.Error())
	}
	return out
}

// CanAdd returns false if Add would panic if passed the same argument.
func (t LocalTime) CanAdd(v Extent) bool {
	_, err := t.add(v)
	return err == nil
}

func (t LocalTime) add(v Extent) (LocalTime, error) {
	if v > maxLocalTime {
		return LocalTime{}, fmt.Errorf("invalid duration v")
	}

	out := t.v + v
	if out > maxLocalTime {
		return LocalTime{}, fmt.Errorf("invalid time t+v")
	}

	if out < 0 {
		return LocalTime{v: 24*Hour + (out % (24 * Hour))}, nil
	}
	return LocalTime{v: out}, nil
}

// Compare compares t with t2. If t is before t2, it returns -1;
// if t is after t2, it returns 1; if they're the same, it returns 0.
func (t LocalTime) Compare(t2 LocalTime) int {
	switch {
	case t.v < t2.v:
		return -1
	case t.v > t2.v:
		return 1
	default:
		return 0
	}
}

func (t LocalTime) String() string {
	hour, min, sec, nsec := fromLocalTime(int64(t.v))
	out := fmt.Sprintf("%02d:%02d:%02d", hour, min, sec)
	if nsec != 0 {
		out += fmt.Sprintf(".%09d", nsec)
	}
	return out
}

const (
	maxLocalTime Extent = 359999999999999
)
