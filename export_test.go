package chrono

func SetupCenturyParsing(v int) {
	overrideCentury = new(int)
	*overrideCentury = v
}

func TearDownCenturyParsing() {
	overrideCentury = nil
}

var DivideAndRoundIntFunc = divideAndRoundInt
