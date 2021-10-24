package chrono

import "math"

// addInt64 attempts to add v1 to v2 but reports if the operation would underflow or overflow int64.
func addInt64(v1, v2 int64) (sum int64, underflows, overflows bool) {
	if v2 > 0 {
		v := math.MaxInt64 - v1
		if v < 0 {
			v = -v
		}

		if v < v2 {
			return 0, false, true
		}
	} else if v2 < 0 {
		v := math.MinInt64 + v1
		if v < 0 {
			v = -v
		}

		if -v > v2 { // v < -v2 can't be used because -math.MinInt64 > math.MaxInt64
			return 0, true, false
		}
	}
	return v1 + v2, false, false
}
