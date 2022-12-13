// Package chronotest provides functionality useful for testing chrono.
// It should not be imported for normal usage of chrono.
//
package chronotest

import (
	_ "unsafe" // for go:linkname

	"github.com/go-chrono/chrono"
)

// InstantOf creates a new Instant with the supplied nanoseconds.
func InstantOf(t int64) chrono.Instant {
	return instant(t)
}

//go:linkname instant github.com/go-chrono/chrono.instant
func instant(int64) chrono.Instant
