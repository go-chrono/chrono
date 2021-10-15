package chrono

import (
	"math"
	"strconv"
)

// Duration represents a period of time with nanosecond precision,
// with a range of approximately ±292,300,000,000 years.
type Duration struct {
	secs int64
	nsec uint32
}

// DurationOf creates a new duration from the supplied extent.
// Durations and extents are semantically equivalent, except that durations exceed,
// and can therefore not be converted to, Go's basic types. Extents are represented as a single integer.
func DurationOf(v Extent) Duration {
	return Duration{
		secs: int64(v / 1e9),
		nsec: uint32(math.Mod(float64(v), 1e9)),
	}
}

// Equal reports whether d and d2 represent the same duration of time.
func (d Duration) Equal(d2 Duration) bool {
	return d2.secs == d.secs && d2.nsec == d.nsec
}

// Add returns the duration d+d2.
func (d Duration) Add(d2 Duration) Duration {
	nsec := float64(d2.secs*1e9) + float64(d2.nsec) + float64(d.secs*1e9) + float64(d.nsec) // TODO fix me - float
	if nsec > math.MaxUint64 {
		panic("d2 + d overflows Duration")
	}

	return Duration{
		secs: int64(nsec / 1e9),
		nsec: uint32(math.Mod(nsec, 1e9)),
	}
}

// Nanoseconds returns the duration as a floating point number of nanoseconds.
func (d Duration) Nanoseconds() float64 {
	return float64(d.secs*1e9) + float64(d.nsec)
}

// Microseconds returns the duration as a floating point number of microseconds.
func (d Duration) Microseconds() float64 {
	return float64(d.secs*1e6) + float64(d.nsec/1e3)
}

// Milliseconds returns the duration as a floating point number of milliseconds.
func (d Duration) Milliseconds() float64 {
	return float64(d.secs*1e3) + float64(d.nsec/1e6)
}

// Seconds returns the duration as a floating point number of seconds.
func (d Duration) Seconds() float64 {
	return float64(d.secs) + float64(d.nsec/1e9)
}

// Minutes returns the duration as a floating point number of minutes.
func (d Duration) Minutes() float64 {
	return float64(d.secs/60) + float64(d.nsec/1e9*60)
}

// Hours returns the duration as a floating point number of hours.
func (d Duration) Hours() float64 {
	return float64(d.secs)/(60*60) + float64(d.nsec)/(1e9*60*60)
}

// String returns a string formatted according to ISO 8601.
// It is equivalent to calling Format with no arguments.
func (d Duration) String() string {
	return d.Format()
}

// Designator of date and time elements present in ISO 8601.
type Designator rune

// Designators.
const (
	Hours   Designator = 'H'
	Minutes Designator = 'M'
	Seconds Designator = 'S'
)

// Format the duration according to ISO 8601.
// The output consists of only the time component - the period component is never included.
// Thus the output always consists of "PT" followed by at least one unit of the time component (hours, minutes, seconds).
//
// The default format, obtained by calling the function with no arguments, consists of the most significant non-zero units,
// presented non-breaking but trimmed.
// 5 minutes is formatted as PT5M, trimming 0-value hours and seconds.
// 1 hour and 5 seconds is formatted at PT1H0M5S, ensuring the sequence of units is not broken.
//
// A list of designators can be optionally passed to the function in order to control which units are included.
// When passed, only those specified units are included in the formatted string, and are present regardless of whether their value is 0.
//
// Fractional values are automatically applied to the least significant unit, if applicable.
// In order to format only integers, the round functions should be used before calling this function.
func (d Duration) Format(exclusive ...Designator) string {
	values := make(map[Designator]float64, 3)
	if len(exclusive) >= 1 {
		for _, d := range exclusive {
			values[d] = 0
		}
	}

	_, h := values[Hours]
	_, m := values[Minutes]
	_, s := values[Seconds]

	switch {
	case len(exclusive) == 0:
		if v := float64(d.secs / 3600); v != 0 {
			values[Hours] = v
			h = true
		}
	case h && (m || s):
		values[Hours] = float64(d.secs / 3600)
	case h:
		values[Hours] = (float64(d.secs) / 3600) + (float64(d.nsec) / 3.6e12)
	}

	switch {
	case len(exclusive) == 0:
		if v := float64((d.secs % 3600) / 60); v != 0 {
			values[Minutes] = v
			m = true
		}
	case m && s && h:
		values[Minutes] = float64((d.secs % 3600) / 60)
	case m && s:
		values[Minutes] = float64(d.secs / 60)
	case m && h:
		values[Minutes] = (float64(d.secs%3600) / 60) + (float64(d.nsec) / 6e10)
	case m:
		values[Minutes] = (float64(d.secs) / 60) + (float64(d.nsec) / 6e10)
	}

	switch {
	case len(exclusive) == 0:
		if v := float64(d.secs%60) + (float64(d.nsec) / 1e9); v != 0 {
			values[Seconds] = v
			if h && !m {
				values[Minutes] = 0
			}
		}
	case s && m:
		values[Seconds] = float64(d.secs%60) + (float64(d.nsec) / 1e9)
	case s && h:
		values[Seconds] = float64(d.secs%3600) + (float64(d.nsec) / 1e9)
	case s:
		values[Seconds] = float64(d.secs) + (float64(d.nsec) / 1e9)
	}

	out := "PT"
	if v, ok := values[Hours]; ok {
		out += strconv.FormatFloat(v, 'f', -1, 64) + "H"
	}

	if v, ok := values[Minutes]; ok {
		out += strconv.FormatFloat(v, 'f', -1, 64) + "M"
	}

	if v, ok := values[Seconds]; ok {
		out += strconv.FormatFloat(v, 'f', -1, 64) + "S"
	}

	return out
}

// Parse an ISO 8601 duration.
func (d *Duration) Parse(s string) error {
	secs, nsec, err := parseDuration(s, false)
	if err != nil {
		return err
	}

	d.secs = secs
	d.nsec = nsec
	return nil
}

// MinDuration returns the minimum allowed Duration of 0 ns.
func MinDuration() *Duration {
	return &minDuration
}

// MaxDuration returns the maximum allowed Duration of 18446744073709551615 s and 999999999 ns (5.846e11 years).
func MaxDuration() *Duration {
	return &maxDuration
}

var minDuration = Duration{secs: 0, nsec: 0}

var maxDuration = Duration{secs: math.MaxInt64, nsec: 1e9 - 1}
