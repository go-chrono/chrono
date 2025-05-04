package chrono

import (
	"math/big"
)

// LocalDateTime is a date and time without a time zone or time component.
// This is a combination of a LocalDate and LocalTime.
type LocalDateTime struct {
	v big.Int
}

// LocalDateTimeOf returns the LocalDateTime that stores the specified year, month, day,
// hour, minute, second, and nanosecond offset within the specified second.
// The same range of values as supported by OfLocalDate and OfLocalTime are allowed here.
func LocalDateTimeOf(year int, month Month, day, hour, min, sec, nsec int) LocalDateTime {
	date, err := makeDate(year, int(month), day)
	if err != nil {
		panic(err.Error())
	}

	time, err := makeTime(hour, min, sec, nsec)
	if err != nil {
		panic(err.Error())
	}

	return LocalDateTime{v: makeDateTime(date, time)}
}

// OfLocalDateTime combines the supplied LocalDate and LocalTime into a single LocalDateTime.
func OfLocalDateTime(date LocalDate, time LocalTime) LocalDateTime {
	return LocalDateTime{v: makeDateTime(int64(date), time.v)}
}

// Unix returns the LocalDateTime that is represented by the supplied Unix time
// (seconds and/or nanoseconds elapsed since 1st January 1970).
// nsecs may be outside of the range [0, 999999999].
func Unix(secs, nsecs int64) LocalDateTime {
	return LocalDateTime{v: unixToDateTime(secs, nsecs)}
}

// UnixMilli returns the LocalDateTime represented by the supplied Unix time
// (milliseconds elapsed since 1st January 1970).
func UnixMilli(msecs int64) LocalDateTime {
	return LocalDateTime{v: unixMilliToDateTime(msecs)}
}

// UnixMicro returns the LocalDateTime represented by the supplied Unix time
// (microseconds elapsed since 1st January 1970).
func UnixMicro(usecs int64) LocalDateTime {
	return LocalDateTime{v: unixMicroToDateTime(usecs)}
}

// Unix returns the Unix time (seconds elapsed since 1st January 1970).
func (d LocalDateTime) Unix() int64 {
	return dateTimeToUnix(d.v)
}

// UnixMilli returns the Unix time (milliseconds elapsed since 1st January 1970).
func (d LocalDateTime) UnixMilli() int64 {
	return dateTimeToUnixMilli(d.v)
}

// UnixMicro returns the Unix time (microseconds elapsed since 1st January 1970).
func (d LocalDateTime) UnixMicro() int64 {
	return dateTimeToUnixMicro(d.v)
}

// UnixNano returns the Unix time (nanoseconds elapsed since 1st January 1970).
func (d LocalDateTime) UnixNano() int64 {
	return dateTimeToUnixNano(d.v)
}

// Compare compares d with d2. If d is before d2, it returns -1;
// if d is after d2, it returns 1; if they're the same, it returns 0.
func (d LocalDateTime) Compare(d2 LocalDateTime) int {
	return d.v.Cmp(&d2.v)
}

// Split returns separate a LocalDate and LocalTime that together represent d.
func (d LocalDateTime) Split() (LocalDate, LocalTime) {
	date, time := splitDateAndTime(d.v)
	return LocalDate(date), LocalTime{v: time}
}

// In returns the OffsetDateTime represeting d with the specified offset.
func (d LocalDateTime) In(offset Offset) OffsetDateTime {
	return OffsetDateTime{v: d.v, o: int64(offset)}
}

// UTC returns the OffsetDateTime represeting d at the UTC offset.
func (d LocalDateTime) UTC() OffsetDateTime {
	return OffsetDateTime{v: d.v}
}

// Add returns the datetime d+v.
// This function panics if the resulting datetime would fall outside of the allowed range.
func (d LocalDateTime) Add(v Duration) LocalDateTime {
	out, err := addDurationToBigDate(d.v, v)
	if err != nil {
		panic(err.Error())
	}
	return LocalDateTime{v: out}
}

// CanAdd returns false if Add would panic if passed the same arguments.
func (d LocalDateTime) CanAdd(v Duration) bool {
	_, err := addDurationToBigDate(d.v, v)
	return err == nil
}

// AddDate returns the datetime corresponding to adding the given number of years, months, and days to d.
// This function panic if the resulting datetime would fall outside of the allowed date range.
func (d LocalDateTime) AddDate(years, months, days int) LocalDateTime {
	out, err := addDateToBigDate(d.v, years, months, days)
	if err != nil {
		panic(err.Error())
	}
	return LocalDateTime{v: out}
}

// CanAddDate returns false if AddDate would panic if passed the same arguments.
func (d LocalDateTime) CanAddDate(years, months, days int) bool {
	_, err := addDateToBigDate(d.v, years, months, days)
	return err == nil
}

// Sub returns the duration d-u.
func (d LocalDateTime) Sub(u LocalDateTime) Duration {
	out := new(big.Int).Set(&d.v)
	out.Sub(out, &u.v)
	return Duration{v: *out}
}

func (d LocalDateTime) String() string {
	date, time := splitDateAndTime(d.v)
	hour, min, sec, nsec := fromTime(time)
	year, month, day, err := fromDate(date)
	if err != nil {
		panic(err.Error())
	}
	return simpleDateStr(year, month, day) + " " + simpleTimeStr(hour, min, sec, nsec, nil)
}

// Format returns a textual representation of the date-time value formatted according to the layout defined by the argument.
// See the constants section of the documentation to see how to represent the layout format.
func (d LocalDateTime) Format(layout string) string {
	date, time := d.Split()
	out, err := formatDateTimeOffset(layout, (*int32)(&date), &time.v, nil)
	if err != nil {
		panic(err.Error())
	}
	return out
}

// Parse a formatted string and store the value it represents in d.
// See the constants section of the documentation to see how to represent the layout format.
func (d *LocalDateTime) Parse(layout, value string) error {
	dv, tv := splitDateAndTime(d.v)
	if err := parseDateAndTime(layout, value, &dv, &tv, nil); err != nil {
		return err
	}

	d.v = makeDateTime(dv, tv)
	return nil
}

// MinLocalDateTime returns the earliest supported datetime.
func MinLocalDateTime() LocalDateTime {
	return minLocalDateTime
}

// MaxLocalDateTime returns the latest supported datetime.
func MaxLocalDateTime() LocalDateTime {
	return maxLocalDateTime
}
