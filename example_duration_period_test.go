package chrono_test

import (
	"fmt"

	"github.com/go-chrono/chrono"
)

func ExampleParseDuration() {
	p, d, _ := chrono.ParseDuration("P3Y6M4DT1M5S")

	fmt.Println(p.Years, "years;", p.Months, "months;", p.Weeks, "weeks;", p.Days, "days;", d.Seconds(), "seconds")
	// Output: 3 years; 6 months; 0 weeks; 4 days; 65 seconds
}

func ExampleFormatDuration() {
	p := chrono.Period{Years: 3, Months: 6, Days: 4}
	d := chrono.DurationOf(1*chrono.Hour + 30*chrono.Minute + 5*chrono.Second)

	fmt.Println(chrono.FormatDuration(p, d))
	// Output: P3Y6M4DT1H30M5S
}

func ExamplePeriod_Parse() {
	var p chrono.Period
	_ = p.Parse("P3Y6M4D")

	fmt.Println(p.Years, "years;", p.Months, "months;", p.Weeks, "weeks;", p.Days, "days")
	// Output: 3 years; 6 months; 0 weeks; 4 days
}

func ExamplePeriod_String() {
	p := chrono.Period{Years: 3, Months: 6, Days: 4}

	fmt.Println(p.String())
	// Output: P3Y6M4D
}

func ExampleDuration_Parse() {
	var d chrono.Duration
	_ = d.Parse("PT1M5S")

	fmt.Println(d.Seconds())
	// Output: 65
}

func ExampleDuration_Format() {
	d := chrono.DurationOf(1*chrono.Hour + 30*chrono.Minute + 5*chrono.Second)

	fmt.Println(d.Format())
	// Output: PT1H30M5S
}
