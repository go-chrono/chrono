package chrono_test

import (
	"fmt"

	"github.com/go-chrono/chrono"
)

func ExampleLocalDateTimeOf() {
	dt := chrono.LocalDateTimeOf(2007, chrono.May, 20, 12, 30, 15, 0)

	fmt.Println(dt)
	// Output: 2007-05-20 12:30:15
}

func ExampleOfLocalDateTime() {
	d := chrono.LocalDateOf(2007, chrono.May, 20)
	t := chrono.LocalTimeOf(12, 30, 15, 0)

	dt := chrono.OfLocalDateTime(d, t)

	fmt.Println(dt)
	// Output: 2007-05-20 12:30:15
}

func ExampleLocalDateTime_Split() {
	dt := chrono.LocalDateTimeOf(2007, chrono.May, 20, 12, 30, 15, 0)
	d, t := dt.Split()

	fmt.Printf("date = %s, time = %s", d, t)
	// Output: date = 2007-05-20, time = 12:30:15
}

func ExampleLocalDateTime_Compare() {
	dt1 := chrono.LocalDateTimeOf(2007, chrono.May, 26, 12, 30, 15, 0)
	dt2 := chrono.LocalDateTimeOf(2007, chrono.May, 20, 12, 30, 15, 0)

	if dt2.Compare(dt1) == -1 {
		fmt.Println(dt2, "is before", dt1)
	}
	// Output: 2007-05-20 12:30:15 is before 2007-05-26 12:30:15
}

func ExampleLocalDateTime_Add() {
	dt := chrono.LocalDateTimeOf(2007, chrono.May, 20, 12, 30, 15, 0)

	fmt.Println(dt.Add(chrono.DurationOf(26 * chrono.Hour)))
	// Output: 2007-05-21 14:30:15
}

func ExampleLocalDateTime_AddDate() {
	dt := chrono.LocalDateTimeOf(2007, chrono.May, 20, 12, 30, 15, 0)

	fmt.Println(dt.AddDate(2, 6, 8))
	// Output: 2009-11-28 12:30:15
}

func ExampleLocalDateTime_Format() {
	dt := chrono.LocalDateTimeOf(2007, chrono.May, 20, 12, 30, 15, 0)

	fmt.Println(dt.Format(chrono.ISO8601DateTimeExtended))
	// Output: 2007-05-20T12:30:15
}

func ExampleLocalDateTime_Parse() {
	var dt chrono.LocalDateTime
	_ = dt.Parse(chrono.ISO8601DateTimeExtended, "2007-05-20T12:30:15")

	fmt.Println(dt)
	// Output: 2007-05-20 12:30:15
}
