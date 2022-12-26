package chrono

import "fmt"

// Weekday specifies the day of the week (Monday = 0, ...).
// Not compatible standard library's time.Weekday (in which Sunday = 0, ...).
type Weekday int

// The days of the week.
const (
	Monday Weekday = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

func (d Weekday) String() string {
	return longWeekdayName(int(d))
}

func longWeekdayName(d int) string {
	if d > int(Sunday) {
		return fmt.Sprintf("%%!Weekday(%d)", d)
	}
	return longDayNames[d]
}

var longDayNames = [7]string{
	Monday:    "Monday",
	Tuesday:   "Tuesday",
	Wednesday: "Wednesday",
	Thursday:  "Thursday",
	Friday:    "Friday",
	Saturday:  "Saturday",
	Sunday:    "Sunday",
}

// Month specifies the month of the year (January = 1, ...).
type Month int

// The months of the year.
const (
	January Month = iota + 1
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)

func (m Month) String() string {
	return longMonthName(int(m))
}

func longMonthName(m int) string {
	if m < int(January) || m > int(December) {
		return fmt.Sprintf("%%!Month(%d)", m)
	}
	return longMonthNames[m-1]
}

var longMonthNames = [12]string{
	January - 1:   "January",
	February - 1:  "February",
	March - 1:     "March",
	April - 1:     "April",
	May - 1:       "May",
	June - 1:      "June",
	July - 1:      "July",
	August - 1:    "August",
	September - 1: "September",
	October - 1:   "October",
	November - 1:  "November",
	December - 1:  "December",
}
