package chrono

import (
	"fmt"
	"strconv"
	"strings"
)

// These are predefined layouts used for the parsing and formatting of dates, times and date-times.
// Additional layouts can be composed using the specifiers detailed below:
//
//	 %Y: The ISO 8601 year as a decimal number, padded to 4 digits with leading 0s.
//	%EY: The year in the era as a decimal number, padded to 4 digits with leading 0s.
//	 %y: The ISO 8601 year without a century as a decimal number, padded to 2 digits with a leading 0, in the range 00 to 99. See note (1).
//	%Ey: The year in the era without a century as a decimal number, padded to 2 digits with a leading 0, in the range 00 to 99. See note (1).
//	%EC: The name of the era, either "CE" (for Common Era) "BCE" (for Before the Common Era).
//	 %j: The day of the year as a decimal number, padded to 3 digits with leading 0s, in the range 001 to 366. See note (2).
//	 %m: The month as a decimal number, padded to 2 digits with a leading 0, in the range 01 to 12.
//	 %B: The full month name, e.g. January, February, etc.
//	 %b: The abbreviated month name, e.g. Jan, Feb, etc.
//	 %d: The day of the month as a decimal number, padded to 2 digits with a leading 0, in the range 01 to 31.
//
//	 %u: The day of the week as a decimal number, e.g. 1 for Monday, 2 for Tuesday, etc. See note (3).
//	 %A: The full name of the day of the week, e.g. Monday, Tuesday, etc. See note (3).
//	 %a: The abbreviated name of the day of the week, e.g. Mon, Tue, etc. See note (3).
//
//	 %G: The ISO 8601 week-based year, padded to 4 digits with leading 0s. This may differ by ±1 to the actual calendar year. See note (2).
//	 %V: The ISO week number, padded to 2 digits with a leading 0, in the range 01 to 53. See note (2).
//
//	 %P: Either "am" or "pm", where noon is "pm" and midnight is "am".
//	 %p: Either "AM" or "PM", where noon is "PM" and midnight is "AM".
//	 %I: The hour of the day using the 12-hour clock as a decimal number, padded to 2 digits with a leading 0, in the range 01 to 12. See note (4).
//
//	 %H: The hour of the day using the 24-hour clock as a decimal number, padded to 2 digits with a leading 0, in the range 00 to 23. See note (5).
//	 %M: The minute as a decimal number, padded to 2 digits with a leading 0, in the range 00 to 59.
//	 %S: The second as a decimal number, padded to 2 digits with a leading 0, in the range 00 to 59.
//
//	 %f: Equivalent to %6f.
//	%3f: The millisecond offset within the represented second, rounded either up or down and padded to 3 digits with leading 0s.
//	%6f: The microsecond offset within the represented second, rounded either up or down and padded to 6 digits with leading 0s.
//	%9f: The nanosecond offset within the represented second, padded to 9 digits with leading 0s.
//
//	 %z: The UTC offset in the format ±HHMM, preceded always by the sign ('+' or '-'), and padded to 4 digits with leading zeros. See notes (6) and (7).
//	%Ez: Equivalent to %z, except that an offset of +0000 is formatted at 'Z', and other offsets as ±HH:MM. See notes (6) and (7).
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
//
//		(1) When 2-digit years are parsed, they are converted according to the POSIX and ISO C standards:
//		    values 69–99 are mapped to 1969–1999, and values 0–68 are mapped to 2000–2068.
//		(2) When a date is parsed in combination with a day of year (%j), and/or an ISO week-based date (%G and/or %V),
//		    an error will be returned if the represented dates to not match.
//		(3) When a date is parsed in combination with a day of the week (%a, %A and/or %u),
//		    an error will be returned if it does not match the day represented by the parsed date.
//		    The day of the week is otherwise ignored - it does not have any effect on the result.
//		(4) When a time represented in the 12-hour clock format (%I) is parsed, and no time of day (%P or %p) is present,
//		    the time of day is assumed to be before noon, i.e. am or AM.
//		(5) When a time is parsed that contains the time of day (%P or %p), any hour (%H) that is present must be valid
//		    on the 12-hour clock.
//		(6) When UTC offsets are parsed into a type which do not include a time offset element, the offset present in the string is ignored.
//		    When UTC offsets are formatted from a type which does not include a time offset element,
//	        the offset will not be present in the returned string.
//		(7) When UTC offsets are parsed (%z or %Ez), the shorted form of ±HH is accepted.
//		    However, when formatted, only the full forms are returned (either ±HHMM or ±HH:MM).
const (
	// ISO 8601.
	ISO8601                          = ISO8601DateTimeExtended
	ISO8601DateSimple                = "%Y%m%d"                                  // 20060102
	ISO8601DateExtended              = "%Y-%m-%d"                                // 2006-01-02
	ISO8601DateTruncated             = "%Y-%m"                                   // 2006-01
	ISO8601TimeSimple                = "T%H%M%S%z"                               // T030405-0700
	ISO8601TimeExtended              = "T%H:%M:%S%Ez"                            // T03:04:05-07:00
	ISO8601TimeMillisSimple          = "T%H%M%S.%3f%z"                           // T030405.000-0700
	ISO8601TimeMillisExtended        = "T%H:%M:%S.%3f%Ez"                        // T03:04:05.000-07:00
	ISO8601TimeTruncatedMinsSimple   = "T%H%M"                                   // T0304
	ISO8601TimeTruncatedMinsExtended = "T%H:%M"                                  // T03:04
	ISO8601TimeTruncatedHours        = "T%H"                                     // T03
	ISO8601DateTimeSimple            = ISO8601DateSimple + ISO8601TimeSimple     // 20060102T030405-0700
	ISO8601DateTimeExtended          = ISO8601DateExtended + ISO8601TimeExtended // 2006-01-02T03:04:05-07:00
	ISO8601WeekSimple                = "%GW%V"                                   // 2006W01
	ISO8601WeekExtended              = "%G-W%V"                                  // 2006-W01
	ISO8601WeekDaySimple             = "%GW%V%u"                                 // 2006W011
	ISO8601WeekDayExtended           = "%G-W%V-%u"                               // 2006-W01-1
	ISO8601OrdinalDateSimple         = "%Y%j"                                    // 2006002
	ISO8601OrdinalDateExtended       = "%Y-%j"                                   // 2006-002
	// Layouts defined by the time package.
	ANSIC   = "%a %b %d %H:%M:%S %Y" // Mon Jan 02 15:04:05 2006
	Kitchen = "%I:%M%p"              // 3:04PM
)

