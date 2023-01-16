package chrono

import (
	"fmt"
	"math"
)

// Extent represents a period of time measured in nanoseconds.
// The represented value is exactly equivalent to the standard library's time.Duration.
type Extent int64

// Common time-based durations relative to 1 nanosecond.
const (
	Nanosecond  Extent = 1
	Microsecond        = 1000 * Nanosecond
	Millisecond        = 1000 * Microsecond
	Second             = 1000 * Millisecond
	Minute             = 60 * Second
	Hour               = 60 * Minute
)

// Nanoseconds returns the extent as an integer nanosecond count.
func (e Extent) Nanoseconds() int64 {
	return int64(e)
}

// Microseconds returns the duration as a floating point number of microseconds.
func (e Extent) Microseconds() float64 {
	micros := e / Microsecond
	nsec := e % micros
	return float64(micros) + float64(nsec)/1e3
}

// Milliseconds returns the duration as a floating point number of milliseconds.
func (e Extent) Milliseconds() float64 {
	millis := e / Millisecond
	nsec := e % millis
	return float64(millis) + float64(nsec)/1e6
}

// Seconds returns the duration as a floating point number of seconds.
func (e Extent) Seconds() float64 {
	secs := e / Second
	nsec := e % secs
	return float64(secs) + float64(nsec)/1e9
}

// Minutes returns the duration as a floating point number of minutes.
func (e Extent) Minutes() float64 {
	mins := e / Minute
	nsec := e % mins
	return float64(mins) + float64(nsec)/(60*1e9)
}

// Hours returns the duration as a floating point number of hours.
func (e Extent) Hours() float64 {
	hours := e / Hour
	nsec := e % hours
	return float64(hours) + float64(nsec)/(60*60*1e9)
}

// Units returns the whole numbers of hours, minutes, seconds, and nanosecond offset represented by e.
func (e Extent) Units() (hours, mins, secs, nsec int) {
	hours = int(e / Hour)
	mins = int(e/Minute) % 60
	secs = int(e/Second) % 60
	nsec = int(e % Second)
	return
}

// Truncate returns the result of rounding e toward zero to a multiple of m.
func (e Extent) Truncate(m Extent) Extent {
	if m <= 0 {
		return e
	}
	return e - e%m
}

// String returns a string formatted according to ISO 8601.
// It is equivalent to calling Format with no arguments.
func (e Extent) String() string {
	return e.Format()
}

// Format the extent according to ISO 8601.
// Behaves the same as Duration.Format.
func (e Extent) Format(exclusive ...Designator) string {
	abs := e.abs()
	out, neg := formatDuration(int64(abs/Second), uint32(abs%Second), e < 0, exclusive...)
	out = "P" + out
	if neg {
		out = "-" + out
	}
	return out
}

// Parse the time portion of an ISO 8601 duration.
// Behaves the same as Duration.Parse.
func (e *Extent) Parse(s string) error {
	_, _, _, _, secs, nsec, neg, err := parseDuration(s, false, true)
	if err != nil {
		return err
	}

	if err := checkExtentRange(secs, nsec, neg); err != nil {
		return err
	}

	*e = Extent(secs*int64(Second) + int64(nsec))
	if neg {
		*e *= -1
	}

	return nil
}

func (e Extent) abs() Extent {
	if e < 0 {
		return e * -1
	}
	return e
}

func checkExtentRange(secs int64, nsec uint32, neg bool) error {
	if neg {
		switch {
		case -secs < minSeconds:
			return fmt.Errorf("seconds underflow")
		case -secs == minSeconds && nsec > uint32(maxNegNanos):
			return fmt.Errorf("nanoseconds underflow")
		}
	} else {
		switch {
		case secs > maxSeconds:
			return fmt.Errorf("seconds overflow")
		case secs == maxSeconds && nsec > uint32(maxPosNanos):
			return fmt.Errorf("nanoseconds overflow")
		}
	}
	return nil
}

const (
	minSeconds  = int64(math.MinInt64) / int64(Second)
	maxNegNanos = -(int64(math.MinInt64) % -minSeconds)
	maxSeconds  = int64(math.MaxInt64) / int64(Second)
	maxPosNanos = int64(math.MaxInt64) % maxSeconds
)
