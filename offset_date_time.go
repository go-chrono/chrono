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
func OfLocalDateTimeOffset(date LocalDate, time LocalTime, offset Offset) OffsetDateTime {
	return OffsetDateTime{
		v: makeDateTime(int64(date), time.v),
		o: int64(offset),
	}
}

// Split returns separate a LocalDate and OffsetTime that together represent d.
func (d OffsetDateTime) Split() (LocalDate, OffsetTime) {
	date, time := splitDateAndTime(d.v)
	return LocalDate(date), OffsetTime{v: time, o: d.o}
}

// Local returns the LocalDateTime represented by t.
func (d OffsetDateTime) Local() LocalDateTime {
	return LocalDateTime{d.v}
}
