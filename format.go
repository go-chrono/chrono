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
// For specifiers that represent padded decimals, leading 0s can be omitted using the '-' character after the '%'.
// For example, '%m' may represent the string '04' (for March), but '%-m' represents '4'.
//
// Depending on the context in which the layout is used, only a subset of specifiers may be supported by a particular function.
// For example, %H is not supported when parsing or formatting a date.
//
// If a specifier is encountered which is not recognized (defined in the list above), or no supported by a particular function,
// the function will panic with a message that includes the unrecognized sequence.
//
// Any other text is enchoed verbatim when formatting, and is expected to appear verbatim in the parsed text.
// In order to print the '%' verbatim character (which normally signifies a specifier), the sequence '%%' can be used.
//
// For familiarity, the examples below use the time package's reference time of "2nd Jan 2006 15:04:05 -0700" (Unix time 1136239445).
// But note that this reference format is not relevant at all to the functioning of this package.
//
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

	var buf []rune
NextChar:
	for _, c := range layout {
		buf = append(buf, c)

		if len(buf) >= 2 && buf[0] == '%' {
			if c == '-' || c == 'E' {
				continue NextChar
			}

			var nopad, localed bool
			if len(buf) == 3 {
				switch buf[1] {
				case '-':
					nopad = true
				case 'E':
					localed = true
				default:
					panic(fmt.Sprintf("unsupported modifier '%c'", buf[1]))
				}
			} else if len(buf) == 4 {
				switch buf[1] {
				case '-':
					nopad = true
				default:
					panic(fmt.Sprintf("unsupported modifier '%c'", buf[1]))
				}

				switch buf[2] {
				case 'E':
					localed = true
				default:
					panic(fmt.Sprintf("unsupported modifier '%c'", buf[1]))
				}
			}

			decimal := func(v int, len int) string {
				if nopad {
					return strconv.Itoa(v)
				}
				return fmt.Sprintf("%0*d", len, v)
			}

			main := buf[len(buf)-1]

			switch {
			case date != nil && main == 'a':
				out = append(out, []rune(date.Weekday().short())...)
			case date != nil && main == 'A':
				out = append(out, []rune(date.Weekday().String())...)
			case date != nil && main == 'b':
				out = append(out, []rune(month.short())...)
			case date != nil && main == 'B':
				out = append(out, []rune(month.String())...)
			case date != nil && localed && main == 'C':
				if year > 1 {
					out = append(out, []rune("CE")...)
				} else {
					out = append(out, []rune("BCE")...)
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
				y := year
				if y < 0 {
					y = y*-1 + 1
				}
				out = append(out, []rune(fmt.Sprintf("%02d", y%100))...)
			case date != nil && main == 'Y':
				out = append(out, []rune(decimal(year, 4))...)
			case date != nil && localed && main == 'Y':
				y := year
				if y < 0 {
					y = y*-1 + 1
				}
				out = append(out, []rune(decimal(y, 4))...)
			case main == '%':
				out = append(out, '%')
			default:
				panic("unsupported sequence " + string(buf))
			}

			buf = nil
		} else if len(buf) == 1 && buf[0] != '%' {
			out = append(out, buf...)
			buf = nil
		}
	}

	return string(out)
}

func parse(layout, value string, date *LocalDate, time *LocalTime) error {
	var pos int

	var buf []byte
NextChar:
	for i := 0; i <= len(layout); i++ {
		// Some short-hands.
		var (
			valid = i < len(layout)
			last  = i == len(layout)
		)

		if valid {
			c := layout[i]

			if len(buf) == 0 {
				goto AppendToBuffer
			} else if len(buf) >= 2 && buf[0] == '%' { // Process a specifier.
				switch c {
				case 'E':
					if last {
						// TODO error
					} else {
						goto AppendToBuffer
					}
				case '%':
					goto VerifyText
				}

				localed := len(buf) == 3 && buf[1] == 'E'
				switch {
				// handle specifiers
				}

				_ = localed // TODO delete

				buf = nil
				continue NextChar
			}

		AppendToBuffer:
			buf = append(buf, c)
		}

	VerifyText:
		if len(buf) == 0 {
			continue NextChar
		}

		verify := func() error {
			if string(buf) > value[pos:] {
				return fmt.Errorf("parsing time \"%s\" as \"%s\": cannot parse \"%s\" as \"%s\"", value, layout, value[pos:], buf)
			}
			return nil
		}

		if buf[len(buf)-1] == '%' {
			if err := verify(); err != nil {
				return err
			}
			buf = []byte{'%'}
		} else {
			if err := verify(); err != nil {
				return err
			}
		}
	}

	if pos < len(value) {
		return fmt.Errorf("parsing time \"%s\": extra text: \"%s\"", value, value[len(buf):])
	}

	return nil
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
