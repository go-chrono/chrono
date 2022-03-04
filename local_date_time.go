package chrono

import (
	"fmt"
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
	date, err := makeLocalDate(year, month, day)
	if err != nil {
		panic(err.Error())
	}

	time, err := makeLocalTime(hour, min, sec, nsec)
	if err != nil {
		panic(err.Error())
	}

	nanos := big.NewInt(date)
	nanos.Mul(nanos, dayExtent)
	nanos.Add(nanos, big.NewInt(time))

	return LocalDateTime{v: *nanos}
}

// OfLocalDateAndTime combines the supplied LocalDate and LocalTime into a single LocalDateTime.
func OfLocalDateAndTime(date LocalDate, time LocalTime) LocalDateTime {
	nanos := big.NewInt(int64(date))
	nanos.Mul(nanos, dayExtent)
	nanos.Add(nanos, big.NewInt(int64(time.v)))

	return LocalDateTime{v: *nanos}
}

// Compare compares d with d2. If d is before d2, it returns -1;
// if d is after d2, it returns 1; if they're the same, it returns 0.
func (d LocalDateTime) Compare(d2 LocalDateTime) int {
	return d.v.Cmp(&d2.v)
}

// Split returns separate LocalDate and LocalTime that together represent d.
func (d LocalDateTime) Split() (LocalDate, LocalTime) {
	date, time := d.split()
	return LocalDate(date), LocalTime{v: Extent(time)}
}

func (d LocalDateTime) String() string {
	date, time := d.split()
	year, month, day := fromLocalDate(date)
	hour, min, sec, nsec := fromLocalTime(time)

	out := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, min, sec)
	if nsec != 0 {
		out += fmt.Sprintf(".%09d", nsec)
	}
	return out
}

func (d LocalDateTime) split() (date, time int64) {
	v := new(big.Int)
	v.Set(&d.v)

	var _time big.Int
	_date, _ := v.DivMod(v, dayExtent, &_time)
	return _date.Int64(), _time.Int64()
}

var (
	dayExtent = big.NewInt(24 * int64(Hour))
)
