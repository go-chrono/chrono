package chrono

import "strconv"

// Instant represents an instantaneous point in time with nanosecond resolution.
type Instant struct {
	*int64
}

// Now returns the Instant that represents the current point in time.
func Now() Instant {
	now := monotime()
	return Instant{
		int64: &now,
	}
}

// Elapsed is shorthand for i.Until(chrono.Now()).
func (i Instant) Elapsed() Duration {
	return i.Until(Now())
}

func (i Instant) String() string {
	return strconv.FormatInt(*i.int64, 10)
}

// Until returns the Duration that represents the elapsed time from i to v.
func (i Instant) Until(v Instant) Duration {
	switch {
	case i.int64 == nil:
		panic("i is not initialized")
	case v.int64 == nil:
		panic("v is not initialized")
	}

	iv, vv := *i.int64, *v.int64
	if vv < iv {
		panic("v is smaller than i")
	}

	d := vv - iv
	return Duration{
		secs: uint64(d / 1e9),
		nsec: uint32(d % 1e9),
	}
}

func instant(t int64) Instant {
	return Instant{
		int64: &t,
	}
}
