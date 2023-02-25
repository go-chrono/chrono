package chrono

import (
	"fmt"
)

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
	if hours == 0 {
		return Offset(Extent(mins) * Minute)
	}

	if mins < 0 {
		mins = -mins
	}

	if hours < 0 {
		return Offset((Extent(hours) * Hour) - (Extent(mins) * Minute))
	}
	return Offset((Extent(hours) * Hour) + (Extent(mins) * Minute))
}

// String returns the time zone designator according to ISO 8601, truncating first to the minute.
// If o == 0, String returns "Z" for the UTC offset.
// In all other cases, a string in the format of ±hh:mm is returned.
// Note that the sign and number of minutes is always included, even if 0.
func (o Offset) String() string {
	e := Extent(o).Truncate(Minute)
	if e == 0 {
		return "Z"
	}

	sign := "+"
	if e < 0 {
		sign = "-"
	}

	hours, mins, _, _ := e.abs().Units()
	return fmt.Sprintf("%s%02d:%02d", sign, hours, mins)
}

// func Parse
