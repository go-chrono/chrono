package chrono

import (
	"fmt"
	"strconv"
)

// These are predefined layouts used for the parsing and formatting of dates, times and date-times.
// Additional layouts can be composed using the specifiers detailed below:
//
//   %Y: The ISO 8601 year as a decimal number, padded to 4 digits with leading 0s.
//  %EY: The year in the era as a decimal number, padded to 4 digits with leading 0s.
//   %y: The ISO 8601 year without a century as a decimal number, padded to 2 digits with a leading 0, in the range 00 to 99.
//  %Ey: The year in the era without a century as a decimal number, padded to 2 digits with a leading 0, in the range 00 to 99.
//  %EC: The name of the era, either "CE" (for Common Era) "BCE" (for Before the Common Era).
//   %j: The day of the year as a decimal number, padded to 3 digits with leading 0s, in the range 001 to 366.
//   %m: The month as a decimal number, padded to 2 digits with a leading 0, in the range 01 to 12.
//   %B: The full month name, e.g. January, February, etc.
//   %b: The abbreviated month name, e.g. Jan, Feb, etc.
//   %d: The day of the month as a decimal number, padded to 2 digits with a leading 0, in the range 01 to 31.
//
//   %u: The day of the week as a decimal number, e.g. 1 for Monday, 2 for Tuesday, etc.
//   %A: The full name of the day of the week, e.g. Monday, Tuesday, etc.
//   %a: The abbreviated name of the day of the week, e.g. Mon, Tue, etc.
//
//   %G: The ISO 8601 week-based year, which may differ by ±1 to the actual calendar year.
//   %V: The ISO week number, padded to 2 digits with a leading 0, in the range 01 to 53.
//
//   %P: Either "am" or "pm", where noon is "pm" and midnight is "am".
//   %p: Either "AM" or "PM", where noon is "PM" and midnight is "AM".
//   %I: The hour of the day using the 12-hour clock as a decimal number, padded to 2 digits with a leading 0, in the range 01 to 12.
//
//   %H: The hour of the day using the 24-hour clock as a decimal number, padded to 2 digits with a leading 0, in the range 00 to 23.
//   %M: The minute as a decimal number, padded to 2 digits with a leading 0, in the range 00 to 59.
//   %S: The second as a decimal number, padded to 2 digits with a leading 0, in the range 00 to 59.
//
// Depending on the context in which the layout is used, only a subset may be supported.
// For example, %H is not supported when parsing or formatting a date.
//
// For familiarity, the examples below use the time package's reference time of "2nd Jan 2006 15:04:05 -0700" (Unix time 1136239445).
// But note that this reference format is not relevant at all to the functioning of this package.
//
const (
	ISO8601Date             = "%Y%m%d"                                  // 20060102
	ISO8601DateExtended     = "%Y-%m-%d"                                // 2006-01-02
	ISO8601Time             = "T%H%M%S"                                 // T030405
	ISO8601TimeExtended     = "T%H:%M:%S"                               // T03:04:05
	ISO8601DateTime         = ISO8601Date + ISO8601Time                 // 20060102T030405
	ISO8601DateTimeExtended = ISO8601DateExtended + ISO8601TimeExtended // 2006-01-02T03:04:05
	Kitchen                 = "%I:%M%p"                                 // 3:04PM
)

