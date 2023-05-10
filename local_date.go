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
// In order to format a string according to the Gregorian calendar, use Format("%EY %EC").
type LocalDate int32

// LocalDateOf returns the LocalDate that represents the specified year, month and day.
// This function panics if the provided date would overflow the internal type,
// or if it is earlier than the first date that can be represented by this type - 24th November -4713 (4714 BCE).
func LocalDateOf(year int, month Month, day int) LocalDate {
	if !isDateValid(year, int(month), day) {
		panic("invalid date")
	}

	out, err := makeDate(year, int(month), day)
	if err != nil {
		panic(err.Error())
	}
	return LocalDate(out)
}

// OfDayOfYear returns the LocalDate that represents the specified day of the year.
// This function panics if the provided date would overflow the internal type,
// or if it is earlier than the first date that can be represented by this type - 24th November -4713 (4714 BCE).
func OfDayOfYear(year, day int) LocalDate {
	d, err := ofDayOfYear(year, day)
	if err != nil {
		panic(err.Error())
	}
	return LocalDate(d)
}

func ofDayOfYear(year, day int) (int64, error) {
	isLeap := isLeapYear(year)
	if (!isLeap && day > 365) || day > 366 {
		return 0, fmt.Errorf("invalid date")
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

	return makeDate(year, int(month), day)
}

// OfFirstWeekday returns the LocalDate that represents the first of the specified weekday of the supplied month and year.
// This function panics if the provided date would overflow the internal type,
// or if it earlier than the first date that can be represented by this type - 24th November -4713 (4714 BCE).
//
// By providing January as the month, the result is therefore also the first specified weekday of the year.
// And by adding increments of 7 to the result, it is therefore possible to find the nth instance of a particular weekday.
func OfFirstWeekday(year int, month Month, weekday Weekday) LocalDate {
	v := makeJDN(int64(year), int64(month), 1)
	wd := (v + unixEpochJDN) % 7

	if diff := int64(weekday) - wd - 1; diff < 0 {
		v += 7 - (diff * -1)
	} else if diff > 0 {
		v += diff
	}

	if v < minJDN || v > maxJDN {
		panic("invalid date")
	}
	return LocalDate(v)
}

// OfISOWeek returns the LocalDate that represents the supplied ISO 8601 year, week number, and weekday.
// See LocalDate.ISOWeek for further explanation of ISO week numbers.
func OfISOWeek(year, week int, day Weekday) (LocalDate, error) {
	out, err := ofISOWeek(year, week, int(day))
	return LocalDate(out), err
}

func ofISOWeek(year, week, day int) (int64, error) {
	if week < 1 || week > 53 {
		return 0, fmt.Errorf("invalid week number")
	}

	jan4th, err := makeDate(year, int(January), 4)
	if err != nil {
		return 0, err
	}

	v := week*7 + int(day) - int(getWeekday(int32(jan4th))+3)

	daysThisYear := getDaysInYear(year)
	switch {
	case v <= 0: // Date is in previous year.
		return ofDayOfYear(year-1, v+getDaysInYear(year-1))
	case v > daysThisYear: // Date is in next year.
		return ofDayOfYear(year+1, v-daysThisYear)
	default: // Date is in this year.
		return ofDayOfYear(year, v)
	}
}

func getDaysInYear(year int) int {
	if isLeapYear(year) {
		return 366
	}
	return 365
}

func makeDate(year, month, day int) (int64, error) {
	if !isDateInBounds(year, month, day) {
		return 0, fmt.Errorf("date out of bounds")
	}
	return makeJDN(int64(year), int64(month), int64(day)), nil
}

func makeJDN(y, m, d int64) int64 {
	return (1461*(y+4800+(m-14)/12))/4 + (367*(m-2-12*((m-14)/12)))/12 - (3*((y+4900+(m-14)/12)/100))/4 + d - 32075 - unixEpochJDN
}

// Date returns the ISO 8601 year, month and day represented by d.
func (d LocalDate) Date() (year int, month Month, day int) {
	year, _month, day, err := fromLocalDate(int64(d))
	if err != nil {
		panic(err.Error())
	}
	return year, Month(_month), day
}

func fromLocalDate(v int64) (year, month, day int, err error) {
	if v < minJDN || v > maxJDN {
		return 0, 0, 0, fmt.Errorf("invalid date")
	}

	dd := int64(v + unixEpochJDN)

	f := dd + 1401 + ((((4*dd + 274277) / 146097) * 3) / 4) - 38
	e := 4*f + 3
	g := (e % 1461) / 4
	h := 5*g + 2

	day = int((h%153)/5) + 1
	month = int((h/153+2)%12) + 1
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
	return Weekday(getWeekday(int32(d)))
}

func getWeekday(ordinal int32) int {
	return int((ordinal+int32(unixEpochJDN))%7) + 1
}

// YearDay returns the day of the year specified by d, in the range [1,365] for non-leap years, and [1,366] in leap years.
func (d LocalDate) YearDay() int {
	out, err := getYearDay(int64(d))
	if err != nil {
		panic(err.Error())
	}
	return out
}

func getYearDay(v int64) (int, error) {
	year, month, day, err := fromLocalDate(v)
	if err != nil {
		return 0, err
	}
	return getOrdinalDate(year, int(month), day), nil
}

// ISOWeek returns the ISO 8601 year and week number in which d occurs.
// Week ranges from 1 to 53 (even for years that are not themselves leap years).
// Jan 01 to Jan 03 of year n might belong to week 52 or 53 of year n-1, and Dec 29 to Dec 31 might belong to week 1 of year n+1.
func (d LocalDate) ISOWeek() (isoYear, isoWeek int) {
	var err error
	if isoYear, isoWeek, err = getISOWeek(int64(d)); err != nil {
		panic(err.Error())
	}
	return
}

func getISOWeek(v int64) (isoYear, isoWeek int, err error) {
	year, month, day, err := fromLocalDate(v)
	if err != nil {
		return 0, 0, err
	}

	isoYear = year
	isoWeek = int((10 + getOrdinalDate(isoYear, int(month), day) - getWeekday(int32(v))) / 7)
	if isoWeek == 0 {
		if isLeapYear(isoYear - 1) {
			return isoYear - 1, 53, nil
		}
		return isoYear - 1, 52, nil
	}

	if isoWeek == 53 && !isLeapYear(year) {
		return isoYear + 1, 1, nil
	}

	return isoYear, isoWeek, nil
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

	out, err := makeDate(year+years, int(month)+months, day+days)
	return LocalDate(out), err
}

func (d LocalDate) String() string {
	year, month, day := d.Date()
	return getDateSimpleStr(year, int(month), day)
}

func getDateSimpleStr(year, month, day int) string {
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

func getISODateSimpleStr(year, week, day int) string {
	return fmt.Sprintf("%04d-W%02d-%d", year, week, day)
}

// Format returns a textual representation of the date value formatted according to the layout defined by the argument.
// See the constants section of the documentation to see how to represent the layout format.
// Time format specifiers encountered in the layout results in a panic.
func (d LocalDate) Format(layout string) string {
	out, err := formatDateTimeOffset(layout, (*int32)(&d), nil, 0)
	if err != nil {
		panic(err.Error())
	}
	return out
}

// Parse a formatted string and store the value it represents in d.
// See the constants section of the documentation to see how to represent the layout format.
// Time format specifiers encountered in the layout results in a panic.
func (d *LocalDate) Parse(layout, value string) error {
	v := int64(*d)
	if err := parseDateAndTime(layout, value, &v, nil); err != nil {
		return err
	}

	*d = LocalDate(v)
	return nil
}

// MinLocalDate returns the earliest supported date.
func MinLocalDate() LocalDate {
	return LocalDate(minJDN)
}

// MaxLocalDate returns the latest supported date.
func MaxLocalDate() LocalDate {
	return LocalDate(maxJDN)
}

func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}

func getOrdinalDate(year, month, day int) int {
	var out int
	for i := int(January); i <= month; i++ {
		if i == month {
			out += int(day)
		} else {
			out += int(daysInMonths[i-1])
		}
	}

	if isLeapYear(year) && month > int(February) {
		out++
	}
	return out
}

func isDateInBounds(year, month, day int) bool {
	if year < minYear {
		return false
	} else if year == minYear {
		if month < minMonth {
			return false
		} else if month == minMonth && day < minDay {
			return false
		}
	}

	if year > maxYear {
		return false
	} else if year == maxYear {
		if month > maxMonth {
			return false
		} else if month == maxMonth && day > maxDay {
			return false
		}
	}

	return true
}

func isDateValid(year, month, date int) bool {
	if month < int(January) || month > int(December) {
		return false
	}

	if isLeapYear(year) && month == int(February) {
		return date > 0 && date <= 29
	}
	return date > 0 && date <= daysInMonths[month-1]
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
	minMonth = int(November)
	minDay   = 24
	minJDN   = -unixEpochJDN

	// The maximum representable date must fit into an int32.
	maxYear  = 5874898
	maxMonth = int(June)
	maxDay   = 3
	maxJDN   = math.MaxInt32 - unixEpochJDN
)
