package chrono

import (
	"math/big"
)

// OffsetDateTime has the same semantics as LocalDateTime, but with the addition of a timezone offset.
type OffsetDateTime struct {
	v big.Int
	o int64
}

// OffsetDateTimeOf returns an OffsetDateTime that represents the specified year, month, day,
// hour, minute, second, and nanosecond offset within the specified second.
// The supplied offset is applied to the returned OffsetDateTime in the same manner as OffsetOf.
// The same range of values as supported by OfLocalDate and OfLocalTime are allowed here.
func OffsetDateTimeOf(year int, month Month, day, hour, min, sec, nsec, offsetHours, offsetMins int) OffsetDateTime {
	date, err := makeDate(year, int(month), day)
	if err != nil {
		panic(err.Error())
	}

	time, err := makeTime(hour, min, sec, nsec)
	if err != nil {
		panic(err.Error())
	}

	return OffsetDateTime{
		v: makeDateTime(date, time),
		o: makeOffset(offsetHours, offsetMins),
	}
}

// OfLocalDateOffsetTime combines a LocalDate and OffsetTime into an OffsetDateTime.
func OfLocalDateOffsetTime(date LocalDate, time OffsetTime) OffsetDateTime {
	return OffsetDateTime{
		v: makeDateTime(int64(date), time.v),
		o: time.o,
	}
}

// OfLocalDateTimeOffset combines a LocalDate, LocalTime, and Offset into an OffsetDateTime.
func OfLocalDateTimeOffset(date LocalDate, time LocalTime, offset Extent) OffsetDateTime {
	return OffsetDateTime{
		v: makeDateTime(int64(date), time.v),
		o: int64(offset),
	}
}

// Compare compares d with d2. If d is before d2, it returns -1;
// if d is after d2, it returns 1; if they're the same, it returns 0.
func (d OffsetDateTime) Compare(d2 OffsetDateTime) int {
	return d.v.Cmp(&d2.v)
}

// Offset returns the offset of d.
func (d OffsetDateTime) Offset() Offset {
	return Offset(d.o)
}

// Split returns separate a LocalDate and OffsetTime that together represent d.
func (d OffsetDateTime) Split() (LocalDate, OffsetTime) {
	date, time := splitDateAndTime(d.v)
	return LocalDate(date), OffsetTime{v: time, o: d.o}
}

// In returns a copy of t, adjusted to the supplied offset.
func (d OffsetDateTime) In(offset Offset) OffsetDateTime {
	return OffsetDateTime{
		v: bigDateToOffset(d.v, d.o, int64(offset)),
		o: int64(offset),
	}
}

// UTC is a shortcut for t.In(UTC).
func (d OffsetDateTime) UTC() OffsetDateTime {
	return OffsetDateTime{v: bigDateToOffset(d.v, d.o, 0)}
}

// Local returns the LocalDateTime represented by d.
func (d OffsetDateTime) Local() LocalDateTime {
	return LocalDateTime{d.v}
}

// Add returns the datetime d+v.
// This function panics if the resulting datetime would fall outside of the allowed range.
func (d OffsetDateTime) Add(v Duration) OffsetDateTime {
	out, err := addDurationToBigDate(d.v, v)
	if err != nil {
		panic(err.Error())
	}
	return OffsetDateTime{v: out, o: d.o}
}

// CanAdd returns false if Add would panic if passed the same arguments.
func (d OffsetDateTime) CanAdd(v Duration) bool {
	_, err := addDurationToBigDate(d.v, v)
	return err == nil
}

// AddDate returns the datetime corresponding to adding the given number of years, months, and days to d.
// This function panic if the resulting datetime would fall outside of the allowed date range.
func (d OffsetDateTime) AddDate(years, months, days int) OffsetDateTime {
	out, err := addDateToBigDate(d.v, years, months, days)
	if err != nil {
		panic(err.Error())
	}
	return OffsetDateTime{v: out, o: d.o}
}

// CanAddDate returns false if AddDate would panic if passed the same arguments.
func (d OffsetDateTime) CanAddDate(years, months, days int) bool {
	_, err := addDateToBigDate(d.v, years, months, days)
	return err == nil
}

// Sub returns the duration d-u.
func (d OffsetDateTime) Sub(u OffsetDateTime) Duration {
	out := new(big.Int).Set(&d.v)
	out.Add(out, big.NewInt(d.o))
	out.Sub(out, &u.v)
	out.Sub(out, big.NewInt(u.o))
	return Duration{v: *out}
}

func (d OffsetDateTime) String() string {
	date, time := splitDateAndTime(d.v)
	hour, min, sec, nsec := fromTime(time)
	year, month, day, err := fromDate(date)
	if err != nil {
		panic(err.Error())
	}
	return simpleDateStr(year, month, day) + " " + simpleTimeStr(hour, min, sec, nsec, &d.o)
}

// Format returns a textual representation of the date-time value formatted according to the layout defined by the argument.
// See the constants section of the documentation to see how to represent the layout format.
func (d OffsetDateTime) Format(layout string) string {
	date, time := d.Split()
	out, err := formatDateTimeOffset(layout, (*int32)(&date), &time.v, &d.o)
	if err != nil {
		panic(err.Error())
	}
	return out
}

// Parse a formatted string and store the value it represents in d.
// See the constants section of the documentation to see how to represent the layout format.
func (d *OffsetDateTime) Parse(layout, value string) error {
	dv, tv := splitDateAndTime(d.v)
	var ov int64
	if err := parseDateAndTime(layout, value, &dv, &tv, &ov); err != nil {
		return err
	}

	d.v = makeDateTime(dv, tv)
	d.o = ov
	return nil
}
