package chrono

func SetupCenturyParsing(v int) {
	overrideCentury = new(int)
	*overrideCentury = v
}
