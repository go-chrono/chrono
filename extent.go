package chrono

// Extent represents a period of time measured in nanoseconds.
// The represented value is exactly equivalent to the standard library's time.Duration.
type Extent int

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

// Truncate returns the result of rounding e toward zero to a multiple of m.
func (e Extent) Truncate(m Extent) Extent {
	if m <= 0 {
		return e
	}
	return e - e%m
}
