package chrono

import (
	"fmt"
	"strconv"
	"strings"
)

// These are predefined layouts used for the parsing and formatting of dates, times and date-times.
// Additional layouts can be composed using the specifiers detailed below:
//
//   %Y: The ISO 8601 year as a decimal number, padded to 4 digits with leading 0s.
//  %EY: The year in the era as a decimal number, padded to 4 digits with leading 0s.
//   %y: The ISO 8601 year without a century as a decimal number, padded to 2 digits with a leading 0, in the range 00 to 99. See note (1).
//  %Ey: The year in the era without a century as a decimal number, padded to 2 digits with a leading 0, in the range 00 to 99. See note (1).
//  %EC: The name of the era, either "CE" (for Common Era) "BCE" (for Before the Common Era).
//   %j: The day of the year as a decimal number, padded to 3 digits with leading 0s, in the range 001 to 366. See note (2).
//   %m: The month as a decimal number, padded to 2 digits with a leading 0, in the range 01 to 12.
//   %B: The full month name, e.g. January, February, etc.
//   %b: The abbreviated month name, e.g. Jan, Feb, etc.
//   %d: The day of the month as a decimal number, padded to 2 digits with a leading 0, in the range 01 to 31.
//
//   %u: The day of the week as a decimal number, e.g. 1 for Monday, 2 for Tuesday, etc. See note (3).
//   %A: The full name of the day of the week, e.g. Monday, Tuesday, etc. See note (3).
//   %a: The abbreviated name of the day of the week, e.g. Mon, Tue, etc. See note (3).
//
//   %G: The ISO 8601 week-based year, which may differ by ±1 to the actual calendar year. See note (2).
//   %V: The ISO week number, padded to 2 digits with a leading 0, in the range 01 to 53. See note (2).
//
//   %P: Either "am" or "pm", where noon is "pm" and midnight is "am".
//   %p: Either "AM" or "PM", where noon is "PM" and midnight is "AM".
//   %I: The hour of the day using the 12-hour clock as a decimal number, padded to 2 digits with a leading 0, in the range 01 to 12. See note (4).
//
//   %H: The hour of the day using the 24-hour clock as a decimal number, padded to 2 digits with a leading 0, in the range 00 to 23. See note (5).
//   %M: The minute as a decimal number, padded to 2 digits with a leading 0, in the range 00 to 59.
//   %S: The second as a decimal number, padded to 2 digits with a leading 0, in the range 00 to 59.
//
// When formatting using specifiers that represent padded decimals, leading 0s can be omitted using the '-' character after the '%'.
// For example, '%m' may produce the string '04' (for March), but '%-m' produces '4'.
// However, when parsing using these specifiers, it is not required that the input string contains any leading zeros.
//
// When parsing using specifier that represent textual values (month names, etc.), the input text is treated case insensitively.
//
// Depending on the context in which the layout is used, only a subset of specifiers may be supported by a particular function.
// For example, %H is not supported when parsing or formatting a date.
//
// When parsing, if multiple instances of the same specifier, or multiple instances of a specifier that represent the same value,
// are encountered, only the instance will be considered. See note (2).
//
// If a specifier is encountered which is not recognized (defined in the list above), or not supported by a particular function,
// the function will panic with a message that includes the unrecognized sequence.
//
// Any other text is enchoed verbatim when formatting, and is expected to appear verbatim in the parsed text.
// In order to print the '%' character verbatim (which normally signifies a specifier), the sequence '%%' can be used.
//
// For familiarity, the examples below use the time package's reference time of "2nd Jan 2006 15:04:05 -0700" (Unix time 1136239445).
// But note that this reference format is not relevant at all to the functioning of this package.
//
// Notes:
//   (1) When 2-digit years are parsed, they are converted according to the POSIX and ISO C standards:
//       values 69–99 are mapped to 1969–1999, and values 0–68 are mapped to 2000–2068.
//   (2) When a date is parsed in combination with a day of year (%j), and/or an ISO week-based date (%G and/or %V),
//       an error will be returned if the represented dates to not match.
//   (3) When a date is parsed in combination with a day of the week (%a, %A and/or %u),
//       an error will be returned if it does not match the day represented by the parsed date.
//       The day of the week is otherwise ignored - it does not have any effect on the result.
//   (4) When a time represented in the 12-hour clock format (%I) is parsed, and no time of day (%P or %p) is present,
//       the time of day is assumed to be before noon, i.e. am or AM.
//   (5) When a time is parsed that contains the time of day (%P or %p), any hour (%H) that is present must be valid
//       on the 12-hour clock.
const (
	// ISO 8601.
	ISO8601Date             = "%Y%m%d"                                  // 20060102
	ISO8601DateExtended     = "%Y-%m-%d"                                // 2006-01-02
	ISO8601Time             = "T%H%M%S"                                 // T030405
	ISO8601TimeExtended     = "T%H:%M:%S"                               // T03:04:05
	ISO8601DateTime         = ISO8601Date + ISO8601Time                 // 20060102T030405
	ISO8601DateTimeExtended = ISO8601DateExtended + ISO8601TimeExtended // 2006-01-02T03:04:05
	// Layouts defined by the time package.
	ANSIC   = "%a %b %d %H:%M:%S %Y" // Mon Jan 02 15:04:05 2006
	Kitchen = "%I:%M%p"              // 3:04PM
)

