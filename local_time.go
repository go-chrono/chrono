package chrono

import "fmt"

type LocalTime struct {
	v Extent
}

func LocalTimeOf(hour, min, sec, nsec int) LocalTime {
	if hour < 0 || hour > 99 || min < 0 || min > 59 || sec < 0 || sec > 59 || nsec < 0 || nsec > 999999999 {
		panic("invalid time")
	}
	return LocalTime{v: Extent(hour)*Hour + Extent(min)*Minute + Extent(sec)*Second + Extent(nsec)}
}

func (t LocalTime) Hour() int {
	return int(t.v / Hour)
}

func (t LocalTime) Minute() int {
	return int(t.v % Hour / Minute)
}

func (t LocalTime) Second() int {
	return int(t.v % Minute / Second)
}

func (t LocalTime) Nanosecond() int {
	return int(t.v % Second)
}

func (t LocalTime) String() string {
	out := fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())
	if nsec := t.Nanosecond(); nsec != 0 {
		out += fmt.Sprintf(".%09d", nsec)
	}
	return out
}