func formatDateTimeOffset(layout string, date *int32, time *int64, offset *int64) (string, error) {
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
		if year, month, day, err = fromDate(v); err != nil {
			return "", err
		}
	}

	if time != nil {
		v := int64(*time)
		hour, min, sec, _ = fromTime(v)
	}

	var buf, out []rune
NextChar:
	for _, c := range layout {
		buf = append(buf, c)

		if len(buf) >= 2 && buf[0] == '%' {
			if c == '-' || c == 'E' || (c >= '0' && c <= '9') {
				continue NextChar
			}

			nopad, localed, precision, main, err := parseSpecifier(buf)
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
			case date != nil && main == 'a': // %a
				out = append(out, []rune(shortWeekdayName(getWeekday(*date)))...)
			case date != nil && main == 'A': // %A
				out = append(out, []rune(longWeekdayName(getWeekday(*date)))...)
			case date != nil && main == 'b': // %b
				out = append(out, []rune(shortMonthName(month))...)
			case date != nil && main == 'B': // %B
				out = append(out, []rune(longMonthName(month))...)
			case date != nil && main == 'C':
				if localed { // %EC
					if _, isBCE := convertISOToGregorianYear(year); isBCE {
						out = append(out, []rune("BCE")...)
					} else {
						out = append(out, []rune("CE")...)
					}
				} else { // %C
					panic("unsupported specifier 'C'")
				}
			case date != nil && main == 'd': // %d
				out = append(out, []rune(decimal(day, 2))...)
			case time != nil && main == 'f': // %f
				if precision == 0 {
					precision = 6
				}

				nanos := timeNanoseconds(*time)
				switch precision {
				case 3: // %3f
					out = append(out, []rune(decimal(divideAndRoundInt(nanos, 1000000), 3))...)
				case 6: // %6f
					out = append(out, []rune(decimal(divideAndRoundInt(nanos, 1000), 6))...)
				case 9: // %9f
					out = append(out, []rune(decimal(nanos, 9))...)
				default:
					panic(fmt.Sprintf("unsupported specifier '%df'", precision))
				}
			case date != nil && main == 'G': // %G
				v := int64(*date)
				y, _, err := getISOWeek(v)
				if err != nil {
					panic(err.Error())
				}
				out = append(out, []rune(decimal(y, 4))...)
			case time != nil && main == 'H': // %H
				out = append(out, []rune(decimal(hour, 2))...)
			case time != nil && main == 'I': // %I
				h, _ := convert24To12HourClock(hour)
				out = append(out, []rune(decimal(h, 2))...)
			case date != nil && main == 'j': // %j
				v := int64(*date)
				d, err := getYearDay(v)
				if err != nil {
					panic(err.Error())
				}
				out = append(out, []rune(decimal(d, 3))...)
			case date != nil && main == 'm': // %m
				out = append(out, []rune(decimal(int(month), 2))...)
			case time != nil && main == 'M': // %M
				out = append(out, []rune(decimal(min, 2))...)
			case time != nil && main == 'p': // %p
				if _, isAfternoon := convert24To12HourClock(hour); !isAfternoon {
					out = append(out, []rune("AM")...)
				} else {
					out = append(out, []rune("PM")...)
				}
			case time != nil && main == 'P': // %P
				if _, isAfternoon := convert24To12HourClock(hour); !isAfternoon {
					out = append(out, []rune("am")...)
				} else {
					out = append(out, []rune("pm")...)
				}
			case time != nil && main == 'S': // %S
				out = append(out, []rune(decimal(sec, 2))...)
			case date != nil && main == 'u': // %u
				out = append(out, []rune(strconv.Itoa(getWeekday(*date)))...)
			case date != nil && main == 'V': // %V
				v := int64(*date)
				_, w, err := getISOWeek(v)
				if err != nil {
					panic(err.Error())
				}
				out = append(out, []rune(decimal(w, 2))...)
			case date != nil && main == 'y': // %y
				y := year
				if localed { // %Ey
					y, _ = convertISOToGregorianYear(y)
				}
				out = append(out, []rune(decimal(y%100, 2))...)
			case date != nil && main == 'Y': // %Y
				y := year
				if localed { // %EY
					y, _ = convertISOToGregorianYear(y)
				}
				out = append(out, []rune(decimal(y, 4))...)
			case time != nil && main == 'z':
				// Formatting %z from a type that contains no offset (e.g. LocalTime, LocalDateTime)
				// is valid, although it will not be printed.
				if offset == nil {
					break
				}

				if localed { // %Ez
					out = append(out, []rune(offsetString(*offset, ":"))...)
				} else { // %z
					out = append(out, []rune(offsetString(*offset, ""))...)
				}
			case main == '%': // %%
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

const (
	extraTextErrMsg   = "parsing time \"%s\": extra text: \"%s\""
	endOfStringErrMsg = "parsing time \"%s\": end of string"
)

// parseDateAndTime parses the supplied value according to the specified layout.
// date, time and offset must be provided in order for those components to be parsed.
// If not provided, and the specifiers that pertain to those components are
// encountered in the supplied layout, then an error is returned.
// If non-zero, date, time, and offset and taken as starting points, where the individual values
// that they represent are replaced only if present in the supplied layout.
func parseDateAndTime(layout, value string, date, time, offset *int64) error {
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
		if year, month, day, err = fromDate(*date); err != nil {
			return err
		}

		if isoYear, isoWeek, err = getISOWeek(*date); err != nil {
			return err
		}
	}

	if time != nil {
		hour, min, sec, nsec = fromTime(*time)
		_, isAfternoon = convert24To12HourClock(hour)
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
				var neg bool

				str := value[pos:]
				if len(str) >= 1 {
					switch str[0] {
					case '-':
						neg = true
						fallthrough
					case '+':
						str = str[1:]
						pos++
					}
				}

				if l := len(str); l == 0 {
					return 0, fmt.Errorf(endOfStringErrMsg, value)
				} else if l < maxLen {
					maxLen = l
				}
				str = str[:maxLen]

				var i int
				for _, char := range str {
					if (char < '0' || char > '9') && char != '.' && char != ',' {
						break
					}
					i++
				}
				pos += i

				if i == 0 {
					return 0, fmt.Errorf(extraTextErrMsg, value, str)
				}

				out, err := strconv.Atoi(str[:i])
				if err != nil {
					return 0, fmt.Errorf(extraTextErrMsg, value, str)
				}

				if neg {
					return out * -1, nil
				}
				return out, nil
			}

			hasMore := func() bool {
				return len(value[pos:]) > 0
			}

			casedAlpha := func(char rune) (rune, bool) {
				str := value[pos:]
				if len(str) != 0 {
					r := rune(str[0])
					if r == char {
						pos++
						return r, true
					}
				}
				return ' ', false
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

			_, localed, precision, main, err := parseSpecifier(buf)
			if err != nil {
				return err
			}

			switch {
			case date != nil && main == 'a': // %a
				lower, original := alphas(3)
				var ok bool
				if dayOfWeek, ok = shortDayNameLookup[lower]; !ok {
					return fmt.Errorf("unrecognized short day name %q", original)
				}
			case date != nil && main == 'A': // %A
				lower, original := alphas(9)
				var ok bool
				if dayOfWeek, ok = longDayNameLookup[lower]; !ok {
					return fmt.Errorf("unrecognized day name %q", original)
				}
			case date != nil && main == 'b': // %b
				lower, original := alphas(3)
				var ok bool
				if month, ok = shortMonthNameLookup[lower]; !ok {
					return fmt.Errorf("unrecognized short month name %q", original)
				}
			case date != nil && main == 'B': // %B
				lower, original := alphas(9)
				var ok bool
				if month, ok = longMonthNameLookup[lower]; !ok {
					return fmt.Errorf("unrecognized month name %q", original)
				}
			case date != nil && main == 'C':
				if localed { // %EC
					haveGregorianYear = true
					lower, original := alphas(3)
					switch lower {
					case "ce":
					case "bce":
						isBCE = true
					default:
						return fmt.Errorf("unrecognized era %q", original)
					}
				} else { // %C
					return fmt.Errorf("unsupported specifier 'C'")
				}
			case date != nil && main == 'd': // %d
				haveDate = true
				if day, err = integer(2); err != nil {
					return err
				}
			case time != nil && main == 'f': // %f
				if precision == 0 {
					precision = 6
				}

				switch precision {
				case 3: // %3f
					millis, err := integer(3)
					if err != nil {
						return err
					}
					nsec = millis * 1000000
				case 6: // %6f
					micros, err := integer(6)
					if err != nil {
						return err
					}
					nsec = micros * 1000
				case 9: // %9f
					if nsec, err = integer(9); err != nil {
						return err
					}
				default:
				}
			case date != nil && main == 'G': // %G
				haveISODate = true
				if isoYear, err = integer(4); err != nil {
					return err
				}
			case time != nil && main == 'H': // %H
				if hour, err = integer(2); err != nil {
					return err
				}
			case time != nil && main == 'I': // %I
				have12HourClock = true
				if hour, err = integer(2); err != nil {
					return err
				}
			case date != nil && main == 'j': // %j
				if dayOfYear, err = integer(3); err != nil {
					return err
				}
			case date != nil && main == 'm': // %m
				if month, err = integer(2); err != nil {
					return err
				}
			case time != nil && main == 'M': // %M
				if min, err = integer(2); err != nil {
					return err
				}
			case time != nil && main == 'p': // %p
				lower, original := alphas(2)
				switch strings.ToUpper(lower) {
				case "AM":
				case "PM":
					isAfternoon = true
				default:
					return fmt.Errorf("failed to parse time of day %q", original)
				}
			case time != nil && main == 'P': // %P
				lower, original := alphas(2)
				switch lower {
				case "am":
				case "pm":
					isAfternoon = true
				default:
					return fmt.Errorf("failed to parse time of day %q", original)
				}
			case time != nil && main == 'S': // %S
				if sec, err = integer(2); err != nil {
					return err
				}
			case date != nil && main == 'u': // %u
				if dayOfWeek, err = integer(1); err != nil {
					return err
				}
			case date != nil && main == 'V': // %V
				haveISODate = true
				if isoWeek, err = integer(2); err != nil {
					return err
				}
			case date != nil && main == 'y': // %y
				if localed { // %Ey
					haveGregorianYear = true
				}

				if year, err = integer(2); err != nil {
					return err
				}
				year += getCentury(year)
			case date != nil && main == 'Y': // %Y
				if localed { // %EY
					haveGregorianYear = true
				}

				if year, err = integer(4); err != nil {
					return err
				}
			case time != nil && main == 'z': // %z
				// Parsing %z into a type that contains no offset (e.g. LocalTime, LocalDateTime)
				// is valid, although the value itself is ignored.
				if offset == nil {
					break
				}

				// Catch the 'Z' case, which is valid for both %z and %Ez.
				if _, ok := casedAlpha('Z'); ok {
					*offset = 0
					break
				}

				h, err := integer(2)
				if err != nil {
					return err
				}

				var m int
				if !hasMore() {
					goto CalculateOffset
				}

				if localed { // %Ez
					if actual, ok := casedAlpha(':'); !ok {
						return fmt.Errorf(extraTextErrMsg, value, string(actual))
					}
				}

				if m, err = integer(2); err != nil {
					return err
				}

			CalculateOffset:
				if h >= 0 {
					*offset = int64(h)*oneHour + int64(m)*oneMinute
				} else {
					*offset = int64(h)*oneHour - int64(m)*oneMinute
				}
			case main == '%': // %%
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
			specifierComplete = isSpecifier && (buf[len(buf)-1] != '-' && buf[len(buf)-1] != 'E' && (buf[len(buf)-1] < '0' || buf[len(buf)-1] > '9'))
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
		return fmt.Errorf(extraTextErrMsg, value, value[pos:])
	}

	if date != nil {
		if haveGregorianYear {
			if year, err = convertGregorianToISOYear(year, isBCE); err != nil {
				return err
			}
		}

		if !isDateValid(year, month, day) {
			return fmt.Errorf("invalid date %q", simpleDateStr(year, month, day))
		}

		_date, err := makeDate(year, month, day)
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
					simpleDateStr(year, month, day),
				)
			}

			*date = doyDate
		}

		// Check ISO week-year according to note (2).
		if haveISODate {
			weekday := dayOfWeek
			if dayOfWeek == 0 {
				weekday = int(Monday)
			}

			isoDate, err := ofISOWeek(isoYear, isoWeek, weekday)
			if err != nil {
				return fmt.Errorf("invalid ISO week-year date %q", getISODateSimpleStr(isoYear, isoWeek, day))
			}

			if haveDate && (isoDate != _date) {
				return fmt.Errorf("ISO week-year date %q does not agree with date %q",
					getISODateSimpleStr(isoYear, isoWeek, day),
					simpleDateStr(year, month, day),
				)
			}

			*date = isoDate
		}

		// Check day of week according to note (3).
		haveDate = haveDate || dayOfYear != 0
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
		// Check validity of hour on 12-hour clock according to note (5).
		if have12HourClock {
			if hour < 1 || hour > 12 {
				return fmt.Errorf("hour %d is not valid on the 12-hour clock", hour)
			}
			hour = convert12To24HourClock(hour, isAfternoon)
		}

		v, err := makeTime(hour, min, sec, nsec)
		if err != nil {
			return err
		}
		*time = v
	}

	return nil
}