func format(layout string, date *LocalDate, time *LocalTime) (string, error) {
	var (
		year  int
		month int
		day   int
		hour  int
		min   int
		sec   int
	)

	var err error
	if date != nil {
		v := int64(*date)
		if year, month, day, err = fromLocalDate(v); err != nil {
			return "", err
		}
	}

	if time != nil {
		v := int64(time.v)
		hour, min, sec, _ = fromLocalTime(v)
	}

	var buf, out []rune
NextChar:
	for _, c := range layout {
		buf = append(buf, c)

		if len(buf) >= 2 && buf[0] == '%' {
			if c == '-' || c == 'E' {
				continue NextChar
			}

			nopad, localed, main, err := parseSpecifier(buf)
			if err != nil {
				return "", err
			}

			decimal := func(v int, len int) string {
				if nopad {
					return strconv.Itoa(v)
				}
				return fmt.Sprintf("%0*d", len, v)
			}

			switch {
			case date != nil && main == 'a':
				out = append(out, []rune(date.Weekday().short())...)
			case date != nil && main == 'A':
				out = append(out, []rune(date.Weekday().String())...)
			case date != nil && main == 'b':
				out = append(out, []rune(shortMonthName(month))...)
			case date != nil && main == 'B':
				out = append(out, []rune(longMonthName(month))...)
			case date != nil && localed && main == 'C':
				_, isBCE := convertISOToGregorianYear(year)
				if isBCE {
					out = append(out, []rune("BCE")...)
				} else {
					out = append(out, []rune("CE")...)
				}
			case date != nil && main == 'd':
				out = append(out, []rune(decimal(day, 2))...)
			case date != nil && main == 'G':
				y, _ := date.ISOWeek()
				out = append(out, []rune(decimal(y, 4))...)
			case time != nil && main == 'H':
				out = append(out, []rune(decimal(hour, 2))...)
			case time != nil && main == 'I':
				if hour <= 12 {
					out = append(out, []rune(decimal(hour, 2))...)
				} else {
					out = append(out, []rune(decimal(hour%12, 2))...)
				}
			case date != nil && main == 'j':
				d := date.YearDay()
				out = append(out, []rune(decimal(d, 3))...)
			case date != nil && main == 'm':
				out = append(out, []rune(decimal(int(month), 2))...)
			case time != nil && main == 'M':
				out = append(out, []rune(decimal(min, 2))...)
			case time != nil && main == 'p':
				if hour < 12 {
					out = append(out, []rune("AM")...)
				} else {
					out = append(out, []rune("PM")...)
				}
			case time != nil && main == 'P':
				if hour < 12 {
					out = append(out, []rune("am")...)
				} else {
					out = append(out, []rune("pm")...)
				}
			case time != nil && main == 'S':
				out = append(out, []rune(decimal(sec, 2))...)
			case date != nil && main == 'u':
				out = append(out, []rune(strconv.Itoa(int(date.Weekday())+1))...)
			case date != nil && main == 'V':
				_, w := date.ISOWeek()
				out = append(out, []rune(decimal(w, 2))...)
			case date != nil && main == 'y':
				out = append(out, []rune(decimal(year%100, 2))...)
			case date != nil && localed && main == 'y':
				y, _ := convertISOToGregorianYear(year)
				out = append(out, []rune(fmt.Sprintf("%02d", y%100))...)
			case date != nil && main == 'Y':
				out = append(out, []rune(decimal(year, 4))...)
			case date != nil && localed && main == 'Y':
				y, _ := convertISOToGregorianYear(year)
				out = append(out, []rune(decimal(y, 4))...)
			case main == '%':
				out = append(out, '%')
			default:
				return "", fmt.Errorf("unsupported sequence %q", string(buf))
			}

			buf = nil
		} else if len(buf) == 1 && buf[0] != '%' {
			out = append(out, buf...)
			buf = nil
		}
	}

	return string(out), nil
}

