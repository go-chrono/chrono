package chrono

import "fmt"

const (
	oneHour   = int64(Hour)
	oneMinute = int64(Minute)
	oneSecond = int64(Second)

	maxTime int64 = 359999999999999
)

func makeTime(hour, min, sec, nsec int) (int64, error) {
	if hour < 0 || hour > 99 || min < 0 || min > 59 || sec < 0 || sec > 59 || nsec < 0 || nsec > 999999999 {
		return 0, fmt.Errorf("invalid time")
	}

	h, m, s, n := int64(hour), int64(min), int64(sec), int64(nsec)
	return h*oneHour + m*oneMinute + s*oneSecond + n, nil
}

func fromTime(v int64) (hour, min, sec, nsec int) {
	nsec = int(v) % int(oneSecond)
	sec = int(v) / int(oneSecond)

	hour = (sec / (60 * 60)) % 24
	sec -= hour * (60 * 60)

	min = sec / 60
	sec -= min * 60
	return
}

func addTime(t, v int64) (int64, error) {
	if v > maxTime {
		return 0, fmt.Errorf("invalid duration v")
	}

	out := t + v
	if out > int64(maxTime) {
		return 0, fmt.Errorf("invalid time t+v")
	}

	if out < 0 {
		return 24*oneHour + (out % (24 * oneHour)), nil
	}
	return out, nil
}

func simpleTimeStr(hour, min, sec, nsec int, offset *int64) string {
	out := fmt.Sprintf("%02d:%02d:%02d", hour, min, sec)
	if nsec != 0 {
		out += fmt.Sprintf(".%09d", nsec)
	}

	if offset == nil {
		return out
	}
	return out + offsetString(*offset, ":")
}

func timeBusinessHour(t int64) int {
	return int(t / oneHour)
}

func timeNanoseconds(t int64) int {
	return int(t % oneSecond)
}

func compareTimes(t, t2 int64) int {
	switch {
	case t < t2:
		return -1
	case t > t2:
		return 1
	default:
		return 0
	}
}
