package chrono

// OffsetTime has the same semantics as LocalTime, but with the addition of a timezone offset.
type OffsetTime struct {
	v, o int64
}

// OffsetTimeOf returns an OffsetTime that represents the specified hour, minute, second, and nanosecond offset within the specified second.
// The supplied offset is applied to the returned OffsetTime in the same manner as OffsetOf.
// A valid time is between 00:00:00 and 99:59:59.999999999. If an invalid time is specified, this function panics.
func OffsetTimeOf(hour, min, sec, nsec, offsetHours, offsetMins int) OffsetTime {
	v, err := makeTime(hour, min, sec, nsec)
	if err != nil {
		panic(err.Error())
	}
	return OffsetTime{
		v: v,
		o: makeOffset(offsetHours, offsetMins),
	}
}

// OfTimeOffset combines a LocalTime and Offset into an OffsetTime.
func OfTimeOffset(time LocalTime, offset Offset) OffsetTime {
	return OffsetTime{
		v: time.v,
		o: int64(offset),
	}
}

// BusinessHour returns the hour specified by t.
// If the hour is greater than 23, that hour is returned without normalization.
func (t OffsetTime) BusinessHour() int {
	return timeBusinessHour(t.v)
}

// Clock returns the hour, minute and second represented by t.
// If hour is greater than 23, the returned value is normalized so as to fit within
// the 24-hour clock as specified by ISO 8601, e.g. 25 is returned as 01.
func (t OffsetTime) Clock() (hour, min, sec int) {
	hour, min, sec, _ = fromTime(int64(t.v))
	return
}

// Nanosecond returns the nanosecond offset within the second specified by t, in the range [0, 999999999].
func (t OffsetTime) Nanosecond() int {
	return timeNanoseconds(t.v)
}

// Sub returns the duration t-u.
func (t OffsetTime) Sub(u OffsetTime) Duration {
	return DurationOf(Extent(t.utc() - u.utc()))
}

// Add returns the time t+v, maintaining the offset of t.
// If the result exceeds the maximum representable time, this function panics.
// If the result would be earlier than 00:00:00, the returned time is moved to the previous day,
// e.g. 01:00:00 - 2 hours = 23:00:00.
func (t OffsetTime) Add(v Extent) OffsetTime {
	out, err := addTime(t.v, int64(v))
	if err != nil {
		panic(err.Error())
	}
	return OffsetTime{
		v: out,
		o: t.o,
	}
}

// CanAdd returns false if Add would panic if passed the same argument.
func (t OffsetTime) CanAdd(v Extent) bool {
	_, err := addTime(t.v, int64(v))
	return err == nil
}

// Compare compares t with t2. If t is before t2, it returns -1;
// if t is after t2, it returns 1; if they're the same, it returns 0.
func (t OffsetTime) Compare(t2 OffsetTime) int {
	return compareTimes(t.utc(), t2.utc())
}

func (t OffsetTime) String() string {
	return timeString(t.v) + offsetString(t.o, ":")
}

// In returns a copy of t, adjusted to the supplied offset.
func (t OffsetTime) In(offset Offset) OffsetTime {
	return OffsetTime{
		v: t.v - int64(t.o) + int64(offset),
		o: int64(offset),
	}
}

// UTC is a shortcut for t.In(UTC).
func (t OffsetTime) UTC() OffsetTime {
	return OffsetTime{v: t.utc()}
}

// Local returns the LocalTime represented by t.
func (t OffsetTime) Local() LocalTime {
	return LocalTime{v: t.v}
}

// Offset returns the offset of t.
func (t OffsetTime) Offset() Offset {
	return Offset(t.o)
}

func (t OffsetTime) utc() int64 {
	return t.v - int64(t.o)
}

// Format returns a textual representation of the time value formatted according to the layout defined by the argument.
// See the constants section of the documentation to see how to represent the layout format.
// Date format specifiers encountered in the layout results in a panic.
func (t OffsetTime) Format(layout string) string {
	out, err := formatDateTimeOffset(layout, nil, &t.v, t.o)
	if err != nil {
		panic(err.Error())
	}
	return out
}

// Parse a formatted string and store the value it represents in t.
// See the constants section of the documentation to see how to represent the layout format.
// Date format specifiers encountered in the layout results in a panic.
func (t *OffsetTime) Parse(layout, value string) error {
	v := t.v
	if err := parseDateAndTime(layout, value, nil, &v); err != nil {
		return err
	}

	t.v = v
	return nil
}
