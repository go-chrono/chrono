package chrono

type ZonedDateTime struct {
	zone Zone
	secs int64
	nsec int32
}

func Current() ZonedDateTime {
	secs, nsec := walltime()
	return ZonedDateTime{
		zone: Local(),
		secs: secs,
		nsec: nsec,
	}
}