func parseSpecifier(buf []rune) (nopad, localed bool, precision uint, main rune, err error) {
	if len(buf) == 3 {
		switch {
		case buf[1] == '-':
			nopad = true
		case buf[1] == 'E':
			localed = true
		case buf[1] >= '0' && buf[1] <= '9':
			precision = uint(buf[1] - 48)
		default:
			return false, false, 0, 0, fmt.Errorf("unsupported modifier '%c'", buf[1])
		}
	} else if len(buf) == 4 {
		switch buf[1] {
		case '-':
			nopad = true
		default:
			return false, false, 0, 0, fmt.Errorf("unsupported modifier '%c'", buf[1])
		}

		switch buf[2] {
		case 'E':
			localed = true
		default:
			return false, false, 0, 0, fmt.Errorf("unsupported modifier '%c'", buf[1])
		}
	}
	return nopad, localed, precision, buf[len(buf)-1], nil
}

func convert12To24HourClock(hour12 int, isAfternoon bool) (hour24 int) {
	if isAfternoon && hour12 == 12 {
		return 12
	} else if isAfternoon {
		return hour12 + 12
	} else if hour12 == 12 {
		return 0
	}
	return hour12
}

func convert24To12HourClock(hour24 int) (hour12 int, isAfternoon bool) {
	if hour24 == 0 {
		return 12, false
	} else if hour24 == 12 {
		return 12, true
	} else if hour24 < 12 {
		return hour24, false
	}
	return hour24 % 12, true
}

func convertGregorianToISOYear(gregorianYear int, isBCE bool) (isoYear int, err error) {
	if gregorianYear == 0 {
		return 0, fmt.Errorf("invalid Gregorian year %04d", gregorianYear)
	}

	if isBCE {
		return (gregorianYear * -1) + 1, nil
	}
	return gregorianYear, nil
}

func convertISOToGregorianYear(isoYear int) (gregorianYear int, isBCE bool) {
	if isoYear <= 0 {
		return (isoYear * -1) + 1, true
	}
	return isoYear, false
}

func shortWeekdayName(d int) string {
	if d < int(Monday) || d > int(Sunday) {
		return fmt.Sprintf("%%!Weekday(%d)", d)
	}
	return shortDayNames[d-1]
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
	Monday - 1:    "Mon",
	Tuesday - 1:   "Tue",
	Wednesday - 1: "Wed",
	Thursday - 1:  "Thu",
	Friday - 1:    "Fri",
	Saturday - 1:  "Sat",
	Sunday - 1:    "Sun",
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
