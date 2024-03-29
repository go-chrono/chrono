package chrono

// LocalTime is a time without a time zone or date component.
// It represents a time within the 24-hour clock system with nanosecond precision, according to ISO 8601.
//
// Additional flexibility is provided whereby times after 23:59:59.999999999 are also considered valid.
// This feature supports various usecases where times such as 25:00 (instead of 01:00) represent
// business hours that extend beyond midnight. LocalTime supports a maximum hour of 99.
type LocalTime struct {
	v int64
}

// LocalTimeOf returns a LocalTime that represents the specified hour, minute, second, and nanosecond offset within the specified second.
// A valid time is between 00:00:00 and 99:59:59.999999999. If an invalid time is specified, this function panics.
func LocalTimeOf(hour, min, sec, nsec int) LocalTime {
	out, err := makeTime(hour, min, sec, nsec)
	if err != nil {
		panic(err.Error())
	}
	return LocalTime{v: out}
}

// BusinessHour returns the hour specified by t.
// If the hour is greater than 23, that hour is returned without normalization.
func (t LocalTime) BusinessHour() int {
	return timeBusinessHour(t.v)
}

// Clock returns the hour, minute and second represented by t.
// If hour is greater than 23, the returned value is normalized so as to fit within
// the 24-hour clock as specified by ISO 8601, e.g. 25 is returned as 01.
func (t LocalTime) Clock() (hour, min, sec int) {
	hour, min, sec, _ = fromTime(t.v)
	return
}

// Nanosecond returns the nanosecond offset within the second specified by t, in the range [0, 999999999].
func (t LocalTime) Nanosecond() int {
	return timeNanoseconds(t.v)
}

// Sub returns the duration t-u.
func (t LocalTime) Sub(u LocalTime) Extent {
	return Extent(t.v - u.v)
}

// Add returns the time t+v. If the result exceeds the maximum representable time, this function panics.
// If the result would be earlier than 00:00:00, the returned time is moved to the previous day,
// e.g. 01:00:00 - 2 hours = 23:00:00.
func (t LocalTime) Add(v Extent) LocalTime {
	out, err := addTime(t.v, int64(v))
	if err != nil {
		panic(err.Error())
	}
	return LocalTime{v: out}
}

// CanAdd returns false if Add would panic if passed the same argument.
func (t LocalTime) CanAdd(v Extent) bool {
	_, err := addTime(t.v, int64(v))
	return err == nil
}

// Compare compares t with t2. If t is before t2, it returns -1;
// if t is after t2, it returns 1; if they're the same, it returns 0.
func (t LocalTime) Compare(t2 LocalTime) int {
	return compareTimes(t.v, t2.v)
}

func (t LocalTime) String() string {
	hour, min, sec, nsec := fromTime(t.v)
	return simpleTimeStr(hour, min, sec, nsec, nil)
}

// In returns the OffsetTime represeting t with the specified offset.
func (t LocalTime) In(offset Offset) OffsetTime {
	return OffsetTime{v: t.v, o: int64(offset)}
}

// UTC returns the OffsetTime represeting t at the UTC offset.
func (t LocalTime) UTC() OffsetTime {
	return OffsetTime{v: t.v}
}

// Format returns a textual representation of the time value formatted according to the layout defined by the argument.
// See the constants section of the documentation to see how to represent the layout format.
// Date format specifiers encountered in the layout results in a panic.
func (t LocalTime) Format(layout string) string {
	out, err := formatDateTimeOffset(layout, nil, &t.v, nil)
	if err != nil {
		panic(err.Error())
	}
	return out
}

// Parse a formatted string and store the value it represents in t.
// See the constants section of the documentation to see how to represent the layout format.
// Date format specifiers encountered in the layout results in a panic.
func (t *LocalTime) Parse(layout, value string) error {
	v := t.v
	if err := parseDateAndTime(layout, value, nil, &v, nil); err != nil {
		return err
	}

	t.v = v
	return nil
}
