package chrono

import (
	"fmt"
	"math"
)

// LocalDate is a date without a time zone or time component, according to ISO 8601.
// It represents a year-month-day in the proleptic Gregorian calendar,
// but cannot represent an instant on a timeline without additional time offset information.
//
// The date is encoded as a Julian Day Number (JDN). Since LocalDate is stored as an integer,
// any two LocalDates can be compared to each other to determine their relationship,
// and the difference can be calculated to determine the number of days between them.
// Additionally, standard addition and subtraction operators can be used to shift a LocalDate by a number of days.
//
// To make usage of LocalDate easier, the default value, 0, represents the date of the Unix epoch, 1st January 1970.
// This differs from the Richards interpretation of JDNs, where 0 represents the date 24th November 4714 BCE (24/11/-4713).
//
// According to ISO 8601, 0000 is a valid year, whereas in the Gregorian calendar, year 0 does not exist.
// The user must be aware of this difference when interfacing between LocalDate and the BCE/CE notation.
// Thus, when using LocalDate, year 0 is intepreted to mean 1 BCE, and year -1 is 2 BCE, and so on.
type LocalDate int32

// LocalDateOf returns the LocalDate that represents the specified year, month and day.
// This function panics if the provided date would overflow the internal type,
// or if it earlier than the first date that can be represented by this type - 24th November -4713 (4714 BCE).
func LocalDateOf(year int, month Month, day int) LocalDate {
	if !dateIsValid(year, month, day) {
		panic("invalid date")
	}

	out, err := makeLocalDate(year, month, day)
	if err != nil {
		panic(err.Error())
	}
	return LocalDate(out)
}

// OfDayOfYear returns the LocalDate that represents the specified day of the year.
// This function panics if the provided date would overflow the internal type,
// or if it earlier than the first date that can be represented by this type - 24th November -4713 (4714 BCE).
func OfDayOfYear(year, day int) LocalDate {
	isLeap := isLeapYear(year)
	if (!isLeap && day > 365) || day > 366 {
		panic("invalid date")
	}

	var month Month

	var total int
	for m, n := range daysInMonths {
		if isLeap && m == 1 {
			n = 29
		}

		if total+n >= day {
			day = day - total
			month = Month(m + 1)
			break
		}
		total += n
	}

	out, err := makeLocalDate(year, month, day)
	if err != nil {
		panic(err.Error())
	}
	return LocalDate(out)
}

func makeLocalDate(year int, month Month, day int) (int64, error) {
	if !dateInBounds(year, month, day) {
		return 0, fmt.Errorf("date out of bounds")
	}

	y, m, d := int64(year), int64(month), int64(day)
	return (1461*(y+4800+(m-14)/12))/4 + (367*(m-2-12*((m-14)/12)))/12 - (3*((y+4900+(m-14)/12)/100))/4 + d - 32075 - unixEpochJDN, nil
}

// Date returns the ISO 8601 year, month and day represented by d.
func (d LocalDate) Date() (year int, month Month, day int) {
	var err error
	if year, month, day, err = fromLocalDate(int64(d)); err != nil {
		panic(err.Error())
	}
	return
}

func fromLocalDate(v int64) (year int, month Month, day int, err error) {
	if v < minJDN || v > maxJDN {
		return 0, 0, 0, fmt.Errorf("invalid date")
	}

	dd := int64(v + unixEpochJDN)

	f := dd + 1401 + ((((4*dd + 274277) / 146097) * 3) / 4) - 38
	e := 4*f + 3
	g := (e % 1461) / 4
	h := 5*g + 2

	day = int((h%153)/5) + 1
	month = Month((h/153+2)%12) + 1
	year = int(e/1461 - 4716 + (14-int64(month))/12)
	return
}

// IsLeapYear reports whether d is a leap year (contains 29th February, and thus 266 days instead of 265).
func (d LocalDate) IsLeapYear() bool {
	year, _, _ := d.Date()
	return isLeapYear(year)
}

// Weekday returns the day of the week specified by d.
func (d LocalDate) Weekday() Weekday {
	return Weekday((d + unixEpochJDN) % 7)
}

