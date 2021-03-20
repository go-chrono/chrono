package chrono

import "strconv"

// Instant represents an instantaneous point in time with nanosecond resolution.
type Instant struct {
	v *int64
}

// Now returns the Instant that represents the current point in time.
func Now() Instant {
	now := monotime()
	return Instant{
		v: &now,
	}
}

// Elapsed is shorthand for i.Until(chrono.Now()).
func (i Instant) Elapsed() Duration {
	return i.Until(Now())
}

func (i Instant) String() string {
	return strconv.FormatInt(*i.v, 10)
}

// Until returns the Duration that represents the elapsed time from i to v.
func (i Instant) Until(v Instant) Duration {
	switch {
	case i.v == nil:
		panic("i is not initialized")
	case v.v == nil:
		panic("v is not initialized")
	}

	iv, vv := *i.v, *v.v
	if vv < iv {
		panic("v is smaller than i")
	}

	d := vv - iv
	return Duration{
		secs: d / 1e9,
		nsec: uint32(d % 1e9),
	}
}

func instant(t int64) Instant {
	return Instant{
		v: &t,
	}
}
