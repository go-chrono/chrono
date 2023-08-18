package chrono_test

import (
	"fmt"

	"github.com/go-chrono/chrono"
)

func ExampleOffsetTimeOf() {
	t := chrono.OffsetTimeOf(12, 30, 15, 0, 2, 0)

	fmt.Println(t)
	// Output: 12:30:15+02:00
}

func ExampleOffsetTime_BusinessHour() {
	t := chrono.OffsetTimeOf(24, 15, 0, 0, 2, 0)

	fmt.Println(t.BusinessHour())
	// Output: 24
}

func ExampleOffsetTime_Sub() {
	t1 := chrono.OffsetTimeOf(12, 30, 0, 0, 1, 0)
	t2 := chrono.OffsetTimeOf(12, 15, 0, 0, 2, 0)

	fmt.Println(t1.Sub(t2))
	// Output: PT1H15M
}

func ExampleOffsetTime_Add() {
	t := chrono.OffsetTimeOf(12, 30, 0, 0, 2, 0)

	fmt.Println(t.Add(4 * chrono.Hour))
	// Output: 16:30:00+02:00
}

func ExampleOffsetTime_Compare() {
	t1 := chrono.OffsetTimeOf(12, 30, 0, 0, 1, 0)
	t2 := chrono.OffsetTimeOf(12, 30, 0, 0, 2, 0)

	if t2.Compare(t1) == -1 {
		fmt.Println(t2, "is before", t1)
	}
	// Output: 12:30:00+02:00 is before 12:30:00+01:00
}

func ExampleOffsetTime_Format() {
	t := chrono.OffsetTimeOf(12, 30, 15, 0, 2, 30)

	fmt.Println(t.Format(chrono.ISO8601TimeExtended))
	// Output: T12:30:15+02:30
}

func ExampleOffsetTime_Parse() {
	var t chrono.OffsetTime
	t.Parse(chrono.ISO8601TimeExtended, "T12:30:15+02:30")

	fmt.Println(t)
	// Output: 12:30:15+02:30
}
