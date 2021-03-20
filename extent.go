package chrono

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

func (e Extent) Nanoseconds() int64 {
	return int64(e)
}

func (e Extent) Microseconds() float64 {
	micros := e / Second
	nsec := e % Second
	return float64(micros) + float64(nsec)/1e3
}

func (e Extent) Milliseconds() float64 {
	millis := e / Second
	nsec := e % Second
	return float64(millis) + float64(nsec)/1e6
}

func (e Extent) Seconds() float64 {
	sec := e / Second
	nsec := e % Second
	return float64(sec) + float64(nsec)/1e9
}

func (e Extent) Minutes() float64 {
	min := e / Minute
	nsec := e % Minute
	return float64(min) + float64(nsec)/(60*1e9)
}

func (e Extent) Hours() float64 {
	hour := e / Hour
	nsec := e % Hour
	return float64(hour) + float64(nsec)/(60*60*1e9)
}

func (e Extent) Truncate(m Extent) Extent {
	if e <= 0 {
		return e + e%m
	}
	return e - e%m
}
