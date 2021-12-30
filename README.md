[![Go Reference](https://pkg.go.dev/badge/github.com/go-chrono/chrono.svg)](https://pkg.go.dev/github.com/go-chrono/chrono)
[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://github.com/go-chrono/chrono/graphs/commit-activity)
[![GoReportCard example](https://goreportcard.com/badge/github.com/go-chrono/chrono)](https://goreportcard.com/report/github.com/go-chrono/chrono)

# `chrono` - supplementary time and date module

`chrono` provides additional functionality and improved ergonomics to complement the Go standard library's `time` package. It is not a replacement for, nor an extension of, the `time` package, but for certain use cases for which it was not explicitly designed to support, `chrono` can help to simplify and clarify.

`chrono` is also designed to look and feel like Go. Many of the ideas and much of the API is inspired by `time`, and should therefore feel familiar. That said, capable time and date libraries exist for most mainstream languages, and `chrono` has taken inspiration from several besides Go's `time` package, including Rust, Java and Python.

---

**Not all features are complete yet. See the [roadmap](https://github.com/go-chrono/chrono/projects/1) for the current state. If in doubt, [create an issue](https://github.com/go-chrono/chrono/issues) to ask a question or open a feature request.**

---

## Use cases

<table>
<tr>
<td style="width:33%">Use case</td>
<td style="width:33%">Using <code>time</code></td>
<td style="width:33%">Using <code>chrono</code></td>
</tr>
<tr>
<td style="vertical-align:top">Parse or format a duration according to ISO 8601. When interfacing with systems where the <code>time</code> package's duration formatting is not understood, ISO 8601 is a commonly-adopted standard.</td>
<td style="vertical-align:top">

`time` doesn't support ISO 8601 durations notation. A simple one-liner that uses only the seconds component is possible, but this is neither readable nor solves the problem of parsing such strings:

```go
var d time.Duration
fmt.Printf("PT%dS", int(d.Seconds()))
```

</td>
<td style="vertical-align:top">

`chrono` supports both formatting and parsing of ISO 8601 strings. Relevant functions:

* [`chrono.ParseDuration`](https://pkg.go.dev/github.com/go-chrono/chrono#ParseDuration)
* [`chrono.FormatDuration`](https://pkg.go.dev/github.com/go-chrono/chrono#FormatDuration)

✅ [See some examples](examples/duration_period_test.go).

</td>
</tr>
</table>

## Concepts

### Durations, extents, and periods

`Extent` is equivalent to `time.Duration`. It is an `int64` that represents a value in nanoseconds, and can therefore represent approximately ±292 years.

`chrono` also introduces the distinct `Duration` type, which is semantically equivalent to `Extent` but can represent a much larger period of time of approximately ±292 billion years. `Duration` also forms part of the support for ISO 8601 durations.

The other half of ISO 8601 duration support comes from the `Period` type. It consists of values for years, months, weeks and days, and has no equivalent in the standard library.

When combined, `Period` and `Duration`, an ISO 8601 duration can be formatted to, and parsed from, a string, e.g. `P3Y6M4DT1M5S`.
