package chrono

type LocalTime int64

func LocalTimeOf(hour, min, sec, nsec int) LocalTime {
	return LocalTime(Extent(hour)*Hour+Extent(min)*Minute+Extent(sec)*Second) + LocalTime(nsec)
}

func (t LocalTime) Hour() int {
	return int(Extent(t) / Hour)
}

func (t LocalTime) Minute() int {
	return int(Extent(t) % Hour / Minute)
}

func (t LocalTime) Second() int {
	return int(Extent(t) % Minute / Second)
}

func (t LocalTime) Nanosecond() int {
	return int(Extent(t) % Second)
}
