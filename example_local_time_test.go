package chrono_test

import (
	"fmt"

	"github.com/go-chrono/chrono"
)

func ExampleLocalTimeOf() {
	t := chrono.LocalTimeOf(12, 30, 15, 0)

	fmt.Println(t)
	// Output: 12:30:15
}

func ExampleLocalTime_BusinessHour() {
	t := chrono.LocalTimeOf(24, 15, 0, 0)

	fmt.Println(t.BusinessHour())
	// Output: 24
}

func ExampleLocalTime_Sub() {
	t1 := chrono.LocalTimeOf(12, 30, 0, 0)
	t2 := chrono.LocalTimeOf(12, 15, 0, 0)

	fmt.Println(t1.Sub(t2))
	// Output: PT15M
}

func ExampleLocalTime_Add() {
	t := chrono.LocalTimeOf(12, 30, 0, 0)

	fmt.Println(t.Add(4 * chrono.Hour))
	// Output: 16:30:00
}

func ExampleLocalTime_Compare() {
	t1 := chrono.LocalTimeOf(12, 30, 0, 0)
	t2 := chrono.LocalTimeOf(12, 15, 0, 0)

	if t2.Compare(t1) == -1 {
		fmt.Println(t2, "is before", t1)
	}
	// Output: 12:15:00 is before 12:30:00
}

func ExampleLocalTime_Format() {
	t := chrono.LocalTimeOf(12, 30, 15, 0)

	fmt.Println(t.Format(chrono.ISO8601TimeExtended))
	// Output: T12:30:15
}
