package examples

import (
	"fmt"

	"github.com/go-chrono/chrono"
)

func ExampleLocalDateOf() {
	d := chrono.LocalDateOf(2007, chrono.May, 20)
	fmt.Println(d)
	// Output: 2007-05-20
}

func ExampleLocalDate_Weekday() {
	d := chrono.LocalDateOf(2007, chrono.May, 20)
	fmt.Println(d.Weekday())
	// Output: Sunday
}

func ExampleLocalDate_compare() {
	d1 := chrono.LocalDateOf(2007, chrono.May, 20)
	d2 := chrono.LocalDateOf(2009, chrono.June, 5)

	if d2 > d1 {
		fmt.Println(d2, "is after", d1)
	}
	// Output: 2009-06-05 is after 2007-05-20
}

func ExampleLocalDate_difference() {
	d1 := chrono.LocalDateOf(2007, chrono.May, 20)
	d2 := chrono.LocalDateOf(2007, chrono.May, 25)

	fmt.Printf("There are %d days from %s to %s\n", d2-d1, d1, d2)
	// Output: There are 5 days from 2007-05-20 to 2007-05-25
}

func ExampleLocalDate_add_subtract() {
	d := chrono.LocalDateOf(2007, chrono.May, 20)

	d += 8
	d -= 3

	fmt.Println(d)
	// Output: 2007-05-25
}