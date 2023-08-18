package chrono

import (
	"fmt"
	"math"
	"math/big"
)

func getISOWeek(v int64) (isoYear, isoWeek int, err error) {
	year, month, day, err := fromDate(v)
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

func getWeekday(ordinal int32) int {
	return int((ordinal+int32(unixEpochJDN))%7) + 1
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

func fromDate(v int64) (year, month, day int, err error) {
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

func getYearDay(v int64) (int, error) {
	year, month, day, err := fromDate(v)
	if err != nil {
		return 0, err
	}
	return getOrdinalDate(year, int(month), day), nil
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

func addDateToDate(d int64, years, months, days int) (int64, error) {
	year, month, day, err := fromDate(d)
	if err != nil {
		return 0, err
	}

	out, err := makeDate(year+years, int(month)+months, day+days)
	return out, err
}

func simpleDateStr(year, month, day int) string {
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

func getISODateSimpleStr(year, week, day int) string {
	return fmt.Sprintf("%04d-W%02d-%d", year, week, day)
}

func makeDateTime(date, time int64) big.Int {
	out := big.NewInt(date)
	out.Mul(out, bigIntDayExtent)
	out.Add(out, big.NewInt(time))
	return *out
}

func addDurationToBigDate(d big.Int, v Duration) (big.Int, error) {
	out := new(big.Int).Set(&d)
	out.Add(out, &v.v)

	if out.Cmp(&minLocalDateTime.v) == -1 || out.Cmp(&maxLocalDateTime.v) == 1 {
		return big.Int{}, fmt.Errorf("datetime out of range")
	}
	return *out, nil
}

func bigDateToOffset(d big.Int, o1, o2 int64) big.Int {
	out := new(big.Int).Set(&d)
	out.Sub(out, big.NewInt(o1))
	out.Add(out, big.NewInt(o2))
	return *out
}

func addDateToBigDate(d big.Int, years, months, days int) (big.Int, error) {
	date, _ := splitDateAndTime(d)

	added, err := addDateToDate(date, years, months, days)
	if err != nil {
		return big.Int{}, err
	}

	if added < minJDN || added > maxJDN {
		return big.Int{}, fmt.Errorf("date out of bounds")
	}

	diff := big.NewInt(int64(added - date))
	diff.Mul(diff, bigIntDayExtent)

	out := new(big.Int).Set(&d)
	out.Add(out, diff)

	return *out, nil
}

func splitDateAndTime(v big.Int) (date, time int64) {
	vv := new(big.Int).Set(&v)

	var _time big.Int
	_date, _ := vv.DivMod(vv, bigIntDayExtent, &_time)
	return _date.Int64(), _time.Int64()
}

var (
	bigIntDayExtent = big.NewInt(24 * int64(Hour))

	minLocalDateTime = OfLocalDateTime(MinLocalDate(), LocalTimeOf(0, 0, 0, 0))
	maxLocalDateTime = OfLocalDateTime(MaxLocalDate(), LocalTimeOf(99, 59, 59, 999999999))
)
