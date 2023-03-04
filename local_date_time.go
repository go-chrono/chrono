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
	date, err := makeLocalDate(year, int(month), day)
	if err != nil {
		panic(err.Error())
	}

	time, err := makeTime(hour, min, sec, nsec)
	if err != nil {
		panic(err.Error())
	}
	return makeLocalDateTime(date, time)
}

// OfLocalDateAndTime combines the supplied LocalDate and LocalTime into a single LocalDateTime.
func OfLocalDateAndTime(date LocalDate, time LocalTime) LocalDateTime {
	return makeLocalDateTime(int64(date), int64(time.v))
}

func makeLocalDateTime(date, time int64) LocalDateTime {
	out := big.NewInt(date)
	out.Mul(out, bigIntDayExtent)
	out.Add(out, big.NewInt(time))
	return LocalDateTime{v: *out}
}

// Compare compares d with d2. If d is before d2, it returns -1;
// if d is after d2, it returns 1; if they're the same, it returns 0.
func (d LocalDateTime) Compare(d2 LocalDateTime) int {
	return d.v.Cmp(&d2.v)
}

// Split returns separate LocalDate and LocalTime that together represent d.
func (d LocalDateTime) Split() (LocalDate, LocalTime) {
	date, time := d.split()
	return LocalDate(date), LocalTime{v: time}
}

// Add returns the datetime d+v.
// This function panics if the resulting datetime would fall outside of the allowed range.
func (d LocalDateTime) Add(v Duration) LocalDateTime {
	out, err := d.add(v)
	if err != nil {
		panic(err.Error())
	}
	return out
}

// CanAdd returns false if Add would panic if passed the same arguments.
func (d LocalDateTime) CanAdd(v Duration) bool {
	_, err := d.add(v)
	return err == nil
}

func (d LocalDateTime) add(v Duration) (LocalDateTime, error) {
	out := new(big.Int).Set(&d.v)
	out.Add(out, &v.v)

	if out.Cmp(&minLocalDateTime.v) == -1 || out.Cmp(&maxLocalDateTime.v) == 1 {
		return LocalDateTime{}, fmt.Errorf("datetime out of range")
	}
	return LocalDateTime{v: *out}, nil
}

// AddDate returns the datetime corresponding to adding the given number of years, months, and days to d.
// This function panic if the resulting datetime would fall outside of the allowed date range.
func (d LocalDateTime) AddDate(years, months, days int) LocalDateTime {
	out, err := d.addDate(years, months, days)
	if err != nil {
		panic(err.Error())
	}
	return out
}

// CanAddDate returns false if AddDate would panic if passed the same arguments.
func (d LocalDateTime) CanAddDate(years, months, days int) bool {
	_, err := d.addDate(years, months, days)
	return err == nil
}

func (d LocalDateTime) addDate(years, months, days int) (LocalDateTime, error) {
	date, _ := d.Split()

	added, err := date.add(years, months, days)
	if err != nil {
		return LocalDateTime{}, err
	}

	if added < minJDN || added > maxJDN {
		return LocalDateTime{}, fmt.Errorf("date out of bounds")
	}

	diff := big.NewInt(int64(added - date))
	diff.Mul(diff, bigIntDayExtent)

	out := new(big.Int).Set(&d.v)
	out.Add(out, diff)

	return LocalDateTime{v: *out}, nil
}

func (d LocalDateTime) String() string {
	date, time := d.split()
	hour, min, sec, nsec := fromTime(time)
	year, month, day, err := fromLocalDate(date)
	if err != nil {
		panic(err.Error())
	}

	out := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, min, sec)
	if nsec != 0 {
		out += fmt.Sprintf(".%09d", nsec)
	}
	return out
}

// Format returns a textual representation of the date-time value formatted according to the layout defined by the argument.
// See the constants section of the documentation to see how to represent the layout format.
func (d LocalDateTime) Format(layout string) string {
	date, time := d.Split()
	out, err := formatDateAndTime(layout, (*int32)(&date), &time.v)
	if err != nil {
		panic(err.Error())
	}
	return out
}

// Parse a formatted string and store the value it represents in d.
// See the constants section of the documentation to see how to represent the layout format.
func (d *LocalDateTime) Parse(layout, value string) error {
	dv, tv := d.split()
	if err := parseDateAndTime(layout, value, &dv, &tv); err != nil {
		return err
	}

	*d = makeLocalDateTime(dv, tv)
	return nil
}

func (d LocalDateTime) split() (date, time int64) {
	v := new(big.Int).Set(&d.v)

	var _time big.Int
	_date, _ := v.DivMod(v, bigIntDayExtent, &_time)
	return _date.Int64(), _time.Int64()
}

// MinLocalDateTime returns the earliest supported datetime.
func MinLocalDateTime() LocalDateTime {
	return minLocalDateTime
}

// MaxLocalDateTime returns the latest supported datetime.
func MaxLocalDateTime() LocalDateTime {
	return maxLocalDateTime
}

var (
	bigIntDayExtent = big.NewInt(24 * int64(Hour))

	minLocalDateTime = OfLocalDateAndTime(MinLocalDate(), LocalTimeOf(0, 0, 0, 0))
	maxLocalDateTime = OfLocalDateAndTime(MaxLocalDate(), LocalTimeOf(99, 59, 59, 999999999))
)
