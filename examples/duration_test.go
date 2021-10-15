package examples

import (
	"fmt"

	"github.com/go-chrono/chrono"
)

func ExampleDuration_Format() {
	d := chrono.DurationOf(1*chrono.Hour + 30*chrono.Minute + 5*chrono.Second)
	fmt.Println(d.Format())
	// Output: PT1H30M5S
}

func ExampleDuration_Parse() {
	var d chrono.Duration
	_ = d.Parse("PT1M5S")
	fmt.Println(d.Seconds())
	// Output: 65
}
