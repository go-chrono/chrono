[![Go Reference](https://pkg.go.dev/badge/github.com/go-chrono/chrono.svg)](https://pkg.go.dev/github.com/go-chrono/chrono)
[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://github.com/go-chrono/chrono/graphs/commit-activity)
[![GoReportCard example](https://goreportcard.com/badge/github.com/go-chrono/chrono)](https://goreportcard.com/report/github.com/go-chrono/chrono)

# `chrono` - supplementary time and date module

`chrono` provides additional functionality and improved ergonomics to complement the Go standard library's `time` package. It is not a replacement for, nor an extension of, the `time` package, but for certain use cases for which it was not explicitly designed to support, `chrono` can help to simplify and clarify.

`chrono` is also designed to look and feel like Go. Many of the ideas and much of the API is inspired by `time`, and should therefore feel familiar. That said, capable time and date libraries exist for most mainstream languages, and `chrono` has taken inspiration from several besides Go's `time` package, including Rust, Java and Python.

---

**Not all features are complete yet. See the [roadmap](https://github.com/go-chrono/chrono/projects/1) for the current state. If in doubt, [create an issue](https://github.com/go-chrono/chrono/issues) to ask a question or open a feature request.**

---

# Use cases

## Parse and format ISO 8601 durations

When interfacing with systems where the <code>time</code> package's duration formatting is not understood, ISO 8601 is a commonly-adopted standard.

`time` doesn't support ISO 8601 durations notation. A simple one-liner that uses only the seconds component is possible, but this is neither readable nor solves the problem of parsing such strings:

```go
var d time.Duration
fmt.Printf("PT%dS", int(d.Seconds()))
```

`chrono` supports both [parsing](https://pkg.go.dev/github.com/go-chrono/chrono#ParseDuration) and [formatting](https://pkg.go.dev/github.com/go-chrono/chrono#FormatDuration) of ISO 8601 strings:

```go
period, duration, _ := chrono.ParseDuration("P3Y6M4DT1M5S")
fmt.Println(chrono.FormatDuration(period, duration))
```

Alternatively, a [`Period`](https://pkg.go.dev/github.com/go-chrono/chrono#Period) and [`Duration`](https://pkg.go.dev/github.com/go-chrono/chrono#Duration) can be initialized with numeric values:

```go
period := chrono.Period{Years: 3, Months: 6, Days: 4}
duration := chrono.DurationOf(1*chrono.Hour + 30*chrono.Minute + 5*chrono.Second)
fmt.Println(chrono.FormatDuration(period, duration))
```

✅ [See more examples](examples/duration_period_test.go).

## Local (or "civil") dates and times

### Dates

Often it's necessary to represent a date with no time component, and no time zone or time offset. For example, you might want to represent a birthday - an event that happens on a particular date, outside of the context of a physical location or time zone.

Usually, this is achieved using the standard library's `time.Date` function and setting the time values to `0`. Alternatively, some people use the `Date` type of Google's [`civil`](https://pkg.go.dev/cloud.google.com/go/civil) package

`chrono` has a dedicated type to describe a local date. Since the date is represented as a single integer, a [`LocalDate`](https://pkg.go.dev/github.com/go-chrono/chrono#LocalDate) is sortable and comparable with built-in operators:

```go
chrono.LocalDateOf(2007, chrono.May, 20)
```

✅ [See more examples](examples/local_date_test.go).

### Times

Where is is necessary to represent a time with no zone or offset, and no date component, `chrono` provides the [`LocalTime`](https://pkg.go.dev/github.com/go-chrono/chrono#LocalTime) type:

```go
chrono.LocalTimeOf(12, 30, 0, 0)
```

✅ [See more examples](examples/local_time_test.go).
