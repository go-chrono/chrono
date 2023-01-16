package chrono

import (
	"fmt"
	"math/big"
	"strconv"
)

// Duration represents a period of time with nanosecond precision,
// with a range of approximately ±292,300,000,000 years.
type Duration struct {
	v big.Int
}

// DurationOf creates a new duration from the supplied extent.
// Durations and extents are semantically equivalent, except that durations exceed,
// and can therefore not be converted to, Go's basic types. Extents are represented as a single integer.
func DurationOf(v Extent) Duration {
	return Duration{v: *big.NewInt(int64(v))}
}

// Compare compares d with d2. If d is less than d2, it returns -1;
// if d is greater than d2, it returns 1; if they're equal, it returns 0.
func (d Duration) Compare(d2 Duration) int {
	return d.v.Cmp(&d2.v)
}

// Add returns the duration d+d2.
// If the operation would overflow the maximum duration, or underflow the minimum duration, it panics.
// Use CanAdd to test whether aa panic would occur.
func (d Duration) Add(d2 Duration) Duration {
	out, err := d.add(d2)
	if err != nil {
		panic(err.Error())
	}
	return out
}

// CanAdd returns false if Add would panic if passed the same argument.
func (d Duration) CanAdd(d2 Duration) bool {
	_, err := d.add(d2)
	return err == nil
}

func (d Duration) add(d2 Duration) (Duration, error) {
	out := new(big.Int).Set(&d.v)
	out.Add(out, &d2.v)

	if out.Cmp(bigIntMinInt64) == -1 || out.Cmp(bigIntMaxInt64) == 1 {
		return Duration{}, fmt.Errorf("duration out of range")
	}
	return Duration{v: *out}, nil
}

// Nanoseconds returns the duration as a floating point number of nanoseconds.
func (d Duration) Nanoseconds() float64 {
	out, _ := new(big.Float).SetInt(&d.v).Float64()
	return out
}

// Microseconds returns the duration as a floating point number of microseconds.
func (d Duration) Microseconds() float64 {
	out, _ := new(big.Float).Quo(new(big.Float).SetInt(&d.v), bigFloatMicrosecondExtent).Float64()
	return out
}

// Milliseconds returns the duration as a floating point number of milliseconds.
func (d Duration) Milliseconds() float64 {
	out, _ := new(big.Float).Quo(new(big.Float).SetInt(&d.v), bigFloatMillisecondExtent).Float64()
	return out
}

// Seconds returns the duration as a floating point number of seconds.
func (d Duration) Seconds() float64 {
	out, _ := new(big.Float).Quo(new(big.Float).SetInt(&d.v), bigFloatSecondExtent).Float64()
	return out
}

// Minutes returns the duration as a floating point number of minutes.
func (d Duration) Minutes() float64 {
	out, _ := new(big.Float).Quo(new(big.Float).SetInt(&d.v), bigFloatMinuteExtent).Float64()
	return out
}

// Hours returns the duration as a floating point number of hours.
func (d Duration) Hours() float64 {
	out, _ := new(big.Float).Quo(new(big.Float).SetInt(&d.v), bigFloatHourExtent).Float64()
	return out
}

// String returns a string formatted according to ISO 8601.
// It is equivalent to calling Format with no arguments.
func (d Duration) String() string {
	return d.Format()
}

// Units returns the whole numbers of hours, minutes, seconds, and nanosecond offset represented by d.
func (d Duration) Units() (hours, mins, secs, nsec int) {
	_secs, _nsec, _ := d.integers()
	hours = int(_secs / 3600)
	mins = int((_secs / 60) % 60)
	secs = int(_secs % 60)
	nsec = int(_nsec)
	return
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
	out, neg := d.format(exclusive...)
	out = "P" + out
	if neg {
		out = "-" + out
	}
	return out
}

func (d Duration) format(exclusive ...Designator) (_ string, neg bool) {
	secs, nsec, neg := d.integers()
	return formatDuration(secs, nsec, neg, exclusive...)
}