var overrideCentury *int

func getCentury(year int) int {
	switch {
	case overrideCentury != nil:
		return *overrideCentury
	case year >= 69 && year <= 99:
		return 1900
	default:
		return 2000
	}
}

func parse(layout, value string, date, time *int64) error {
	var (
		haveDate          bool
		haveGregorianYear bool
		isBCE             bool
		year              int
		month             int
		day               int

		dayOfWeek int

		dayOfYear int

		haveISODate bool
		isoYear     int
		isoWeek     int

		have12HourClock bool
		isAfternoon     bool
		hour            int
		min             int
		sec             int
		nsec            int
	)

	var err error
	if date != nil {
		if year, month, day, err = fromLocalDate(*date); err != nil {
			return err
		}

		if isoYear, isoWeek, err = getISOWeek(*date); err != nil {
			return err
		}
	}

	if time != nil {
		hour, min, sec, _ = fromLocalTime(*time)
	}

	var pos int
	var buf []rune
	for i := 0; i <= len(layout); i++ {
		verifyText := func() error {
			if len(buf) == 0 {
				return nil
			}

			verify := func() error {
				if !strings.HasPrefix(value[pos:], string(buf)) {
					return fmt.Errorf("parsing time \"%s\" as \"%s\": cannot parse \"%s\" as \"%s\"", value, layout, value[pos:], string(buf))
				}
				return nil
			}

			if buf[len(buf)-1] == '%' {
				if err := verify(); err != nil {
					return err
				}
				buf = []rune{'%'}
			} else {
				if err := verify(); err != nil {
					return err
				}
			}

			pos += len(buf)
			buf = nil
			return nil
		}

		processSpecifier := func() error {
			integer := func(maxLen int) (int, error) {
				str := value[pos:]
				if len(str) >= 1 && (str[0] == '+' || str[0] == '-') {
					maxLen++
				}

				if l := len(str); l < maxLen {
					maxLen = l
				}
				str = value[pos : pos+maxLen]

				var i int
				for _, char := range str {
					if (char < '0' || char > '9') && char != '.' && char != ',' {
						break
					}
					i++
				}
				pos += i

				return strconv.Atoi(str[:i])
			}

			alphas := func(maxLen int) (lower, original string) {
				str := value[pos:]

				if l := len(str); l < maxLen {
					maxLen = l
				}
				str = value[pos : pos+maxLen]

				_lower := make([]rune, maxLen)
				_original := make([]rune, maxLen)

				var i int
				for _, char := range str {
					if char >= 'a' && char <= 'z' {
						_lower[i] = char
						_original[i] = char
						i++
					} else if char >= 'A' && char <= 'Z' {
						_lower[i] = char + 32
						_original[i] = char
						i++
					} else {
						break
					}
				}
				pos += i

				return string(_lower[:i]), string(_original[:i])
			}

			_, localed, main, err := parseSpecifier(buf)
			if err != nil {
				return err
			}

			switch {
			case date != nil && main == 'a':
				lower, original := alphas(3)

				var ok bool
				if dayOfWeek, ok = shortDayNameLookup[lower]; !ok {
					return fmt.Errorf("unrecognized short day name %q", original)
				}
			case date != nil && main == 'A':
				lower, original := alphas(9)

				var ok bool
				if dayOfWeek, ok = longDayNameLookup[lower]; !ok {
					return fmt.Errorf("unrecognized day name %q", original)
				}
			case date != nil && main == 'b':
				haveDate = true

				lower, original := alphas(3)

				var ok bool
				if month, ok = shortMonthNameLookup[lower]; !ok {
					return fmt.Errorf("unrecognized short month name %q", original)
				}
			case date != nil && main == 'B':
				haveDate = true

				lower, original := alphas(9)

				var ok bool
				if month, ok = longMonthNameLookup[lower]; !ok {
					return fmt.Errorf("unrecognized month name %q", original)
				}
			case date != nil && main == 'C':
				if localed { // 'EC'
					haveGregorianYear = true

					lower, original := alphas(2)

					switch lower {
					case "CE":
					case "BCE":
						isBCE = true
					default:
						return fmt.Errorf("unrecognized era %q", original)
					}
				}
			case date != nil && main == 'd':
				haveDate = true
				if day, err = integer(2); err != nil {
					return err
				}
			case date != nil && main == 'G':
				haveISODate = true

				if isoYear, err = integer(4); err != nil {
					return err
				}
			case time != nil && main == 'H':
				if hour, err = integer(2); err != nil {
					return err
				}
			case time != nil && main == 'I':
				have12HourClock = true

				if hour, err = integer(2); err != nil {
					return err
				}
			case date != nil && main == 'j':
				if dayOfYear, err = integer(3); err != nil {
					return err
				}
			case date != nil && main == 'm':
				haveDate = true
				if month, err = integer(2); err != nil {
					return err
				}
			case time != nil && main == 'M':
				if min, err = integer(2); err != nil {
					return err
				}
			case time != nil && main == 'p':
				lower, original := alphas(2)

				switch strings.ToUpper(lower) {
				case "AM":
				case "PM":
					isAfternoon = true
				default:
					return fmt.Errorf("failed to parse time of day %q", original)
				}
			case time != nil && main == 'P':
				lower, original := alphas(2)

				switch lower {
				case "am":
				case "pm":
					isAfternoon = true
				default:
					return fmt.Errorf("failed to parse time of day %q", original)
				}
			case time != nil && main == 'S':
				if sec, err = integer(2); err != nil {
					return err
				}
			case date != nil && main == 'u':
				if dayOfWeek, err = integer(1); err != nil {
					return err
				}
			case date != nil && main == 'V':
				haveISODate = true

				if isoWeek, err = integer(2); err != nil {
					return err
				}
			case date != nil && main == 'y':
				haveDate = true

				if localed { // 'Ey'
					haveGregorianYear = true
				}

				if year, err = integer(2); err != nil {
					return err
				}
				year += getCentury(year)
			case date != nil && main == 'Y':
				haveDate = true

				if localed { // 'EY'
					haveGregorianYear = true
				}

				if year, err = integer(4); err != nil {
					return err
				}
			case main == '%':
			default:
				return fmt.Errorf("unsupported sequence %q", string(buf))
			}

			buf = nil
			return nil
		}

		// Some short-hands.
		var (
			valid             = i < len(layout)
			isSpecifier       = len(buf) >= 2 && buf[0] == '%'
			specifierComplete = isSpecifier && (buf[len(buf)-1] != '-' && buf[len(buf)-1] != 'E')
			isText            = len(buf) >= 1 && buf[0] != '%'
		)

		if valid {
			c := layout[i]
			if len(buf) == 0 {
				goto AppendToBuffer
			} else if isSpecifier {
				if !specifierComplete {
					goto AppendToBuffer
				}

				if err := processSpecifier(); err != nil {
					return err
				}
				goto AppendToBuffer
			} else if isText && c == '%' {
				if err := verifyText(); err != nil {
					return err
				}
				goto AppendToBuffer
			}

		AppendToBuffer:
			buf = append(buf, rune(c))
		} else if isSpecifier {
			if err := processSpecifier(); err != nil {
				return err
			}
		} else if isText {
			if err := verifyText(); err != nil {
				return err
			}
		}
	}

	if pos < len(value) {
		return fmt.Errorf("parsing time \"%s\": extra text: \"%s\"", value, value[pos:])
	}

	if date != nil {
		if haveGregorianYear {
			if year, err = convertGregorianToISOYear(year, isBCE); err != nil {
				return err
			}
		}

		if !isDateValid(year, month, day) {
			return fmt.Errorf("invalid date %q", getDateSimpleStr(year, month, day))
		}

		_date, err := makeLocalDate(year, month, day)
		if err != nil {
			return err
		}

		*date = _date

		// Check day of year according to note (2).
		if dayOfYear != 0 {
			doyDate, err := ofDayOfYear(year, dayOfYear)
			if err != nil {
				return err
			}

			if haveDate && (doyDate != _date) {
				return fmt.Errorf("day-of-year date %q does not agree with date %q",
					LocalDate(doyDate).String(),
					getDateSimpleStr(year, month, day),
				)
			}

			*date = doyDate
		}

		// Check ISO week-year according to note (2).
		if haveISODate {
			isoDate, err := ofISOWeek(isoYear, isoWeek, day)
			if err != nil {
				return fmt.Errorf("invalid ISO week-year date %q", getISODateSimpleStr(isoYear, isoWeek, day))
			}

			if haveDate && (isoDate != _date) {
				return fmt.Errorf("ISO week-year date %q does not agree with date %q",
					getISODateSimpleStr(isoYear, isoWeek, day),
					getDateSimpleStr(year, month, day),
				)
			}

			*date = isoDate
		}

		// Check day of week according to note (3).
		haveDate = haveDate || dayOfYear != 0 || haveISODate
		if dayOfWeek != 0 && haveDate {
			if actual := getWeekday(int32(*date)); dayOfWeek != actual {
				return fmt.Errorf("day of week %q does not agree with actual day of week %q",
					longWeekdayName(dayOfWeek),
					longWeekdayName(actual),
				)
			}
		}
	}

	if time != nil {
		if have12HourClock {
			if hour < 1 || hour > 12 {
				return fmt.Errorf("hour %d is not valid on the 12-hour clock", hour)
			}
			hour = convert12To24HourClock(hour, isAfternoon)
		}

		v, err := makeLocalTime(hour, min, sec, nsec)
		if err != nil {
			return err
		}
		*time = v
	}

	return nil
}

