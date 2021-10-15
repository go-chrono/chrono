package chrono

import (
	"fmt"
	"math"
	"strconv"
)

func parseDuration(s string, periodAllowed bool) (secs int64, nsec uint32, err error) {
	if len(s) == 0 || s[0] != 'P' {
		return 0, 0, fmt.Errorf("expecting 'P'")
	}

	var value int
	var time bool

	for i := 1; i < len(s); i++ {
		digit := (s[i] >= '0' && s[i] <= '9') || s[i] == '.'

		if value == 0 {
			if digit {
				value = i
			} else if s[i] == 'T' {
				if !time {
					time = true
				} else {
					return 0, 0, fmt.Errorf("unexpected '%c', expecting digit", s[i])
				}
			} else {
				return 0, 0, fmt.Errorf("unexpected '%c', expecting digit or 'T'", s[i])
			}
		} else {
			if !time {
				if !periodAllowed {
					return 0, 0, fmt.Errorf("cannot parse duration as Duration")
				} else if digit {
					continue
				}

				switch s[i] {
				case 'Y':
					fmt.Printf("P%sY\n", s[value:i])
				case 'M':
					fmt.Printf("P%sM\n", s[value:i])
				case 'D':
					fmt.Printf("P%sD\n", s[value:i])
				default:
					return 0, 0, fmt.Errorf("unexpected '%c', expecting 'Y', 'M' or 'D'", s[i])
				}

				value = 0
			} else {
				if digit {
					continue
				}

				v, err := strconv.ParseFloat(s[value:i], 64)
				if err != nil {
					return 0, 0, err
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
					return 0, 0, fmt.Errorf("unexpected '%c', expecting 'H', 'M' or 'S'", s[i])
				}

				_secsInt := int64(_secs)

				// TODO check overflow
				secs += _secsInt
				nsec += _nsec

				if nsec >= 1e9 {
					secs++
					nsec -= 1e9
				}

				value = 0
			}
		}
	}

	return
}