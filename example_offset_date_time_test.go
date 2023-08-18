package chrono_test

import (
	"fmt"

	"github.com/go-chrono/chrono"
)

func ExampleOffsetDateTimeOf() {
	dt := chrono.OffsetDateTimeOf(2007, chrono.May, 20, 12, 30, 15, 0, 2, 30)

	fmt.Println(dt)
	// Output: 2007-05-20 12:30:15+02:30
}

func ExampleOfLocalDateTimeOffset() {
	d := chrono.LocalDateOf(2007, chrono.May, 20)
	t := chrono.LocalTimeOf(12, 30, 15, 0)

	dt := chrono.OfLocalDateTimeOffset(d, t, 2*chrono.Hour+30*chrono.Minute)

	fmt.Println(dt)
	// Output: 2007-05-20 12:30:15+02:30
}

func ExampleOffsetDateTime_Split() {
	dt := chrono.OffsetDateTimeOf(2007, chrono.May, 20, 12, 30, 15, 0, 2, 30)
	d, t := dt.Split()

	fmt.Printf("date = %s, time = %s", d, t)
	// Output: date = 2007-05-20, time = 12:30:15+02:30
}

func ExampleOffsetDateTime_Compare() {
	dt1 := chrono.OffsetDateTimeOf(2007, chrono.May, 26, 12, 30, 15, 0, 2, 30)
	dt2 := chrono.OffsetDateTimeOf(2007, chrono.May, 20, 12, 30, 15, 0, 2, 30)

	if dt2.Compare(dt1) == -1 {
		fmt.Println(dt2, "is before", dt1)
	}
	// Output: 2007-05-20 12:30:15+02:30 is before 2007-05-26 12:30:15+02:30
}

func ExampleOffsetDateTime_Add() {
	dt := chrono.OffsetDateTimeOf(2007, chrono.May, 20, 12, 30, 15, 0, 2, 30)

	fmt.Println(dt.Add(chrono.DurationOf(26 * chrono.Hour)))
	// Output: 2007-05-21 14:30:15+02:30
}

func ExampleOffsetDateTime_AddDate() {
	dt := chrono.OffsetDateTimeOf(2007, chrono.May, 20, 12, 30, 15, 0, 2, 30)

	fmt.Println(dt.AddDate(2, 6, 8))
	// Output: 2009-11-28 12:30:15+02:30
}

func ExampleOffsetDateTime_Format() {
	dt := chrono.OffsetDateTimeOf(2007, chrono.May, 20, 12, 30, 15, 0, 2, 30)

	fmt.Println(dt.Format(chrono.ISO8601DateTimeExtended))
	// Output: 2007-05-20T12:30:15+02:30
}

func ExampleOffsetDateTime_Parse() {
	var dt chrono.OffsetDateTime
	_ = dt.Parse(chrono.ISO8601DateTimeExtended, "2007-05-20T12:30:15+02:30")

	fmt.Println(dt)
	// Output: 2007-05-20 12:30:15+02:30
}