func format(layout string, date *LocalDate, time *LocalTime) string {
	var out []rune

	var (
		year  int
		month Month
		day   int
		hour  int
		min   int
		sec   int
	)

	if date != nil {
		year, month, day = date.Date()
	}

	if time != nil {
		hour, min, sec = time.Clock()
	}

	var lit []rune
NextChar:
	for _, c := range layout {
		lit = append(lit, c)

		if len(lit) >= 2 && lit[0] == '%' {
			if c == 'E' {
				continue NextChar
			}

			switch {
			case date != nil && lit[1] == 'a':
				out = append(out, []rune(date.Weekday().short())...)
			case date != nil && lit[1] == 'A':
				out = append(out, []rune(date.Weekday().String())...)
			case date != nil && lit[1] == 'b':
				out = append(out, []rune(month.short())...)
			case date != nil && lit[1] == 'B':
				out = append(out, []rune(month.String())...)
			case date != nil && lit[1] == 'E' && lit[2] == 'C':
				if year > 1 {
					out = append(out, []rune("CE")...)
				} else {
					out = append(out, []rune("BCE")...)
				}
			case date != nil && lit[1] == 'd':
				out = append(out, []rune(fmt.Sprintf("%02d", day))...)
			case date != nil && lit[1] == 'G':
				y, _ := date.ISOWeek()
				out = append(out, []rune(fmt.Sprintf("%04d", y))...)
			case time != nil && lit[1] == 'H':
				out = append(out, []rune(fmt.Sprintf("%02d", hour))...)
			case time != nil && lit[1] == 'I':
				out = append(out, []rune(fmt.Sprintf("%02d", hour%12))...)
			case date != nil && lit[1] == 'j':
				d := date.YearDay()
				out = append(out, []rune(fmt.Sprintf("%03d", d))...)
			case date != nil && lit[1] == 'm':
				out = append(out, []rune(fmt.Sprintf("%02d", month))...)
			case time != nil && lit[1] == 'M':
				out = append(out, []rune(fmt.Sprintf("%02d", min))...)
			case time != nil && lit[1] == 'p':
				if hour < 12 {
					out = append(out, []rune("AM")...)
				} else {
					out = append(out, []rune("PM")...)
				}
			case time != nil && lit[1] == 'P':
				if hour < 12 {
					out = append(out, []rune("am")...)
				} else {
					out = append(out, []rune("pm")...)
				}
			case time != nil && lit[1] == 'S':
				out = append(out, []rune(fmt.Sprintf("%02d", sec))...)
			case date != nil && lit[1] == 'u':
				out = append(out, []rune(strconv.Itoa(int(date.Weekday())+1))...)
			case date != nil && lit[1] == 'V':
				_, week := date.ISOWeek()
				out = append(out, []rune(fmt.Sprintf("%02d", week))...)
			case date != nil && lit[1] == 'y':
				out = append(out, []rune(fmt.Sprintf("%02d", year%100))...)
			case date != nil && lit[1] == 'E' && lit[2] == 'y':
				y := year
				if y < 0 {
					y = y*-1 + 1
				}
				out = append(out, []rune(fmt.Sprintf("%02d", y%100))...)
			case date != nil && lit[1] == 'Y':
				out = append(out, []rune(fmt.Sprintf("%04d", year))...)
			case date != nil && lit[1] == 'E' && lit[2] == 'Y':
				y := year
				if y < 0 {
					y = y*-1 + 1
				}
				out = append(out, []rune(fmt.Sprintf("%04d", y))...)
			case lit[1] == '%':
				out = append(out, '%')
			default:
				panic("unsupported sequence " + string(lit))
			}

			lit = nil
		} else if len(lit) == 1 && lit[0] != '%' {
			out = append(out, lit...)
			lit = nil
		}
	}

	return string(out)
}

func (w Weekday) short() string {
	if w > Sunday {
		return fmt.Sprintf("%%!Weekday(%d)", w)
	}
	return shortDayNames[w]
}

var shortDayNames = [7]string{
	Monday:    "Mon",
	Tuesday:   "Tue",
	Wednesday: "Wed",
	Thursday:  "Thu",
	Friday:    "Fri",
	Saturday:  "Sat",
	Sunday:    "Sun",
}

func (m Month) short() string {
	if m < January || m > December {
		return fmt.Sprintf("%%!Month(%d)", m)
	}
	return shortMonthNames[m-1]
}

var shortMonthNames = [12]string{
	January - 1:   "Jan",
	February - 1:  "Feb",
	March - 1:     "Mar",
	April - 1:     "Apr",
	May - 1:       "May",
	June - 1:      "Jun",
	July - 1:      "Jul",
	August - 1:    "Aug",
	September - 1: "Sep",
	October - 1:   "Oct",
	November - 1:  "Nov",
	December - 1:  "Dec",
}