func parseSpecifier(buf []rune) (nopad, localed bool, main rune, err error) {
	if len(buf) == 3 {
		switch buf[1] {
		case '-':
			nopad = true
		case 'E':
			localed = true
		default:
			return false, false, 0, fmt.Errorf("unsupported modifier '%c'", buf[1])
		}
	} else if len(buf) == 4 {
		switch buf[1] {
		case '-':
			nopad = true
		default:
			return false, false, 0, fmt.Errorf("unsupported modifier '%c'", buf[1])
		}

		switch buf[2] {
		case 'E':
			localed = true
		default:
			return false, false, 0, fmt.Errorf("unsupported modifier '%c'", buf[1])
		}
	}
	return nopad, localed, buf[len(buf)-1], nil
}

func convert12To24HourClock(hour int, isAfternoon bool) int {
	if isAfternoon && hour == 12 {
		return 12
	} else if isAfternoon {
		return hour + 12
	} else if hour == 12 {
		return 0
	}
	return hour
}

func convertGregorianToISOYear(gregorianYear int, isBCE bool) (isoYear int, err error) {
	if gregorianYear == 0 {
		return 0, fmt.Errorf("invalid Gregorian year %04d", gregorianYear)
	}

	if isBCE {
		return (gregorianYear * -1) - 1, nil
	}
	return gregorianYear, nil
}