func formatDuration(secs int64, nsec uint32, neg bool, exclusive ...Designator) (_ string, isNeg bool) {
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
		if v := float64(secs / 3600); v != 0 {
			values[Hours] = v
			h = true
		}
	case h && (m || s):
		values[Hours] = float64(secs / 3600)
	case h:
		values[Hours] = (float64(secs) / 3600) + (float64(nsec) / 3.6e12)
	}

	switch {
	case len(exclusive) == 0:
		if v := float64((secs % 3600) / 60); v != 0 {
			values[Minutes] = v
			m = true
		}
	case m && s && h:
		values[Minutes] = float64((secs % 3600) / 60)
	case m && s:
		values[Minutes] = float64(secs / 60)
	case m && h:
		values[Minutes] = (float64(secs%3600) / 60) + (float64(nsec) / 6e10)
	case m:
		values[Minutes] = (float64(secs) / 60) + (float64(nsec) / 6e10)
	}

	switch {
	case len(exclusive) == 0:
		if v := float64(secs%60) + (float64(nsec) / 1e9); v != 0 {
			values[Seconds] = v
			if h && !m {
				values[Minutes] = 0
			}
		} else if !h && !m {
			values[Seconds] = 0
		}
	case s && m:
		values[Seconds] = float64(secs%60) + (float64(nsec) / 1e9)
	case s && h:
		values[Seconds] = float64(secs%3600) + (float64(nsec) / 1e9)
	case s:
		values[Seconds] = float64(secs) + (float64(nsec) / 1e9)
	}

	out := "T"
	if v, ok := values[Hours]; ok {
		out += strconv.FormatFloat(v, 'f', -1, 64) + "H"
	}

	if v, ok := values[Minutes]; ok {
		out += strconv.FormatFloat(v, 'f', -1, 64) + "M"
	}

	if v, ok := values[Seconds]; ok {
		out += strconv.FormatFloat(v, 'f', -1, 64) + "S"
	}
	return out, neg
}

func (d Duration) integers() (secs int64, nsec uint32, neg bool) {
	v := new(big.Int).Abs(&d.v)
	var _nsec big.Int
	_secs, _ := new(big.Int).DivMod(v, bigIntSecondExtent, &_nsec)
	return _secs.Int64(), uint32(_nsec.Uint64()), d.v.Cmp(bigIntNegOne) == -1
}

// Parse the time portion of an ISO 8601 duration.
func (d *Duration) Parse(s string) error {
	_, _, _, _, secs, nsec, neg, err := parseDuration(s, false, true)
	if err != nil {
		return err
	}

	*d = makeDuration(secs, nsec, neg)
	return nil
}

// MinDuration returns the minimum supported duration.
func MinDuration() Duration {
	return Duration{v: *bigIntMinInt64}
}

// MaxDuration returns the maximum supported duration.
func MaxDuration() Duration {
	return Duration{v: *bigIntMaxInt64}
}

func makeDuration(secs int64, nsec uint32, neg bool) Duration {
	out := new(big.Int).Mul(big.NewInt(secs), bigIntSecondExtent)
	out.Add(out, big.NewInt(int64(nsec)))
	if neg {
		out.Neg(out)
	}
	return Duration{v: *out}
}

var (
	bigIntNegOne = big.NewInt(-1)

	bigIntMinInt64 = new(big.Int).Lsh(big.NewInt(int64(-Second)), 63)
	bigIntMaxInt64 = new(big.Int).Add(new(big.Int).Lsh(big.NewInt(int64(Second)), 63), big.NewInt(-1))

	bigIntSecondExtent = big.NewInt(int64(Second))

	bigFloatMicrosecondExtent = big.NewFloat(float64(Microsecond))
	bigFloatMillisecondExtent = big.NewFloat(float64(Millisecond))
	bigFloatSecondExtent      = big.NewFloat(float64(Second))
	bigFloatMinuteExtent      = big.NewFloat(float64(Minute))
	bigFloatHourExtent        = big.NewFloat(float64(Hour))
)
