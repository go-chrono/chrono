package chrono

import "math"

type Duration struct {
	secs uint64
	nsec uint32
}

type Extent float64

// Common time-based durations relative to 1 nanosecond.
const (
	Nanosecond  Extent = 1
	Microsecond        = 1000 * Nanosecond
	Millisecond        = 1000 * Microsecond
	Second             = 1000 * Millisecond
	Minute             = 60 * Second
	Hour               = 60 * Minute
)

func DurationOf(v Extent) Duration {
	secs := v / 1e9
	if secs > math.MaxUint64 {
		panic("v overflows Duration")
	}

	return Duration{
		secs: uint64(v / 1e9),
		nsec: uint32(math.Mod(float64(v), 1e9)),
	}
}

func (d Duration) Add(d2 Duration) Duration {
	nsec := float64(d2.secs*1e9) + float64(d2.nsec) + float64(d.secs*1e9) + float64(d.nsec)
	if nsec > math.MaxUint64 {
		panic("d2 + d overflows Duration")
	}

	return Duration{
		secs: uint64(nsec / 1e9),
		nsec: uint32(math.Mod(nsec, 1e9)),
	}
}

func (d Duration) Nanoseconds() float64 {
	return float64(d.secs*1e9) + float64(d.nsec)
}

func (d Duration) Microseconds() float64 {
	return float64(d.secs*1e6) + float64(d.nsec/1e3)
}

func (d Duration) Milliseconds() float64 {
	return float64(d.secs*1e3) + float64(d.nsec/1e6)
}

func (d Duration) Seconds() float64 {
	return float64(d.secs) + float64(d.nsec/1e9)
}

func (d Duration) Minutes() float64 {
	return float64(d.secs/60) + float64(d.nsec/1e9*60)
}

func (d Duration) Hours() float64 {
	return float64(d.secs)/(60*60) + float64(d.nsec)/(1e9*60*60)
}

// MinDuration returns the minimum allowed Duration of 0 ns.
func MinDuration() *Duration {
	return &minDuration
}

// MaxDuration returns the maximum allowed Duration of 18446744073709551615 s and 999999999 ns.
func MaxDuration() *Duration {
	return &maxDuration
}

var minDuration = Duration{secs: 0, nsec: 0}

var maxDuration = Duration{secs: math.MaxUint64, nsec: 1e9 - 1}
