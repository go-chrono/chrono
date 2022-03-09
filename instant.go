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

// Compare compares i with i2. If i is before i2, it returns -1;
// if i is after i2, it returns 1; if they're the same, it returns 0.
func (i Instant) Compare(i2 Instant) int {
	switch {
	case i.v == nil:
		panic("i is not initialized")
	case i2.v == nil:
		panic("i2 is not initialized")
	case *i.v < *i2.v:
		return -1
	case *i.v > *i2.v:
		return 1
	default:
		return 0
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
func (i Instant) Until(i2 Instant) Duration {
	switch {
	case i.v == nil:
		panic("i is not initialized")
	case i2.v == nil:
		panic("i2 is not initialized")
	}

	iv, vv := *i.v, *i2.v
	if vv < iv {
		panic("v is smaller than i")
	}

	d := vv - iv
	return DurationOf(Extent(d))
}

// instant is used by the chronotest package.
func instant(t int64) Instant {
	return Instant{
		v: &t,
	}
}
