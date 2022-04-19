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

func (w Weekday) String() string {
	if w > Sunday {
		return fmt.Sprintf("%%!Weekday(%d)", w)
	}
	return longDayNames[w]
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
	if m < January || m > December {
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

var longMonthNameLookup = map[string]Month{
	"January":   January,
	"February":  February,
	"March":     March,
	"April":     April,
	"May":       May,
	"June":      June,
	"July":      July,
	"August":    August,
	"September": September,
	"October":   October,
	"November":  November,
	"December":  December,
}
