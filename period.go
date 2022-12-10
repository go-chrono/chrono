package chrono

import (
	"fmt"
	"math"
	"strconv"
)

// Period represents an amount of time in years, months, weeks and days.
// A period is not a measurable quantity since the lengths of these components is ambiguous.
type Period struct {
	Years  float32
	Months float32
	Weeks  float32
	Days   float32
}

// Equal reports whether p and p2 represent the same period of time.
func (p Period) Equal(p2 Period) bool {
	return p2.Years == p.Years && p2.Months == p.Months && p2.Weeks == p.Weeks && p2.Days == p.Days
}

// Format the duration according to ISO 8601.
// The output consists of only the period component - the time component is never included.
func (p Period) Format() string {
	if p.Years == 0 && p.Months == 0 && p.Weeks == 0 && p.Days == 0 {
		return "P0D"
	}

	out := "P"
	if p.Years != 0 {
		out += strconv.FormatFloat(math.Abs(float64(p.Years)), 'f', -1, 32) + "Y"
	}

	if p.Months != 0 {
		out += strconv.FormatFloat(math.Abs(float64(p.Months)), 'f', -1, 32) + "M"
	}

	if p.Weeks != 0 {
		out += strconv.FormatFloat(math.Abs(float64(p.Weeks)), 'f', -1, 32) + "W"
	}

	if p.Days != 0 {
		out += strconv.FormatFloat(math.Abs(float64(p.Days)), 'f', -1, 32) + "D"
	}
	return out
}

// Parse the period portion of an ISO 8601 duration.
func (p *Period) Parse(s string) error {
	years, months, weeks, days, _, _, err := parseDuration(s, true, false)
	if err != nil {
		return err
	}

	p.Years = years
	p.Months = months
	p.Weeks = weeks
	p.Days = days
	return nil
}

// ParseDuration parses a complete ISO 8601 duration.
func ParseDuration(s string) (Period, Duration, error) {
	years, months, weeks, days, secs, nsec, err := parseDuration(s, true, true)
	return Period{
			Years:  years,
			Months: months,
			Weeks:  weeks,
			Days:   days,
		},
		Duration{
			secs: secs,
			nsec: nsec,
		}, err
}

func parseDuration(s string, parsePeriod, parseTime bool) (years, months, weeks, days float32, secs int64, nsec uint32, err error) {
	if len(s) == 0 || s[0] != 'P' {
		return 0, 0, 0, 0, 0, 0, fmt.Errorf("expecting 'P'")
	}

	var value int
	var onTime bool
	var haveUnit bool
	var haveWeeks int // 0 = undecided, 1 = W, -1 = YMD

	for i := 1; i < len(s); i++ {
		digit := (s[i] >= '0' && s[i] <= '9') || s[i] == '.'

		if value == 0 {
			if digit {
				value = i
			} else if s[i] == 'T' {
				if !onTime {
					onTime = true
				} else {
					return 0, 0, 0, 0, 0, 0, fmt.Errorf("unexpected '%c', expecting digit", s[i])
				}
			} else {
				return 0, 0, 0, 0, 0, 0, fmt.Errorf("unexpected '%c', expecting digit or 'T'", s[i])
			}
		} else {
			if !onTime {
				if !parsePeriod {
					return 0, 0, 0, 0, 0, 0, fmt.Errorf("cannot parse duration as Duration")
				} else if digit {
					continue
				}

				v, err := strconv.ParseFloat(s[value:i], 32)
				if err != nil {
					return 0, 0, 0, 0, 0, 0, err
				}

				if s[i] == 'W' {
					switch {
					case haveWeeks == 0:
						haveWeeks++
					case haveWeeks < 0:
						return 0, 0, 0, 0, 0, 0, fmt.Errorf("cannot mix 'W' with 'Y'/'M'/'D'")
					}

					weeks = float32(v)
				} else {
					switch {
					case haveWeeks == 0:
						haveWeeks--
					case haveWeeks > 0:
						return 0, 0, 0, 0, 0, 0, fmt.Errorf("cannot mix 'Y'/'M'/'D' with 'W'")
					}

					switch s[i] {
					case 'Y':
						years = float32(v)
					case 'M':
						months = float32(v)
					case 'D':
						days = float32(v)
					default:
						return 0, 0, 0, 0, 0, 0, fmt.Errorf("unexpected '%c', expecting 'Y', 'M' or 'D'", s[i])
					}
				}

				value = 0
				haveUnit = true
			} else {
				if digit {
					continue
				}

				v, err := strconv.ParseFloat(s[value:i], 64)
				if err != nil {
					return 0, 0, 0, 0, 0, 0, err
				}

				var _secs float64
				var _nsec uint32
				switch s[i] {
				case 'H':
					_secs = math.Floor(v * 3600)
					_nsec = uint32((v * 3.6e12) - (_secs * 1e9))
				case 'M':
					_secs = math.Floor(v * 60)
					_nsec = uint32((v * 6e10) - (_secs * 1e9))
				case 'S':
					_secs = math.Floor(v)
					_nsec = uint32((v * 1e9) - (_secs * 1e9))
				default:
					return 0, 0, 0, 0, 0, 0, fmt.Errorf("unexpected '%c', expecting 'H', 'M' or 'S'", s[i])
				}

				if _secs < math.MinInt64 {
					return 0, 0, 0, 0, 0, 0, fmt.Errorf("seconds underflow")
				} else if _secs > math.MaxInt64 {
					return 0, 0, 0, 0, 0, 0, fmt.Errorf("seconds overflow")
				}

				var under, over bool
				if secs, under, over = addInt64(secs, int64(_secs)); under {
					return 0, 0, 0, 0, 0, 0, fmt.Errorf("seconds underflow")
				} else if over {
					return 0, 0, 0, 0, 0, 0, fmt.Errorf("seconds overflow")
				}

				if secs, under, over = addInt64(secs, int64(_nsec/1e9)); under {
					return 0, 0, 0, 0, 0, 0, fmt.Errorf("seconds underflow")
				} else if over {
					return 0, 0, 0, 0, 0, 0, fmt.Errorf("seconds overflow")
				}
				nsec = _nsec % 1e9

				value = 0
				haveUnit = true
			}
		}
	}

	if !haveUnit {
		return 0, 0, 0, 0, 0, 0, fmt.Errorf("expecting at least one unit")
	}

	return
}
