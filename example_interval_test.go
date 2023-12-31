package chrono_test

import (
	"fmt"

	"github.com/go-chrono/chrono"
)

func ExampleParseInterval() {
	i, _ := chrono.ParseInterval("R5/2007-03-01T13:00:00Z/P1Y2M10DT2H30M")

	s, _ := i.Start()
	p, d, _ := i.Duration()
	e, _ := i.End()

	fmt.Printf("start: %v; duration: %v; end: %v; repetitions: %v", s, chrono.FormatDuration(p, d), e, i.Repetitions())
	// Output: start: 2007-03-01 13:00:00Z; duration: P1Y2M10DT2H30M; end: 2008-05-11 15:30:00Z; repetitions: 5
}