func convertISOToGregorianYear(isoYear int) (gregorianYear int, isBCE bool) {
	if isoYear <= 0 {
		return (isoYear * -1) + 1, true
	}
	return isoYear, false
}

func (d Weekday) short() string {
	return shortWeekdayName(int(d))
}

func shortWeekdayName(d int) string {
	if d > int(Sunday) {
		return fmt.Sprintf("%%!Weekday(%d)", d)
	}
	return shortDayNames[d]
}

var longDayNameLookup = map[string]int{
	"monday":    int(Monday),
	"tuesday":   int(Tuesday),
	"wednesday": int(Wednesday),
	"thursday":  int(Thursday),
	"friday":    int(Friday),
	"saturday":  int(Saturday),
	"sunday":    int(Sunday),
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

var shortDayNameLookup = map[string]int{
	"mon": int(Monday),
	"tue": int(Tuesday),
	"wed": int(Wednesday),
	"thu": int(Thursday),
	"fri": int(Friday),
	"sat": int(Saturday),
	"sun": int(Sunday),
}

func (m Month) short() string {
	return shortMonthName(int(m))
}

func shortMonthName(m int) string {
	if m < int(January) || m > int(December) {
		return fmt.Sprintf("%%!Month(%d)", m)
	}
	return shortMonthNames[m-1]
}

var longMonthNameLookup = map[string]int{
	"january":   int(January),
	"february":  int(February),
	"march":     int(March),
	"april":     int(April),
	"may":       int(May),
	"june":      int(June),
	"july":      int(July),
	"august":    int(August),
	"september": int(September),
	"october":   int(October),
	"november":  int(November),
	"december":  int(December),
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

var shortMonthNameLookup = map[string]int{
	"jan": int(January),
	"feb": int(February),
	"mar": int(March),
	"apr": int(April),
	"may": int(May),
	"jun": int(June),
	"jul": int(July),
	"aug": int(August),
	"sep": int(September),
	"oct": int(October),
	"nov": int(November),
	"dec": int(December),
}
