package chrono

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

// Date returns the ISO 8601 year, month and day represented by d.
func (d LocalDate) Date() (year int, month Month, day int) {
	year, _month, day, err := fromDate(int64(d))
	if err != nil {
		panic(err.Error())
	}
	return year, Month(_month), day
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

// YearDay returns the day of the year specified by d, in the range [1,365] for non-leap years, and [1,366] in leap years.
func (d LocalDate) YearDay() int {
	out, err := getYearDay(int64(d))
	if err != nil {
		panic(err.Error())
	}
	return out
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

// AddDate returns the date corresponding to adding the given number of years, months, and days to d.
func (d LocalDate) AddDate(years, months, days int) LocalDate {
	out, err := addDateToDate(int64(d), years, months, days)
	if err != nil {
		panic(err.Error())
	}
	return LocalDate(out)
}

// CanAddDate returns false if AddDate would panic if passed the same arguments.
func (d LocalDate) CanAddDate(years, months, days int) bool {
	_, err := addDateToDate(int64(d), years, months, days)
	return err == nil
}

func (d LocalDate) String() string {
	year, month, day := d.Date()
	return simpleDateStr(year, int(month), day)
}

// Format returns a textual representation of the date value formatted according to the layout defined by the argument.
// See the constants section of the documentation to see how to represent the layout format.
// Time format specifiers encountered in the layout results in a panic.
func (d LocalDate) Format(layout string) string {
	out, err := formatDateTimeOffset(layout, (*int32)(&d), nil, nil)
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
	if err := parseDateAndTime(layout, value, &v, nil, nil); err != nil {
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
