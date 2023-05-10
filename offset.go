package chrono

import "fmt"

// UTC represents Universal Coordinated Time (UTC).
const UTC = Offset(0)

// Offset represents a time zone offset from UTC with precision to the minute.
type Offset Extent

// OffsetOf returns the Offset represented by a number of hours and minutes.
// If hours is non-zero, the sign of minutes is ignored, e.g.:
//   * OffsetOf(-2, 30) = -02h:30m
//   * OffsetOf(2, -30) = 02h:30m
//   * OffsetOf(0, 30) = 00h:30m
//   * OffsetOf(0, -30) = -00h:30m
func OffsetOf(hours, mins int) Offset {
	return Offset(makeOffset(hours, mins))
}

func makeOffset(hours, mins int) int64 {
	if hours == 0 {
		return int64(mins) * oneMinute
	}

	if mins < 0 {
		mins = -mins
	}

	if hours < 0 {
		return (int64(hours) * oneHour) - (int64(mins) * oneMinute)
	}
	return (int64(hours) * oneHour) + (int64(mins) * oneMinute)
}

// String returns the time zone designator according to ISO 8601, truncating first to the minute.
// If o == 0, String returns "Z" for the UTC offset.
// In all other cases, a string in the format of ±hh:mm is returned.
// Note that the sign and number of minutes is always included, even if 0.
func (o Offset) String() string {
	return offsetString(int64(o), ":")
}

func offsetString(o int64, sep string) string {
	e := truncateExtent(o, oneMinute)
	if e == 0 {
		return "Z"
	}

	sign := "+"
	if e < 0 {
		sign = "-"
	}

	hours, mins, _, _ := extentUnits(extentAbs(e))
	return fmt.Sprintf("%s%02d%s%02d", sign, hours, sep, mins)
}