// YearDay returns the day of the year specified by d, in the range [1,365] for non-leap years, and [1,366] in leap years.
func (d LocalDate) YearDay() int {
	year, month, day := d.Date()
	return ordinalDate(year, month, day)
}

// ISOWeek returns the ISO 8601 year and week number in which d occurs.
// Week ranges from 1 to 53 (even for years that are not themselves leap years).
// Jan 01 to Jan 03 of year n might belong to week 52 or 53 of year n-1, and Dec 29 to Dec 31 might belong to week 1 of year n+1.
func (d LocalDate) ISOWeek() (year int, week int) {
	year, month, day := d.Date()
	week = int((10 + ordinalDate(year, month, day) - int(d.Weekday()) - 1) / 7)

	if week == 0 {
		if isLeapYear(year - 1) {
			return year - 1, 53
		}
		return year - 1, 52
	}

	if week == 53 && !d.IsLeapYear() {
		return year + 1, 1
	}

	return year, week
}

// Add returns the date corresponding to adding the given number of years, months, and days to d.
func (d LocalDate) Add(years, months, days int) LocalDate {
	out, err := d.add(years, months, days)
	if err != nil {
		panic(err.Error())
	}
	return out
}

// CanAdd returns false if Add would panic if passed the same arguments.
func (d LocalDate) CanAdd(years, months, days int) bool {
	_, err := d.add(years, months, days)
	return err == nil
}

func (d LocalDate) add(years, months, days int) (LocalDate, error) {
	year, month, day, err := fromLocalDate(int64(d))
	if err != nil {
		return 0, err
	}

	out, err := makeLocalDate(year+years, month+Month(months), day+days)
	return LocalDate(out), err
}

func (d LocalDate) String() string {
	year, month, day := d.Date()
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

// Format returns a textual representation of the date value formatted according to the layout defined by the argument.
// See the constants section of the documentation to see how to represent the layout format.
func (d LocalDate) Format(layout string) string {
	return format(layout, &d, nil)
}

// MinLocalDate returns the earliest supported date.
func MinLocalDate() LocalDate {
	return LocalDate(minJDN)
}

// MaxLocalDate returns the latest supported date.
func MaxLocalDate() LocalDate {
	return LocalDate(maxJDN)
}

func isLeapYear(y int) bool {
	return (y%4 == 0 && y%100 != 0) || y%400 == 0
}

func ordinalDate(y int, m Month, d int) int {
	var out int
	for i := January; i <= m; i++ {
		if i == m {
			out += int(d)
		} else {
			out += int(daysInMonths[i-1])
		}
	}

	if isLeapYear(y) && m > February {
		out++
	}
	return out
}

func dateInBounds(y int, m Month, d int) bool {
	if y < minYear {
		return false
	} else if y == minYear {
		if m < minMonth {
			return false
		} else if m == minMonth && d < minDay {
			return false
		}
	}

	if y > maxYear {
		return false
	} else if y == maxYear {
		if m > maxMonth {
			return false
		} else if m == maxMonth && d > maxDay {
			return false
		}
	}

	return true
}

func dateIsValid(y int, m Month, d int) bool {
	if m < January || m > December {
		return false
	}

	if isLeapYear(y) && m == February {
		return d > 0 && d <= 29
	}
	return d > 0 && d <= daysInMonths[m-1]
}

var daysInMonths = [12]int{
	January - 1:   31,
	February - 1:  28,
	March - 1:     31,
	April - 1:     30,
	May - 1:       31,
	June - 1:      30,
	July - 1:      31,
	August - 1:    31,
	September - 1: 30,
	October - 1:   31,
	November - 1:  30,
	December - 1:  31,
}

const (
	// unixEpochJDN is the JDN that corresponds to 1st January 1970 (Gregorian).
	unixEpochJDN = 2440588

	// The minimum representable date is JDN 0.
	minYear  = -4713
	minMonth = November
	minDay   = 24
	minJDN   = -unixEpochJDN

	// The maximum representable date must fit into an int32.
	maxYear  = 5874898
	maxMonth = June
	maxDay   = 3
	maxJDN   = math.MaxInt32 - unixEpochJDN
)
